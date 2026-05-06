package payment

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"askoc-mvp/internal/fixtures"
)

func TestHandlerReturnsPaidAndUnpaidTranscriptPaymentStatus(t *testing.T) {
	handler := testHandler(t)

	tests := []struct {
		id            string
		status        string
		amountDue     float64
		transactionID string
	}{
		{id: "S100001", status: "paid", amountDue: 0, transactionID: "SYNTH-PAY-100001"},
		{id: "S100002", status: "unpaid", amountDue: 15, transactionID: "SYNTH-PAY-100002"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/students/"+tt.id+"/payment-status", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			var got PaymentStatus
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode payment status: %v", err)
			}
			if got.StudentID != tt.id || got.Status != tt.status || got.AmountDue != tt.amountDue || got.Currency != "CAD" || got.TransactionID != tt.transactionID {
				t.Fatalf("payment status = %+v", got)
			}
		})
	}
}

func TestHandlerReturnsSafeNotFoundForUnknownPayment(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S999999/payment-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	assertErrorCode(t, rec.Body.Bytes(), "payment_not_found")
}

func TestPaymentStatusContractIncludesSyntheticTransactionOnly(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S100002/payment-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode payment response: %v", err)
	}
	for _, field := range []string{"student_id", "item", "amount_due", "currency", "status", "transaction_id", "synthetic"} {
		if _, ok := got[field]; !ok {
			t.Fatalf("field %q missing from response %v", field, got)
		}
	}
	if got["transaction_id"] != "SYNTH-PAY-100002" {
		t.Fatalf("transaction_id = %v, want synthetic ID", got["transaction_id"])
	}
	if _, ok := got["card_number"]; ok {
		t.Fatalf("payment response exposed card_number: %v", got)
	}
}

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	records, err := fixtures.Load(context.Background(), "../../../data/synthetic-students.json")
	if err != nil {
		t.Fatalf("load fixtures: %v", err)
	}
	return NewHandler(records)
}

func assertErrorCode(t *testing.T, body []byte, want string) {
	t.Helper()
	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode error response: %v; body=%s", err, body)
	}
	if got.Error.Code != want {
		t.Fatalf("error code = %q, want %q; body=%s", got.Error.Code, want, body)
	}
}
