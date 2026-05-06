package orchestrator

import (
	"context"
	"strings"

	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/rag"
)

type groundingStatus string

const (
	groundingReady       groundingStatus = "ready"
	groundingUnavailable groundingStatus = "unavailable"
	groundingCaution     groundingStatus = "caution"
)

func (o *Orchestrator) handleGroundedAnswer(ctx context.Context, req domain.ChatRequest, resp domain.ChatResponse) domain.ChatResponse {
	sources, status := o.groundedSources(ctx, req.Message)
	switch status {
	case groundingReady:
		resp.Sources = sources
		resp.Actions = append(resp.Actions, o.action(ctx, "rag_sources_retrieved", domain.ActionStatusCompleted, "Approved public source chunks met the retrieval confidence threshold.", ""))
		if answer, ok := o.generateGroundedAnswer(ctx, req.Message, sources); ok {
			resp.Answer = answer
			resp.Actions = append(resp.Actions, o.action(ctx, "llm_answer_generated", domain.ActionStatusCompleted, "Guarded LLM answer generated from approved sources.", ""))
			return resp
		}
		resp.Answer = "I found approved public source guidance for this demo. Use the linked approved public source for official transcript request steps, and ask staff for account-specific, fee, deadline, or eligibility details."
	case groundingCaution:
		resp.Sources = sources
		resp.Answer = "I found an approved public source, but this transcript guidance needs staff confirmation before I present it as authoritative."
		resp.Actions = append(resp.Actions, o.action(ctx, "source_confirmation_required", domain.ActionStatusPending, "Retrieved source is stale or high-risk for the current confidence level.", ""))
	default:
		resp.Answer = "verified information is unavailable in the demo knowledge base for that policy question. I can help route this to staff instead of guessing."
		resp.Actions = append(resp.Actions, o.action(ctx, "rag_sources_retrieved", domain.ActionStatusPending, "No approved source chunk met the retrieval confidence threshold.", ""))
	}
	return resp
}

func (o *Orchestrator) attachGroundingIfAvailable(ctx context.Context, req domain.ChatRequest, resp domain.ChatResponse) domain.ChatResponse {
	sources, status := o.groundedSources(ctx, req.Message)
	if len(sources) == 0 {
		return resp
	}
	resp.Sources = sources
	switch status {
	case groundingCaution:
		resp.Actions = append(resp.Actions, o.action(ctx, "source_confirmation_required", domain.ActionStatusPending, "Retrieved source is stale or high-risk for the current confidence level.", ""))
	default:
		resp.Actions = append(resp.Actions, o.action(ctx, "rag_sources_retrieved", domain.ActionStatusCompleted, "Approved public source chunks met the retrieval confidence threshold.", ""))
	}
	return resp
}

func (o *Orchestrator) groundedSources(ctx context.Context, query string) ([]domain.Source, groundingStatus) {
	retrieved, err := o.retriever.RetrieveSources(ctx, query)
	if err != nil {
		return nil, groundingUnavailable
	}
	deduped := make([]domain.Source, 0, len(retrieved))
	seen := map[string]int{}
	status := groundingReady
	for _, source := range retrieved {
		key := sourceKey(source)
		if key == "" {
			continue
		}
		if source.Confidence < rag.DefaultMinimumConfidence {
			continue
		}
		source = withSourceCaution(source)
		if source.Caution != "" {
			status = groundingCaution
		}
		if existing, ok := seen[key]; ok {
			if source.Confidence > deduped[existing].Confidence {
				deduped[existing] = source
			}
			continue
		}
		seen[key] = len(deduped)
		deduped = append(deduped, source)
	}
	if len(deduped) == 0 {
		return nil, groundingUnavailable
	}
	return deduped, status
}

func sourceKey(source domain.Source) string {
	if strings.TrimSpace(source.URL) != "" {
		return strings.TrimSpace(strings.ToLower(source.URL))
	}
	return strings.TrimSpace(strings.ToLower(source.ID + ":" + source.ChunkID))
}

func withSourceCaution(source domain.Source) domain.Source {
	if strings.EqualFold(source.FreshnessStatus, string(rag.FreshnessStale)) {
		source.Caution = "source is stale; staff confirmation required"
		return source
	}
	if strings.EqualFold(source.RiskLevel, string(rag.RiskHigh)) && source.Confidence < rag.HighRiskConfidenceThreshold {
		source.Caution = "high-risk source confidence is below staff-confirmation threshold"
	}
	return source
}
