package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"askoc-mvp/internal/middleware"
)

func TestBannerClientGetsStudentAndSendsTraceID(t *testing.T) {
	var gotTrace string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTrace = r.Header.Get(middleware.TraceHeader)
		if r.URL.Path != "/api/v1/students/S100002" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		writeJSON(t, w, BannerStudentProfile{
			StudentID:     "S100002",
			PreferredName: "Demo Learner Two",
			Status:        "active",
			Program:       "Business Administration Demo Program",
			Synthetic:     true,
		})
	}))
	defer server.Close()

	client := NewBannerClient(server.URL, server.Client())
	ctx := middleware.WithTraceID(context.Background(), "trace-banner-client")

	got, err := client.GetStudent(ctx, "S100002")
	if err != nil {
		t.Fatalf("GetStudent() error = %v", err)
	}
	if got.StudentID != "S100002" || !got.Synthetic {
		t.Fatalf("profile = %+v", got)
	}
	if gotTrace != "trace-banner-client" {
		t.Fatalf("trace header = %q, want trace-banner-client", gotTrace)
	}
}

func TestBannerClientMapsNotFoundAndMalformedJSON(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		wantKind  ErrorKind
		callTrans bool
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"error":{"code":"student_not_found","message":"missing"}}`, wantKind: KindNotFound},
		{name: "malformed", status: http.StatusOK, body: `{`, wantKind: KindParse},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewBannerClient(server.URL, server.Client())
			_, err := client.GetTranscriptStatus(context.Background(), "S100999")
			if !IsKind(err, tt.wantKind) {
				t.Fatalf("error = %v, want kind %s", err, tt.wantKind)
			}
		})
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, payload any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("encode JSON: %v", err)
	}
}
