package eval

import (
	"context"
	"errors"
	"testing"
	"time"

	"askoc-mvp/internal/domain"
)

func TestRunnerCallsChatClientForEachCase(t *testing.T) {
	client := &recordingChatClient{
		response: domain.ChatResponse{
			Intent:  domain.IntentResult{Name: domain.IntentTranscriptRequest, Confidence: 0.92},
			Actions: []domain.Action{{Type: "grounded_answer_returned", Status: domain.ActionStatusCompleted}},
		},
	}
	runner := Runner{Client: client, Timeout: time.Second}
	cases := []Case{
		{ID: "D01", Prompt: "How do I order my transcript?", ExpectedIntent: domain.IntentTranscriptRequest, ExpectedActions: []string{"grounded_answer_returned"}},
		{ID: "D02", Prompt: "How do I request a transcript copy?", ExpectedIntent: domain.IntentTranscriptRequest, ExpectedActions: []string{"grounded_answer_returned"}},
	}

	report, err := runner.Run(context.Background(), cases)
	if err != nil {
		t.Fatalf("Run error = %v", err)
	}
	if len(client.prompts) != 2 {
		t.Fatalf("client prompts = %+v, want two calls", client.prompts)
	}
	if report.Summary.TotalCases != 2 || report.Summary.Passed != 2 {
		t.Fatalf("summary = %+v, want two passed cases", report.Summary)
	}
	for _, result := range report.Results {
		if result.LatencyMS <= 0 {
			t.Fatalf("LatencyMS = %d, want recorded latency", result.LatencyMS)
		}
	}
}

func TestRunnerCapturesTimeoutAsFailure(t *testing.T) {
	runner := Runner{
		Client:  timeoutChatClient{},
		Timeout: 10 * time.Millisecond,
	}
	cases := []Case{{
		ID:              "TIMEOUT",
		Prompt:          "How do I order my transcript?",
		ExpectedIntent:  domain.IntentTranscriptRequest,
		ExpectedActions: []string{"grounded_answer_returned"},
		Critical:        true,
	}}

	report, err := runner.Run(context.Background(), cases)
	if err != nil {
		t.Fatalf("Run error = %v", err)
	}
	if report.Summary.Passed != 0 || report.Summary.Failed != 1 {
		t.Fatalf("summary = %+v, want one failed timeout", report.Summary)
	}
	if report.Results[0].Error == "" || report.Results[0].Score.Passed {
		t.Fatalf("result = %+v, want captured timeout failure", report.Results[0])
	}
}

type recordingChatClient struct {
	prompts  []string
	response domain.ChatResponse
}

func (c *recordingChatClient) Chat(_ context.Context, req domain.ChatRequest) (domain.ChatResponse, error) {
	c.prompts = append(c.prompts, req.Message)
	return c.response, nil
}

type timeoutChatClient struct{}

func (timeoutChatClient) Chat(ctx context.Context, _ domain.ChatRequest) (domain.ChatResponse, error) {
	<-ctx.Done()
	return domain.ChatResponse{}, ctx.Err()
}

var errFakeChat = errors.New("fake chat error")
