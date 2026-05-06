# Changelog

All notable MVP task changes are recorded here with what changed, where it changed, when it changed, why it changed, and how it was completed.

## 2026-05-06 - P2 Chat API And UI Skeleton

### P2-T01 - Define chat domain models

- What: added provider-neutral chat request/response domain models for intents, sentiment, source citations, action traces, and escalation metadata.
- Where: `internal/domain/chat.go`, `internal/domain/chat_test.go`, `docs/api-spec.md`, `docs/architecture.md`.
- When: 2026-05-06.
- Why: establish the stable API contract before real retrieval, AI orchestration, and enterprise integrations are added.
- How: wrote JSON round-trip and validation-facing tests first, confirmed the missing package failure, then implemented typed Go structs and constants.

### P2-T02 - Implement `POST /api/v1/chat` handler with deterministic placeholder response

- What: added the chat HTTP handler, chat service interface, deterministic placeholder service, trace ID propagation, strict JSON decoding, and safe service-error handling.
- Where: `internal/handlers/chat.go`, `internal/handlers/chat_test.go`, `cmd/api/main.go`, `docs/api-spec.md`, `docs/golang-implementation.md`.
- When: 2026-05-06.
- Why: expose the public chat contract without relying on live AI, RAG, workflow, or mock enterprise services.
- How: wrote `httptest` coverage for success, invalid JSON, missing message, unsupported method, validation errors, and safe service failure before implementing the handler and placeholder.

### P2-T03 - Serve minimal Go web chat UI

- What: added a Go-rendered chat page, static JavaScript, static CSS, and routes for `/chat`, `/`, and `/static/`.
- Where: `web/templates/chat.html`, `web/static/app.js`, `web/static/app.css`, `internal/handlers/ui.go`, `internal/handlers/ui_test.go`, `cmd/api/main.go`, `README.md`.
- When: 2026-05-06.
- Why: provide an interview-friendly usable chat surface while keeping frontend complexity low.
- How: wrote rendering and static-content tests first, confirmed missing handler/static asset failures, then added the template, static assets, and route registration.

### P2-T04 - Add request validation and safe client-facing errors

- What: added chat request validation for missing, whitespace-only, oversized messages and synthetic-only student ID shape.
- Where: `internal/validation/chat.go`, `internal/validation/chat_test.go`, `internal/handlers/chat.go`, `internal/handlers/chat_test.go`, `docs/api-spec.md`, `README.md`.
- When: 2026-05-06.
- Why: keep malformed requests and accidental non-demo identifiers from destabilizing the local demo or leaking raw input.
- How: added table-driven validation tests and handler assertions that error responses do not echo request bodies, then implemented safe validation codes and messages.

### P2-T05 - Add in-memory conversation session store

- What: added a concurrency-safe in-memory session store with configurable TTL, create/append/read/expire behavior, and redaction before message persistence.
- Where: `internal/session/store.go`, `internal/session/store_test.go`, `internal/handlers/chat.go`, `Makefile`, `docs/test-plan.md`, `README.md`.
- When: 2026-05-06.
- Why: track short synthetic demo conversations so follow-up behavior can be layered in later phases without storing raw PII.
- How: wrote create/append/read, expiration, redaction, and concurrent access tests first, confirmed missing package failures, then implemented the mutex-protected store and `make test-race` target.

### P2 review evidence

- What: completed P2 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API surface, local commands, and verification expectations aligned with the implemented chat skeleton.
- How: marked P2 tasks and gates complete after the red/green cycle and verified `go test ./internal/domain ./internal/validation ./internal/handlers ./internal/session`, `go test -race ./internal/session`, `go test ./...`, `go test -race ./...`, `make test`, `make test-race`, and `go vet ./...` pass.

## 2026-05-06 - P1 Go Project Foundation

### P1-T01 - Initialize Go module, repository layout, and developer commands

- What: created the initial Go module, API skeleton, build smoke test, and developer commands.
- Where: `go.mod`, `Makefile`, `cmd/api/main.go`, `internal/build/smoke_test.go`.
- When: 2026-05-06.
- Why: establish a compilable Go workspace and predictable local commands before feature work.
- How: added failing compile/Makefile smoke tests first, confirmed `go test ./...` failed for missing implementation, then added the minimal module, `make test`, and `make dev` targets.

### P1-T02 - Build typed configuration loader

- What: added typed environment configuration with defaults, overrides, validation, and redacted string output.
- Where: `internal/config/config.go`, `internal/config/config_test.go`, `README.md`.
- When: 2026-05-06.
- Why: make HTTP, auth, logging, workflow, and provider settings explicit without leaking secrets.
- How: added table-driven tests for defaults, overrides, invalid booleans/timeouts/log levels, and secret redaction, then implemented `Load`, `LoadFromEnv`, and safe formatting.

### P1-T03 - Implement health and readiness endpoints

