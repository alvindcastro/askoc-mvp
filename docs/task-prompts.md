# Detailed Task Prompts

Use these prompts with a coding assistant or as self-instructions while implementing the Go-based AskOC AI Concierge. Each prompt is scoped to one task from [Phases and Tickable Tasks](phases-and-tasks.md).

For code tasks, the prompt explicitly enforces strict TDD. Do not remove that section when using the prompt.

## Global context to include with every prompt

```text
Project: AskOC AI Concierge.
Language: Go/Golang.
Purpose: Applicant portfolio MVP for an AI/Automation Solutions Developer role in higher education.
Core demo: AI learner-service assistant answers transcript questions, checks synthetic Banner/payment records, triggers a workflow reminder, and escalates urgent or unresolved cases to a mock CRM.
Data rule: Use synthetic learner data only. Do not use real student data, real payment data, private portals, or secrets.
Architecture rule: Prefer small Go packages, interfaces at boundaries, context-aware clients, safe errors, redacted logs, and tests that do not call live external services.
TDD rule: For code tasks, write failing tests first, verify failure, implement the smallest production code, run tests, then refactor.
```

---


## P0-T01 — Write the applicant story and MVP thesis

**Type:** Documentation  
**Goal:** State why the project exists, what learner pain point it addresses, and how it maps to the AI/Automation Solutions Developer role.

### Copy/paste prompt

```text
You are helping implement task P0-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Write the applicant story and MVP thesis
Task goal: State why the project exists, what learner pain point it addresses, and how it maps to the AI/Automation Solutions Developer role.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `README.md`
- `docs/mvp-scope.md`

Required tests or review checks:
- Peer-read the one-sentence pitch and confirm it mentions Go, RAG, workflow automation, mock enterprise integrations, and privacy.

Acceptance criteria:
- README contains a one-sentence pitch.
- MVP scope states the primary transcript/payment workflow.
- The story avoids claiming access to real OC systems.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P0-T02 — Define synthetic data and privacy boundary

**Type:** Documentation  
**Goal:** Make it explicit that the MVP uses synthetic learners, synthetic IDs, mock payments, and mock CRM cases only.

### Copy/paste prompt

```text
You are helping implement task P0-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Define synthetic data and privacy boundary
Task goal: Make it explicit that the MVP uses synthetic learners, synthetic IDs, mock payments, and mock CRM cases only.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/privacy-impact-lite.md`
- `data/synthetic-students.json`

Required tests or review checks:
- Review all examples for real names, real student IDs, private portal data, or secrets.

Acceptance criteria:
- Synthetic data policy is visible.
- PII boundary is documented.
- Demo records are clearly marked as fake.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P0-T03 — Confirm source allowlist and knowledge-domain limits

**Type:** Documentation  
**Goal:** Define which public pages or manually curated content can be ingested and which content must never be ingested.

### Copy/paste prompt

```text
You are helping implement task P0-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Confirm source allowlist and knowledge-domain limits
Task goal: Define which public pages or manually curated content can be ingested and which content must never be ingested.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `data/seed-sources.json`
- `docs/source-references.md`
- `docs/privacy-impact-lite.md`

Required tests or review checks:
- Manually verify every seed source is public and relevant to learner services.

Acceptance criteria:
- Allowlist exists.
- Private portal scraping is explicitly out of scope.
- Stale-source handling is documented.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P0-T04 — Create demo acceptance matrix

**Type:** Documentation  
**Goal:** Turn the interview demo into measurable scenarios before implementation starts.

### Copy/paste prompt

```text
You are helping implement task P0-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Create demo acceptance matrix
Task goal: Turn the interview demo into measurable scenarios before implementation starts.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/demo-script.md`
- `docs/model-evaluation.md`

Required tests or review checks:
- Walk through the golden path and confirm each expected action has an observable output.

Acceptance criteria:
- Golden path includes transcript answer, unpaid payment workflow, financial-hold escalation, and urgent sentiment escalation.
- Each scenario has expected intent, source, action, and handoff behavior.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P0-T05 — Freeze MVP scope and defer nice-to-haves

**Type:** Documentation  
**Goal:** Prevent overbuilding by separating must-have demo features from optional stretch work.

### Copy/paste prompt

```text
You are helping implement task P0-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Freeze MVP scope and defer nice-to-haves
Task goal: Prevent overbuilding by separating must-have demo features from optional stretch work.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/mvp-scope.md`
- `docs/implementation-roadmap.md`

Required tests or review checks:
- Check that every Phase 1-9 task supports the transcript/payment MVP or demo operations.

Acceptance criteria:
- Nice-to-have features are explicitly deferred.
- No real authentication, payments, or Banner integrations are listed as MVP work.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P1-T01 — Initialize Go module, repository layout, and developer commands

**Type:** Code  
**Goal:** Create a compilable Go workspace with predictable local commands before adding business logic.

### Copy/paste prompt

