package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PowerAutomateClientConfig struct {
	WebhookURL      string
	HTTPClient      *http.Client
	Signature       string
	SignatureHeader string
	MaxRetries      int
}

type PowerAutomateClient struct {
	webhookURL      string
	httpClient      *http.Client
	signature       string
	signatureHeader string
	maxRetries      int
}

func NewPowerAutomateClient(cfg PowerAutomateClientConfig) (*PowerAutomateClient, error) {
	if strings.TrimSpace(cfg.WebhookURL) == "" {
		return nil, errors.New("workflow webhook URL is required for Power Automate client")
	}
	client := cfg.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	signatureHeader := strings.TrimSpace(cfg.SignatureHeader)
	if signatureHeader == "" {
		signatureHeader = "X-AskOC-Workflow-Signature"
	}
	maxRetries := cfg.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}
	return &PowerAutomateClient{
		webhookURL:      strings.TrimSpace(cfg.WebhookURL),
		httpClient:      client,
		signature:       cfg.Signature,
		signatureHeader: signatureHeader,
		maxRetries:      maxRetries,
	}, nil
}

func (c *PowerAutomateClient) SendPaymentReminder(ctx context.Context, req PaymentReminderRequest) (PaymentReminderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return PaymentReminderResponse{}, fmt.Errorf("encode workflow request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		attemptCount := attempt + 1
		if err := ctx.Err(); err != nil {
			return PaymentReminderResponse{}, err
		}
		resp, err := c.sendOnce(ctx, req, body)
		if err == nil {
			resp.AttemptCount = attemptCount
			return resp, nil
		}
		lastErr = err
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return PaymentReminderResponse{}, err
		}
		var statusErr *StatusError
		if errors.As(err, &statusErr) && statusErr.Transient {
			statusErr.AttemptCount = attemptCount
			continue
		}
		if errors.As(err, &statusErr) {
			statusErr.AttemptCount = attemptCount
		}
		return PaymentReminderResponse{}, err
	}
	return PaymentReminderResponse{}, lastErr
}

func (c *PowerAutomateClient) sendOnce(ctx context.Context, reminder PaymentReminderRequest, body []byte) (PaymentReminderResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return PaymentReminderResponse{}, errors.New("create workflow request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(reminder.TraceID) != "" {
		httpReq.Header.Set("X-Trace-ID", reminder.TraceID)
	}
	if strings.TrimSpace(reminder.IdempotencyKey) != "" {
		httpReq.Header.Set("Idempotency-Key", reminder.IdempotencyKey)
	}
	if c.signature != "" {
		httpReq.Header.Set(c.signatureHeader, c.signature)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return PaymentReminderResponse{}, ctxErr
		}
		return PaymentReminderResponse{}, errors.New("workflow webhook request failed")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		io.Copy(io.Discard, httpResp.Body)
		return PaymentReminderResponse{}, &StatusError{
			StatusCode: httpResp.StatusCode,
			Transient:  httpResp.StatusCode >= 500,
		}
	}

	var workflowResp PaymentReminderResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&workflowResp); err != nil {
		return PaymentReminderResponse{}, errors.New("decode workflow response")
	}
	return workflowResp, nil
}

type StatusError struct {
	StatusCode   int
	Transient    bool
	AttemptCount int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("workflow webhook returned HTTP %d", e.StatusCode)
}
