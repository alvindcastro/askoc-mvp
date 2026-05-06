package orchestrator

import (
	"context"
	"fmt"
	"strings"

	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/tools"
)

type handoffRequest struct {
	queue       string
	priority    string
	reason      string
	status      tools.BannerTranscriptStatus
	includeFlow bool
}

func (o *Orchestrator) createHandoff(ctx context.Context, req domain.ChatRequest, result classifier.Result, resp domain.ChatResponse, handoff handoffRequest) (domain.ChatResponse, error) {
	if strings.TrimSpace(handoff.queue) == "" {
		handoff.queue = "learner_support"
	}
	if strings.TrimSpace(handoff.priority) == "" {
		handoff.priority = "normal"
	}
	caseReq := tools.CRMCaseRequest{
		StudentID:      strings.TrimSpace(req.StudentID),
		ConversationID: resp.ConversationID,
		Intent:         string(result.Intent),
		Priority:       handoff.priority,
		Queue:          handoff.queue,
		Summary:        o.redact(o.summary(req, result, handoff)),
		SourceTraceID:  traceID(ctx),
	}
	created, err := o.crm.CreateCase(ctx, caseReq)
	if err != nil {
		resp.Answer = "I could not create the mock staff handoff, but no internal CRM details were exposed."
		resp.Actions = append(resp.Actions, o.action(ctx, "crm_case_created", domain.ActionStatusPending, "Mock CRM case could not be confirmed.", ""))
		resp.Escalation = &domain.Escalation{
			Required: true,
			Status:   domain.HandoffPending,
			Queue:    handoff.queue,
			Priority: handoff.priority,
			Reason:   handoff.reason,
		}
		return resp, nil
	}

	resp.Actions = append(resp.Actions, o.action(ctx, "crm_case_created", domain.ActionStatusCompleted, "Mock CRM case created for staff review.", created.CaseID))
	resp.Escalation = &domain.Escalation{
		Required: true,
		Status:   domain.HandoffPending,
		Queue:    created.Queue,
		Priority: created.Priority,
		CaseID:   created.CaseID,
		Reason:   handoff.reason,
	}
	if strings.TrimSpace(resp.Answer) == "" {
		resp.Answer = "I created a mock staff handoff so a learner-services team member can review this synthetic demo case."
	}
	return resp, nil
}

func (o *Orchestrator) summary(req domain.ChatRequest, result classifier.Result, handoff handoffRequest) string {
	parts := []string{
		fmt.Sprintf("Intent: %s", result.Intent),
		fmt.Sprintf("Confidence: %.2f", result.Confidence),
		fmt.Sprintf("Sentiment: %s", result.Sentiment),
		"Queue reason: " + handoff.reason,
	}
	if req.StudentID != "" {
		parts = append(parts, "Synthetic student ID: "+req.StudentID)
	}
	if handoff.status.TranscriptRequestID != "" {
		parts = append(parts, "Transcript request: "+handoff.status.TranscriptRequestID)
		parts = append(parts, "Transcript status: "+handoff.status.TranscriptRequestStatus)
	}
	if strings.TrimSpace(req.Message) != "" {
		parts = append(parts, "Learner message: "+req.Message)
	}
	return strings.Join(parts, ". ")
}

func needsEscalation(result classifier.Result) bool {
	return result.Intent == domain.IntentHumanHandoff ||
		result.Intent == domain.IntentEscalationRequest ||
		result.NeedsHandoff ||
		result.Sentiment == domain.SentimentNegative ||
		result.Sentiment == domain.SentimentUrgentNegative ||
		result.Confidence < classifier.SensitiveToolConfidence
}

func queueForResult(result classifier.Result) string {
	if result.Sentiment == domain.SentimentUrgentNegative {
		return "priority_staff_queue"
	}
	return "learner_support"
}

func priorityForResult(result classifier.Result) string {
	if result.Sentiment == domain.SentimentUrgentNegative {
		return "high"
	}
	return "normal"
}
