package orchestrator

import (
	"context"
	"errors"
	"strings"
	"testing"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/tools"
	"askoc-mvp/internal/workflow"
)

var _ handlers.ChatService = (*Orchestrator)(nil)
var _ IntentClassifier = (*fakeClassifier)(nil)
var _ BannerTool = (*fakeBanner)(nil)
var _ PaymentTool = (*fakePayment)(nil)
var _ WorkflowTool = (*fakeWorkflow)(nil)
var _ CRMTool = (*fakeCRM)(nil)
var _ AuditRecorder = (*fakeAudit)(nil)

func TestNewRequiresDependencies(t *testing.T) {
	if _, err := New(completeDeps()); err != nil {
		t.Fatalf("New returned error with complete deps: %v", err)
	}

	tests := []struct {
		name string
		edit func(*Dependencies)
		want string
	}{
		{name: "classifier", edit: func(d *Dependencies) { d.Classifier = nil }, want: "classifier"},
		{name: "retriever", edit: func(d *Dependencies) { d.Retriever = nil }, want: "retriever"},
		{name: "llm", edit: func(d *Dependencies) { d.LLM = nil }, want: "llm"},
		{name: "banner", edit: func(d *Dependencies) { d.Banner = nil }, want: "banner"},
		{name: "payment", edit: func(d *Dependencies) { d.Payment = nil }, want: "payment"},
		{name: "workflow", edit: func(d *Dependencies) { d.Workflow = nil }, want: "workflow"},
		{name: "crm", edit: func(d *Dependencies) { d.CRM = nil }, want: "crm"},
		{name: "audit", edit: func(d *Dependencies) { d.Audit = nil }, want: "audit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := completeDeps()
			tt.edit(&deps)
			_, err := New(deps)
			if err == nil {
				t.Fatal("New returned nil error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want it to mention %q", err.Error(), tt.want)
			}
		})
	}
}

func TestTranscriptStatusPaidRecordSkipsWorkflow(t *testing.T) {
	deps := completeDeps()
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-paid"), domain.ChatRequest{
		ConversationID: "conv-paid",
		Channel:        "web",
		Message:        "I ordered my transcript but it has not been processed.",
		StudentID:      "S100001",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if !strings.Contains(resp.Answer, "ready for processing") {
		t.Fatalf("answer = %q, want ready status", resp.Answer)
	}
	if resp.Escalation != nil {
		t.Fatalf("escalation = %+v, want nil", resp.Escalation)
	}
	if len(workflowFake(t, deps).calls) != 0 {
		t.Fatalf("workflow calls = %+v, want none", workflowFake(t, deps).calls)
	}
	if len(crmFake(t, deps).requests) != 0 {
		t.Fatalf("crm requests = %+v, want none", crmFake(t, deps).requests)
	}
	assertAction(t, resp, "banner_status_checked", domain.ActionStatusCompleted)
	assertAction(t, resp, "payment_status_checked", domain.ActionStatusCompleted)
	assertAction(t, resp, "payment_reminder_skipped", domain.ActionStatusSkipped)
	assertActionsHaveTrace(t, resp, "trace-paid")
}

func TestTranscriptStatusUnpaidRecordTriggersWorkflowAndAudit(t *testing.T) {
	deps := completeDeps()
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-unpaid"), domain.ChatRequest{
		ConversationID: "conv-unpaid",
		Channel:        "web",
		Message:        "I ordered my transcript but it has not been processed.",
		StudentID:      "S100002",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(workflowFake(t, deps).calls) != 1 {
		t.Fatalf("workflow calls = %+v, want exactly one", workflowFake(t, deps).calls)
	}
	call := workflowFake(t, deps).calls[0]
	if call.StudentID != "S100002" || call.Item != "official_transcript" {
		t.Fatalf("workflow call = %+v", call)
	}
	wantKey := workflow.PaymentReminderKey("trace-unpaid", "S100002", "official_transcript")
	if call.IdempotencyKey != wantKey {
		t.Fatalf("idempotency key = %q, want %q", call.IdempotencyKey, wantKey)
	}
	action := assertAction(t, resp, "payment_reminder_triggered", domain.ActionStatusCompleted)
	if action.ReferenceID != "WF-ACCEPTED-1" {
		t.Fatalf("workflow action reference = %q", action.ReferenceID)
	}
	if action.IdempotencyKey != wantKey {
		t.Fatalf("action idempotency key = %q, want %q", action.IdempotencyKey, wantKey)
	}
	if !auditFake(t, deps).has("workflow_payment_reminder", "attempted") || !auditFake(t, deps).has("workflow_payment_reminder", "completed") {
		t.Fatalf("audit events = %+v, want workflow attempted and completed", auditFake(t, deps).events)
	}
	assertActionsHaveTrace(t, resp, "trace-unpaid")
}

