# Go Implementation Guide

## Goal

This guide translates the AskOC AI Concierge MVP into a practical Go build plan. It is written for an applicant portfolio project where the reviewer should quickly see strong backend, integration, automation, privacy, and AI engineering skills.

## Recommended Go style

Use a boring, production-like Go structure:

- `cmd/` for runnable services,
- `internal/` for application packages,
- typed domain models,
- interfaces at orchestration boundaries,
- context-aware outbound calls,
- table-driven tests,
- structured logging,
- explicit errors and safe fallbacks.

## Service commands

```text
cmd/api
cmd/mock-banner
cmd/mock-payment
cmd/mock-crm
cmd/mock-lms
cmd/workflow-sim   # later P8
cmd/ingest         # P5 local RAG ingestion
cmd/eval           # later P9
```

Each command should have a small `main.go` that loads configuration, creates dependencies, registers routes, and starts an HTTP server.

## Core packages

| Package | Responsibility |
|---|---|
| `internal/config` | Environment variables and typed config |
| `internal/domain` | Shared request, response, intent, source, action, case, and student models |
| `internal/handlers` | HTTP handlers and JSON encoding/decoding |
| `internal/session` | In-memory demo conversation sessions with TTL and redaction before persistence |
| `internal/validation` | Chat request validation and safe validation error codes |
| `internal/orchestrator` | P6 guarded chat decision workflow, prompt templates, source packaging, and dependency ports |
| `internal/rag` | P5 allowlist parsing, ingestion, chunking, local retrieval, and source freshness metadata |
| `internal/llm` | P6 provider-neutral request/response types and Azure/OpenAI-compatible REST client |
| `internal/classifier` | P6 deterministic fallback, strict JSON parser, and fixture-backed intent/sentiment tests |
| `internal/tools` | Banner, payment, CRM, LMS, notification clients |
| `internal/workflow` | P4 in-process idempotent workflow port; later P8 webhook/simulator clients |
| `internal/privacy` | Later P7 PII redaction, prompt-injection checks, safe summaries |
| `internal/audit` | P4 audit event port types; later P7 event store and dashboard summaries |
| `internal/middleware` | Trace IDs, recovery, logging, auth, rate limits |

## Dependency rule

Keep dependencies pointing inward:

```text
handlers -> orchestrator -> classifier/rag/llm/tools/workflow/audit/privacy
```

Avoid letting tool clients call handlers or UI code. This keeps tests simple.

## Interfaces

### Orchestrator dependencies

```go
type Retriever interface {
    Retrieve(ctx context.Context, query string, limit int) ([]domain.SourceChunk, error)
}

type LLM interface {
    GenerateGroundedAnswer(ctx context.Context, req llm.AnswerRequest) (llm.AnswerResult, error)
    Classify(ctx context.Context, message string) (classifier.Result, error)
}

type StudentClient interface {
    GetStudent(ctx context.Context, studentID string) (domain.Student, error)
    GetTranscriptStatus(ctx context.Context, studentID string) (domain.TranscriptStatus, error)
}

type PaymentClient interface {
    GetTranscriptPayment(ctx context.Context, studentID string) (domain.PaymentStatus, error)
}

type CRMClient interface {
    CreateCase(ctx context.Context, req domain.CreateCaseRequest) (domain.Case, error)
}

type WorkflowClient interface {
    TriggerPaymentReminder(ctx context.Context, req domain.PaymentReminderRequest) (domain.WorkflowResult, error)
}

type Auditor interface {
    Record(ctx context.Context, event domain.AuditEvent) error
}
```

This lets you test the orchestrator without real network calls.

## HTTP server pattern

Use standard Go HTTP patterns or a small router.

