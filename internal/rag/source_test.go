package rag

import (
	"context"
	"strings"
	"testing"
)

func TestParseAllowlistParsesValidSourceConfig(t *testing.T) {
	allowlist, err := ParseAllowlist([]byte(validAllowlistJSON("https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards")))
	if err != nil {
		t.Fatalf("ParseAllowlist returned error: %v", err)
	}

	if !allowlist.IsURLAllowlisted("https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards") {
		t.Fatal("transcript source URL was not allowlisted")
	}
	source, ok := allowlist.SourceByID("oc-transcript-request-2005-onwards")
	if !ok {
		t.Fatal("transcript source not found by ID")
	}
	if source.Title != "Transcript Request Guidance" || source.Department != "Registrar" {
		t.Fatalf("source metadata = %+v", source)
	}
	if source.RiskLevel != RiskHigh {
		t.Fatalf("risk level = %q, want %q", source.RiskLevel, RiskHigh)
	}
	if source.RetrievedAt.IsZero() || source.FreshnessStatus != FreshnessFresh {
		t.Fatalf("freshness metadata = retrieved_at %s status %q", source.RetrievedAt, source.FreshnessStatus)
	}
}

func TestParseAllowlistRejectsUnsafeSources(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "non HTTPS URL",
			body: validAllowlistJSON("http://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards"),
			want: "https",
		},
		{
			name: "missing title",
			body: strings.Replace(validAllowlistJSON("https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards"), `"title":"Transcript Request Guidance"`, `"title":""`, 1),
			want: "title",
		},
		{
			name: "missing department",
			body: strings.Replace(validAllowlistJSON("https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards"), `"department":"Registrar"`, `"department":""`, 1),
			want: "department",
		},
		{
			name: "private portal marker",
			body: strings.Replace(validAllowlistJSON("https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards"), `"private_portal":false`, `"private_portal":true`, 1),
			want: "private",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseAllowlist([]byte(tt.body))
			if err == nil {
				t.Fatal("ParseAllowlist returned nil error")
			}
			if !strings.Contains(strings.ToLower(err.Error()), tt.want) {
				t.Fatalf("error = %q, want it to mention %q", err.Error(), tt.want)
			}
		})
	}
}

func TestLoadSeedSourcesFixtureParsesP5FreshnessMetadata(t *testing.T) {
	allowlist, err := LoadAllowlist(context.Background(), "../../data/seed-sources.json")
	if err != nil {
		t.Fatalf("LoadAllowlist returned error: %v", err)
	}

	if len(allowlist.Sources) < 5 {
		t.Fatalf("sources = %d, want at least 5 approved public sources", len(allowlist.Sources))
	}
	for _, source := range allowlist.Sources {
		if source.RetrievedAt.IsZero() {
			t.Fatalf("source %s missing retrieved_at", source.ID)
		}
		if source.FreshnessStatus == "" {
			t.Fatalf("source %s missing freshness_status", source.ID)
		}
		if !allowlist.IsURLAllowlisted(source.URL) {
			t.Fatalf("source URL %q was not allowlisted", source.URL)
		}
	}
	if allowlist.IsURLAllowlisted("https://www.okanagancollege.ca/private/student-account") {
		t.Fatal("unlisted private-looking URL was allowlisted")
	}
}

func validAllowlistJSON(sourceURL string) string {
	return `{
  "fixture_name":"test-allowlist",
  "public_sources_only":true,
  "ingestion_policy":{
    "allow_private_portal_scraping":false,
    "allow_authenticated_pages":false,
    "allow_personal_account_pages":false,
    "allow_unlisted_urls":false,
    "stale_after_days_default":30
  },
  "sources":[
    {
      "id":"oc-transcript-request-2005-onwards",
      "title":"Transcript Request Guidance",
      "url":"` + sourceURL + `",
      "department":"Registrar",
      "content_type":"public_web_page",
      "allowlisted":true,
      "risk_level":"high",
      "requires_freshness_check":true,
      "stale_after_days":14,
      "retrieved_at":"2026-05-06",
      "freshness_status":"fresh",
      "knowledge_domains":["transcript_request_guidance"],
      "private_portal":false
    }
  ]
}`
}