func TestTranscriptStatusWorkflowFailureIsSafeAndAudited(t *testing.T) {
	deps := completeDeps()
	workflowFake(t, deps).err = errors.New("private webhook token leaked")
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-workflow-fail"), domain.ChatRequest{
		ConversationID: "conv-workflow-fail",
		Channel:        "web",
		Message:        "I ordered my transcript but it has not been processed.",
		StudentID:      "S100002",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(workflowFake(t, deps).calls) != 1 {
		t.Fatalf("workflow calls = %+v, want exactly one", workflowFake(t, deps).calls)
	}
	if strings.Contains(resp.Answer, "webhook token") {
		t.Fatalf("answer leaked workflow error: %q", resp.Answer)
	}
	action := assertAction(t, resp, "payment_reminder_triggered", domain.ActionStatusPending)
	if strings.Contains(action.Message, "webhook token") {
		t.Fatalf("action leaked workflow error: %+v", action)
	}
	if !auditFake(t, deps).has("workflow_payment_reminder", "failed") {
		t.Fatalf("audit events = %+v, want failed workflow event", auditFake(t, deps).events)
	}
	for _, event := range auditFake(t, deps).events {
		if strings.Contains(event.Message, "webhook token") {
			t.Fatalf("audit leaked workflow error: %+v", event)
		}
	}
}

func TestTranscriptStatusFinancialHoldCreatesCRMCase(t *testing.T) {
	deps := completeDeps()
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-hold"), domain.ChatRequest{
		ConversationID: "conv-hold",
		Channel:        "web",
		Message:        "Can you check my transcript?",
		StudentID:      "S100003",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(workflowFake(t, deps).calls) != 0 {
		t.Fatalf("workflow calls = %+v, want none for financial hold", workflowFake(t, deps).calls)
	}
	if len(crmFake(t, deps).requests) != 1 {
		t.Fatalf("crm requests = %+v, want one", crmFake(t, deps).requests)
	}
	gotCase := crmFake(t, deps).requests[0]
	if gotCase.Queue != "registrar_student_accounts" || gotCase.Priority != "high" {
		t.Fatalf("crm route = queue %q priority %q", gotCase.Queue, gotCase.Priority)
	}
	if gotCase.Intent != string(domain.IntentTranscriptStatus) || gotCase.SourceTraceID != "trace-hold" {
		t.Fatalf("crm request = %+v", gotCase)
	}
	if resp.Escalation == nil || resp.Escalation.CaseID != "MOCK-CRM-1" {
		t.Fatalf("escalation = %+v, want CRM case", resp.Escalation)
	}
	assertAction(t, resp, "financial_hold_detected", domain.ActionStatusCompleted)
	assertAction(t, resp, "crm_case_created", domain.ActionStatusCompleted)
	assertActionsHaveTrace(t, resp, "trace-hold")
}

func TestTranscriptStatusUnknownSyntheticRecordCreatesNormalHandoff(t *testing.T) {
	deps := completeDeps()
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-unknown-record"), domain.ChatRequest{
		ConversationID: "conv-unknown-record",
		Channel:        "web",
		Message:        "Can you check my transcript?",
		StudentID:      "S100004",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(crmFake(t, deps).requests) != 1 {
		t.Fatalf("crm requests = %+v, want one", crmFake(t, deps).requests)
	}
	if crmFake(t, deps).requests[0].Queue != "learner_support" || crmFake(t, deps).requests[0].Priority != "normal" {
		t.Fatalf("crm request = %+v, want normal learner support route", crmFake(t, deps).requests[0])
	}
	if resp.Escalation == nil || resp.Escalation.Status != domain.HandoffPending {
		t.Fatalf("escalation = %+v, want pending handoff", resp.Escalation)
	}
	assertAction(t, resp, "transcript_status_unresolved", domain.ActionStatusPending)
}

