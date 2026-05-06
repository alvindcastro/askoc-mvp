package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaymentClientGetsPaymentStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/students/S100002/payment-status" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		writeJSON(t, w, PaymentStatus{
			StudentID:     "S100002",
			Item:          "official_transcript",
			AmountDue:     15,
			Currency:      "CAD",
			Status:        "unpaid",
			TransactionID: "SYNTH-PAY-100002",
			Synthetic:     true,
		})
	}))
	defer server.Close()

	client := NewPaymentClient(server.URL, server.Client())
	got, err := client.GetPaymentStatus(context.Background(), "S100002")
	if err != nil {
		t.Fatalf("GetPaymentStatus() error = %v", err)
	}
	if got.Status != "unpaid" || got.TransactionID != "SYNTH-PAY-100002" {
		t.Fatalf("payment status = %+v", got)
	}
}

func TestPaymentClientMapsCanceledContextToTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("server should not be called for canceled context")
	}))
	defer server.Close()

	client := NewPaymentClient(server.URL, server.Client())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetPaymentStatus(ctx, "S100002")
	if !IsKind(err, KindTimeout) {
		t.Fatalf("error = %v, want timeout kind", err)
	}
}