- What: added `GET /healthz` and `GET /readyz` handlers with JSON responses, method checks, trace IDs, and safe dependency status reporting.
- Where: `internal/handlers/health.go`, `internal/handlers/health_test.go`, `docs/api-spec.md`, `README.md`.
- When: 2026-05-06.
- Why: provide operational endpoints for local demo startup checks and later Docker/deployment health checks.
- How: added `httptest` coverage for healthy responses, readiness dependency success/failure, 405 responses, trace propagation, and non-leaky dependency failures before implementing handlers.

### P1-T04 - Add middleware for trace ID, panic recovery, request logging, and mock auth

- What: added middleware chaining, trace ID propagation, panic recovery, mock bearer auth, request logging, and basic redaction.
- Where: `internal/middleware/chain.go`, `internal/middleware/trace.go`, `internal/middleware/recover.go`, `internal/middleware/auth.go`, `internal/middleware/logging.go`, `internal/middleware/*_test.go`, `docs/task-prompts.md`, `docs/phases-and-tasks.md`.
- When: 2026-05-06.
- Why: establish API hygiene before chat, orchestration, and mock integration endpoints are added.
- How: wrote failing tests for trace preservation/generation, safe panic conversion, auth enabled/disabled behavior, and logging redaction hooks, then implemented standard-library middleware.

### P1-T05 - Create JSON response and error helpers

- What: added response helpers for successful JSON payloads and stable safe API errors.
- Where: `internal/handlers/respond.go`, `internal/handlers/respond_test.go`, `internal/handlers/health.go`.
- When: 2026-05-06.
- Why: standardize response shape for later handlers and avoid leaking raw Go errors to clients.
- How: added tests for success headers/body, `{error:{code,message,trace_id}}` error shape, and unsupported-value fallback before implementing helper functions.

### P1 review evidence

- What: completed P1 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/task-prompts.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API surface, and local commands aligned with the implemented Go foundation.
- How: marked P1 tasks and gates complete after the red/green cycle and verified `go test ./...`, `make test`, and `go vet ./...` pass.

## 2026-05-06 - P0 Product Framing And Applicant Strategy

### P0-T01 - Write the applicant story and MVP thesis

- What: tightened the applicant story and one-sentence pitch around a Go-based learner-service agent using RAG, workflow automation, mock enterprise integrations, and privacy-first synthetic data controls.
- Where: `README.md`, `docs/mvp-scope.md`.
- When: 2026-05-06.
- Why: make the project purpose, learner pain point, and AI/Automation Solutions Developer role mapping explicit before implementation starts.
- How: updated the README pitch and scope language, froze the primary transcript/payment workflow, and removed wording that could imply access to real OC systems.

### P0-T02 - Define synthetic data and privacy boundary

- What: documented the synthetic-data-only boundary and created the synthetic learner fixture.
- Where: `docs/privacy-impact-lite.md`, `data/synthetic-students.json`, `README.md`.
- When: 2026-05-06.
- Why: make it clear that learner records, IDs, payments, transcript states, and CRM cases are invented demo artifacts only.
- How: added fake-data markers, fixture rules, synthetic ID patterns, and four visibly fake demo records using `S10000X`, `SYNTH-*`, and `MOCK-*` identifiers.

### P0-T03 - Confirm source allowlist and knowledge-domain limits

- What: created the public source allowlist and documented retrieval boundaries.
- Where: `data/seed-sources.json`, `docs/source-references.md`, `docs/privacy-impact-lite.md`.
- When: 2026-05-06.
- Why: prevent private portal scraping, unapproved source ingestion, stale-source overconfidence, and learner-specific data leakage.
- How: listed approved public Okanagan College URLs already present in repo docs, spot-checked them as accessible public pages on 2026-05-06, separated implementation references from learner-service RAG sources, added freshness metadata, and defined fallback behavior for stale or missing sources.

### P0-T04 - Create demo acceptance matrix

- What: turned the interview demo into measurable acceptance scenarios.
- Where: `docs/demo-script.md`, `docs/model-evaluation.md`.
- When: 2026-05-06.
- Why: ensure the transcript answer, unpaid payment workflow, financial-hold escalation, and urgent sentiment escalation have observable pass criteria before code implementation.
- How: added D01-D05 demo cases with expected intent, source, action, handoff behavior, and pass evidence; aligned source checks with the source allowlist fixture.

### P0-T05 - Freeze MVP scope and defer nice-to-haves

- What: froze the MVP around transcript/payment support and deferred nonessential workflows and real integrations.
- Where: `docs/mvp-scope.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/phases-and-tasks.md`.
- When: 2026-05-06.
- Why: keep the MVP narrow enough for strict TDD delivery and avoid overbuilding beyond the applicant demo.
- How: marked P0 complete, updated phase gates, clarified that non-transcript learner-service topics use fallback or handoff, and documented that real authentication, Banner, payment, CRM, LMS, and private portal integrations are out of scope.

### Review evidence

- What: completed documentation review checks for all P0 tasks.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `data/synthetic-students.json`, `data/seed-sources.json`, and the updated P0 Markdown files.
- When: 2026-05-06.
- Why: P0 contains documentation tasks only, so Go failing-test evidence and `go test ./...` are not applicable.
- How: used JSON validation, public URL spot-checks, source/fixture consistency checks, targeted content searches for required and prohibited terms, and Markdown whitespace checks.