```go
func main() {
    cfg := config.MustLoad()
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    app := bootstrap.NewApp(cfg, logger)

    mux := http.NewServeMux()
    handlers.RegisterRoutes(mux, app)

    srv := &http.Server{
        Addr:              cfg.HTTPAddr,
        Handler:           middleware.Chain(mux, middleware.TraceID, middleware.Recover(logger), middleware.Auth(cfg)),
        ReadHeaderTimeout: 5 * time.Second,
    }

    logger.Info("api starting", "addr", cfg.HTTPAddr)
    if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        logger.Error("api failed", "error", err)
        os.Exit(1)
    }
}
```

## Chat handler pattern

```go
func ChatHandler(service ChatService) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
            return
        }

        var req domain.ChatRequest
        decoder := json.NewDecoder(r.Body)
        decoder.DisallowUnknownFields()
        if err := decoder.Decode(&req); err != nil {
            WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
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
```

P6 wires this handler to the guarded orchestrator at `POST /api/v1/chat`. By default the orchestrator returns deterministic fallback intent/sentiment classification, local RAG source packaging for transcript-request answers, mock Banner/payment/CRM actions, idempotent payment-reminder workflow results, and optional handoff metadata without calling live AI. When `ASKOC_PROVIDER=openai-compatible` is explicitly configured, the same orchestrator uses the tested REST LLM gateway behind strict JSON parsing, prompt-version metadata, source-only answer validation, and deterministic fallback on timeout or unsafe output. The Go UI is served at `/chat`.

## Orchestrator decision flow

```text
1. Create trace ID.
2. Redact message for logging.
3. Classify intent and sentiment with deterministic fallback or guarded LLM JSON output.
4. Reject invalid, unknown, or low-confidence classification before sensitive tool checks.
5. Retrieve approved source chunks for transcript-request answers and attach source confidence/risk/freshness metadata.
6. If LLM mode is enabled, generate answer text only after source guardrails pass.
7. If transcript-status intent and student ID exists, call mock Banner and payment APIs.
8. If unpaid, trigger payment reminder workflow.
9. If hold, negative sentiment, urgent context, or low confidence, create CRM case.
10. Record workflow and guardrail audit-port events.
11. Return structured response.
```

## RAG ingestion command

P5 `cmd/ingest`:

1. Load `data/seed-sources.json`.
2. Fetch allowlisted public URLs.
3. Strip navigation and boilerplate where possible.
4. Chunk text with configurable `-max-words` and `-overlap-words`.
5. Store metadata: source ID, URL, title, retrieved date, content hash, risk level, and freshness status.
6. Write chunks to local JSON at `data/rag-chunks.json` by default.
7. Leave embeddings, PostgreSQL/pgvector, and Azure AI Search for later implementation.

Example source config:

```json
[
  {
    "url": "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
    "department": "Registrar",
    "risk_level": "medium",
    "requires_freshness_check": true
  }
]
```

## LLM structured output

For classification, P6 requires a strict JSON response:

```json
{
  "intent": "transcript_status",
  "intent_confidence": 0.91,
  "sentiment": "negative",
  "urgency": "high",
  "needs_handoff": true,
  "reason": "Learner reports blocked transcript and urgent job deadline."
}
```

Validate model output before using it:

- reject invalid JSON,
- reject out-of-range confidence values,
- default unknown intent to safe fallback,
- never execute tools based only on low-confidence classification.

The fixture target for P6 classification is 100% intent accuracy on the synthetic demo set in `data/classification-fixtures.jsonl`; `internal/classifier/e2e_test.go` fails with the fixture ID and expected intent when a regression occurs.

## Tool calling without magic

Do not let the LLM call arbitrary tools directly. The Go orchestrator should decide based on typed logic.

```go
switch result.Intent {
case domain.IntentTranscriptStatus:
    return o.handleTranscriptStatus(ctx, req, result)
case domain.IntentHumanHandoff:
    return o.createHandoff(ctx, req, result)
default:
    return o.handleGroundedAnswer(ctx, req, result)
}
```

## Mock Banner endpoint

```http
GET /api/v1/students/S100002/transcript-status
```

Response:

```json
{
  "student_id": "S100002",
  "status": "requested",
  "hold": "none",
  "eligible": true
}
```

## Mock Payment endpoint

