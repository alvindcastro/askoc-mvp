package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"askoc-mvp/internal/domain"
)

type ChatClient interface {
	Chat(context.Context, domain.ChatRequest) (domain.ChatResponse, error)
}

type Runner struct {
	Client  ChatClient
	Timeout time.Duration
}

func (r Runner) Run(ctx context.Context, cases []Case) (Report, error) {
	if r.Client == nil {
		return Report{}, fmt.Errorf("eval runner chat client is required")
	}
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	results := make([]CaseResult, 0, len(cases))
	for _, tc := range cases {
		req := domain.ChatRequest{
			ConversationID: "eval-" + tc.ID,
			Channel:        channelForCase(tc),
			Message:        tc.Prompt,
			StudentID:      tc.StudentID,
		}
		caseCtx, cancel := context.WithTimeout(ctx, timeout)
		start := time.Now()
		resp, err := r.Client.Chat(caseCtx, req)
		latency := time.Since(start)
		cancel()
		if err != nil && caseCtx.Err() == context.DeadlineExceeded {
			err = context.DeadlineExceeded
		}
		score := ScoreCase(tc, resp, latency, err)
		result := CaseResult{
			ID:        tc.ID,
			Prompt:    tc.Prompt,
			Critical:  tc.Critical,
			Response:  resp,
			LatencyMS: durationMillis(latency),
			Score:     score,
		}
		if err != nil {
			result.Error = safeError(err)
		}
		results = append(results, result)
	}
	return Report{
		Summary: BuildSummary(cases, results),
		Results: results,
	}, nil
}

type HTTPChatClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c HTTPChatClient) Chat(ctx context.Context, req domain.ChatRequest) (domain.ChatResponse, error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	endpoint := chatEndpoint(c.BaseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("encode chat request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("create chat request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("call chat API: %w", err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return domain.ChatResponse{}, fmt.Errorf("chat API returned HTTP %d", httpResp.StatusCode)
	}
	var resp domain.ChatResponse
	decoder := json.NewDecoder(httpResp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&resp); err != nil {
		return domain.ChatResponse{}, fmt.Errorf("decode chat response: %w", err)
	}
	return resp, nil
}

func chatEndpoint(baseURL string) string {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasSuffix(baseURL, "/chat") {
		return baseURL
	}
	return baseURL + "/chat"
}

func channelForCase(tc Case) string {
	if got := strings.TrimSpace(tc.Channel); got != "" {
		return got
	}
	return "eval"
}

func durationMillis(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	ms := d.Milliseconds()
	if ms == 0 {
		return 1
	}
	return ms
}
