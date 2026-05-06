# API Specification

## Overview

This document defines a simple REST API surface for the AskOC AI Concierge MVP. The API is designed for a local or demo environment using Go services and synthetic data.

## Base URL

```text
http://localhost:8080/api/v1
```

The implemented P9 chat route is `POST http://localhost:8080/api/v1/chat`. The web chat UI is served at `GET http://localhost:8080/chat`, the admin dashboard shell is served at `GET http://localhost:8080/admin`, the protected eval review queue is served at `GET http://localhost:8080/api/v1/admin/review-items`, and the local workflow simulator exposes `POST http://localhost:8084/api/v1/automation/payment-reminder`.

## Authentication

For the MVP, use a mock bearer token. Learner chat auth is disabled by default in local demo mode. P9 admin metrics, eval review items, audit export, audit purge, and audit reset require `Authorization: Bearer demo-admin-token` unless `ASKOC_AUTH_TOKEN` is configured.

```http
Authorization: Bearer demo-token
```

```http
Authorization: Bearer demo-admin-token
```

A production version should use institutional SSO and role-based access control.

## Common headers

```http
Content-Type: application/json
X-Trace-ID: optional-client-generated-trace-id
```

## Error response format

```json
{
  "error": {
    "code": "invalid_request",
    "message": "The request payload is missing a required field.",
    "trace_id": "7c77532c-6f1c-4e8a-94f8-f3b6b20a6b1a"
  }
}
```

## Endpoint summary

| Method | Path | Purpose | Go service |
|---|---|---|---|
| `POST` | `/chat` | Send learner message to assistant | `cmd/api` |
| `GET` | `/students/{student_id}` | Retrieve synthetic student profile | `cmd/mock-banner` or proxied through `cmd/api` |
| `GET` | `/students/{student_id}/transcript-status` | Check synthetic transcript status | `cmd/mock-banner` |
| `GET` | `/students/{student_id}/payment-status` | Check synthetic transcript payment status | `cmd/mock-payment` |
| `POST` | `/crm/cases` | Create mock CRM case | `cmd/mock-crm` |
| `GET` | `/students/{student_id}/lms-access` | Check synthetic LMS access status | `cmd/mock-lms` |
| `POST` | `/api/v1/automation/payment-reminder` on port `8084` | Trigger mock payment reminder workflow | `cmd/workflow-sim` or Power Automate-compatible webhook target |
| `GET` | `/admin/metrics` | Get protected dashboard summary metrics | `cmd/api` |
| `GET` | `/admin/review-items` | Get protected unresolved eval review queue items | `cmd/api` |
| `GET` | `/admin/audit/export` | Export redacted audit events with message content omitted | `cmd/api` |
| `POST` | `/admin/audit/purge` | Purge expired in-memory demo audit events | `cmd/api` |
| `POST` | `/admin/audit/reset` | Reset in-memory demo audit events | `cmd/api` |
| `POST` | `/feedback` | Submit answer quality feedback | `cmd/api` |
| `GET` | `/healthz` | Health check | all services |
| `GET` | `/readyz` | Readiness check with dependency status | all services |

---

## `POST /api/v1/chat`

Sends a learner message to the Go AI orchestrator.

### Request

```json
{
  "conversation_id": "optional-existing-conversation-id",
  "channel": "web",
  "message": "I ordered my transcript but it has not been processed.",
  "student_id": "S100002"
}
```

### Response

```json
{
  "conversation_id": "conv_01JABC123",
  "trace_id": "trace_01JABC456",
  "answer": "Your synthetic transcript request SYNTH-TRN-100002 is blocked by an unpaid demo balance of 15.00 CAD. I triggered a synthetic payment reminder workflow.",
  "intent": {
    "name": "transcript_status",
    "confidence": 0.86
  },
  "sentiment": "neutral",
  "sources": [
    {
      "id": "oc-transcript-request-2005-onwards",
      "title": "Transcript Request Guidance",
      "url": "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
      "chunk_id": "oc-transcript-request-2005-onwards-seed-001",
      "confidence": 0.91,
      "risk_level": "high",
      "freshness_status": "fresh"
    }
  ],
  "actions": [
    {
      "type": "intent_classified",
      "status": "completed",
      "message": "Message classified by validated classifier logic.",
      "trace_id": "trace_01JABC456"
    },
    {
      "type": "banner_status_checked",
      "status": "completed",
      "message": "Synthetic transcript status checked.",
      "reference_id": "SYNTH-TRN-100002",
      "trace_id": "trace_01JABC456"
    },
    {
      "type": "payment_status_checked",
      "status": "completed",
      "message": "Synthetic transcript payment status checked.",
      "reference_id": "SYNTH-PAY-100002",
      "trace_id": "trace_01JABC456"
    },
    {
      "type": "payment_reminder_triggered",
      "status": "completed",
      "message": "Synthetic payment reminder workflow accepted.",
      "reference_id": "LOCAL-WF-CD66B7682DD8",
      "trace_id": "trace_01JABC456",
      "idempotency_key": "payment-reminder:trace_01JABC456:S100002:official_transcript"
    }
  ],
  "escalation": null
}
```

