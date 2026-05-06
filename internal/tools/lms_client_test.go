package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLMSClientGetsAccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("course_id"); got != "DEMO-LMS-101" {
			t.Fatalf("course_id = %q", got)
		}
		writeJSON(t, w, LMSAccessStatus{
			StudentID:       "S100001",
			AccountStatus:   "active",
			CourseID:        "DEMO-LMS-101",
			CourseName:      "Online Learning Orientation",
			AccessStatus:    "available",
			Synthetic:       true,
			ContentIncluded: false,
		})
	}))
	defer server.Close()

	client := NewLMSClient(server.URL, server.Client())
	got, err := client.GetAccessStatus(context.Background(), "S100001", "DEMO-LMS-101")
	if err != nil {
		t.Fatalf("GetAccessStatus() error = %v", err)
	}
	if got.AccessStatus != "available" || got.ContentIncluded {
		t.Fatalf("LMS access = %+v", got)
	}
}

func TestLMSClientMapsCanceledContextToTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("server should not be called for canceled context")
	}))
	defer server.Close()

	client := NewLMSClient(server.URL, server.Client())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetAccessStatus(ctx, "S100001", "DEMO-LMS-101")
	if !IsKind(err, KindTimeout) {
		t.Fatalf("error = %v, want timeout kind", err)
	}
}
