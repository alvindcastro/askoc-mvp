package tools

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type PaymentClient struct {
	baseURL    string
	httpClient *http.Client
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

func NewPaymentClient(baseURL string, httpClient *http.Client) *PaymentClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &PaymentClient{baseURL: strings.TrimRight(baseURL, "/"), httpClient: httpClient}
}

func (c *PaymentClient) GetPaymentStatus(ctx context.Context, studentID string) (PaymentStatus, error) {
	var out PaymentStatus
	err := doJSON(ctx, c.httpClient, "payment", http.MethodGet, c.baseURL+"/api/v1/students/"+url.PathEscape(studentID)+"/payment-status", nil, &out)
	return out, err
}