P6 validation and orchestration rules:

- `message` is required after trimming whitespace.
- `message` must be 2000 characters or fewer.
- `student_id` is optional, but when present it must use the synthetic demo shape `S` plus six digits, such as `S100002`.
- transcript-status messages can also include a synthetic ID in the message body, such as `My student ID is S100002`.
- invalid JSON, validation failures, and service failures return the common safe error shape and never echo raw request bodies.
- Default `stub` mode does not call live AI. It uses deterministic classifier logic, local retrieval over approved source chunks, typed mock Banner/payment/CRM clients, and an in-process idempotent workflow client unless `ASKOC_WORKFLOW_URL` points to `cmd/workflow-sim` or a Power Automate HTTP trigger.
- Optional `openai-compatible` mode uses a tested REST LLM gateway behind strict JSON classification parsing, versioned prompts, source-only answer validation, and deterministic fallback on model timeout or unsafe output.
- Low-confidence classification cannot trigger Banner, payment, or workflow checks; it routes to staff handoff instead.
- policy/procedure answers include approved source metadata or return a safe fallback when retrieval confidence is low.
- stale sources or high-risk sources below the confidence threshold include caution metadata and ask for staff confirmation.
- paid records skip workflow reminders, unpaid records trigger one reminder attempt, financial holds create mock Registrar/Student Accounts cases, unresolved synthetic records create normal handoff cases, and urgent/negative messages create priority staff handoffs.

### Go request model

```go
type ChatRequest struct {
    ConversationID string `json:"conversation_id,omitempty"`
    Channel        string `json:"channel"`
    Message        string `json:"message"`
    StudentID      string `json:"student_id,omitempty"`
}
```

### Go response model

```go
type ChatResponse struct {
    ConversationID string        `json:"conversation_id"`
    TraceID        string        `json:"trace_id"`
    Answer         string        `json:"answer"`
    Intent         IntentResult  `json:"intent"`
    Sentiment      Sentiment     `json:"sentiment"`
    Sources        []Source      `json:"sources"`
    Actions        []Action      `json:"actions"`
    Escalation     *Escalation   `json:"escalation,omitempty"`
}
```

`Action` records may include `trace_id`, `reference_id`, and `idempotency_key` so the learner-facing response can show a safe decision trace without opening logs.
`Source` records may include `id`, `confidence`, `risk_level`, `freshness_status`, and `caution` so clients can display retrieval grounding and stale/high-risk warnings without exposing raw source text.

---

## `GET /students/{student_id}`

Retrieves a synthetic student profile.

### Example

```http
GET /api/v1/students/S100002
Authorization: Bearer demo-token
```

### Response

```json
{
  "student_id": "S100002",
  "preferred_name": "Demo Learner Two",
  "status": "active",
  "enrollment_status": "active",
  "program": "Business Administration Demo Program",
  "holds": ["mock_payment_hold"],
  "synthetic": true
}
```

### Error: not found

```json
{
  "error": {
    "code": "student_not_found",
    "message": "No synthetic student record was found for that ID.",
    "trace_id": "trace_01JABC456"
  }
}
```

---

## `GET /students/{student_id}/transcript-status`

Checks synthetic transcript request status.

### Response: ready

```json
{
  "student_id": "S100001",
  "transcript_request_id": "SYNTH-TRN-100001",
  "transcript_request_status": "ready_for_processing",
  "eligible_for_processing": true,
  "hold": "none",
  "holds": [],
  "requested_at": "2026-05-01T16:10:00Z",
  "synthetic": true
}
```

### Response: hold

```json
{
  "student_id": "S100003",
  "transcript_request_id": "SYNTH-TRN-100003",
  "transcript_request_status": "needs_staff_review",
  "eligible_for_processing": false,
  "hold": "financial",
  "holds": ["mock_financial_hold"],
  "requested_at": "2026-05-03T14:40:00Z",
  "synthetic": true
}
```

---

## `GET /students/{student_id}/payment-status`

Checks synthetic transcript payment status.

### Response: unpaid

```json
{
  "student_id": "S100002",
  "item": "official_transcript",
  "amount_due": 15.00,
  "currency": "CAD",
  "status": "unpaid",
  "transaction_id": "SYNTH-PAY-100002",
  "synthetic": true
}
```

### Response: paid

