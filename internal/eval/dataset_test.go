package eval

import (
	"context"
	"os"
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
)

func TestParseDatasetHandlesValidJSONL(t *testing.T) {
	input := strings.NewReader(`
{"id":"D01","prompt":"How do I order my official transcript?","expected_intent":"transcript_request","expected_source_contains":"oc-transcript-request-2005-onwards","expected_actions":["grounded_answer_returned"],"expected_handoff":"none","critical":true}
{"id":"D02","prompt":"I ordered my transcript but it has not been processed. My student ID is S100002.","expected_intent":"transcript_status","expected_source_contains":"transcript","expected_actions":["banner_status_checked","payment_status_checked","payment_reminder_triggered"],"expected_handoff":"none","critical":true}
`)

	cases, err := ParseDataset(context.Background(), input)
	if err != nil {
		t.Fatalf("ParseDataset error = %v", err)
	}
	if len(cases) != 2 {
		t.Fatalf("len(cases) = %d, want 2", len(cases))
	}
	if cases[0].ID != "D01" || cases[0].ExpectedIntent != domain.IntentTranscriptRequest || !cases[0].Critical {
		t.Fatalf("first case = %+v, want parsed transcript request", cases[0])
	}
	if got := cases[1].ExpectedActions; len(got) != 3 || got[2] != "payment_reminder_triggered" {
		t.Fatalf("ExpectedActions = %+v, want workflow actions", got)
	}
}

func TestParseDatasetInvalidRowIncludesLineNumber(t *testing.T) {
	input := strings.NewReader("{\"id\":\"D01\",\"prompt\":\"ok\",\"expected_intent\":\"unknown\",\"expected_actions\":[\"safe_fallback\"]}\n{\"id\":")

	_, err := ParseDataset(context.Background(), input)
	if err == nil {
		t.Fatal("ParseDataset error = nil, want invalid JSON error")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("error = %q, want line number", err.Error())
	}
}

func TestParseDatasetValidatesRequiredExpectedFields(t *testing.T) {
	input := strings.NewReader(`{"id":"bad","prompt":"How do I register?","expected_actions":["safe_fallback"]}`)

	_, err := ParseDataset(context.Background(), input)
	if err == nil {
		t.Fatal("ParseDataset error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "expected_intent") {
		t.Fatalf("error = %q, want expected_intent validation", err.Error())
	}
}

func TestEvalDatasetFixtureIncludesThirtyCasesAndCriticalSafety(t *testing.T) {
	file, err := os.Open("../../data/eval-questions.jsonl")
	if err != nil {
		t.Fatalf("open eval dataset: %v", err)
	}
	defer file.Close()

	cases, err := ParseDataset(context.Background(), file)
	if err != nil {
		t.Fatalf("ParseDataset fixture error = %v", err)
	}
	if len(cases) < 30 {
		t.Fatalf("len(cases) = %d, want at least 30", len(cases))
	}

	var safetyCases int
	for _, tc := range cases {
		if tc.MustRefuse || tc.MustRedact != "" || tc.MustWarnPassword || tc.MustRefuseRecordAccess {
			safetyCases++
		}
		if len(tc.ExpectedActions) == 0 && tc.ExpectedSourceContains == "" && tc.ExpectedHandoff == "" && !tc.MustRefuse && tc.MustRedact == "" && !tc.MustWarnPassword && !tc.MustRefuseRecordAccess {
			t.Fatalf("case %s has no expected action/source/handoff/safety behavior", tc.ID)
		}
	}
	if safetyCases < 4 {
		t.Fatalf("safetyCases = %d, want at least 4", safetyCases)
	}
}