func TestTranscriptStatusMissingStudentIDPromptsWithoutToolCalls(t *testing.T) {
	deps := completeDeps()
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-missing-id"), domain.ChatRequest{
		ConversationID: "conv-missing-id",
		Channel:        "web",
		Message:        "I ordered my transcript but it has not been processed.",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(bannerFake(t, deps).calls) != 0 || len(paymentFake(t, deps).calls) != 0 || len(workflowFake(t, deps).calls) != 0 || len(crmFake(t, deps).requests) != 0 {
		t.Fatalf("tool calls should be empty: banner=%v payment=%v workflow=%v crm=%v", bannerFake(t, deps).calls, paymentFake(t, deps).calls, workflowFake(t, deps).calls, crmFake(t, deps).requests)
	}
	if !strings.Contains(resp.Answer, "synthetic demo student ID") {
		t.Fatalf("answer = %q, want synthetic ID prompt", resp.Answer)
	}
	assertAction(t, resp, "student_id_required", domain.ActionStatusPending)
}

func TestLowConfidenceCreatesNormalHandoffWithoutSensitiveToolCalls(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentUnknown,
		Confidence: 0.30,
		Sentiment:  domain.SentimentNeutral,
	}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-low-confidence"), domain.ChatRequest{
		ConversationID: "conv-low-confidence",
		Channel:        "web",
		Message:        "Can you solve something unrelated?",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(bannerFake(t, deps).calls) != 0 || len(paymentFake(t, deps).calls) != 0 || len(workflowFake(t, deps).calls) != 0 {
		t.Fatalf("sensitive tool calls should be empty: banner=%v payment=%v workflow=%v", bannerFake(t, deps).calls, paymentFake(t, deps).calls, workflowFake(t, deps).calls)
	}
	if len(crmFake(t, deps).requests) != 1 {
		t.Fatalf("crm requests = %+v, want normal handoff", crmFake(t, deps).requests)
	}
	if crmFake(t, deps).requests[0].Priority != "normal" {
		t.Fatalf("crm priority = %q, want normal", crmFake(t, deps).requests[0].Priority)
	}
	if resp.Escalation == nil || resp.Escalation.Queue != "learner_support" {
		t.Fatalf("escalation = %+v, want learner support handoff", resp.Escalation)
	}
}

func TestUrgentEscalationRedactsCRMSummary(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentHumanHandoff,
		Confidence: 0.88,
		Sentiment:  domain.SentimentUrgentNegative,
	}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-urgent"), domain.ChatRequest{
		ConversationID: "conv-urgent",
		Channel:        "web",
		Message:        "This is extremely frustrating. Email me at learner@example.com or call 250-555-0199 today.",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(crmFake(t, deps).requests) != 1 {
		t.Fatalf("crm requests = %+v, want one", crmFake(t, deps).requests)
	}
	got := crmFake(t, deps).requests[0]
	if got.Queue != "priority_staff_queue" || got.Priority != "high" {
		t.Fatalf("crm route = %+v, want priority staff queue", got)
	}
	if strings.Contains(got.Summary, "learner@example.com") || strings.Contains(got.Summary, "250-555-0199") {
		t.Fatalf("summary was not redacted: %q", got.Summary)
	}
	if !strings.Contains(got.Summary, "[REDACTED_EMAIL]") || !strings.Contains(got.Summary, "[REDACTED_PHONE]") {
		t.Fatalf("summary missing redaction markers: %q", got.Summary)
	}
	if resp.Escalation == nil || resp.Escalation.Priority != "high" {
		t.Fatalf("escalation = %+v, want high priority", resp.Escalation)
	}
}

