package rag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type FreshnessStatus string

const (
	FreshnessFresh   FreshnessStatus = "fresh"
	FreshnessStale   FreshnessStatus = "stale"
	FreshnessUnknown FreshnessStatus = "unknown"
)

type Source struct {
	ID                     string
	Title                  string
	URL                    string
	Department             string
	ContentType            string
	Allowlisted            bool
	RiskLevel              RiskLevel
	RequiresFreshnessCheck bool
	StaleAfterDays         int
	RetrievedAt            time.Time
	FreshnessStatus        FreshnessStatus
	KnowledgeDomains       []string
	PrivatePortal          bool
}

type IngestionPolicy struct {
	AllowPrivatePortalScraping bool
	AllowAuthenticatedPages    bool
	AllowPersonalAccountPages  bool
	AllowUnlistedURLs          bool
	StaleAfterDaysDefault      int
}

type Allowlist struct {
	FixtureName       string
	PublicSourcesOnly bool
	SourceReviewedAt  time.Time
	IngestionPolicy   IngestionPolicy
	Sources           []Source

	byID  map[string]Source
	byURL map[string]Source
}

func LoadAllowlist(ctx context.Context, path string) (Allowlist, error) {
	if err := ctx.Err(); err != nil {
		return Allowlist{}, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return Allowlist{}, fmt.Errorf("load source allowlist: %w", err)
	}
	return ParseAllowlist(body)
}

func ParseAllowlist(body []byte) (Allowlist, error) {
	var raw rawAllowlist
	if err := json.Unmarshal(body, &raw); err != nil {
		return Allowlist{}, fmt.Errorf("parse source allowlist: %w", err)
	}

	reviewedAt, err := parseOptionalDate(raw.SourceReviewedAt)
	if err != nil {
		return Allowlist{}, fmt.Errorf("source_reviewed_at: %w", err)
	}
	allowlist := Allowlist{
		FixtureName:       strings.TrimSpace(raw.FixtureName),
		PublicSourcesOnly: raw.PublicSourcesOnly,
		SourceReviewedAt:  reviewedAt,
		IngestionPolicy: IngestionPolicy{
			AllowPrivatePortalScraping: raw.IngestionPolicy.AllowPrivatePortalScraping,
			AllowAuthenticatedPages:    raw.IngestionPolicy.AllowAuthenticatedPages,
			AllowPersonalAccountPages:  raw.IngestionPolicy.AllowPersonalAccountPages,
			AllowUnlistedURLs:          raw.IngestionPolicy.AllowUnlistedURLs,
			StaleAfterDaysDefault:      raw.IngestionPolicy.StaleAfterDaysDefault,
		},
		byID:  map[string]Source{},
		byURL: map[string]Source{},
	}
	if !allowlist.PublicSourcesOnly {
		return Allowlist{}, errors.New("source allowlist must be public_sources_only")
	}
	if allowlist.IngestionPolicy.AllowPrivatePortalScraping ||
		allowlist.IngestionPolicy.AllowAuthenticatedPages ||
		allowlist.IngestionPolicy.AllowPersonalAccountPages ||
		allowlist.IngestionPolicy.AllowUnlistedURLs {
		return Allowlist{}, errors.New("source allowlist policy must reject private, authenticated, personal, and unlisted URLs")
	}

	for i, rawSource := range raw.Sources {
		source, err := convertSource(rawSource, allowlist.IngestionPolicy.StaleAfterDaysDefault)
		if err != nil {
			return Allowlist{}, fmt.Errorf("source %d: %w", i, err)
		}
		if _, exists := allowlist.byID[source.ID]; exists {
			return Allowlist{}, fmt.Errorf("source %q is duplicated", source.ID)
		}
		normalizedURL := normalizeURL(source.URL)
		if _, exists := allowlist.byURL[normalizedURL]; exists {
			return Allowlist{}, fmt.Errorf("source URL %q is duplicated", source.URL)
		}
		allowlist.Sources = append(allowlist.Sources, source)
		allowlist.byID[source.ID] = source
		allowlist.byURL[normalizedURL] = source
	}
	return allowlist, nil
}

func (a Allowlist) IsURLAllowlisted(rawURL string) bool {
	_, ok := a.SourceByURL(rawURL)
	return ok
}

func (a Allowlist) SourceByID(id string) (Source, bool) {
	if a.byID != nil {
		source, ok := a.byID[strings.TrimSpace(id)]
		return source, ok
	}
	for _, source := range a.Sources {
		if source.ID == strings.TrimSpace(id) {
			return source, true
		}
	}
	return Source{}, false
}