```http
GET /api/v1/students/S100002/payment-status
```

Response:

```json
{
  "student_id": "S100002",
  "item": "official_transcript",
  "amount_due": 15.00,
  "currency": "CAD",
  "status": "unpaid"
}
```

## Mock LMS endpoint

```http
GET /api/v1/students/S100001/lms-access?course_id=DEMO-LMS-101
```

Response:

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

## Workflow client

```go
type PaymentReminderRequest struct {
    StudentID      string  `json:"student_id"`
    ConversationID string  `json:"conversation_id,omitempty"`
    TraceID        string  `json:"trace_id,omitempty"`
    Item           string  `json:"item"`
    AmountDue      float64 `json:"amount_due,omitempty"`
    Currency       string  `json:"currency,omitempty"`
    Reason         string  `json:"reason"`
    IdempotencyKey string  `json:"idempotency_key"`
}
```

Use idempotency keys to avoid duplicate reminders:

```text
payment-reminder:trace_01JABC456:S100002:official_transcript
```

## Privacy redaction

Use redaction before logging:

```go
func Redact(input string) string {
    input = emailRegex.ReplaceAllString(input, "[REDACTED_EMAIL]")
    input = phoneRegex.ReplaceAllString(input, "[REDACTED_PHONE]")
    input = passwordRegex.ReplaceAllString(input, "[REDACTED_SECRET]")
    return input
}
```

For the MVP, retain only redacted conversation text and minimal synthetic identifiers.

## Testing strategy

Use table-driven tests for orchestration:

```go
func TestTranscriptStatusWorkflow(t *testing.T) {
    tests := []struct {
        name        string
        studentID   string
        payment     string
        hold        string
        wantAction  string
        wantEscalate bool
    }{
        {"paid no hold", "S100001", "paid", "none", "payment_reminder_skipped", false},
        {"unpaid with payment hold", "S100002", "unpaid", "mock_payment_hold", "payment_reminder_triggered", false},
        {"financial hold", "S100003", "review_required", "mock_financial_hold", "crm_case_created", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange fake clients.
            // Act by calling orchestrator.HandleChat.
            // Assert actions, escalation, and audit events.
        })
    }
}
```

## Local commands

```bash
go test ./...
go test ./internal/classifier ./internal/workflow ./internal/orchestrator
go test ./internal/rag ./internal/orchestrator
go test -race ./internal/session
go run ./cmd/api
go run ./cmd/ingest -sources data/seed-sources.json -out data/rag-chunks.json
go run ./cmd/mock-banner
go run ./cmd/mock-payment
go run ./cmd/mock-crm
go run ./cmd/mock-lms
```

`cmd/workflow-sim` and `cmd/eval` are later-phase commands.

## Makefile targets

```makefile
.PHONY: dev test test-race

dev:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./internal/session
```

## What to show in GitHub

A reviewer should see:

- idiomatic Go folder structure,
- clean interfaces,
- context-aware HTTP clients,
- safe error handling,
- redacted structured logs,
- table-driven tests,
- Docker Compose demo,
- concise docs and diagrams.

## MVP build shortcut

For the fastest credible build:

1. Start with one Go binary serving API + UI.
2. Implement mock Banner/payment/CRM as in-memory Go services.
3. Store RAG chunks in local JSON first.
4. Add PostgreSQL later if time allows.
5. Use LLM structured output for intent/sentiment.
6. Add Power Automate webhook only after local workflow simulator works.


## Strict TDD workflow

Implementation should follow [TDD Policy](tdd-policy.md) and the task-level prompts in [Task Prompts](task-prompts.md). For every Go code task:

- [ ] create failing unit/handler/client/orchestrator tests first,
- [ ] verify the red state,
- [ ] implement minimal code,
- [ ] run the narrow package test,
- [ ] run `go test ./...`,
- [ ] refactor only while green.

Use `httptest.Server` for clients, fakes for orchestrator dependencies, table-driven tests for pure logic, and redaction tests for any log/audit/session behavior.
