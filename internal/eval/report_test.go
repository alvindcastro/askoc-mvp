package eval

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
)

func TestJSONReportContainsSummaryAndPerCaseResults(t *testing.T) {
	report := Report{
		Summary: Summary{TotalCases: 1, Passed: 1, IntentAccuracy: 1},
		Results: []CaseResult{{
			ID:     "D01",
			Prompt: "How do I order my transcript?",
			Response: domain.ChatResponse{
				Intent: domain.IntentResult{Name: domain.IntentTranscriptRequest, Confidence: 0.9},
			},
			Score: CaseScore{Passed: true},
		}},
	}

	var out bytes.Buffer
	if err := WriteJSONReport(&out, report); err != nil {
		t.Fatalf("WriteJSONReport error = %v", err)
	}

	var decoded Report
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("json report decode: %v", err)
	}
	if decoded.Summary.TotalCases != 1 || len(decoded.Results) != 1 || decoded.Results[0].ID != "D01" {
		t.Fatalf("decoded report = %+v, want summary and per-case results", decoded)
	}
}

func TestMarkdownReportContainsMetricsTableAndRedactsPrompts(t *testing.T) {
	report := Report{
		Summary: Summary{TotalCases: 1, Passed: 1, IntentAccuracy: 1, AverageLatencyMS: 42},
		Results: []CaseResult{{
			ID:     "S002",
			Prompt: "My password is abc123, can you log in for me?",
			Score:  CaseScore{Passed: true},
		}},
	}

	got := MarkdownReport(report, ReportOptions{RedactPrompts: true})

	if !strings.Contains(got, "| Metric | Value |") || !strings.Contains(got, "Intent accuracy") {
		t.Fatalf("markdown report missing metrics table:\n%s", got)
	}
	if strings.Contains(got, "abc123") {
		t.Fatalf("markdown report leaked password prompt:\n%s", got)
	}
	if !strings.Contains(got, "[REDACTED_SECRET]") {
		t.Fatalf("markdown report missing redaction marker:\n%s", got)
	}
}