```text
You are helping implement task P1-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Initialize Go module, repository layout, and developer commands
Task goal: Create a compilable Go workspace with predictable local commands before adding business logic.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `go.mod`
- `Makefile`
- `cmd/api/main.go`
- `internal/config`

Required tests or review checks:
- Add a failing smoke test that runs `go test ./...` and proves the module compiles once minimal code exists.
- Add a Makefile test target that executes `go test ./...`.

Acceptance criteria:
- `go test ./...` passes.
- `make test` passes.
- `make dev` has a documented target, even if it only runs the API skeleton.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P1-T02 — Build typed configuration loader

**Type:** Code  
**Goal:** Load HTTP addresses, auth mode, log level, workflow URL, and provider settings from environment variables with safe defaults.

### Copy/paste prompt

```text
You are helping implement task P1-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Build typed configuration loader
Task goal: Load HTTP addresses, auth mode, log level, workflow URL, and provider settings from environment variables with safe defaults.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/config/config.go`
- `internal/config/config_test.go`

Required tests or review checks:
- Table-driven tests for defaults.
- Table-driven tests for overridden env vars.
- Test invalid numeric/boolean values return clear errors.

Acceptance criteria:
- Config has no global mutable state.
- Secrets are not printed.
- Invalid config fails fast with actionable error messages.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P1-T03 — Implement health and readiness endpoints

**Type:** Code  
**Goal:** Expose simple operational endpoints for local demo, Docker health checks, and future deployment.

### Copy/paste prompt

```text
You are helping implement task P1-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Implement health and readiness endpoints
Task goal: Expose simple operational endpoints for local demo, Docker health checks, and future deployment.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/handlers/health.go`
- `internal/handlers/health_test.go`

Required tests or review checks:
- `httptest` verifies `GET /healthz` returns 200 JSON.
- `httptest` verifies `GET /readyz` returns dependency status.
- Invalid methods return 405.

Acceptance criteria:
- Health response includes status and trace ID when middleware is enabled.
- No external dependencies are required for `/healthz`.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P1-T04 — Add middleware for trace ID, panic recovery, request logging, and mock auth

**Type:** Code  
**Goal:** Create enterprise-style API hygiene before feature work begins.

### Copy/paste prompt

```text
You are helping implement task P1-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Add middleware for trace ID, panic recovery, request logging, and mock auth
Task goal: Create enterprise-style API hygiene before feature work begins.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/middleware/trace.go`
- `internal/middleware/recover.go`
- `internal/middleware/auth.go`
- `internal/middleware/*_test.go`

Required tests or review checks:
- Trace middleware adds or preserves trace ID.
- Recovery middleware converts panic to 500 without exposing internals.
- Mock auth rejects missing token when auth is enabled.
- Logging test verifies redaction hooks are used.

Acceptance criteria:
- All handlers receive context with trace ID.
- Panic stack traces are not returned to clients.
- Auth can be disabled for local demo.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P1-T05 — Create JSON response and error helpers

**Type:** Code  
**Goal:** Standardize API responses so later handlers are easy to test and document.

### Copy/paste prompt

```text
You are helping implement task P1-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Create JSON response and error helpers
Task goal: Standardize API responses so later handlers are easy to test and document.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/handlers/respond.go`
- `internal/handlers/respond_test.go`

Required tests or review checks:
- Test successful JSON response headers/body.
- Test error response shape.
- Test unsupported values are handled safely.

Acceptance criteria:
- Errors follow a stable `{error:{code,message,trace_id}}` shape.
- No raw Go error details leak to user responses.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P2-T01 — Define chat domain models

**Type:** Code  
**Goal:** Create typed request/response structs for messages, intents, sources, actions, and handoff status.

### Copy/paste prompt

```text
You are helping implement task P2-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Define chat domain models
Task goal: Create typed request/response structs for messages, intents, sources, actions, and handoff status.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/domain/chat.go`
- `internal/domain/chat_test.go`

Required tests or review checks:
- JSON marshal/unmarshal tests for `ChatRequest` and `ChatResponse`.
- Validation tests for missing message and invalid student ID shape.

Acceptance criteria:
- Domain structs are provider-neutral.
- Response can include source citations, tool actions, and escalation metadata.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P2-T02 — Implement `POST /api/v1/chat` handler with deterministic placeholder response

**Type:** Code  
**Goal:** Create the public chat API contract before real AI orchestration exists.

### Copy/paste prompt

```text
You are helping implement task P2-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Implement `POST /api/v1/chat` handler with deterministic placeholder response
Task goal: Create the public chat API contract before real AI orchestration exists.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/handlers/chat.go`
- `internal/handlers/chat_test.go`

Required tests or review checks:
- `httptest` verifies valid request returns 200.
- Invalid JSON returns 400.
- Missing message returns 400.
- Unsupported method returns 405.

Acceptance criteria:
- Handler depends on an interface, not concrete orchestrator.
- Trace ID appears in response.
- Tests do not call network or model APIs.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P2-T03 — Serve minimal Go web chat UI

**Type:** Code  
**Goal:** Provide an interview-friendly UI without making frontend complexity the focus.

### Copy/paste prompt

```text
You are helping implement task P2-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Serve minimal Go web chat UI
Task goal: Provide an interview-friendly UI without making frontend complexity the focus.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `web/templates/chat.html`
- `web/static/app.js`
- `web/static/app.css`
- `internal/handlers/ui.go`
- `internal/handlers/ui_test.go`

Required tests or review checks:
- `httptest` verifies chat page renders.
- Template test verifies API endpoint path is present.
- Static file route returns expected content type.

Acceptance criteria:
- User can type a message and see a response.
- UI clearly marks demo/synthetic mode.
- UI can be replaced later without changing orchestrator.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P2-T04 — Add request validation and safe client-facing errors

**Type:** Code  
**Goal:** Keep malformed requests and accidental PII from destabilizing the demo.

### Copy/paste prompt

```text
You are helping implement task P2-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Add request validation and safe client-facing errors
Task goal: Keep malformed requests and accidental PII from destabilizing the demo.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/validation`
- `internal/handlers/chat_test.go`

Required tests or review checks:
- Table-driven tests for empty, oversized, and whitespace-only messages.
- Student ID validation accepts synthetic IDs only.
- Error body never includes raw request body.

Acceptance criteria:
- Oversized messages are rejected.
- Synthetic ID rule is explicit.
- Error responses are consistent.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P2-T05 — Add in-memory conversation session store

**Type:** Code  
**Goal:** Track a short demo conversation so follow-up questions can reference prior transcript/payment context.

### Copy/paste prompt

```text
You are helping implement task P2-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Add in-memory conversation session store
Task goal: Track a short demo conversation so follow-up questions can reference prior transcript/payment context.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/session/store.go`
- `internal/session/store_test.go`

Required tests or review checks:
- Create, append, read, and expire session tests.
- Concurrent access test with `-race` target.
- PII redaction before persistence test.

Acceptance criteria:
- Store is concurrency-safe.
- TTL is configurable.
- Only redacted or synthetic data is stored.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P3-T01 — Create synthetic fixture loader

**Type:** Code  
**Goal:** Load deterministic demo records for students, transcripts, payments, LMS access, and CRM examples.

### Copy/paste prompt

```text
You are helping implement task P3-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Create synthetic fixture loader
Task goal: Load deterministic demo records for students, transcripts, payments, LMS access, and CRM examples.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `data/synthetic-students.json`
- `internal/fixtures/loader.go`
- `internal/fixtures/loader_test.go`

Required tests or review checks:
- Valid fixture loads all expected students.
- Duplicate IDs fail.
- Missing required fields fail.
- Fixture contains only synthetic IDs.

Acceptance criteria:
- S100001-S100004 exist with expected states.
- Fixtures can be reused by mock services and tests.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P3-T02 — Build mock Banner-style student API

**Type:** Code  
**Goal:** Simulate student profile, enrollment status, transcript request status, and holds.

### Copy/paste prompt

```text
You are helping implement task P3-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Build mock Banner-style student API
Task goal: Simulate student profile, enrollment status, transcript request status, and holds.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/mock-banner/main.go`
- `internal/mock/banner`
- `internal/tools/banner_client.go`

Required tests or review checks:
- Handler tests for known student, unknown student, transcript status, and financial hold.
- Client tests against `httptest.Server`.
- Contract tests validate response schema.

Acceptance criteria:
- Known IDs return deterministic data.
- Unknown IDs return 404.
- No real Banner naming/secrets are used beyond mock labels.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P3-T03 — Build mock payment API

**Type:** Code  
**Goal:** Simulate transcript payment status without processing real payments.

### Copy/paste prompt

```text
You are helping implement task P3-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Build mock payment API
Task goal: Simulate transcript payment status without processing real payments.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/mock-payment/main.go`
- `internal/mock/payment`
- `internal/tools/payment_client.go`

Required tests or review checks:
- Paid student returns paid status.
- Unpaid student returns unpaid status.
- Unknown payment returns clear safe error.
- Client timeout test.

Acceptance criteria:
- Payment response includes status, amount, currency, and synthetic transaction ID only.
- The UI never accepts real card/payment data.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P3-T04 — Build mock CRM API for case creation and queue routing

**Type:** Code  
**Goal:** Show human handoff and staff-facing case summaries.

### Copy/paste prompt

```text
You are helping implement task P3-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Build mock CRM API for case creation and queue routing
Task goal: Show human handoff and staff-facing case summaries.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/mock-crm/main.go`
- `internal/mock/crm`
- `internal/tools/crm_client.go`

Required tests or review checks:
- Create case success.
- Priority case success.
- Validation rejects empty summary.
- Client retries or returns typed error on 5xx.

Acceptance criteria:
- Case response includes case ID, queue, priority, and summary.
- Conversation summary is redacted before storage.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P3-T05 — Build mock LMS API

**Type:** Code  
**Goal:** Support basic Moodle/Brightspace-style learner access questions without real LMS access.

### Copy/paste prompt

```text
You are helping implement task P3-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Build mock LMS API
Task goal: Support basic Moodle/Brightspace-style learner access questions without real LMS access.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/mock-lms/main.go`
- `internal/mock/lms`
- `internal/tools/lms_client.go`

Required tests or review checks:
- Known synthetic student returns LMS access status.
- Unknown course returns safe fallback.
- Client handles timeout.

Acceptance criteria:
- LMS API is clearly marked synthetic.
- It supports only demo access status, not course content.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P3-T06 — Add typed enterprise clients with shared error model

**Type:** Code  
**Goal:** Let the orchestrator call mock services through interfaces that resemble production integrations.

### Copy/paste prompt

```text
You are helping implement task P3-T06 for the Go-based AskOC AI Concierge MVP.

Task title: Add typed enterprise clients with shared error model
Task goal: Let the orchestrator call mock services through interfaces that resemble production integrations.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/tools/errors.go`
- `internal/tools/*_client.go`
- `internal/tools/*_test.go`

Required tests or review checks:
- Timeout tests with context cancellation.
- 404 maps to typed not-found error.
- 5xx maps to retryable or external-service error.
- Malformed JSON maps to typed parse error.

Acceptance criteria:
- No orchestrator code depends on raw HTTP details.
- All tool calls accept context and emit trace ID headers.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P4-T01 — Define orchestrator ports and dependency injection

**Type:** Code  
**Goal:** Make the orchestrator testable by depending on interfaces for retrieval, classification, LLM, tools, workflow, and audit.

### Copy/paste prompt

```text
You are helping implement task P4-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Define orchestrator ports and dependency injection
Task goal: Make the orchestrator testable by depending on interfaces for retrieval, classification, LLM, tools, workflow, and audit.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/orchestrator.go`
- `internal/orchestrator/orchestrator_test.go`

Required tests or review checks:
- Compile-time interface fakes.
- Test orchestrator can be constructed with fake dependencies.
- Nil dependency validation test.

Acceptance criteria:
- No live network dependencies in orchestrator tests.
- Every external dependency has a small interface.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P4-T02 — Implement fallback intent and sentiment classifier

**Type:** Code  
**Goal:** Provide deterministic behavior for demos and tests before using an LLM.

### Copy/paste prompt

```text
You are helping implement task P4-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Implement fallback intent and sentiment classifier
Task goal: Provide deterministic behavior for demos and tests before using an LLM.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/classifier/fallback.go`
- `internal/classifier/fallback_test.go`

Required tests or review checks:
- Table-driven messages map to expected intents.
- Frustration/urgency phrases map to negative/high urgency.
- Unknown questions map to unknown with low confidence.

Acceptance criteria:
- Classifier returns typed confidence.
- Low confidence cannot trigger sensitive tools.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P4-T03 — Implement transcript-status decision flow

**Type:** Code  
**Goal:** Handle the core Tier 1 scenario using synthetic Banner and payment data.

### Copy/paste prompt

```text
You are helping implement task P4-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Implement transcript-status decision flow
Task goal: Handle the core Tier 1 scenario using synthetic Banner and payment data.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/transcript.go`
- `internal/orchestrator/transcript_test.go`

Required tests or review checks:
- S100001 -> ready/no workflow.
- S100002 -> unpaid/payment reminder.
- S100003 -> hold/CRM escalation.
- S100004 -> unknown/handoff.
- Missing student ID prompts for synthetic ID.

Acceptance criteria:
- Response contains user-friendly answer plus machine-readable actions.
- No payment reminder is sent for paid records.
- Financial holds route to staff, not self-service.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P4-T04 — Trigger payment reminder workflow from orchestrator

**Type:** Code  
**Goal:** Connect the transcript decision flow to workflow automation.

### Copy/paste prompt

```text
You are helping implement task P4-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Trigger payment reminder workflow from orchestrator
Task goal: Connect the transcript decision flow to workflow automation.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/workflow/client.go`
- `internal/orchestrator/transcript_test.go`

Required tests or review checks:
- Unpaid transcript calls workflow exactly once.
- Idempotency key is passed.
- Workflow failure is reported safely and audited.
- Paid transcript does not call workflow.

Acceptance criteria:
- Workflow action appears in chat response.
- Failures do not create duplicate reminders.
- Audit event records workflow attempt.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P4-T05 — Create CRM escalation summary and priority routing

**Type:** Code  
**Goal:** Turn unresolved, urgent, or sensitive conversations into staff-ready mock cases.

### Copy/paste prompt

```text
You are helping implement task P4-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Create CRM escalation summary and priority routing
Task goal: Turn unresolved, urgent, or sensitive conversations into staff-ready mock cases.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/escalation.go`
- `internal/orchestrator/escalation_test.go`

Required tests or review checks:
- Negative sentiment creates priority case.
- Low confidence creates normal handoff.
- Financial hold routes to Finance/Registrar queue.
- Summary redacts emails/phone numbers.

Acceptance criteria:
- CRM case includes intent, queue, priority, trace ID, and redacted summary.
- Learner receives a clear handoff message.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P4-T06 — Return action trace in chat responses

**Type:** Code  
**Goal:** Make the demo transparent by showing what was checked, triggered, and escalated.

### Copy/paste prompt

```text
You are helping implement task P4-T06 for the Go-based AskOC AI Concierge MVP.

Task title: Return action trace in chat responses
Task goal: Make the demo transparent by showing what was checked, triggered, and escalated.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/domain/chat.go`
- `internal/orchestrator/orchestrator_test.go`

Required tests or review checks:
- Response includes tool action names and statuses.
- Internal errors are not exposed.
- Trace ID is present on all action results.

Acceptance criteria:
- Interview demo can show decision flow without opening logs.
- Action trace is safe for learner-facing display.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P5-T01 — Define source allowlist schema

**Type:** Code  
**Goal:** Represent approved public learner-service sources with freshness and risk metadata.

### Copy/paste prompt

```text
You are helping implement task P5-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Define source allowlist schema
Task goal: Represent approved public learner-service sources with freshness and risk metadata.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `data/seed-sources.json`
- `internal/rag/source.go`
- `internal/rag/source_test.go`

Required tests or review checks:
- Valid source config parses.
- Non-HTTPS URL fails.
- Missing title/department fails.
- Private-domain marker fails.

Acceptance criteria:
- Sources include URL, title, department, risk level, retrieved date, and freshness flag.
- No source is ingested unless allowlisted.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P5-T02 — Implement ingestion fetcher with allowlist and cleaning boundaries

**Type:** Code  
**Goal:** Fetch only approved public content and prepare text for chunking.

### Copy/paste prompt

```text
You are helping implement task P5-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Implement ingestion fetcher with allowlist and cleaning boundaries
Task goal: Fetch only approved public content and prepare text for chunking.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/ingest/main.go`
- `internal/rag/ingest.go`
- `internal/rag/ingest_test.go`

Required tests or review checks:
- Fetcher rejects URL not in allowlist.
- HTML cleaner removes nav/script/style.
- Network failure returns typed error.
- Test uses `httptest.Server`, not live internet.

Acceptance criteria:
- Ingestion is deterministic in tests.
- Content hash is stored.
- Private pages cannot be fetched accidentally.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P5-T03 — Implement chunking with metadata preservation

**Type:** Code  
**Goal:** Split retrieved content into searchable chunks while preserving source URL and title.

### Copy/paste prompt

```text
You are helping implement task P5-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Implement chunking with metadata preservation
Task goal: Split retrieved content into searchable chunks while preserving source URL and title.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/rag/chunk.go`
- `internal/rag/chunk_test.go`

Required tests or review checks:
- Short content produces one chunk.
- Long content produces bounded chunks.
- Chunk IDs are stable for same input.
- Metadata is copied to every chunk.

Acceptance criteria:
- Chunk size is configurable.
- No empty chunks are stored.
- Chunk IDs can be cited in responses.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P5-T04 — Create local retrieval implementation

**Type:** Code  
**Goal:** Provide a demo-safe retrieval path before wiring Azure AI Search or pgvector.

### Copy/paste prompt

```text
You are helping implement task P5-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Create local retrieval implementation
Task goal: Provide a demo-safe retrieval path before wiring Azure AI Search or pgvector.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/rag/retrieve.go`
- `internal/rag/local_retriever.go`
- `internal/rag/retrieve_test.go`

Required tests or review checks:
- Transcript query ranks transcript chunks first.
- Unrelated query returns low confidence.
- Retrieval limit is respected.
- Source metadata is returned.

Acceptance criteria:
- Retriever interface can be swapped later.
- No hallucinated source links are created.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P5-T05 — Add grounded answer source packaging

**Type:** Code  
**Goal:** Attach source citations and confidence to chat responses regardless of LLM provider.

### Copy/paste prompt

```text
You are helping implement task P5-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Add grounded answer source packaging
Task goal: Attach source citations and confidence to chat responses regardless of LLM provider.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/grounded_answer.go`
- `internal/orchestrator/grounded_answer_test.go`

Required tests or review checks:
- Answer with sufficient chunks includes citations.
- Low retrieval confidence returns safe fallback.
- Duplicate sources are de-duplicated.
- Source risk level is included for internal response.

Acceptance criteria:
- Every policy/procedure answer includes at least one source or a fallback.
- Unsupported claims are not invented.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P5-T06 — Flag stale or high-risk sources

**Type:** Code  
**Goal:** Avoid confidently answering from outdated or sensitive content.

### Copy/paste prompt

```text
You are helping implement task P5-T06 for the Go-based AskOC AI Concierge MVP.

Task title: Flag stale or high-risk sources
Task goal: Avoid confidently answering from outdated or sensitive content.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/rag/freshness.go`
- `internal/rag/freshness_test.go`

Required tests or review checks:
- Stale source triggers caution flag.
- High-risk source requires stronger confidence or handoff.
- Fresh source passes without warning.

Acceptance criteria:
- Assistant can say it needs staff confirmation for stale/high-risk policy questions.
- Dashboard can show stale-source warnings.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P6-T01 — Define LLM provider interface and request/response types

**Type:** Code  
**Goal:** Hide provider details behind typed Go interfaces.

### Copy/paste prompt

```text
You are helping implement task P6-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Define LLM provider interface and request/response types
Task goal: Hide provider details behind typed Go interfaces.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/llm/types.go`
- `internal/llm/types_test.go`

Required tests or review checks:
- JSON schema for answer request/response marshals correctly.
- Provider-neutral error type works.
- Timeout field validation.

Acceptance criteria:
- Orchestrator imports interfaces/types only.
- Provider can be replaced without handler changes.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P6-T02 — Implement OpenAI/Azure-compatible REST client

**Type:** Code  
**Goal:** Call an OpenAI-compatible chat/completions endpoint through a testable Go client.

### Copy/paste prompt

```text
You are helping implement task P6-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Implement OpenAI/Azure-compatible REST client
Task goal: Call an OpenAI-compatible chat/completions endpoint through a testable Go client.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/llm/openai_client.go`
- `internal/llm/openai_client_test.go`

Required tests or review checks:
- `httptest.Server` verifies request payload.
- Client parses successful response.
- Client handles 429/500/timeout.
- No live API call in tests.

Acceptance criteria:
- API key is read from config only.
- Logs never include prompts containing PII.
- Client supports context cancellation.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P6-T03 — Parse strict JSON classification output

**Type:** Code  
**Goal:** Convert model output into trusted typed classification only after validation.

### Copy/paste prompt

```text
You are helping implement task P6-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Parse strict JSON classification output
Task goal: Convert model output into trusted typed classification only after validation.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/classifier/llm_parser.go`
- `internal/classifier/llm_parser_test.go`

Required tests or review checks:
- Valid JSON parses.
- Malformed JSON fails safely.
- Unknown intent maps to `unknown`.
- Out-of-range confidence is rejected or clamped.
- Tool-triggering is disabled below threshold.

Acceptance criteria:
- Invalid model output never panics.
- Low-confidence classification returns safe fallback.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P6-T04 — Create prompt templates for classification and grounded answers

**Type:** Code  
**Goal:** Make prompts versioned, testable, and aligned with privacy/safety rules.

### Copy/paste prompt

```text
You are helping implement task P6-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Create prompt templates for classification and grounded answers
Task goal: Make prompts versioned, testable, and aligned with privacy/safety rules.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/prompts.go`
- `internal/orchestrator/prompts_test.go`

Required tests or review checks:
- Prompt contains strict JSON instruction.
- Prompt includes source-only answer rule.
- Prompt includes privacy/no-real-data rule.
- Golden test catches accidental prompt drift.

Acceptance criteria:
- Prompts are plain text constants or embedded templates.
- Prompt version is included in audit metadata.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P6-T05 — Add LLM fallback and guardrail behavior

**Type:** Code  
**Goal:** Keep the demo reliable if the model is unavailable or produces unsafe output.

### Copy/paste prompt

```text
You are helping implement task P6-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Add LLM fallback and guardrail behavior
Task goal: Keep the demo reliable if the model is unavailable or produces unsafe output.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/ai_guardrails.go`
- `internal/orchestrator/ai_guardrails_test.go`

Required tests or review checks:
- Model timeout uses deterministic fallback.
- Unsafe answer without sources is rejected.
- Tool calls require validated classification.
- Sensitive/unsupported requests escalate.

Acceptance criteria:
- No core demo path depends solely on live model availability.
- Guardrail failures are logged and visible in dashboard.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P6-T06 — Add end-to-end classification tests with fixture messages

**Type:** Code  
**Goal:** Measure intent and sentiment behavior before using the assistant in the demo.

### Copy/paste prompt

```text
You are helping implement task P6-T06 for the Go-based AskOC AI Concierge MVP.

Task title: Add end-to-end classification tests with fixture messages
Task goal: Measure intent and sentiment behavior before using the assistant in the demo.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `data/classification-fixtures.jsonl`
- `internal/classifier/e2e_test.go`

Required tests or review checks:
- At least 5 examples per supported intent.
- Negative/urgent sentiment examples pass.
- Unknown/off-topic examples do not trigger actions.

Acceptance criteria:
- Fixture accuracy target is documented.
- Failures identify which intent regressed.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P7-T01 — Implement PII redaction

**Type:** Code  
**Goal:** Redact emails, phone numbers, likely passwords, and non-synthetic IDs before logging or case summaries.

### Copy/paste prompt

```text
You are helping implement task P7-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Implement PII redaction
Task goal: Redact emails, phone numbers, likely passwords, and non-synthetic IDs before logging or case summaries.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/privacy/redact.go`
- `internal/privacy/redact_test.go`

Required tests or review checks:
- Email redaction.
- Phone redaction.
- Password phrase redaction.
- Real-looking student number redaction.
- Synthetic S100001-style IDs are allowed for demo.

Acceptance criteria:
- Redactor is used by logs, audit, session store, and CRM summaries.
- Tests cover false positives and false negatives.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P7-T02 — Create audit event schema and store

**Type:** Code  
**Goal:** Record safe operational events for intent, sources, tool calls, workflow triggers, and escalations.

### Copy/paste prompt

```text
You are helping implement task P7-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Create audit event schema and store
Task goal: Record safe operational events for intent, sources, tool calls, workflow triggers, and escalations.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/audit/event.go`
- `internal/audit/store.go`
- `internal/audit/store_test.go`

Required tests or review checks:
- Record and query events by trace ID.
- Event payload is redacted.
- Concurrent writes are safe.
- Retention timestamp is set.

Acceptance criteria:
- Every orchestrator action can produce an audit event.
- Audit store can run in memory for demo and later be backed by DB.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P7-T03 — Wire redacted structured logging across services

**Type:** Code  
**Goal:** Ensure logs are useful for debugging but safe for learner privacy.

### Copy/paste prompt

```text
You are helping implement task P7-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Wire redacted structured logging across services
Task goal: Ensure logs are useful for debugging but safe for learner privacy.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/middleware/logging.go`
- `internal/privacy/logging_test.go`

Required tests or review checks:
- Chat message is not logged raw.
- Trace ID is logged.
- Tool status is logged without raw payload.
- Secret-like values are redacted.

Acceptance criteria:
- No log line contains real-looking PII from tests.
- Logger uses `slog` fields, not string concatenation.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P7-T04 — Build admin metrics endpoint

**Type:** Code  
**Goal:** Summarize containment, escalation, top intents, confidence, workflow count, and low-confidence questions.

### Copy/paste prompt

```text
You are helping implement task P7-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Build admin metrics endpoint
Task goal: Summarize containment, escalation, top intents, confidence, workflow count, and low-confidence questions.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/handlers/admin.go`
- `internal/handlers/admin_test.go`
- `internal/audit/metrics.go`

Required tests or review checks:
- Metrics calculation from seeded audit events.
- Endpoint requires admin/mock auth.
- Empty store returns zeros safely.

Acceptance criteria:
- Dashboard data can be generated after demo conversation.
- Metrics are aggregate and redacted.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P7-T05 — Create minimal admin dashboard UI

**Type:** Code  
**Goal:** Give stakeholders an at-a-glance view of AI adoption and operational risk.

### Copy/paste prompt

```text
You are helping implement task P7-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Create minimal admin dashboard UI
Task goal: Give stakeholders an at-a-glance view of AI adoption and operational risk.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `web/templates/admin.html`
- `web/static/admin.js`
- `internal/handlers/admin_ui.go`

Required tests or review checks:
- Template renders metric labels.
- Unauthorized request is rejected.
- Review queue renders redacted question text only.

Acceptance criteria:
- Dashboard shows top intents, escalations, workflows, low-confidence items, and stale-source warnings.
- No raw PII appears.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P7-T06 — Add retention and export controls

**Type:** Code  
**Goal:** Document and enforce how demo audit data is kept, deleted, and exported.

### Copy/paste prompt

```text
You are helping implement task P7-T06 for the Go-based AskOC AI Concierge MVP.

Task title: Add retention and export controls
Task goal: Document and enforce how demo audit data is kept, deleted, and exported.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/audit/retention.go`
- `internal/audit/retention_test.go`
- `docs/privacy-impact-lite.md`

Required tests or review checks:
- Expired events are purged.
- Export excludes raw messages.
- Retention config default is short for demo.

Acceptance criteria:
- Privacy doc matches implementation behavior.
- Demo can reset all stored data.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P8-T01 — Build local workflow simulator service

**Type:** Code  
**Goal:** Simulate Power Automate-style HTTP-triggered flows for offline demos.

### Copy/paste prompt

```text
You are helping implement task P8-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Build local workflow simulator service
Task goal: Simulate Power Automate-style HTTP-triggered flows for offline demos.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/workflow-sim/main.go`
- `internal/workflow/sim_handler.go`
- `internal/workflow/sim_handler_test.go`

Required tests or review checks:
- Payment reminder payload accepted.
- Missing idempotency key rejected.
- Duplicate idempotency key returns same workflow ID.
- Invalid payload returns 400.

Acceptance criteria:
- Simulator returns workflow ID and status.
- Workflow events are visible in dashboard.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P8-T02 — Implement workflow idempotency and retry policy

**Type:** Code  
**Goal:** Prevent duplicate payment reminders and handle transient failures safely.

### Copy/paste prompt

```text
You are helping implement task P8-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Implement workflow idempotency and retry policy
Task goal: Prevent duplicate payment reminders and handle transient failures safely.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/workflow/idempotency.go`
- `internal/workflow/client_test.go`

Required tests or review checks:
- Same trace/student/action creates same idempotency key.
- Transient 500 is retried within limit.
- Permanent 400 is not retried.
- Context cancellation stops retry.

Acceptance criteria:
- No duplicate learner reminders in repeated calls.
- Retry behavior is observable in audit logs.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P8-T03 — Add optional Power Automate webhook client

**Type:** Code  
**Goal:** Show RPA/workflow integration relevance while keeping local simulator as default.

### Copy/paste prompt

```text
You are helping implement task P8-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Add optional Power Automate webhook client
Task goal: Show RPA/workflow integration relevance while keeping local simulator as default.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/workflow/powerautomate.go`
- `internal/workflow/powerautomate_test.go`

Required tests or review checks:
- Client sends expected JSON schema to `httptest.Server`.
- Signature/header is included when configured.
- Webhook URL missing falls back to simulator config.
- Secrets are not logged.

Acceptance criteria:
- Same interface supports simulator and Power Automate.
- Docs explain secure webhook storage.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P8-T04 — Audit workflow outcomes and errors

**Type:** Code  
**Goal:** Make workflow automation measurable and debuggable.

### Copy/paste prompt

```text
You are helping implement task P8-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Audit workflow outcomes and errors
Task goal: Make workflow automation measurable and debuggable.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/orchestrator/transcript.go`
- `internal/audit/event.go`
- `internal/workflow/client_test.go`

Required tests or review checks:
- Success event recorded.
- Failure event recorded with safe message.
- Dashboard count increments.
- No raw payload is stored.

Acceptance criteria:
- Every workflow attempt has trace ID, action, status, and idempotency key hash.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P8-T05 — Document Power Automate flow schema and setup

**Type:** Documentation  
**Goal:** Provide enough detail for an interviewer to see how the Go webhook would connect to Power Automate.

### Copy/paste prompt

```text
You are helping implement task P8-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Document Power Automate flow schema and setup
Task goal: Provide enough detail for an interviewer to see how the Go webhook would connect to Power Automate.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/power-automate-flow.md`

Required tests or review checks:
- Review JSON examples against client tests.
- Confirm docs do not include real webhook secrets.

Acceptance criteria:
- Webhook request/response examples are present.
- Security notes mention idempotency, secret storage, and replay protection.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P9-T01 — Create JSONL evaluation dataset

**Type:** Code  
**Goal:** Define repeatable test cases for intents, sources, workflow actions, escalations, and safety behavior.

### Copy/paste prompt

```text
You are helping implement task P9-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Create JSONL evaluation dataset
Task goal: Define repeatable test cases for intents, sources, workflow actions, escalations, and safety behavior.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `data/eval-questions.jsonl`
- `internal/eval/dataset.go`
- `internal/eval/dataset_test.go`

Required tests or review checks:
- Dataset parser handles valid JSONL.
- Invalid rows fail with line number.
- Required expected fields are validated.
- Dataset includes critical safety cases.

Acceptance criteria:
- At least 30 evaluation examples exist.
- Each case has expected intent and expected action/source behavior.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P9-T02 — Build `cmd/eval` runner

**Type:** Code  
**Goal:** Run the assistant against the evaluation dataset and produce repeatable quality results.

### Copy/paste prompt

```text
You are helping implement task P9-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Build `cmd/eval` runner
Task goal: Run the assistant against the evaluation dataset and produce repeatable quality results.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/eval/main.go`
- `internal/eval/runner.go`
- `internal/eval/runner_test.go`

Required tests or review checks:
- Runner calls fake chat client for each case.
- Timeouts are captured as failures.
- Results include latency.
- CLI exits zero when all critical tests pass.

Acceptance criteria:
- `go run ./cmd/eval` works locally.
- Runner can target in-process fake or live local API.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P9-T03 — Implement scoring functions

**Type:** Code  
**Goal:** Score intent accuracy, source grounding, action correctness, escalation precision, safety, and latency.

### Copy/paste prompt

```text
You are helping implement task P9-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Implement scoring functions
Task goal: Score intent accuracy, source grounding, action correctness, escalation precision, safety, and latency.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/eval/score.go`
- `internal/eval/score_test.go`

Required tests or review checks:
- Intent exact match scoring.
- Expected source present scoring.
- Expected action present/absent scoring.
- Critical hallucination scoring.
- Latency threshold scoring.

Acceptance criteria:
- Scores are deterministic.
- Critical policy errors are separated from minor misses.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P9-T04 — Generate JSON and Markdown evaluation reports

**Type:** Code  
**Goal:** Create portfolio-ready evidence of model/system quality.

### Copy/paste prompt

```text
You are helping implement task P9-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Generate JSON and Markdown evaluation reports
Task goal: Create portfolio-ready evidence of model/system quality.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/eval/report.go`
- `internal/eval/report_test.go`
- `reports/eval-summary.md`

Required tests or review checks:
- JSON report contains summary and per-case results.
- Markdown report contains metrics table.
- Report redacts user messages when configured.

Acceptance criteria:
- Reports can be committed as sample outputs.
- Failures are readable by non-technical stakeholders.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P9-T05 — Fail builds on critical evaluation failures

**Type:** Code  
**Goal:** Make responsible AI behavior part of the development gate, not a manual afterthought.

### Copy/paste prompt

```text
You are helping implement task P9-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Fail builds on critical evaluation failures
Task goal: Make responsible AI behavior part of the development gate, not a manual afterthought.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `cmd/eval/main.go`
- `internal/eval/gates.go`
- `internal/eval/gates_test.go`
- `Makefile`

Required tests or review checks:
- Critical hallucination returns non-zero exit.
- Missing required escalation returns non-zero exit.
- Minor accuracy miss can be warning based on threshold.
- Make target propagates exit code.

Acceptance criteria:
- `make eval` fails for critical safety regressions.
- README explains quality gate.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P9-T06 — Create review queue for failed/low-confidence answers

**Type:** Code  
**Goal:** Show how the system would improve over time through human review.

### Copy/paste prompt

```text
You are helping implement task P9-T06 for the Go-based AskOC AI Concierge MVP.

Task title: Create review queue for failed/low-confidence answers
Task goal: Show how the system would improve over time through human review.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `internal/eval/review_queue.go`
- `internal/handlers/admin.go`
- `internal/eval/review_queue_test.go`

Required tests or review checks:
- Low-confidence answer added to queue.
- Failed critical eval added to queue.
- Duplicate review items collapse by normalized question.
- Resolved item no longer appears as open.

Acceptance criteria:
- Dashboard can display unresolved review items.
- Review items contain sources/actions but no raw PII.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P10-T01 — Create Dockerfiles for Go services

**Type:** Code  
**Goal:** Package API and mock services as small reproducible containers.

### Copy/paste prompt

```text
You are helping implement task P10-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Create Dockerfiles for Go services
Task goal: Package API and mock services as small reproducible containers.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `Dockerfile`
- `cmd/*`
- `.dockerignore`

Required tests or review checks:
- Add a shell or Makefile smoke test that builds the API image.
- Unit tests still pass before image build.
- Container starts and responds to `/healthz` in local smoke script.

Acceptance criteria:
- Image does not include source secrets.
- Non-root runtime user is used where practical.
- Build is documented.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P10-T02 — Create Docker Compose local stack

**Type:** Code  
**Goal:** Start API, mock services, workflow simulator, and optional database with one command.

### Copy/paste prompt

```text
You are helping implement task P10-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Create Docker Compose local stack
Task goal: Start API, mock services, workflow simulator, and optional database with one command.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `docker-compose.yml`
- `.env.example`
- `Makefile`

Required tests or review checks:
- `make compose-test` or smoke script waits for health endpoints.
- Missing optional LLM config still allows deterministic demo mode.
- Service URLs match README.

Acceptance criteria:
- `make dev` or `docker compose up` starts the full demo.
- Ports are documented.
- Synthetic mode is default.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P10-T03 — Add CI workflow for Go tests and evaluation gates

**Type:** Code  
**Goal:** Demonstrate professional delivery discipline.

### Copy/paste prompt

```text
You are helping implement task P10-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Add CI workflow for Go tests and evaluation gates
Task goal: Demonstrate professional delivery discipline.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `.github/workflows/ci.yml`
- `Makefile`

Required tests or review checks:
- CI runs `go test ./...`.
- CI runs `go vet ./...`.
- CI runs `make eval` with deterministic/fake providers.
- CI fails on critical eval failure.

Acceptance criteria:
- Pull requests cannot pass with failing tests.
- CI avoids live model/API dependencies.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P10-T04 — Add environment sample and secret-safety checks

**Type:** Code  
**Goal:** Prevent accidental commit of real keys or private data.

### Copy/paste prompt

```text
You are helping implement task P10-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Add environment sample and secret-safety checks
Task goal: Prevent accidental commit of real keys or private data.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `.env.example`
- `.gitignore`
- `internal/config/config_test.go`

Required tests or review checks:
- Config tests prove secrets are optional for deterministic mode.
- Secret scanner or grep script rejects known key patterns in repo.
- `.env` is ignored.

Acceptance criteria:
- Only placeholders appear in docs.
- No real webhook URL or API key is committed.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P10-T05 — Add one-command smoke test

**Type:** Code  
**Goal:** Make the project easy to verify during interview prep.

### Copy/paste prompt

```text
You are helping implement task P10-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Add one-command smoke test
Task goal: Make the project easy to verify during interview prep.

Strict TDD instructions:
1. Before writing or modifying production code, create the failing tests listed below.
2. Run the narrowest relevant `go test` command and confirm failure for the expected missing behavior.
3. Implement only the minimum production code required to pass those tests.
4. Run the narrow package test, then `go test ./...`.
5. Refactor only after tests are green.
6. Do not call live external APIs in tests; use fakes, fixtures, or `httptest.Server`.
7. Keep logs, audit payloads, and test fixtures free of real PII and secrets.

Primary files to create or update:
- `scripts/smoke.sh`
- `Makefile`
- `docs/test-plan.md`

Required tests or review checks:
- Smoke test starts stack or assumes stack is running.
- Calls `/healthz`.
- Calls transcript chat scenario.
- Verifies workflow and CRM outputs for synthetic IDs.

Acceptance criteria:
- `make smoke` gives a clear pass/fail.
- Failures include actionable messages.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the failing tests you will add first.
2. Show or describe the expected initial failure.
3. Implement the minimal code.
4. List commands to run.
5. Summarize changed files and remaining risks.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P11-T01 — Polish README for applicant storytelling

**Type:** Documentation  
**Goal:** Make the project understandable to a hiring manager in under two minutes.

### Copy/paste prompt

```text
You are helping implement task P11-T01 for the Go-based AskOC AI Concierge MVP.

Task title: Polish README for applicant storytelling
Task goal: Make the project understandable to a hiring manager in under two minutes.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `README.md`

Required tests or review checks:
- Read README top-to-bottom and confirm it explains problem, solution, stack, demo, and privacy.

Acceptance criteria:
- README includes screenshot placeholder, quickstart, architecture, success metrics, and TDD quality statement.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P11-T02 — Add architecture diagram and sequence diagram

**Type:** Documentation  
**Goal:** Show enterprise integration thinking visually.

### Copy/paste prompt

```text
You are helping implement task P11-T02 for the Go-based AskOC AI Concierge MVP.

Task title: Add architecture diagram and sequence diagram
Task goal: Show enterprise integration thinking visually.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/architecture.md`

Required tests or review checks:
- Trace golden path through diagram and verify every component exists in tasks or code.

Acceptance criteria:
- Diagram includes chat UI, Go API, orchestrator, RAG, mock Banner/payment/CRM/LMS, workflow, audit, dashboard.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P11-T03 — Finalize 5-7 minute demo script

**Type:** Documentation  
**Goal:** Prepare a concise interview walkthrough that highlights the role requirements.

### Copy/paste prompt

```text
You are helping implement task P11-T03 for the Go-based AskOC AI Concierge MVP.

Task title: Finalize 5-7 minute demo script
Task goal: Prepare a concise interview walkthrough that highlights the role requirements.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/demo-script.md`

Required tests or review checks:
- Dry-run the script and verify it fits the time box.
- Confirm every spoken claim can be shown in UI, logs, docs, or tests.

Acceptance criteria:
- Demo covers Tier 0 answer, Tier 1 workflow, urgent escalation, dashboard, and TDD/evaluation evidence.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P11-T04 — Prepare screenshots or short GIF placeholders

**Type:** Documentation  
**Goal:** Make the GitHub repo visually scannable.

### Copy/paste prompt

```text
You are helping implement task P11-T04 for the Go-based AskOC AI Concierge MVP.

Task title: Prepare screenshots or short GIF placeholders
Task goal: Make the GitHub repo visually scannable.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `README.md`
- `docs/demo-script.md`

Required tests or review checks:
- Check screenshots do not show real student data, real tokens, or private URLs.

Acceptance criteria:
- README has image placeholders or links.
- Captions explain what each screen proves.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.


## P11-T05 — Run final release checklist

**Type:** Documentation  
**Goal:** Confirm the portfolio is coherent, safe, and demonstrable.

### Copy/paste prompt

```text
You are helping implement task P11-T05 for the Go-based AskOC AI Concierge MVP.

Task title: Run final release checklist
Task goal: Confirm the portfolio is coherent, safe, and demonstrable.

Documentation/review instructions:
1. Update only the Markdown or fixture files needed for this task.
2. Keep the project Go-first and synthetic-data-only.
3. Make acceptance criteria observable and easy to verify.
4. Do not claim real OC system access or production deployment.

Primary files to create or update:
- `docs/phases-and-tasks.md`
- `docs/test-plan.md`
- `docs/privacy-impact-lite.md`

Required tests or review checks:
- Run quickstart commands.
- Run tests and eval.
- Review docs for consistency.
- Check links.

Acceptance criteria:
- All must-have tasks are checked or clearly deferred.
- ZIP/repo is ready to share.
- Known limitations are honest and documented.

Implementation constraints:
- Use idiomatic Go with small packages and explicit interfaces at external boundaries.
- Pass `context.Context` into operations that may call services, perform I/O, or take time.
- Return typed errors where the caller needs to branch on behavior.
- Keep user-facing errors safe and non-leaky.
- Keep all demo data synthetic.
- Update related Markdown docs when behavior, commands, API payloads, or assumptions change.

Expected response format:
1. List the review checks you will perform.
2. Update the requested Markdown or fixture files.
3. Summarize changed files.
4. Call out assumptions, deferred items, and any consistency risks.
5. Confirm that no real learner data, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Tests/review checks were created or performed before marking done.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real learner data, secrets, or private URLs were added.
