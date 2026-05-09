# Power Automate / Workflow Design

## Purpose

This document describes the workflow automation layer for the AskOC AI Concierge MVP. The implementation is Go-first, with Power Automate as the realistic institutional automation layer and a local Go workflow simulator for portfolio demos.

## Workflow goal

Automate transcript payment/status follow-up:

```text
Learner asks why transcript has not been processed
→ Go orchestrator checks synthetic Banner/payment records
→ If unpaid, trigger payment reminder workflow
→ If hold or urgent sentiment, create CRM case
→ Notify learner
→ Log auditable workflow event
```

## Implementation options

### Option A: Power Automate webhook

The Go backend sends an HTTP POST request to a Power Automate cloud flow.

```text
cmd/api or internal/workflow
→ POST Power Automate HTTP endpoint
→ flow validates payload
→ flow sends email/Teams-style notification
→ flow updates mock CRM or audit endpoint
```

### Option B: Local Go workflow simulator

For a fully local P8 demo, `cmd/workflow-sim` simulates the Power Automate flow. The API still defaults to the in-process idempotent workflow client when `ASKOC_WORKFLOW_URL` is empty, and can call either `cmd/workflow-sim` or a Power Automate HTTP trigger by setting `ASKOC_WORKFLOW_URL`.

```text
cmd/api
→ POST http://workflow-sim:9084/api/v1/automation/payment-reminder
→ workflow-sim writes event to audit store
→ workflow-sim returns workflow ID
```

The implemented approach keeps Option B as the reliable offline demo path and Option A as a webhook-compatible extension behind the same `PaymentReminderSender` interface.

## Trigger payload

```json
{
  "student_id": "S100002",
  "conversation_id": "conv_01JABC123",
  "trace_id": "trace_01JABC456",
  "item": "official_transcript",
  "amount_due": 15.00,
  "currency": "CAD",
  "reason": "Transcript request cannot be processed until payment is complete.",
  "idempotency_key": "payment-reminder:trace_01JABC456:S100002:official_transcript"
}
```

## Response payload

```json
{
  "workflow_id": "WF-2026-000789",
  "status": "accepted",
  "message": "Payment reminder workflow accepted.",
  "idempotency_key": "payment-reminder:trace_01JABC456:S100002:official_transcript",
  "synthetic": true,
  "attempt_count": 1
}
```

## Go workflow client interface

```go
type PaymentReminderSender interface {
    SendPaymentReminder(ctx context.Context, req PaymentReminderRequest) (PaymentReminderResponse, error)
}

type PaymentReminderRequest struct {
    StudentID       string  `json:"student_id"`
    ConversationID  string  `json:"conversation_id,omitempty"`
    TraceID         string  `json:"trace_id,omitempty"`
    Item            string  `json:"item"`
    AmountDue       float64 `json:"amount_due,omitempty"`
    Currency        string  `json:"currency,omitempty"`
    Reason          string  `json:"reason"`
    IdempotencyKey  string  `json:"idempotency_key"`
}

type PaymentReminderResponse struct {
    WorkflowID     string `json:"workflow_id"`
    Status         string `json:"status"`
    Message        string `json:"message,omitempty"`
    IdempotencyKey string `json:"idempotency_key,omitempty"`
    Synthetic      bool   `json:"synthetic"`
    AttemptCount   int    `json:"attempt_count,omitempty"`
}
```

## Go webhook client behavior

The Go client should:

- use a short timeout,
- include `X-Trace-ID`,
- include an idempotency key,
- include `Idempotency-Key`,
- include `X-AskOC-Workflow-Signature` when `ASKOC_WORKFLOW_SIGNATURE` is configured,
- retry at most once for transient failures,
- never include sensitive raw conversation text,
- return a safe error to the orchestrator if the workflow fails.

Example behavior:

```go
ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()

result, err := workflowClient.TriggerPaymentReminder(ctx, req)
if err != nil {
    // Log failure, then create CRM case or show safe fallback.
}
```

Implemented runtime settings:

