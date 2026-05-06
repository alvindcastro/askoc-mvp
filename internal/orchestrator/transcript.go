package orchestrator

import (
	"context"
	"fmt"
	"strings"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/tools"
	"askoc-mvp/internal/workflow"
)

func (o *Orchestrator) handleTranscriptStatus(ctx context.Context, req domain.ChatRequest, result classifier.Result, resp domain.ChatResponse) (domain.ChatResponse, error) {
	studentID := studentIDFrom(req)
	if studentID == "" {
		resp.Answer = "Please provide a synthetic demo student ID like S100002 so I can check the local transcript/payment demo record."
		resp.Actions = append(resp.Actions, o.action(ctx, "student_id_required", domain.ActionStatusPending, "No synthetic demo student ID was provided.", ""))
		return resp, nil
	}

	status, err := o.banner.GetTranscriptStatus(ctx, studentID)
	if err != nil {
		resp.Answer = "I could not find a matching synthetic transcript record. No real student system was checked."
		resp.Actions = append(resp.Actions, o.action(ctx, "banner_status_checked", domain.ActionStatusSkipped, "Synthetic transcript lookup did not find a demo record.", ""))
		return resp, nil
	}
	resp.Actions = append(resp.Actions, o.action(ctx, "banner_status_checked", domain.ActionStatusCompleted, "Synthetic transcript status checked.", status.TranscriptRequestID))

	if isFinancialHold(status) {
		resp.Actions = append(resp.Actions, o.action(ctx, "financial_hold_detected", domain.ActionStatusCompleted, "Synthetic financial hold requires staff routing.", status.TranscriptRequestID))
		resp.Answer = "Your synthetic transcript request needs staff review because a financial hold is present. I created a mock Registrar/Student Accounts case instead of sending a self-service payment reminder."
		return o.createHandoff(ctx, withStudentID(req, studentID), result, resp, handoffRequest{
			queue:       "registrar_student_accounts",
			priority:    "high",
			reason:      "financial hold detected",
			status:      status,
			includeFlow: true,
		})
	}

	if strings.EqualFold(status.TranscriptRequestStatus, "not_found") {
		resp.Answer = "I could not find an active synthetic transcript request for that demo student. I created a normal mock handoff so staff can review the unresolved status."
		resp.Actions = append(resp.Actions, o.action(ctx, "transcript_status_unresolved", domain.ActionStatusPending, "Synthetic transcript status is unresolved.", status.TranscriptRequestID))
		return o.createHandoff(ctx, withStudentID(req, studentID), result, resp, handoffRequest{
			queue:    "learner_support",
			priority: "normal",
			reason:   "synthetic transcript request not found",
			status:   status,
		})
	}

	payment, err := o.payment.GetPaymentStatus(ctx, studentID)
	if err != nil {
		resp.Answer = "I checked the synthetic transcript status, but the synthetic payment check was unavailable. I created a staff handoff without exposing internal tool details."
		resp.Actions = append(resp.Actions, o.action(ctx, "payment_status_checked", domain.ActionStatusPending, "Synthetic payment status could not be confirmed.", ""))
		return o.createHandoff(ctx, withStudentID(req, studentID), result, resp, handoffRequest{
			queue:    "learner_support",
			priority: "normal",
			reason:   "payment status unavailable",
			status:   status,
		})
	}
	resp.Actions = append(resp.Actions, o.action(ctx, "payment_status_checked", domain.ActionStatusCompleted, "Synthetic transcript payment status checked.", payment.TransactionID))

	if strings.EqualFold(payment.Status, "unpaid") || payment.AmountDue > 0 {
		return o.handleUnpaidTranscript(ctx, withStudentID(req, studentID), resp, status, payment)
	}

	resp.Answer = fmt.Sprintf("Your synthetic transcript request %s is ready for processing and payment is marked paid. No payment reminder is needed.", status.TranscriptRequestID)
	resp.Actions = append(resp.Actions, o.action(ctx, "payment_reminder_skipped", domain.ActionStatusSkipped, "Payment reminder skipped because the synthetic record is not unpaid.", ""))
	return resp, nil
}

