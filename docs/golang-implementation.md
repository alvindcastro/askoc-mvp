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
cmd/workflow-sim
cmd/ingest
cmd/eval
```

Each command should have a small `main.go` that loads configuration, creates dependencies, registers routes, and starts an HTTP server.

## Core packages

| Package | Responsibility |
|---|---|
| `internal/config` | Environment variables and typed config |
| `internal/domain` | Shared request, response, intent, source, action, case, and student models |
| `internal/handlers` | HTTP handlers and JSON encoding/decoding |
| `internal/orchestrator` | Main chat decision workflow |
| `internal/rag` | Ingestion, chunking, retrieval, source metadata |
| `internal/llm` | LLM provider interface and Azure/OpenAI-compatible REST client |
| `internal/classifier` | Intent and sentiment classification |
| `internal/tools` | Banner, payment, CRM, LMS, notification clients |
| `internal/workflow` | Power Automate webhook client and local workflow simulator client |
| `internal/privacy` | PII redaction, prompt-injection checks, safe summaries |
| `internal/audit` | Audit event writing and dashboard summaries |
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
func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST.")
        return
    }

    var req domain.ChatRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON.")
        return
    }

    if strings.TrimSpace(req.Message) == "" {
        writeError(w, http.StatusBadRequest, "missing_message", "Message is required.")
        return
    }

    resp, err := h.Orchestrator.HandleChat(r.Context(), req)
    if err != nil {
        h.Logger.Error("chat failed", "error", err)
        writeError(w, http.StatusInternalServerError, "chat_failed", "The assistant could not process this message.")
        return
    }

    writeJSON(w, http.StatusOK, resp)
}
```

## Orchestrator decision flow

```text
1. Create trace ID.
2. Redact message for logging.
3. Classify intent and sentiment.
4. Retrieve sources for policy/procedure questions.
5. If source confidence is too low, produce fallback and/or escalate.
6. If transcript-status intent and student ID exists, call mock Banner and payment APIs.
7. If unpaid, trigger payment reminder workflow.
8. If hold, negative sentiment, urgent context, or low confidence, create CRM case.
9. Generate grounded response.
10. Write audit events.
11. Return structured response.
```

## RAG ingestion command

`cmd/ingest` should:

1. Load `data/seed-sources.json`.
2. Fetch allowlisted public URLs.
3. Strip navigation and boilerplate where possible.
4. Chunk text into 500–900 token-equivalent chunks.
5. Store metadata: URL, title, retrieved date, content hash, risk level.
6. Generate embeddings if using vector search.
7. Write chunks to PostgreSQL/pgvector, Azure AI Search, or local JSON for the MVP.

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

For classification, require a strict JSON response:

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
- clamp confidence values,
- default unknown intent to safe fallback,
- never execute tools based only on low-confidence classification.

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
  "amount_due": 10.00,
  "currency": "CAD",
  "status": "unpaid"
}
```

## Workflow client

```go
type PaymentReminderRequest struct {
    TraceID       string `json:"trace_id"`
    StudentID     string `json:"student_id"`
    ConversationID string `json:"conversation_id"`
    Reason        string `json:"reason"`
    IdempotencyKey string `json:"idempotency_key"`
}
```

Use idempotency keys to avoid duplicate reminders:

```text
payment-reminder:S100002:official-transcript:2026-05-06
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
        {"paid no hold", "S100001", "paid", "none", "status_checked", false},
        {"unpaid", "S100002", "unpaid", "none", "payment_reminder_sent", false},
        {"financial hold", "S100003", "paid", "financial", "crm_case_created", true},
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
go test -race ./...
go run ./cmd/api
go run ./cmd/mock-banner
go run ./cmd/mock-payment
go run ./cmd/mock-crm
go run ./cmd/workflow-sim
go run ./cmd/ingest -sources data/seed-sources.json
go run ./cmd/eval -input data/eval-questions.jsonl
```

## Makefile targets

```makefile
.PHONY: dev test race ingest eval

dev:
	docker compose up --build

test:
	go test ./...

race:
	go test -race ./...

ingest:
	go run ./cmd/ingest -sources data/seed-sources.json

eval:
	go run ./cmd/eval -input data/eval-questions.jsonl
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