| Variable | Default | Purpose |
|---|---|---|
| `ASKOC_WORKFLOW_URL` | empty | Empty uses the in-process client; set to `http://localhost:9084/api/v1/automation/payment-reminder` for `cmd/workflow-sim` or to a Power Automate trigger URL for webhook mode; redacted from config output |
| `ASKOC_WORKFLOW_TIMEOUT_SECONDS` | `5` | HTTP client timeout |
| `ASKOC_WORKFLOW_SIGNATURE` | empty | Optional shared signature/header value; redacted from config output and never logged |
| `ASKOC_WORKFLOW_SIGNATURE_HEADER` | `X-AskOC-Workflow-Signature` | Header name used with the signature value |
| `ASKOC_WORKFLOW_MAX_RETRIES` | `1` | Number of retries for transient `5xx` webhook responses |

## Workflow decision table

| Transcript status | Payment status | Hold | Sentiment | Workflow action |
|---|---|---|---|---|
| requested | paid | none | neutral | No reminder; explain status |
| requested | unpaid | none | neutral | Trigger payment reminder |
| requested | unpaid | none | negative/urgent | Trigger reminder and create priority CRM case |
| requested | paid | financial | any | Create Registrar/Finance CRM case |
| unknown | unknown | none | any | Human handoff or low-confidence case |
| not found | n/a | n/a | any | Safe not-found response; do not expose data |

## Local workflow simulator endpoints

### `POST /api/v1/automation/payment-reminder`

Accepts the trigger payload and returns `202 Accepted` with a deterministic workflow ID, status, idempotency key, `synthetic: true`, and `attempt_count: 1`. Invalid JSON, missing idempotency keys, missing student IDs, and missing items return safe `400` errors. Duplicate idempotency keys return the same workflow ID.

### `GET /healthz`

Service health check.

### `GET /api/v1/admin/metrics`

Protected with `Authorization: Bearer demo-admin-token`. Shows workflow event counts from the simulator's in-memory audit store.

## Power Automate flow outline

1. Trigger: “When an HTTP request is received.”
2. Validate schema and required fields.
3. Check idempotency key against storage or mock list.
4. Compose learner notification text.
5. Send email/Teams-style notification.
6. Return workflow result with the original idempotency key and a workflow ID.
7. Keep any callback/audit endpoint optional for future hardening; the current Go API audits workflow attempts around the client call.

## Example learner notification

```text
Subject: Transcript request payment reminder

Your transcript request appears to be waiting for payment in this demo scenario.
Please review the transcript payment steps in myOkanagan or contact the Registrar if you believe this is incorrect.

Reference: WF-2026-000789
```

## Audit event

```json
{
  "trace_id": "trace_01JABC456",
  "conversation_id": "conv_01JABC123",
  "student_id": "S100002",
  "type": "workflow",
  "action": "workflow_payment_reminder",
  "status": "completed",
  "reference_id": "WF-2026-000789",
  "metadata": {
    "idempotency_key_hash": "sha256-hex-value",
    "attempt_count": "1"
  },
  "recorded_at": "2026-05-06T12:06:00Z"
}
```

## Privacy controls

- Do not send raw conversation transcripts to Power Automate.
- Do not send passwords, payment card details, government ID values, or personal notes.
- Use synthetic student IDs only in the demo.
- Keep notification content generic.
- Record workflow IDs for audit and troubleshooting.
- Store webhook URLs and signatures in environment variables or a secret manager, never in committed files.
- Treat Power Automate trigger URLs as secrets because they often contain query-string tokens.
- Use idempotency keys and signature headers for replay protection; reject duplicate keys rather than sending duplicate reminders.

## Failure handling

| Failure | Go orchestrator response |
|---|---|
| Workflow endpoint timeout | Log failure, retry once, then offer CRM handoff |
| Invalid workflow payload | Log developer error and show safe fallback |
| Duplicate idempotency key | Do not send duplicate reminder; return existing workflow status |
| Power Automate unavailable | Use local workflow simulator or create CRM case |
| CRM unavailable | Explain that handoff could not be completed in demo and log unresolved case |

## Demo positioning

Even if Power Automate is not fully configured, show the workflow as a webhook-compatible design. This proves that the Go AI solution can integrate with institutional automation tools while remaining testable in a local development environment.