func (o *Orchestrator) handleUnpaidTranscript(ctx context.Context, req domain.ChatRequest, resp domain.ChatResponse, status tools.BannerTranscriptStatus, payment tools.PaymentStatus) (domain.ChatResponse, error) {
	item := payment.Item
	if strings.TrimSpace(item) == "" {
		item = "official_transcript"
	}
	key := workflow.PaymentReminderKey(traceID(ctx), req.StudentID, item)
	reminder := workflow.PaymentReminderRequest{
		StudentID:      req.StudentID,
		ConversationID: resp.ConversationID,
		TraceID:        traceID(ctx),
		Item:           item,
		AmountDue:      payment.AmountDue,
		Currency:       payment.Currency,
		Reason:         "Transcript request cannot be processed until payment is complete.",
		IdempotencyKey: key,
	}

	o.recordAudit(ctx, audit.Event{
		TraceID:        traceID(ctx),
		ConversationID: resp.ConversationID,
		StudentID:      req.StudentID,
		Type:           "workflow",
		Action:         "workflow_payment_reminder",
		Status:         "attempted",
		Message:        "payment reminder workflow attempted",
	})

	workflowResp, err := o.workflow.SendPaymentReminder(ctx, reminder)
	if err != nil {
		o.recordAudit(ctx, audit.Event{
			TraceID:        traceID(ctx),
			ConversationID: resp.ConversationID,
			StudentID:      req.StudentID,
			Type:           "workflow",
			Action:         "workflow_payment_reminder",
			Status:         "failed",
			Message:        "payment reminder workflow failed",
		})
		resp.Answer = "Your synthetic transcript request is blocked by an unpaid demo balance, but I could not confirm the reminder workflow. No duplicate reminder was created."
		resp.Actions = append(resp.Actions, withIdempotency(o.action(ctx, "payment_reminder_triggered", domain.ActionStatusPending, "Payment reminder workflow could not be confirmed.", ""), key))
		return resp, nil
	}

	o.recordAudit(ctx, audit.Event{
		TraceID:        traceID(ctx),
		ConversationID: resp.ConversationID,
		StudentID:      req.StudentID,
		Type:           "workflow",
		Action:         "workflow_payment_reminder",
		Status:         "completed",
		ReferenceID:    workflowResp.WorkflowID,
		Message:        "payment reminder workflow accepted",
	})

	resp.Answer = fmt.Sprintf("Your synthetic transcript request %s is blocked by an unpaid demo balance of %.2f %s. I triggered a synthetic payment reminder workflow.", status.TranscriptRequestID, payment.AmountDue, payment.Currency)
	resp.Actions = append(resp.Actions, withIdempotency(o.action(ctx, "payment_reminder_triggered", domain.ActionStatusCompleted, "Synthetic payment reminder workflow accepted.", workflowResp.WorkflowID), key))
	return resp, nil
}

func (o *Orchestrator) recordAudit(ctx context.Context, event audit.Event) {
	_ = o.audit.Record(ctx, event)
}

func studentIDFrom(req domain.ChatRequest) string {
	if got := strings.TrimSpace(req.StudentID); got != "" {
		return got
	}
	return syntheticIDPattern.FindString(strings.ToUpper(req.Message))
}

func isFinancialHold(status tools.BannerTranscriptStatus) bool {
	if strings.EqualFold(status.Hold, "financial") {
		return true
	}
	for _, hold := range status.Holds {
		if strings.EqualFold(hold, "mock_financial_hold") {
			return true
		}
	}
	return false
}

func withStudentID(req domain.ChatRequest, studentID string) domain.ChatRequest {
	req.StudentID = studentID
	return req
}