func (a Allowlist) SourceByURL(rawURL string) (Source, bool) {
	normalizedURL := normalizeURL(rawURL)
	if a.byURL != nil {
		source, ok := a.byURL[normalizedURL]
		return source, ok
	}
	for _, source := range a.Sources {
		if normalizeURL(source.URL) == normalizedURL {
			return source, true
		}
	}
	return Source{}, false
}

type rawAllowlist struct {
	FixtureName       string `json:"fixture_name"`
	PublicSourcesOnly bool   `json:"public_sources_only"`
	SourceReviewedAt  string `json:"source_reviewed_at"`
	IngestionPolicy   struct {
		AllowPrivatePortalScraping bool `json:"allow_private_portal_scraping"`
		AllowAuthenticatedPages    bool `json:"allow_authenticated_pages"`
		AllowPersonalAccountPages  bool `json:"allow_personal_account_pages"`
		AllowUnlistedURLs          bool `json:"allow_unlisted_urls"`
		StaleAfterDaysDefault      int  `json:"stale_after_days_default"`
	} `json:"ingestion_policy"`
	Sources []rawSource `json:"sources"`
}

type rawSource struct {
	ID                     string   `json:"id"`
	Title                  string   `json:"title"`
	URL                    string   `json:"url"`
	Department             string   `json:"department"`
	ContentType            string   `json:"content_type"`
	Allowlisted            bool     `json:"allowlisted"`
	RiskLevel              string   `json:"risk_level"`
	RequiresFreshnessCheck bool     `json:"requires_freshness_check"`
	StaleAfterDays         int      `json:"stale_after_days"`
	RetrievedAt            string   `json:"retrieved_at"`
	FreshnessStatus        string   `json:"freshness_status"`
	KnowledgeDomains       []string `json:"knowledge_domains"`
	PrivatePortal          bool     `json:"private_portal"`
}

func convertSource(raw rawSource, defaultStaleAfterDays int) (Source, error) {
	source := Source{
		ID:                     strings.TrimSpace(raw.ID),
		Title:                  strings.TrimSpace(raw.Title),
		URL:                    strings.TrimSpace(raw.URL),
		Department:             strings.TrimSpace(raw.Department),
		ContentType:            strings.TrimSpace(raw.ContentType),
		Allowlisted:            raw.Allowlisted,
		RiskLevel:              RiskLevel(strings.TrimSpace(strings.ToLower(raw.RiskLevel))),
		RequiresFreshnessCheck: raw.RequiresFreshnessCheck,
		StaleAfterDays:         raw.StaleAfterDays,
		FreshnessStatus:        FreshnessStatus(strings.TrimSpace(strings.ToLower(raw.FreshnessStatus))),
		KnowledgeDomains:       raw.KnowledgeDomains,
		PrivatePortal:          raw.PrivatePortal,
	}
	if source.ID == "" {
		return Source{}, errors.New("id is required")
	}
	if source.Title == "" {
		return Source{}, errors.New("title is required")
	}
	if source.Department == "" {
		return Source{}, errors.New("department is required")
	}
	if !source.Allowlisted {
		return Source{}, fmt.Errorf("source %q is not allowlisted", source.ID)
	}
	if source.PrivatePortal {
		return Source{}, fmt.Errorf("source %q is marked private portal", source.ID)
	}
	parsedURL, err := url.Parse(source.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return Source{}, fmt.Errorf("url %q is invalid", source.URL)
	}
	if parsedURL.Scheme != "https" {
		return Source{}, fmt.Errorf("url %q must use https", source.URL)
	}
	switch source.RiskLevel {
	case RiskLow, RiskMedium, RiskHigh:
	default:
		return Source{}, fmt.Errorf("risk_level %q is unsupported", raw.RiskLevel)
	}
	if source.StaleAfterDays == 0 {
		source.StaleAfterDays = defaultStaleAfterDays
	}
	if source.StaleAfterDays <= 0 {
		return Source{}, errors.New("stale_after_days must be greater than zero")
	}
	retrievedAt, err := parseRequiredDate(raw.RetrievedAt)
	if err != nil {
		return Source{}, fmt.Errorf("retrieved_at: %w", err)
	}
	source.RetrievedAt = retrievedAt
	switch source.FreshnessStatus {
	case FreshnessFresh, FreshnessStale, FreshnessUnknown:
	default:
		return Source{}, fmt.Errorf("freshness_status %q is unsupported", raw.FreshnessStatus)
	}
	return source, nil
}

func parseRequiredDate(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, errors.New("date is required")
	}
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func parseOptionalDate(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func normalizeURL(rawURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return strings.TrimRight(strings.TrimSpace(rawURL), "/")
	}
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	parsed.Fragment = ""
	return strings.TrimRight(parsed.String(), "/")
}
