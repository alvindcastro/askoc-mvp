package rag

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ErrorKind string

const (
	ErrorKindInvalidAllowlist ErrorKind = "invalid_allowlist"
	ErrorKindNotAllowlisted   ErrorKind = "not_allowlisted"
	ErrorKindFetchFailed      ErrorKind = "fetch_failed"
)

type Error struct {
	Kind    ErrorKind
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func IsErrorKind(err error, kind ErrorKind) bool {
	var ragErr *Error
	if !asError(err, &ragErr) {
		return false
	}
	return ragErr.Kind == kind
}

func asError(err error, target **Error) bool {
	for err != nil {
		if got, ok := err.(*Error); ok {
			*target = got
			return true
		}
		type unwrapper interface {
			Unwrap() error
		}
		next, ok := err.(unwrapper)
		if !ok {
			return false
		}
		err = next.Unwrap()
	}
	return false
}

type Document struct {
	Source      Source    `json:"source"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Text        string    `json:"text"`
	ContentHash string    `json:"content_hash"`
	RetrievedAt time.Time `json:"retrieved_at"`
}

type Fetcher struct {
	allowlist Allowlist
	client    *http.Client
}

func NewFetcher(allowlist Allowlist, client *http.Client) Fetcher {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return Fetcher{allowlist: allowlist, client: client}
}

func (f Fetcher) Fetch(ctx context.Context, rawURL string) (Document, error) {
	source, ok := f.allowlist.SourceByURL(rawURL)
	if !ok {
		return Document{}, &Error{
			Kind:    ErrorKindNotAllowlisted,
			Message: "source URL is not allowlisted for ingestion",
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return Document{}, &Error{
			Kind:    ErrorKindFetchFailed,
			Message: "source URL could not be requested",
			Err:     err,
		}
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return Document{}, &Error{
			Kind:    ErrorKindFetchFailed,
			Message: "allowlisted source fetch failed",
			Err:     err,
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Document{}, &Error{
			Kind:    ErrorKindFetchFailed,
			Message: fmt.Sprintf("allowlisted source returned HTTP %d", resp.StatusCode),
		}
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return Document{}, &Error{
			Kind:    ErrorKindFetchFailed,
			Message: "allowlisted source response could not be read",
			Err:     err,
		}
	}
	text := CleanHTML(string(body))
	hash := sha256.Sum256([]byte(text))
	return Document{
		Source:      source,
		Title:       source.Title,
		URL:         source.URL,
		Text:        text,
		ContentHash: hex.EncodeToString(hash[:]),
		RetrievedAt: time.Now().UTC(),
	}, nil
}

var (
	scriptBlock = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	styleBlock  = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	navBlock    = regexp.MustCompile(`(?is)<nav\b[^>]*>.*?</nav>`)
	tagPattern  = regexp.MustCompile(`(?is)<[^>]+>`)
	spaceRun    = regexp.MustCompile(`\s+`)
)

func CleanHTML(raw string) string {
	cleaned := scriptBlock.ReplaceAllString(raw, " ")
	cleaned = styleBlock.ReplaceAllString(cleaned, " ")
	cleaned = navBlock.ReplaceAllString(cleaned, " ")
	cleaned = tagPattern.ReplaceAllString(cleaned, " ")
	cleaned = html.UnescapeString(cleaned)
	cleaned = spaceRun.ReplaceAllString(cleaned, " ")
	return strings.TrimSpace(cleaned)
}
