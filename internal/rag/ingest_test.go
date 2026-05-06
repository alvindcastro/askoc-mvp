package rag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetcherRejectsURLNotInAllowlist(t *testing.T) {
	allowlist := mustParseAllowlist(t, validAllowlistJSON("https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards"))
	fetcher := NewFetcher(allowlist, &http.Client{Timeout: time.Second})

	_, err := fetcher.Fetch(context.Background(), "https://www.okanagancollege.ca/private/student-account")
	if !IsErrorKind(err, ErrorKindNotAllowlisted) {
		t.Fatalf("Fetch error = %v, want not-allowlisted kind", err)
	}
}

func TestFetcherCleansHTMLAndStoresContentHash(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html>
			<head><style>.hidden{display:none}</style><script>console.log("secret")</script></head>
			<body><nav>portal navigation</nav><main><h1>Transcript Request Guidance</h1><p>Order official transcripts through the approved public guidance page.</p></main></body>
		</html>`))
	}))
	defer server.Close()

	allowlist := mustParseAllowlist(t, validAllowlistJSON(server.URL+"/transcript"))
	fetcher := NewFetcher(allowlist, server.Client())

	doc, err := fetcher.Fetch(context.Background(), server.URL+"/transcript")
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if doc.Source.ID != "oc-transcript-request-2005-onwards" || doc.Title != "Transcript Request Guidance" {
		t.Fatalf("document metadata = %+v", doc)
	}
	if doc.ContentHash == "" {
		t.Fatal("ContentHash was empty")
	}
	if !strings.Contains(doc.Text, "Order official transcripts") {
		t.Fatalf("cleaned text = %q, want main content", doc.Text)
	}
	if strings.Contains(doc.Text, "portal navigation") || strings.Contains(doc.Text, "console.log") || strings.Contains(doc.Text, "display:none") {
		t.Fatalf("cleaned text retained removed HTML regions: %q", doc.Text)
	}
}

func TestFetcherNetworkFailureReturnsTypedError(t *testing.T) {
	allowlist := mustParseAllowlist(t, validAllowlistJSON("https://127.0.0.1:1/transcript"))
	fetcher := NewFetcher(allowlist, &http.Client{Timeout: 10 * time.Millisecond})

	_, err := fetcher.Fetch(context.Background(), "https://127.0.0.1:1/transcript")
	if !IsErrorKind(err, ErrorKindFetchFailed) {
		t.Fatalf("Fetch error = %v, want fetch-failed kind", err)
	}
}

func mustParseAllowlist(t *testing.T, body string) Allowlist {
	t.Helper()
	allowlist, err := ParseAllowlist([]byte(body))
	if err != nil {
		t.Fatalf("ParseAllowlist returned error: %v", err)
	}
	return allowlist
}