func completeDeps() Dependencies {
	return Dependencies{
		Classifier: &fakeClassifier{result: classifier.Result{
			Intent:     domain.IntentTranscriptStatus,
			Confidence: 0.88,
			Sentiment:  domain.SentimentNeutral,
		}},
		Retriever: DisabledRetriever{},
		LLM:       DisabledLLM{},
		Banner: &fakeBanner{statuses: map[string]tools.BannerTranscriptStatus{
			"S100001": {
				StudentID:               "S100001",
				TranscriptRequestID:     "SYNTH-TRN-100001",
				TranscriptRequestStatus: "ready_for_processing",
				EligibleForProcessing:   true,
				Hold:                    "none",
				Synthetic:               true,
			},
			"S100002": {
				StudentID:               "S100002",
				TranscriptRequestID:     "SYNTH-TRN-100002",
				TranscriptRequestStatus: "blocked_by_unpaid_fee",
				EligibleForProcessing:   false,
				Hold:                    "payment",
				Holds:                   []string{"mock_payment_hold"},
				Synthetic:               true,
			},
			"S100003": {
				StudentID:               "S100003",
				TranscriptRequestID:     "SYNTH-TRN-100003",
				TranscriptRequestStatus: "needs_staff_review",
				EligibleForProcessing:   false,
				Hold:                    "financial",
				Holds:                   []string{"mock_financial_hold"},
				Synthetic:               true,
			},
			"S100004": {
				StudentID:               "S100004",
				TranscriptRequestID:     "SYNTH-TRN-100004",
				TranscriptRequestStatus: "not_found",
				EligibleForProcessing:   false,
				Hold:                    "none",
				Synthetic:               true,
			},
		}},
		Payment: &fakePayment{statuses: map[string]tools.PaymentStatus{
			"S100001": {
				StudentID:     "S100001",
				Item:          "official_transcript",
				Status:        "paid",
				AmountDue:     0,
				Currency:      "CAD",
				TransactionID: "SYNTH-PAY-100001",
				Synthetic:     true,
			},
			"S100002": {
				StudentID:     "S100002",
				Item:          "official_transcript",
				Status:        "unpaid",
				AmountDue:     15,
				Currency:      "CAD",
				TransactionID: "SYNTH-PAY-100002",
				Synthetic:     true,
			},
			"S100003": {
				StudentID:     "S100003",
				Item:          "official_transcript",
				Status:        "review_required",
				AmountDue:     40,
				Currency:      "CAD",
				TransactionID: "SYNTH-PAY-100003",
				Synthetic:     true,
			},
			"S100004": {
				StudentID: "S100004",
				Item:      "official_transcript",
				Status:    "not_applicable",
				Currency:  "CAD",
				Synthetic: true,
			},
		}},
		Workflow: &fakeWorkflow{resp: workflow.PaymentReminderResponse{
			WorkflowID: "WF-ACCEPTED-1",
			Status:     "accepted",
			Synthetic:  true,
		}},
		CRM:   &fakeCRM{},
		Audit: &fakeAudit{},
	}
}

func mustNew(t *testing.T, deps Dependencies) *Orchestrator {
	t.Helper()
	orch, err := New(deps)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	return orch
}

func traceContext(traceID string) context.Context {
	return middleware.WithTraceID(context.Background(), traceID)
}

func classifierFake(t *testing.T, deps Dependencies) *fakeClassifier {
	t.Helper()
	got, ok := deps.Classifier.(*fakeClassifier)
	if !ok {
		t.Fatalf("classifier fake has type %T", deps.Classifier)
	}
	return got
}

func bannerFake(t *testing.T, deps Dependencies) *fakeBanner {
	t.Helper()
	got, ok := deps.Banner.(*fakeBanner)
	if !ok {
		t.Fatalf("banner fake has type %T", deps.Banner)
	}
	return got
}

func paymentFake(t *testing.T, deps Dependencies) *fakePayment {
	t.Helper()
	got, ok := deps.Payment.(*fakePayment)
	if !ok {
		t.Fatalf("payment fake has type %T", deps.Payment)
	}
	return got
}

func workflowFake(t *testing.T, deps Dependencies) *fakeWorkflow {
	t.Helper()
	got, ok := deps.Workflow.(*fakeWorkflow)
	if !ok {
		t.Fatalf("workflow fake has type %T", deps.Workflow)
	}
	return got
}

func crmFake(t *testing.T, deps Dependencies) *fakeCRM {
	t.Helper()
	got, ok := deps.CRM.(*fakeCRM)
	if !ok {
		t.Fatalf("crm fake has type %T", deps.CRM)
	}
	return got
}

