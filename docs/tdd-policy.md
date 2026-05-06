# Strict TDD Policy for Go Code Tasks

This project uses a strict test-driven development workflow for every code task. The goal is not only higher code quality; it is also strong applicant evidence for maintainable AI/automation systems.

## Core rule

No production Go code should be written for a task until at least one relevant test fails for the expected reason.

For each code task:

1. **Red:** write the smallest meaningful failing test.
2. **Verify red:** run the narrow test and confirm the failure proves the missing behavior.
3. **Green:** write the smallest production code needed to pass.
4. **Verify green:** run the narrow test and then `go test ./...`.
5. **Refactor:** improve naming, structure, duplication, and interfaces while tests stay green.
6. **Document:** update the related Markdown file when behavior, command usage, or API shape changes.

## Required evidence before checking off a code task

A code task is not done unless all of the following are true:

- [ ] A test exists for the new behavior.
- [ ] The test failed before production code was added or changed.
- [ ] The test passes after implementation.
- [ ] `go test ./...` passes.
- [ ] Any handler/client code has negative-path tests.
- [ ] Any workflow/tool action has audit or traceability tests.
- [ ] Any AI/classifier/RAG behavior has safe fallback tests.
- [ ] Any logging/session/audit behavior has redaction tests where relevant.
- [ ] The task documentation or API spec is updated when behavior changes.

## Test categories by code type

| Code type | Required tests |
|---|---|
| Pure Go logic | Table-driven unit tests |
| HTTP handlers | `httptest` success, validation, auth, method, and error tests |
| HTTP clients | `httptest.Server`, timeout, 4xx, 5xx, malformed JSON, context cancellation |
| Orchestrator flows | Fakes/mocks for every dependency; action assertions; no network calls |
| RAG ingestion | Allowlist, HTML cleaning, chunking, metadata, stale-source behavior |
| Classifier/LLM parsing | Valid JSON, invalid JSON, unknown intent, low confidence, unsafe tool trigger prevention |
| Privacy | Redaction, false positives, false negatives, no raw PII in logs/audit/CRM summaries |
| Workflow automation | Idempotency, retry policy, duplicate prevention, audit event creation |
| Evaluation runner | Dataset parsing, scoring, report generation, non-zero critical failure |
| CLI commands | Argument/config validation, deterministic fake provider mode, exit codes |

## Recommended Go test commands

Use narrow commands during Red/Green, then broader commands before checking off a task.

```bash
# Run all tests
go test ./...

# Run one package
go test ./internal/orchestrator

# Run one test by name
go test ./internal/orchestrator -run TestTranscriptStatus_UnpaidTriggersWorkflow

# Race detection for concurrent code
go test -race ./...

# Coverage snapshot
go test ./... -coverprofile=coverage.out

# Optional vet check
go vet ./...
```

## Table-driven test pattern

```go
func TestClassifyIntent(t *testing.T) {
    tests := []struct {
        name       string
        message    string
        wantIntent Intent
        wantUrgent bool
    }{
        {name: "transcript request", message: "How do I order my transcript?", wantIntent: IntentTranscriptRequest},
        {name: "urgent frustration", message: "This is frustrating and I need it today", wantIntent: IntentTranscriptStatus, wantUrgent: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewFallbackClassifier().Classify(context.Background(), tt.message)
            if got.Intent != tt.wantIntent {
                t.Fatalf("intent = %v, want %v", got.Intent, tt.wantIntent)
            }
            if got.Urgent != tt.wantUrgent {
                t.Fatalf("urgent = %v, want %v", got.Urgent, tt.wantUrgent)
            }
        })
    }
}
```

## HTTP handler test pattern

```go
func TestChatHandlerRejectsMissingMessage(t *testing.T) {
    handler := NewChatHandler(fakeOrchestrator{})

    req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(`{"message":""}`))
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusBadRequest {
        t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
    }
    if strings.Contains(rec.Body.String(), "panic") {
        t.Fatalf("response leaked internal details: %s", rec.Body.String())
    }
}
```

## Orchestrator fake dependency rule

Orchestrator tests must not call real HTTP services, live LLM providers, live search indexes, or Power Automate. Use small fakes that record calls.

Example assertion style:

```go
if fakeWorkflow.calls != 1 {
    t.Fatalf("workflow calls = %d, want 1", fakeWorkflow.calls)
}
if fakeCRM.calls != 0 {
    t.Fatalf("CRM calls = %d, want 0 for unpaid/no-hold case", fakeCRM.calls)
}
```

## AI and RAG safety gates

AI-adjacent code must be tested as deterministic software. Do not rely on manual prompt inspection only.

Required gates:

- [ ] The assistant refuses or escalates when retrieval confidence is too low.
- [ ] Policy/procedure answers include source metadata or a safe fallback.
- [ ] Invalid JSON from the model cannot trigger tools.
- [ ] Tool execution requires validated intent and confidence threshold.
- [ ] Stale or high-risk sources produce caution or escalation behavior.
- [ ] Critical hallucination tests fail the evaluation gate.

## Definition of done for pull requests

Before merging or marking a phase complete:

- [ ] Relevant task checkbox is updated in `docs/phases-and-tasks.md`.
- [ ] Prompt used, tests written, and commands run are summarized in the commit or PR notes.
- [ ] `go test ./...` passes.
- [ ] `go vet ./...` passes or known exceptions are documented.
- [ ] `make eval` passes when the task touches AI, RAG, classifier, orchestration, workflow, privacy, or prompts.
- [ ] No real secrets, real learner records, private URLs, or raw PII were added.

## Anti-patterns to reject

- Writing handlers first and adding tests later.
- Mocking the function under test instead of its dependencies.
- Using live APIs in unit tests.
- Marking tests as skipped to make a phase pass.
- Logging raw learner messages before redaction.
- Letting LLM output directly call tools without validation.
- Adding fine-tuning before proving RAG, routing, and evaluation.
- Hiding known limitations in the portfolio docs.
