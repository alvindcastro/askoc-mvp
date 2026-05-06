package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCRMClientCreatesCase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/crm/cases" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		writeJSON(t, w, CRMCaseResponse{
			CaseID:   "MOCK-CRM-000001",
			Status:   "created",
			Queue:    "registrar_student_accounts",
			Priority: "priority",
			Summary:  "redacted summary",
		})
	}))
	defer server.Close()

	client := NewCRMClient(server.URL, server.Client())
	got, err := client.CreateCase(context.Background(), CRMCaseRequest{
		StudentID:      "S100003",
		ConversationID: "conv_123",
		Intent:         "transcript_status",
		Queue:          "registrar_student_accounts",
		Priority:       "priority",
		Summary:        "Learner needs staff follow-up.",
		SourceTraceID:  "trace-crm-client",
	})
	if err != nil {
		t.Fatalf("CreateCase() error = %v", err)
	}
	if got.CaseID != "MOCK-CRM-000001" || got.Priority != "priority" {
		t.Fatalf("case response = %+v", got)
	}
}

func TestCRMClientMaps5xxToRetryableError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "temporary backend failure", http.StatusBadGateway)
	}))
	defer server.Close()

	client := NewCRMClient(server.URL, server.Client())
	_, err := client.CreateCase(context.Background(), CRMCaseRequest{
		StudentID:      "S100003",
		ConversationID: "conv_123",
		Intent:         "transcript_status",
		Queue:          "registrar_student_accounts",
		Priority:       "normal",
		Summary:        "Learner needs staff follow-up.",
	})
	if !IsKind(err, KindRetryable) {
		t.Fatalf("error = %v, want retryable kind", err)
	}
}