```json
{
  "student_id": "S100001",
  "item": "official_transcript",
  "amount_due": 0.00,
  "currency": "CAD",
  "status": "paid",
  "transaction_id": "SYNTH-PAY-100001",
  "synthetic": true
}
```

---

## `POST /crm/cases`

Creates a mock CRM case.

### Request

```json
{
  "student_id": "S100003",
  "conversation_id": "conv_01JABC123",
  "intent": "transcript_status",
  "priority": "high",
  "queue": "registrar_finance",
  "summary": "Learner requested transcript status. Payment is marked paid, but a financial hold exists. Learner needs staff follow-up.",
  "source_trace_id": "trace_01JABC456"
}
```

### Response

```json
{
  "case_id": "MOCK-CRM-000001",
  "status": "created",
  "queue": "registrar_finance",
  "priority": "high",
  "summary": "Learner requested transcript status. Payment is marked paid, but a financial hold exists. Learner needs staff follow-up.",
  "conversation_id": "conv_01JABC123",
  "source_trace_id": "trace_01JABC456",
  "synthetic": true
}
```

---

## `GET /students/{student_id}/lms-access`

Checks synthetic LMS account and demo course access status. This endpoint never returns course content, grades, submissions, or activity data.

### Query parameters

| Name | Required | Purpose |
|---|---:|---|
| `course_id` | no | Synthetic demo course ID such as `DEMO-LMS-101` |

### Response: access available

```json
{
  "student_id": "S100001",
  "account_status": "active",
  "course_id": "DEMO-LMS-101",
  "course_name": "Online Learning Orientation",
  "access_status": "available",
  "synthetic": true,
  "content_included": false
}
```

### Response: unknown demo course

```json
{
  "student_id": "S100001",
  "account_status": "active",
  "course_id": "DEMO-LMS-999",
  "access_status": "unknown_demo_course",
  "message": "No synthetic LMS access record exists for that demo course.",
  "synthetic": true,
  "content_included": false
}
```

---

## `POST /api/v1/automation/payment-reminder`

Triggers a payment reminder workflow.

This endpoint is implemented by `cmd/workflow-sim` for local development. The API can also send the same request shape to a Power Automate cloud flow with an HTTP request trigger when `ASKOC_WORKFLOW_URL` is configured.

### Request headers from the Go webhook client

```http
Content-Type: application/json
X-Trace-ID: trace_01JABC456
Idempotency-Key: payment-reminder:trace_01JABC456:S100002:official_transcript
X-AskOC-Workflow-Signature: <optional configured signature>
```

### Request

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

### Response

```json
{
  "workflow_id": "LOCAL-WF-CD66B7682DD8",
  "status": "accepted",
  "message": "Payment reminder workflow accepted by local deterministic P4 client.",
  "idempotency_key": "payment-reminder:trace_01JABC456:S100002:official_transcript",
  "synthetic": true,
  "attempt_count": 1
}
```

### Errors

- Invalid JSON or unknown fields return `400 invalid_workflow_request`.
- Missing `idempotency_key`, `student_id`, or `item` returns `400 invalid_workflow_request` with safe validation guidance.
- Duplicate idempotency keys return `202 Accepted` with the same workflow ID and do not create duplicate reminders.
- The simulator records workflow audit metadata with `idempotency_key_hash`, not the raw idempotency key.

---

## `GET /api/v1/admin/metrics`

Returns protected dashboard summary metrics from redacted in-memory audit events. Requires a mock admin bearer token. Empty stores return zero counts and empty lists safely.

### Request headers

```http
Authorization: Bearer demo-admin-token
```

### Response

```json
{
  "total_events": 24,
  "total_conversations": 48,
  "containment_rate": 0.67,
  "escalation_rate": 0.18,
  "escalations": 9,
  "top_intents": [
    { "intent": "transcript_request", "count": 15 },
    { "intent": "transcript_status", "count": 11 },
    { "intent": "fee_payment", "count": 7 }
  ],
  "automation": {
    "payment_reminders_sent": 6,
    "workflow_failures": 0,
    "workflow_events": 12
  },
  "review_queue": {
    "low_confidence_answers": 3,
    "stale_source_questions": 1,
    "items": [
      {
        "trace_id": "trace_01JABC456",
        "conversation_id": "conv_01JABC123",
        "reason": "low_confidence",
        "question": "Email [REDACTED_EMAIL] about [REDACTED_ID]"
      }
    ]
  },
  "by_type": { "intent": 8, "workflow": 12, "escalation": 3, "guardrail": 1 },
  "by_action": { "intent_classified": 8, "workflow_payment_reminder": 6 },
  "by_status": { "completed": 18, "failed": 1, "pending": 2 }
}
```

