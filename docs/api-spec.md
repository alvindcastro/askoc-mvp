# API Specification

## Overview

This document defines a simple REST API surface for the AskOC AI Concierge MVP. The API is designed for a local or demo environment using Go services and synthetic data.

## Base URL

```text
http://localhost:8080/api/v1
```

The implemented P2 chat route is `POST http://localhost:8080/api/v1/chat`. The web chat UI is served at `GET http://localhost:8080/chat`.

## Authentication

For the MVP, use a mock bearer token.

```http
Authorization: Bearer demo-token
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
| `POST` | `/automation/payment-reminder` | Trigger mock payment reminder workflow | `cmd/workflow-sim` or Power Automate |
| `GET` | `/analytics/summary` | Get dashboard summary metrics | `cmd/api` |
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
  "answer": "This P2 demo placeholder received your transcript status question. Synthetic Banner and payment checks will be added in later phases.",
  "intent": {
    "name": "transcript_status",
    "confidence": 0.6
  },
  "sentiment": "neutral",
  "sources": [
    {
      "title": "Transcript Request - 2005 Onwards",
      "url": "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
      "chunk_id": "placeholder-transcript-source"
    }
  ],
  "actions": [
    {
      "type": "placeholder_response",
      "status": "completed",
      "message": "No live AI, retrieval, or enterprise system was called in this P2 placeholder."
    }
  ],
  "escalation": null
}
```

P2 validation rules:

- `message` is required after trimming whitespace.
- `message` must be 2000 characters or fewer.
- `student_id` is optional, but when present it must use the synthetic demo shape `S` plus six digits, such as `S100002`.
- invalid JSON, validation failures, and service failures return the common safe error shape and never echo raw request bodies.
- P2 does not call live AI, retrieval, Banner, payment, CRM, LMS, or workflow services.

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
  "program": "Business Administration Demo Program",
  "holds": [],
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
  "transcript_request_status": "requested",
  "eligible_for_processing": true,
  "hold": "none",
  "last_updated": "2026-05-06T12:00:00Z"
}
```

### Response: hold

```json
{
  "student_id": "S100003",
  "transcript_request_status": "requested",
  "eligible_for_processing": false,
  "hold": "financial",
  "last_updated": "2026-05-06T12:00:00Z"
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
  "amount_due": 10.00,
  "currency": "CAD",
  "status": "unpaid",
  "last_updated": "2026-05-06T12:00:00Z"
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
  "last_updated": "2026-05-06T12:00:00Z"
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
  "case_id": "CASE-2026-000123",
  "status": "created",
  "queue": "registrar_finance",
  "priority": "high",
  "created_at": "2026-05-06T12:05:00Z"
}
```

---

## `POST /automation/payment-reminder`

Triggers a payment reminder workflow.

This endpoint can be implemented by:

1. `cmd/workflow-sim` for local development, or
2. Power Automate cloud flow with an HTTP request trigger.

### Request

```json
{
  "student_id": "S100002",
  "conversation_id": "conv_01JABC123",
  "trace_id": "trace_01JABC456",
  "item": "official_transcript",
  "reason": "Transcript request cannot be processed until payment is complete.",
  "idempotency_key": "payment-reminder:S100002:official-transcript:2026-05-06"
}
```

### Response

```json
{
  "workflow_id": "WF-2026-000789",
  "status": "accepted",
  "message": "Payment reminder workflow accepted.",
  "created_at": "2026-05-06T12:06:00Z"
}
```

---

## `GET /analytics/summary`

Returns dashboard summary metrics.

### Response

```json
{
  "total_conversations": 48,
  "containment_rate": 0.67,
  "escalation_rate": 0.18,
  "average_response_ms": 1840,
  "top_intents": [
    { "intent": "transcript_request", "count": 15 },
    { "intent": "transcript_status", "count": 11 },
    { "intent": "fee_payment", "count": 7 }
  ],
  "automation": {
    "payment_reminders_sent": 6,
    "workflow_failures": 0
  },
  "review_queue": {
    "low_confidence_answers": 3,
    "stale_source_questions": 1
  }
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
              title:
                type: string
              url:
                type: string
              chunk_id:
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
- Redact sensitive values before writing logs or CRM summaries.
