package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/middleware"
)

type fakeChatService struct {
	resp domain.ChatResponse
	err  error
	req  domain.ChatRequest
}

func (f *fakeChatService) HandleChat(ctx context.Context, req domain.ChatRequest) (domain.ChatResponse, error) {
	f.req = req
	if f.err != nil {
		return domain.ChatResponse{}, f.err
	}
	if f.resp.ConversationID == "" {
		f.resp.ConversationID = req.ConversationID
	}
	return f.resp, nil
}

func TestChatHandlerReturnsDeterministicResponseWithTraceID(t *testing.T) {
	service := &fakeChatService{
		resp: domain.ChatResponse{
			ConversationID: "conv_existing",
			Answer:         "deterministic placeholder",
			Intent: domain.IntentResult{
				Name:       domain.IntentTranscriptStatus,
				Confidence: 0.5,
			},
			Sentiment: domain.SentimentNeutral,
			Actions: []domain.Action{
				{Type: "placeholder_response", Status: domain.ActionStatusCompleted},
			},
		},
	}
	handler := ChatHandler(service)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(`{
		"conversation_id":"conv_existing",
		"channel":"web",
		"message":"I ordered my transcript but it has not been processed.",
		"student_id":"S100002"
	}`))
	req = req.WithContext(middleware.WithTraceID(req.Context(), "trace-chat"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var got domain.ChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode chat response: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got.TraceID != "trace-chat" {
		t.Fatalf("trace_id = %q, want trace-chat", got.TraceID)
	}
	if service.req.StudentID != "S100002" || service.req.Message == "" {
		t.Fatalf("service request = %+v", service.req)
	}
}

func TestChatHandlerRejectsInvalidJSON(t *testing.T) {
	handler := ChatHandler(&fakeChatService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(`{"message":`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	assertErrorCode(t, rec.Body.Bytes(), "invalid_request")
}

func TestChatHandlerRejectsInvalidRequestsSafely(t *testing.T) {
	tests := []struct {
		name string
		body string
		code string
	}{
		{name: "missing message", body: `{"channel":"web","student_id":"S100002"}`, code: "missing_message"},
		{name: "whitespace message", body: `{"channel":"web","message":"   ","student_id":"S100002"}`, code: "missing_message"},
		{name: "oversized message", body: `{"channel":"web","message":"` + strings.Repeat("a", 2001) + `"}`, code: "message_too_large"},
		{name: "invalid student ID", body: `{"channel":"web","message":"Please help","student_id":"123456"}`, code: "invalid_student_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ChatHandler(&fakeChatService{})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(tt.body))
			req = req.WithContext(middleware.WithTraceID(req.Context(), "trace-validation"))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			assertErrorCode(t, rec.Body.Bytes(), tt.code)
			if strings.Contains(rec.Body.String(), "123456") || strings.Contains(rec.Body.String(), strings.Repeat("a", 80)) {
				t.Fatalf("error response echoed raw request body: %s", rec.Body.String())
			}
		})
	}
}

func TestChatHandlerRejectsUnsupportedMethod(t *testing.T) {
	handler := ChatHandler(&fakeChatService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	assertErrorCode(t, rec.Body.Bytes(), "method_not_allowed")
}

func TestChatHandlerReturnsSafeErrorWhenServiceFails(t *testing.T) {
	handler := ChatHandler(&fakeChatService{err: errors.New("private orchestrator detail")})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(`{"channel":"web","message":"Please help"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), "private orchestrator detail") {
		t.Fatalf("response leaked service error: %s", rec.Body.String())
	}
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