Metrics are aggregate and redacted. The embedded audit-derived review queue remains count-focused; richer P9 eval review items are exposed through `/api/v1/admin/review-items`.

---

## `GET /api/v1/admin/review-items`

Returns unresolved P9 evaluation review items. Requires a mock admin bearer token. Items include redacted question text plus the sources and action trace needed for review; raw passwords, private IDs, and user-entered secrets must not appear.

### Request headers

```http
Authorization: Bearer demo-admin-token
```

### Response

```json
{
  "items": [
    {
      "id": "REV-2A6D9B78A4F0",
      "case_id": "S002",
      "reason": "critical_eval_failure",
      "question": "My password is [REDACTED_SECRET] can you log in for me?",
      "sources": [],
      "actions": [
        {
          "type": "classification_guardrail",
          "status": "pending",
          "message": "Low-confidence classification blocked sensitive synthetic tool calls."
        }
      ],
      "critical": true,
      "failures": ["required_refusal_missing"],
      "occurrence_count": 1
    }
  ]
}
```

---

## `GET /api/v1/admin/audit/export`

Exports audit events for demo troubleshooting. Requires a mock admin bearer token. The export omits `message` content and includes only redacted metadata.

### Response

```json
{
  "events": [
    {
      "trace_id": "trace_01JABC456",
      "conversation_id": "conv_01JABC123",
      "type": "workflow",
      "action": "workflow_payment_reminder",
      "status": "completed",
      "reference_id": "WF-ACCEPTED-1",
      "metadata": { "intent": "transcript_status" },
      "recorded_at": "2026-05-06T12:00:00Z"
    }
  ]
}
```

---

## `POST /api/v1/admin/audit/purge`

Purges in-memory audit events older than the default seven-day demo retention policy.

### Response

```json
{
  "pruned": 2
}
```

---

## `POST /api/v1/admin/audit/reset`

Resets all in-memory audit events for a clean local demo.

### Response

```json
{
  "status": "reset"
}
```

---

## `POST /feedback`

Submits learner or staff feedback.

### Request

```json
{
  "conversation_id": "conv_01JABC123",
  "message_id": "msg_01JABC789",
  "rating": "helpful",
  "comment": "The answer was clear and the source link helped."
}
```

### Response

```json
{
  "status": "received"
}
```

---

## Health check

```http
GET /healthz
```

Response:

```json
{
  "status": "ok",
  "trace_id": "trace-demo"
}
```

The API echoes an inbound `X-Trace-ID` header or generates one when it is missing. `/healthz` has no external dependencies.

## Readiness check

```http
GET /readyz
```

Response when all registered dependencies are available:

```json
{
  "status": "ready",
  "trace_id": "trace-demo",
  "dependencies": {}
}
```

Response when a dependency is unavailable:

```json
{
  "status": "not_ready",
  "trace_id": "trace-demo",
  "dependencies": {
    "workflow": "unavailable"
  }
}
```

## OpenAPI starter

```yaml
openapi: 3.0.3
info:
  title: AskOC AI Concierge API
  version: 0.1.0
servers:
  - url: http://localhost:8080/api/v1
paths:
  /chat:
    post:
      summary: Send learner message to AI concierge
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ChatRequest'
      responses:
        '200':
          description: Assistant response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ChatResponse'
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
  schemas:
    ChatRequest:
      type: object
      required: [channel, message]
      properties:
        conversation_id:
          type: string
        channel:
          type: string
          example: web
        message:
          type: string
          maxLength: 2000
        student_id:
          type: string
          pattern: '^S[0-9]{6}$'
    ChatResponse:
      type: object
      properties:
        conversation_id:
          type: string
        trace_id:
          type: string
        answer:
          type: string
        intent:
          type: object
          properties:
            name:
              type: string
            confidence:
              type: number
              format: float
        sentiment:
          type: string
        sources:
          type: array
          items:
            type: object
            properties:
              id:
                type: string
              title:
                type: string
              url:
                type: string
              chunk_id:
                type: string
              confidence:
                type: number
                format: float
              risk_level:
                type: string
              freshness_status:
                type: string
              caution:
                type: string
        actions:
          type: array
          items:
            type: object
            properties:
              type:
                type: string
              status:
                type: string
              message:
                type: string
        escalation:
          nullable: true
          type: object
```

## Go implementation notes

- Keep API models in `internal/domain`.
- Use `encoding/json` with strict validation where possible.
- Use `context.Context` for all downstream calls.
- Set HTTP client timeouts for LLM, retrieval, and tool calls.
- Return safe user-facing errors; log internal errors separately.
- Include `trace_id` in every response and audit event.
- Redact sensitive values before writing logs, audit events, session messages, or CRM summaries.
