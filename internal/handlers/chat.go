package handlers

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/session"
	"askoc-mvp/internal/validation"
)

type ChatService interface {
	HandleChat(context.Context, domain.ChatRequest) (domain.ChatResponse, error)
}

func ChatHandler(service ChatService) http.Handler {
	if service == nil {
		service = NewPlaceholderChatService(nil)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		var req domain.ChatRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, validation.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			WriteError(w, r, http.StatusBadRequest, validation.CodeInvalidRequest, "request body must contain one JSON object")
			return
		}
		if err := validation.ValidateChatRequest(req); err != nil {
			WriteError(w, r, http.StatusBadRequest, validation.Code(err), validation.SafeMessage(err))
			return
		}

		resp, err := service.HandleChat(r.Context(), req)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "chat_unavailable", "unable to produce chat response")
			return
		}
		resp.TraceID = middleware.TraceIDFromContext(r.Context())
		WriteJSON(w, r, http.StatusOK, resp)
	})
}

type PlaceholderChatService struct {
	store *session.Store
}

func NewPlaceholderChatService(store *session.Store) *PlaceholderChatService {
	return &PlaceholderChatService{store: store}
}

func (s *PlaceholderChatService) HandleChat(ctx context.Context, req domain.ChatRequest) (domain.ChatResponse, error) {
	conversationID := conversationIDFor(req)
	intent, confidence := classifyPlaceholderIntent(req.Message)
	sentiment := domain.SentimentNeutral
	if containsAny(req.Message, "frustrating", "urgent", "today", "job application") {
		sentiment = domain.SentimentUrgent
	}

	resp := domain.ChatResponse{
		ConversationID: conversationID,
		Answer:         placeholderAnswer(intent),
		Intent: domain.IntentResult{
			Name:       intent,
			Confidence: confidence,
		},
		Sentiment: sentiment,
		Sources:   placeholderSources(intent),
		Actions: []domain.Action{
			{
				Type:    "placeholder_response",
				Status:  domain.ActionStatusCompleted,
				Message: "No live AI, retrieval, or enterprise system was called in this P2 placeholder.",
			},
		},
	}
	if intent == domain.IntentHumanHandoff || sentiment == domain.SentimentUrgent {
		resp.Escalation = &domain.Escalation{
			Required: true,
			Status:   domain.HandoffPending,
			Queue:    "demo_review",
			Priority: "normal",
			Reason:   "placeholder handoff metadata for later CRM routing",
		}
		if sentiment == domain.SentimentUrgent {
			resp.Escalation.Priority = "priority"
		}
	}

	if s.store != nil {
		if _, ok := s.store.Get(conversationID); !ok {
			if err := s.store.Create(conversationID); err != nil {
				return domain.ChatResponse{}, err
			}
		}
		if err := s.store.Append(conversationID, session.Message{Role: session.RoleUser, Content: req.Message}); err != nil {
			return domain.ChatResponse{}, err
		}
		if err := s.store.Append(conversationID, session.Message{Role: session.RoleAssistant, Content: resp.Answer}); err != nil {
			return domain.ChatResponse{}, err
		}
	}

	return resp, nil
}

func conversationIDFor(req domain.ChatRequest) string {
	if got := strings.TrimSpace(req.ConversationID); got != "" {
		return got
	}
	key := strings.TrimSpace(req.StudentID)
	if key == "" {
		key = strings.TrimSpace(req.Channel) + ":" + strings.TrimSpace(req.Message)
	}
	hash := sha1.Sum([]byte(key))
	return "conv_" + hex.EncodeToString(hash[:])[:12]
}

func classifyPlaceholderIntent(message string) (domain.Intent, float64) {
	switch {
	case containsAny(message, "person", "human", "staff", "advisor"):
		return domain.IntentHumanHandoff, 0.6
	case containsAny(message, "payment", "paid", "fee", "balance"):
		return domain.IntentFeePayment, 0.55
	case containsAny(message, "processed", "status", "arrived", "check") && containsAny(message, "transcript"):
		return domain.IntentTranscriptStatus, 0.6
	case containsAny(message, "order", "request") && containsAny(message, "transcript"):
		return domain.IntentTranscriptRequest, 0.6
	default:
		return domain.IntentUnknown, 0.4
	}
}

func placeholderAnswer(intent domain.Intent) string {
	switch intent {
	case domain.IntentTranscriptRequest:
		return "This P2 demo placeholder received your transcript request question. Source-grounded transcript guidance will be added in the RAG phase."
	case domain.IntentTranscriptStatus:
		return "This P2 demo placeholder received your transcript status question. Synthetic Banner and payment checks will be added in later phases."
	case domain.IntentFeePayment:
		return "This P2 demo placeholder received your payment question. Mock payment workflow actions will be added in later phases."
	case domain.IntentHumanHandoff:
		return "This P2 demo placeholder marked the conversation for staff handoff metadata without creating a real case."
	default:
		return "This P2 demo placeholder received your message and returned a deterministic response without calling live AI or enterprise systems."
	}
}

func placeholderSources(intent domain.Intent) []domain.Source {
	if intent != domain.IntentTranscriptRequest && intent != domain.IntentTranscriptStatus {
		return []domain.Source{}
	}
	return []domain.Source{
		{
			Title:   "Transcript Request - 2005 Onwards",
			URL:     "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID: "placeholder-transcript-source",
		},
	}
}

func containsAny(value string, terms ...string) bool {
	value = strings.ToLower(value)
	for _, term := range terms {
		if strings.Contains(value, term) {
			return true
		}
	}
	return false
}
