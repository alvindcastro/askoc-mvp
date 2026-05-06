package orchestrator

import (
	"context"
	"strings"
	"testing"

	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
)

func TestGroundedTranscriptAnswerIncludesDedupedCitations(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentTranscriptRequest,
		Confidence: 0.92,
		Sentiment:  domain.SentimentNeutral,
	}
	deps.Retriever = &fakeRetriever{sources: []domain.Source{
		{
			Title:           "Transcript Request Guidance",
			URL:             "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:         "transcript-chunk-1",
			Confidence:      0.91,
			RiskLevel:       "high",
			FreshnessStatus: "fresh",
		},
		{
			Title:           "Transcript Request Guidance",
			URL:             "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:         "transcript-chunk-duplicate",
			Confidence:      0.88,
			RiskLevel:       "high",
			FreshnessStatus: "fresh",
		},
	}}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-grounded"), domain.ChatRequest{
		ConversationID: "conv-grounded",
		Channel:        "web",
		Message:        "How do I order my official transcript?",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if !strings.Contains(resp.Answer, "approved public source") {
		t.Fatalf("answer = %q, want grounded source wording", resp.Answer)
	}
	if len(resp.Sources) != 1 {
		t.Fatalf("sources = %+v, want duplicate source URLs de-duplicated", resp.Sources)
	}
	if resp.Sources[0].RiskLevel != "high" || resp.Sources[0].Confidence == 0 {
		t.Fatalf("source missing risk/confidence metadata: %+v", resp.Sources[0])
	}
	assertAction(t, resp, "rag_sources_retrieved", domain.ActionStatusCompleted)
}

func TestGroundedTranscriptAnswerLowConfidenceUsesSafeFallback(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentTranscriptRequest,
		Confidence: 0.88,
		Sentiment:  domain.SentimentNeutral,
	}
	deps.Retriever = &fakeRetriever{sources: []domain.Source{
		{
			Title:      "Transcript Request Guidance",
			URL:        "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:    "transcript-low-confidence",
			Confidence: 0.22,
			RiskLevel:  "high",
		},
	}}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-low-rag"), domain.ChatRequest{
		ConversationID: "conv-low-rag",
		Channel:        "web",
		Message:        "How do I order my official transcript?",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if !strings.Contains(resp.Answer, "verified information is unavailable") {
		t.Fatalf("answer = %q, want safe fallback", resp.Answer)
	}
	if len(resp.Sources) != 0 {
		t.Fatalf("sources = %+v, want low-confidence sources withheld", resp.Sources)
	}
	assertAction(t, resp, "rag_sources_retrieved", domain.ActionStatusPending)
}

func TestGroundedTranscriptAnswerStaleHighRiskSourceRequestsStaffConfirmation(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentTranscriptRequest,
		Confidence: 0.90,
		Sentiment:  domain.SentimentNeutral,
	}
	deps.Retriever = &fakeRetriever{sources: []domain.Source{
		{
			Title:           "Transcript Request Guidance",
			URL:             "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:         "transcript-stale-high-risk",
			Confidence:      0.72,
			RiskLevel:       "high",
			FreshnessStatus: "stale",
		},
	}}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-stale-risk"), domain.ChatRequest{
		ConversationID: "conv-stale-risk",
		Channel:        "web",
		Message:        "How do I order my official transcript?",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if !strings.Contains(resp.Answer, "staff confirmation") {
		t.Fatalf("answer = %q, want staff confirmation wording", resp.Answer)
	}
	if len(resp.Sources) != 1 || resp.Sources[0].Caution == "" {
		t.Fatalf("sources = %+v, want stale/high-risk caution metadata", resp.Sources)
	}
	assertAction(t, resp, "source_confirmation_required", domain.ActionStatusPending)
}

func TestTranscriptStatusAttachesGroundingSourcesWhenRetrieverHasApprovedChunk(t *testing.T) {
	deps := completeDeps()
	deps.Retriever = &fakeRetriever{sources: []domain.Source{
		{
			ID:              "oc-transcript-request-2005-onwards",
			Title:           "Transcript Request Guidance",
			URL:             "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:         "transcript-status-chunk",
			Confidence:      0.92,
			RiskLevel:       "high",
			FreshnessStatus: "fresh",
		},
	}}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-status-grounded"), domain.ChatRequest{
		ConversationID: "conv-status-grounded",
		Channel:        "web",
		Message:        "I ordered my transcript but it has not been processed.",
		StudentID:      "S100002",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(resp.Sources) != 1 || resp.Sources[0].ChunkID != "transcript-status-chunk" {
		t.Fatalf("sources = %+v, want transcript status citation", resp.Sources)
	}
	assertAction(t, resp, "rag_sources_retrieved", domain.ActionStatusCompleted)
	assertAction(t, resp, "payment_reminder_triggered", domain.ActionStatusCompleted)
}

type fakeRetriever struct {
	sources []domain.Source
	err     error
}

func (f *fakeRetriever) RetrieveSources(context.Context, string) ([]domain.Source, error) {
	return f.sources, f.err
}