func auditFake(t *testing.T, deps Dependencies) *fakeAudit {
	t.Helper()
	got, ok := deps.Audit.(*fakeAudit)
	if !ok {
		t.Fatalf("audit fake has type %T", deps.Audit)
	}
	return got
}

func assertAction(t *testing.T, resp domain.ChatResponse, actionType string, status domain.ActionStatus) domain.Action {
	t.Helper()
	for _, action := range resp.Actions {
		if action.Type == actionType {
			if action.Status != status {
				t.Fatalf("action %s status = %q, want %q; actions=%+v", actionType, action.Status, status, resp.Actions)
			}
			return action
		}
	}
	t.Fatalf("action %s not found in %+v", actionType, resp.Actions)
	return domain.Action{}
}

func assertActionsHaveTrace(t *testing.T, resp domain.ChatResponse, traceID string) {
	t.Helper()
	for _, action := range resp.Actions {
		if action.TraceID != traceID {
			t.Fatalf("action %s trace_id = %q, want %q; actions=%+v", action.Type, action.TraceID, traceID, resp.Actions)
		}
	}
}

type fakeClassifier struct {
	result classifier.Result
	err    error
}

func (f *fakeClassifier) Classify(context.Context, string) (classifier.Result, error) {
	return f.result, f.err
}

type fakeBanner struct {
	statuses map[string]tools.BannerTranscriptStatus
	calls    []string
	err      error
}

func (f *fakeBanner) GetTranscriptStatus(_ context.Context, studentID string) (tools.BannerTranscriptStatus, error) {
	f.calls = append(f.calls, studentID)
	if f.err != nil {
		return tools.BannerTranscriptStatus{}, f.err
	}
	status, ok := f.statuses[studentID]
	if !ok {
		return tools.BannerTranscriptStatus{}, &tools.ToolError{Kind: tools.KindNotFound, Service: "banner"}
	}
	return status, nil
}

type fakePayment struct {
	statuses map[string]tools.PaymentStatus
	calls    []string
	err      error
}

func (f *fakePayment) GetPaymentStatus(_ context.Context, studentID string) (tools.PaymentStatus, error) {
	f.calls = append(f.calls, studentID)
	if f.err != nil {
		return tools.PaymentStatus{}, f.err
	}
	status, ok := f.statuses[studentID]
	if !ok {
		return tools.PaymentStatus{}, &tools.ToolError{Kind: tools.KindNotFound, Service: "payment"}
	}
	return status, nil
}

type fakeWorkflow struct {
	calls []workflow.PaymentReminderRequest
	resp  workflow.PaymentReminderResponse
	err   error
}

func (f *fakeWorkflow) SendPaymentReminder(_ context.Context, req workflow.PaymentReminderRequest) (workflow.PaymentReminderResponse, error) {
	f.calls = append(f.calls, req)
	if f.err != nil {
		return workflow.PaymentReminderResponse{}, f.err
	}
	resp := f.resp
	resp.IdempotencyKey = req.IdempotencyKey
	return resp, nil
}

type fakeCRM struct {
	requests []tools.CRMCaseRequest
	err      error
}

func (f *fakeCRM) CreateCase(_ context.Context, req tools.CRMCaseRequest) (tools.CRMCaseResponse, error) {
	f.requests = append(f.requests, req)
	if f.err != nil {
		return tools.CRMCaseResponse{}, f.err
	}
	return tools.CRMCaseResponse{
		CaseID:         "MOCK-CRM-1",
		Status:         "created",
		Queue:          req.Queue,
		Priority:       req.Priority,
		Summary:        req.Summary,
		ConversationID: req.ConversationID,
		SourceTraceID:  req.SourceTraceID,
		Synthetic:      true,
	}, nil
}

type fakeAudit struct {
	events []audit.Event
	err    error
}

func (f *fakeAudit) Record(_ context.Context, event audit.Event) error {
	f.events = append(f.events, event)
	return f.err
}

func (f *fakeAudit) has(action, status string) bool {
	for _, event := range f.events {
		if event.Action == action && event.Status == status {
			return true
		}
	}
	return false
}
