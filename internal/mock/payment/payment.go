package payment

import (
	"net/http"
	"strings"

	"askoc-mvp/internal/fixtures"
	"askoc-mvp/internal/handlers"
)

type Handler struct {
	fixture *fixtures.Fixture
}

type PaymentStatus struct {
	StudentID     string  `json:"student_id"`
	Item          string  `json:"item"`
	AmountDue     float64 `json:"amount_due"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id"`
	Synthetic     bool    `json:"synthetic"`
}

func NewHandler(fixture *fixtures.Fixture) http.Handler {
	return &Handler{fixture: fixture}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	studentID, ok := parsePaymentPath(r.URL.Path)
	if !ok {
		handlers.WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
		return
	}
	student, found := h.fixture.StudentByID(studentID)
	if !found {
		handlers.WriteError(w, r, http.StatusNotFound, "payment_not_found", "No synthetic payment record was found for that ID.")
		return
	}

	handlers.WriteJSON(w, r, http.StatusOK, PaymentStatus{
		StudentID:     student.StudentID,
		Item:          "official_transcript",
		AmountDue:     student.Payment.AmountDue,
		Currency:      student.Payment.Currency,
		Status:        student.Payment.Status,
		TransactionID: transactionIDValue(student.Payment.TransactionID),
		Synthetic:     true,
	})
}

func transactionIDValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func parsePaymentPath(path string) (string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 5 {
		return "", false
	}
	if parts[0] != "api" || parts[1] != "v1" || parts[2] != "students" || parts[3] == "" || parts[4] != "payment-status" {
		return "", false
	}
	return parts[3], true
}
