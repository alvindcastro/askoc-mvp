# Phases and Tickable Tasks

This task board turns the Go-based AskOC AI Concierge MVP into an execution plan. It is intentionally strict: **every code task must follow test-driven development (TDD)**. Documentation tasks still require review checks, but they do not require Go tests.

## Status key

- `[ ]` Not started
- `[~]` In progress
- `[x]` Done
- `Blocked:` add a short reason directly under the task

## Non-negotiable code rule

For every task marked **Code**:

1. Write or update the failing test first.
2. Run the narrow test and confirm it fails for the expected reason.
3. Write the smallest production code required to pass.
4. Run the narrow test, then `go test ./...`.
5. Refactor only while tests are green.
6. Do not mark the checkbox done until the acceptance evidence is true.

See [TDD Policy](tdd-policy.md) and [Task Prompts](task-prompts.md) for detailed copy/paste prompts.

## Phase overview

| Phase | Outcome |
|---|---|
| P0 — Product framing and applicant strategy | Applicant story, scope, privacy boundary, source limits, and demo matrix are frozen. |
| P1 — Go project foundation | Go project compiles, serves health endpoints, and has API hygiene. |
| P2 — Chat API and UI skeleton | Chat API and minimal UI exist with deterministic placeholder behavior. |
| P3 — Synthetic enterprise APIs and clients | Synthetic Banner/payment/CRM/LMS APIs and typed clients are working. |
| P4 — Deterministic orchestration before AI | Transcript/payment orchestration works without relying on live AI. |
| P5 — RAG ingestion and source-grounded answers | Approved-source ingestion and retrieval produce grounded responses. |
| P6 — LLM gateway and structured classification | LLM gateway, strict JSON classification, prompts, and guardrails are testable. |
| P7 — Privacy, audit, and dashboard | Privacy redaction, audit events, and dashboard metrics are in place. |
| P8 — Workflow automation and Power Automate option | Workflow simulator and optional Power Automate webhook path are ready. |
| P9 — Evaluation runner and quality gates | Evaluation runner produces quality reports and fails critical regressions. |
| P10 — Docker, CI, and local developer experience | Docker, CI, env safety, and smoke test support repeatable demos. |
| P11 — Portfolio polish and interview readiness | Portfolio materials are polished for interview presentation. |

---

## P0 — Product framing and applicant strategy
**Phase outcome:** Applicant story, scope, privacy boundary, source limits, and demo matrix are frozen.
- [x] **P0-T01 — Write the applicant story and MVP thesis**
  **Type:** Documentation  
  **Goal:** State why the project exists, what learner pain point it addresses, and how it maps to the AI/Automation Solutions Developer role.  
  **Primary files:** `README.md`, `docs/mvp-scope.md`  
  **Prompt:** [`P0-T01` in task-prompts.md](task-prompts.md#p0-t01-write-the-applicant-story-and-mvp-thesis)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Peer-read the one-sentence pitch and confirm it mentions Go, RAG, workflow automation, mock enterprise integrations, and privacy.  
  **Done when:** README contains a one-sentence pitch.; MVP scope states the primary transcript/payment workflow.; The story avoids claiming access to real OC systems.
- [x] **P0-T02 — Define synthetic data and privacy boundary**
  **Type:** Documentation  
  **Goal:** Make it explicit that the MVP uses synthetic learners, synthetic IDs, mock payments, and mock CRM cases only.  
  **Primary files:** `docs/privacy-impact-lite.md`, `data/synthetic-students.json`  
  **Prompt:** [`P0-T02` in task-prompts.md](task-prompts.md#p0-t02-define-synthetic-data-and-privacy-boundary)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Review all examples for real names, real student IDs, private portal data, or secrets.  
  **Done when:** Synthetic data policy is visible.; PII boundary is documented.; Demo records are clearly marked as fake.
- [x] **P0-T03 — Confirm source allowlist and knowledge-domain limits**
  **Type:** Documentation  
  **Goal:** Define which public pages or manually curated content can be ingested and which content must never be ingested.  
  **Primary files:** `data/seed-sources.json`, `docs/source-references.md`, `docs/privacy-impact-lite.md`  
  **Prompt:** [`P0-T03` in task-prompts.md](task-prompts.md#p0-t03-confirm-source-allowlist-and-knowledge-domain-limits)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Manually verify every seed source is public and relevant to learner services.  
  **Done when:** Allowlist exists.; Private portal scraping is explicitly out of scope.; Stale-source handling is documented.
- [x] **P0-T04 — Create demo acceptance matrix**
  **Type:** Documentation  
  **Goal:** Turn the interview demo into measurable scenarios before implementation starts.  
  **Primary files:** `docs/demo-script.md`, `docs/model-evaluation.md`  
  **Prompt:** [`P0-T04` in task-prompts.md](task-prompts.md#p0-t04-create-demo-acceptance-matrix)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Walk through the golden path and confirm each expected action has an observable output.  
  **Done when:** Golden path includes transcript answer, unpaid payment workflow, financial-hold escalation, and urgent sentiment escalation.; Each scenario has expected intent, source, action, and handoff behavior.
- [x] **P0-T05 — Freeze MVP scope and defer nice-to-haves**
  **Type:** Documentation  
  **Goal:** Prevent overbuilding by separating must-have demo features from optional stretch work.  
  **Primary files:** `docs/mvp-scope.md`, `docs/implementation-roadmap.md`  
  **Prompt:** [`P0-T05` in task-prompts.md](task-prompts.md#p0-t05-freeze-mvp-scope-and-defer-nice-to-haves)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Check that every Phase 1-9 task supports the transcript/payment MVP or demo operations.  
  **Done when:** Nice-to-have features are explicitly deferred.; No real authentication, payments, or Banner integrations are listed as MVP work.

### P0 phase gate

- [x] All P0 tasks above are complete or explicitly deferred with a reason.
- [x] No code tasks exist in P0; failing-test evidence is not applicable.
- [x] `go test ./...` is not required for P0 because no Go code was added or changed.
- [x] Relevant docs and fixtures are updated with changed behavior and assumptions.

---

## P1 — Go project foundation
**Phase outcome:** Go project compiles, serves health endpoints, and has API hygiene.
- [x] **P1-T01 — Initialize Go module, repository layout, and developer commands**  
  **Type:** Code  
  **Goal:** Create a compilable Go workspace with predictable local commands before adding business logic.  
  **Primary files:** `go.mod`, `Makefile`, `cmd/api/main.go`, `internal/config`  
  **Prompt:** [`P1-T01` in task-prompts.md](task-prompts.md#p1-t01-initialize-go-module-repository-layout-and-developer-commands)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Add a failing smoke test that runs `go test ./...` and proves the module compiles once minimal code exists.; Add a Makefile test target that executes `go test ./...`.  
  **Done when:** `go test ./...` passes.; `make test` passes.; `make dev` has a documented target, even if it only runs the API skeleton.
- [x] **P1-T02 — Build typed configuration loader**  
  **Type:** Code  
  **Goal:** Load HTTP addresses, auth mode, log level, workflow URL, and provider settings from environment variables with safe defaults.  
  **Primary files:** `internal/config/config.go`, `internal/config/config_test.go`  
  **Prompt:** [`P1-T02` in task-prompts.md](task-prompts.md#p1-t02-build-typed-configuration-loader)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Table-driven tests for defaults.; Table-driven tests for overridden env vars.; Test invalid numeric/boolean values return clear errors.  
  **Done when:** Config has no global mutable state.; Secrets are not printed.; Invalid config fails fast with actionable error messages.
- [x] **P1-T03 — Implement health and readiness endpoints**  
  **Type:** Code  
  **Goal:** Expose simple operational endpoints for local demo, Docker health checks, and future deployment.  
  **Primary files:** `internal/handlers/health.go`, `internal/handlers/health_test.go`  
  **Prompt:** [`P1-T03` in task-prompts.md](task-prompts.md#p1-t03-implement-health-and-readiness-endpoints)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** `httptest` verifies `GET /healthz` returns 200 JSON.; `httptest` verifies `GET /readyz` returns dependency status.; Invalid methods return 405.  
  **Done when:** Health response includes status and trace ID when middleware is enabled.; No external dependencies are required for `/healthz`.
- [x] **P1-T04 — Add middleware for trace ID, panic recovery, request logging, and mock auth**  
  **Type:** Code  
  **Goal:** Create enterprise-style API hygiene before feature work begins.  
  **Primary files:** `internal/middleware/trace.go`, `internal/middleware/recover.go`, `internal/middleware/auth.go`, `internal/middleware/logging.go`, `internal/middleware/*_test.go`  
  **Prompt:** [`P1-T04` in task-prompts.md](task-prompts.md#p1-t04-add-middleware-for-trace-id-panic-recovery-request-logging-and-mock-auth)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Trace middleware adds or preserves trace ID.; Recovery middleware converts panic to 500 without exposing internals.; Mock auth rejects missing token when auth is enabled.; Logging test verifies redaction hooks are used.  
  **Done when:** All handlers receive context with trace ID.; Panic stack traces are not returned to clients.; Auth can be disabled for local demo.
- [x] **P1-T05 — Create JSON response and error helpers**  
  **Type:** Code  
  **Goal:** Standardize API responses so later handlers are easy to test and document.  
  **Primary files:** `internal/handlers/respond.go`, `internal/handlers/respond_test.go`  
  **Prompt:** [`P1-T05` in task-prompts.md](task-prompts.md#p1-t05-create-json-response-and-error-helpers)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Test successful JSON response headers/body.; Test error response shape.; Test unsupported values are handled safely.  
  **Done when:** Errors follow a stable `{error:{code,message,trace_id}}` shape.; No raw Go error details leak to user responses.

### P1 phase gate

- [x] All P1 tasks above are complete or explicitly deferred with a reason.
- [x] All code tasks in P1 have failing-test evidence before implementation.
- [x] `go test ./...` passes after P1 code tasks.
- [x] `go vet ./...` passes after P1 code tasks.
- [x] Relevant docs are updated with any changed behavior or assumptions.

---

## P2 — Chat API and UI skeleton
**Phase outcome:** Chat API and minimal UI exist with deterministic placeholder behavior.
- [x] **P2-T01 — Define chat domain models**  
  **Type:** Code  
  **Goal:** Create typed request/response structs for messages, intents, sources, actions, and handoff status.  
  **Primary files:** `internal/domain/chat.go`, `internal/domain/chat_test.go`  
  **Prompt:** [`P2-T01` in task-prompts.md](task-prompts.md#p2-t01-define-chat-domain-models)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** JSON marshal/unmarshal tests for `ChatRequest` and `ChatResponse`.; Validation tests for missing message and invalid student ID shape.  
  **Done when:** Domain structs are provider-neutral.; Response can include source citations, tool actions, and escalation metadata.
- [x] **P2-T02 — Implement `POST /api/v1/chat` handler with deterministic placeholder response**  
  **Type:** Code  
  **Goal:** Create the public chat API contract before real AI orchestration exists.  
  **Primary files:** `internal/handlers/chat.go`, `internal/handlers/chat_test.go`  
  **Prompt:** [`P2-T02` in task-prompts.md](task-prompts.md#p2-t02-implement-post-api-v1-chat-handler-with-deterministic-placeholder-response)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** `httptest` verifies valid request returns 200.; Invalid JSON returns 400.; Missing message returns 400.; Unsupported method returns 405.  
  **Done when:** Handler depends on an interface, not concrete orchestrator.; Trace ID appears in response.; Tests do not call network or model APIs.
- [x] **P2-T03 — Serve minimal Go web chat UI**  
  **Type:** Code  
  **Goal:** Provide an interview-friendly UI without making frontend complexity the focus.  
  **Primary files:** `web/templates/chat.html`, `web/static/app.js`, `web/static/app.css`, `internal/handlers/ui.go`, `internal/handlers/ui_test.go`  
  **Prompt:** [`P2-T03` in task-prompts.md](task-prompts.md#p2-t03-serve-minimal-go-web-chat-ui)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** `httptest` verifies chat page renders.; Template test verifies API endpoint path is present.; Static file route returns expected content type.  
  **Done when:** User can type a message and see a response.; UI clearly marks demo/synthetic mode.; UI can be replaced later without changing orchestrator.
- [x] **P2-T04 — Add request validation and safe client-facing errors**  
  **Type:** Code  
  **Goal:** Keep malformed requests and accidental PII from destabilizing the demo.  
  **Primary files:** `internal/validation`, `internal/handlers/chat_test.go`  
  **Prompt:** [`P2-T04` in task-prompts.md](task-prompts.md#p2-t04-add-request-validation-and-safe-client-facing-errors)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Table-driven tests for empty, oversized, and whitespace-only messages.; Student ID validation accepts synthetic IDs only.; Error body never includes raw request body.  
  **Done when:** Oversized messages are rejected.; Synthetic ID rule is explicit.; Error responses are consistent.
- [x] **P2-T05 — Add in-memory conversation session store**  
  **Type:** Code  
  **Goal:** Track a short demo conversation so follow-up questions can reference prior transcript/payment context.  
  **Primary files:** `internal/session/store.go`, `internal/session/store_test.go`  
  **Prompt:** [`P2-T05` in task-prompts.md](task-prompts.md#p2-t05-add-in-memory-conversation-session-store)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Create, append, read, and expire session tests.; Concurrent access test with `-race` target.; PII redaction before persistence test.  
  **Done when:** Store is concurrency-safe.; TTL is configurable.; Only redacted or synthetic data is stored.

### P2 phase gate

- [x] All P2 tasks above are complete or explicitly deferred with a reason.
- [x] All code tasks in P2 have failing-test evidence before implementation.
- [x] `go test ./...` passes after P2 code tasks.
- [x] Relevant docs are updated with any changed behavior or assumptions.

---

## P3 — Synthetic enterprise APIs and clients
**Phase outcome:** Synthetic Banner/payment/CRM/LMS APIs and typed clients are working.
- [x] **P3-T01 — Create synthetic fixture loader**
  **Type:** Code  
  **Goal:** Load deterministic demo records for students, transcripts, payments, LMS access, and CRM examples.  
  **Primary files:** `data/synthetic-students.json`, `internal/fixtures/loader.go`, `internal/fixtures/loader_test.go`  
  **Prompt:** [`P3-T01` in task-prompts.md](task-prompts.md#p3-t01-create-synthetic-fixture-loader)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Valid fixture loads all expected students.; Duplicate IDs fail.; Missing required fields fail.; Fixture contains only synthetic IDs.  
  **Done when:** S100001-S100004 exist with expected states.; Fixtures can be reused by mock services and tests.
- [x] **P3-T02 — Build mock Banner-style student API**
  **Type:** Code  
  **Goal:** Simulate student profile, enrollment status, transcript request status, and holds.  
  **Primary files:** `cmd/mock-banner/main.go`, `internal/mock/banner`, `internal/tools/banner_client.go`  
  **Prompt:** [`P3-T02` in task-prompts.md](task-prompts.md#p3-t02-build-mock-banner-style-student-api)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Handler tests for known student, unknown student, transcript status, and financial hold.; Client tests against `httptest.Server`.; Contract tests validate response schema.  
  **Done when:** Known IDs return deterministic data.; Unknown IDs return 404.; No real Banner naming/secrets are used beyond mock labels.
- [x] **P3-T03 — Build mock payment API**
  **Type:** Code  
  **Goal:** Simulate transcript payment status without processing real payments.  
  **Primary files:** `cmd/mock-payment/main.go`, `internal/mock/payment`, `internal/tools/payment_client.go`  
  **Prompt:** [`P3-T03` in task-prompts.md](task-prompts.md#p3-t03-build-mock-payment-api)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Paid student returns paid status.; Unpaid student returns unpaid status.; Unknown payment returns clear safe error.; Client timeout test.  
  **Done when:** Payment response includes status, amount, currency, and synthetic transaction ID only.; The UI never accepts real card/payment data.
- [x] **P3-T04 — Build mock CRM API for case creation and queue routing**
  **Type:** Code  
  **Goal:** Show human handoff and staff-facing case summaries.  
  **Primary files:** `cmd/mock-crm/main.go`, `internal/mock/crm`, `internal/tools/crm_client.go`  
  **Prompt:** [`P3-T04` in task-prompts.md](task-prompts.md#p3-t04-build-mock-crm-api-for-case-creation-and-queue-routing)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Create case success.; Priority case success.; Validation rejects empty summary.; Client retries or returns typed error on 5xx.  
  **Done when:** Case response includes case ID, queue, priority, and summary.; Conversation summary is redacted before storage.
- [x] **P3-T05 — Build mock LMS API**
  **Type:** Code  
  **Goal:** Support basic Moodle/Brightspace-style learner access questions without real LMS access.  
  **Primary files:** `cmd/mock-lms/main.go`, `internal/mock/lms`, `internal/tools/lms_client.go`  
  **Prompt:** [`P3-T05` in task-prompts.md](task-prompts.md#p3-t05-build-mock-lms-api)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Known synthetic student returns LMS access status.; Unknown course returns safe fallback.; Client handles timeout.  
  **Done when:** LMS API is clearly marked synthetic.; It supports only demo access status, not course content.
- [x] **P3-T06 — Add typed enterprise clients with shared error model**
  **Type:** Code  
  **Goal:** Let the orchestrator call mock services through interfaces that resemble production integrations.  
  **Primary files:** `internal/tools/errors.go`, `internal/tools/*_client.go`, `internal/tools/*_test.go`  
  **Prompt:** [`P3-T06` in task-prompts.md](task-prompts.md#p3-t06-add-typed-enterprise-clients-with-shared-error-model)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Timeout tests with context cancellation.; 404 maps to typed not-found error.; 5xx maps to retryable or external-service error.; Malformed JSON maps to typed parse error.  
  **Done when:** No orchestrator code depends on raw HTTP details.; All tool calls accept context and emit trace ID headers.

### P3 phase gate

- [x] All P3 tasks above are complete or explicitly deferred with a reason.
- [x] All code tasks in P3 have failing-test evidence before implementation.
- [x] `go test ./...` passes after P3 code tasks.
- [x] Relevant docs are updated with any changed behavior or assumptions.

---

## P4 — Deterministic orchestration before AI
**Phase outcome:** Transcript/payment orchestration works without relying on live AI.
- [x] **P4-T01 — Define orchestrator ports and dependency injection**
  **Type:** Code  
  **Goal:** Make the orchestrator testable by depending on interfaces for retrieval, classification, LLM, tools, workflow, and audit.  
  **Primary files:** `internal/orchestrator/orchestrator.go`, `internal/orchestrator/orchestrator_test.go`  
  **Prompt:** [`P4-T01` in task-prompts.md](task-prompts.md#p4-t01-define-orchestrator-ports-and-dependency-injection)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Compile-time interface fakes.; Test orchestrator can be constructed with fake dependencies.; Nil dependency validation test.  
  **Done when:** No live network dependencies in orchestrator tests.; Every external dependency has a small interface.
- [x] **P4-T02 — Implement fallback intent and sentiment classifier**
  **Type:** Code  
  **Goal:** Provide deterministic behavior for demos and tests before using an LLM.  
  **Primary files:** `internal/classifier/fallback.go`, `internal/classifier/fallback_test.go`  
  **Prompt:** [`P4-T02` in task-prompts.md](task-prompts.md#p4-t02-implement-fallback-intent-and-sentiment-classifier)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Table-driven messages map to expected intents.; Frustration/urgency phrases map to negative/high urgency.; Unknown questions map to unknown with low confidence.  
  **Done when:** Classifier returns typed confidence.; Low confidence cannot trigger sensitive tools.
- [x] **P4-T03 — Implement transcript-status decision flow**
  **Type:** Code  
  **Goal:** Handle the core Tier 1 scenario using synthetic Banner and payment data.  
  **Primary files:** `internal/orchestrator/transcript.go`, `internal/orchestrator/transcript_test.go`  
  **Prompt:** [`P4-T03` in task-prompts.md](task-prompts.md#p4-t03-implement-transcript-status-decision-flow)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** S100001 -> ready/no workflow.; S100002 -> unpaid/payment reminder.; S100003 -> hold/CRM escalation.; S100004 -> unknown/handoff.; Missing student ID prompts for synthetic ID.  
  **Done when:** Response contains user-friendly answer plus machine-readable actions.; No payment reminder is sent for paid records.; Financial holds route to staff, not self-service.
- [x] **P4-T04 — Trigger payment reminder workflow from orchestrator**
  **Type:** Code  
  **Goal:** Connect the transcript decision flow to workflow automation.  
  **Primary files:** `internal/workflow/client.go`, `internal/orchestrator/transcript_test.go`  
  **Prompt:** [`P4-T04` in task-prompts.md](task-prompts.md#p4-t04-trigger-payment-reminder-workflow-from-orchestrator)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Unpaid transcript calls workflow exactly once.; Idempotency key is passed.; Workflow failure is reported safely and audited.; Paid transcript does not call workflow.  
  **Done when:** Workflow action appears in chat response.; Failures do not create duplicate reminders.; Audit event records workflow attempt.
- [x] **P4-T05 — Create CRM escalation summary and priority routing**
  **Type:** Code  
  **Goal:** Turn unresolved, urgent, or sensitive conversations into staff-ready mock cases.  
  **Primary files:** `internal/orchestrator/escalation.go`, `internal/orchestrator/escalation_test.go`  
  **Prompt:** [`P4-T05` in task-prompts.md](task-prompts.md#p4-t05-create-crm-escalation-summary-and-priority-routing)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Negative sentiment creates priority case.; Low confidence creates normal handoff.; Financial hold routes to Finance/Registrar queue.; Summary redacts emails/phone numbers.  
  **Done when:** CRM case includes intent, queue, priority, trace ID, and redacted summary.; Learner receives a clear handoff message.
- [x] **P4-T06 — Return action trace in chat responses**
  **Type:** Code  
  **Goal:** Make the demo transparent by showing what was checked, triggered, and escalated.  
  **Primary files:** `internal/domain/chat.go`, `internal/orchestrator/orchestrator_test.go`  
  **Prompt:** [`P4-T06` in task-prompts.md](task-prompts.md#p4-t06-return-action-trace-in-chat-responses)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Response includes tool action names and statuses.; Internal errors are not exposed.; Trace ID is present on all action results.  
  **Done when:** Interview demo can show decision flow without opening logs.; Action trace is safe for learner-facing display.

### P4 phase gate

- [x] All P4 tasks above are complete or explicitly deferred with a reason.
- [x] All code tasks in P4 have failing-test evidence before implementation.
- [x] `go test ./...` passes after P4 code tasks.
- [x] Relevant docs are updated with any changed behavior or assumptions.

---

## P5 — RAG ingestion and source-grounded answers
**Phase outcome:** Approved-source ingestion and retrieval produce grounded responses.
- [ ] **P5-T01 — Define source allowlist schema**  
  **Type:** Code  
  **Goal:** Represent approved public learner-service sources with freshness and risk metadata.  
  **Primary files:** `data/seed-sources.json`, `internal/rag/source.go`, `internal/rag/source_test.go`  
  **Prompt:** [`P5-T01` in task-prompts.md](task-prompts.md#p5-t01-define-source-allowlist-schema)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Valid source config parses.; Non-HTTPS URL fails.; Missing title/department fails.; Private-domain marker fails.  
  **Done when:** Sources include URL, title, department, risk level, retrieved date, and freshness flag.; No source is ingested unless allowlisted.
- [ ] **P5-T02 — Implement ingestion fetcher with allowlist and cleaning boundaries**  
  **Type:** Code  
  **Goal:** Fetch only approved public content and prepare text for chunking.  
  **Primary files:** `cmd/ingest/main.go`, `internal/rag/ingest.go`, `internal/rag/ingest_test.go`  
  **Prompt:** [`P5-T02` in task-prompts.md](task-prompts.md#p5-t02-implement-ingestion-fetcher-with-allowlist-and-cleaning-boundaries)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Fetcher rejects URL not in allowlist.; HTML cleaner removes nav/script/style.; Network failure returns typed error.; Test uses `httptest.Server`, not live internet.  
  **Done when:** Ingestion is deterministic in tests.; Content hash is stored.; Private pages cannot be fetched accidentally.
- [ ] **P5-T03 — Implement chunking with metadata preservation**  
  **Type:** Code  
  **Goal:** Split retrieved content into searchable chunks while preserving source URL and title.  
  **Primary files:** `internal/rag/chunk.go`, `internal/rag/chunk_test.go`  
  **Prompt:** [`P5-T03` in task-prompts.md](task-prompts.md#p5-t03-implement-chunking-with-metadata-preservation)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Short content produces one chunk.; Long content produces bounded chunks.; Chunk IDs are stable for same input.; Metadata is copied to every chunk.  
  **Done when:** Chunk size is configurable.; No empty chunks are stored.; Chunk IDs can be cited in responses.
- [ ] **P5-T04 — Create local retrieval implementation**  
  **Type:** Code  
  **Goal:** Provide a demo-safe retrieval path before wiring Azure AI Search or pgvector.  
  **Primary files:** `internal/rag/retrieve.go`, `internal/rag/local_retriever.go`, `internal/rag/retrieve_test.go`  
  **Prompt:** [`P5-T04` in task-prompts.md](task-prompts.md#p5-t04-create-local-retrieval-implementation)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Transcript query ranks transcript chunks first.; Unrelated query returns low confidence.; Retrieval limit is respected.; Source metadata is returned.  
  **Done when:** Retriever interface can be swapped later.; No hallucinated source links are created.
- [ ] **P5-T05 — Add grounded answer source packaging**  
  **Type:** Code  
  **Goal:** Attach source citations and confidence to chat responses regardless of LLM provider.  
  **Primary files:** `internal/orchestrator/grounded_answer.go`, `internal/orchestrator/grounded_answer_test.go`  
  **Prompt:** [`P5-T05` in task-prompts.md](task-prompts.md#p5-t05-add-grounded-answer-source-packaging)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Answer with sufficient chunks includes citations.; Low retrieval confidence returns safe fallback.; Duplicate sources are de-duplicated.; Source risk level is included for internal response.  
  **Done when:** Every policy/procedure answer includes at least one source or a fallback.; Unsupported claims are not invented.
- [ ] **P5-T06 — Flag stale or high-risk sources**  
  **Type:** Code  
  **Goal:** Avoid confidently answering from outdated or sensitive content.  
  **Primary files:** `internal/rag/freshness.go`, `internal/rag/freshness_test.go`  
  **Prompt:** [`P5-T06` in task-prompts.md](task-prompts.md#p5-t06-flag-stale-or-high-risk-sources)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Stale source triggers caution flag.; High-risk source requires stronger confidence or handoff.; Fresh source passes without warning.  
  **Done when:** Assistant can say it needs staff confirmation for stale/high-risk policy questions.; Dashboard can show stale-source warnings.

### P5 phase gate

- [ ] All P5 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P5 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P5 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.

---

## P6 — LLM gateway and structured classification
**Phase outcome:** LLM gateway, strict JSON classification, prompts, and guardrails are testable.
- [ ] **P6-T01 — Define LLM provider interface and request/response types**  
  **Type:** Code  
  **Goal:** Hide provider details behind typed Go interfaces.  
  **Primary files:** `internal/llm/types.go`, `internal/llm/types_test.go`  
  **Prompt:** [`P6-T01` in task-prompts.md](task-prompts.md#p6-t01-define-llm-provider-interface-and-request-response-types)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** JSON schema for answer request/response marshals correctly.; Provider-neutral error type works.; Timeout field validation.  
  **Done when:** Orchestrator imports interfaces/types only.; Provider can be replaced without handler changes.
- [ ] **P6-T02 — Implement OpenAI/Azure-compatible REST client**  
  **Type:** Code  
  **Goal:** Call an OpenAI-compatible chat/completions endpoint through a testable Go client.  
  **Primary files:** `internal/llm/openai_client.go`, `internal/llm/openai_client_test.go`  
  **Prompt:** [`P6-T02` in task-prompts.md](task-prompts.md#p6-t02-implement-openai-azure-compatible-rest-client)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** `httptest.Server` verifies request payload.; Client parses successful response.; Client handles 429/500/timeout.; No live API call in tests.  
  **Done when:** API key is read from config only.; Logs never include prompts containing PII.; Client supports context cancellation.
- [ ] **P6-T03 — Parse strict JSON classification output**  
  **Type:** Code  
  **Goal:** Convert model output into trusted typed classification only after validation.  
  **Primary files:** `internal/classifier/llm_parser.go`, `internal/classifier/llm_parser_test.go`  
  **Prompt:** [`P6-T03` in task-prompts.md](task-prompts.md#p6-t03-parse-strict-json-classification-output)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Valid JSON parses.; Malformed JSON fails safely.; Unknown intent maps to `unknown`.; Out-of-range confidence is rejected or clamped.; Tool-triggering is disabled below threshold.  
  **Done when:** Invalid model output never panics.; Low-confidence classification returns safe fallback.
- [ ] **P6-T04 — Create prompt templates for classification and grounded answers**  
  **Type:** Code  
  **Goal:** Make prompts versioned, testable, and aligned with privacy/safety rules.  
  **Primary files:** `internal/orchestrator/prompts.go`, `internal/orchestrator/prompts_test.go`  
  **Prompt:** [`P6-T04` in task-prompts.md](task-prompts.md#p6-t04-create-prompt-templates-for-classification-and-grounded-answers)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Prompt contains strict JSON instruction.; Prompt includes source-only answer rule.; Prompt includes privacy/no-real-data rule.; Golden test catches accidental prompt drift.  
  **Done when:** Prompts are plain text constants or embedded templates.; Prompt version is included in audit metadata.
- [ ] **P6-T05 — Add LLM fallback and guardrail behavior**  
  **Type:** Code  
  **Goal:** Keep the demo reliable if the model is unavailable or produces unsafe output.  
  **Primary files:** `internal/orchestrator/ai_guardrails.go`, `internal/orchestrator/ai_guardrails_test.go`  
  **Prompt:** [`P6-T05` in task-prompts.md](task-prompts.md#p6-t05-add-llm-fallback-and-guardrail-behavior)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Model timeout uses deterministic fallback.; Unsafe answer without sources is rejected.; Tool calls require validated classification.; Sensitive/unsupported requests escalate.  
  **Done when:** No core demo path depends solely on live model availability.; Guardrail failures are logged and visible in dashboard.
- [ ] **P6-T06 — Add end-to-end classification tests with fixture messages**  
  **Type:** Code  
  **Goal:** Measure intent and sentiment behavior before using the assistant in the demo.  
  **Primary files:** `data/classification-fixtures.jsonl`, `internal/classifier/e2e_test.go`  
  **Prompt:** [`P6-T06` in task-prompts.md](task-prompts.md#p6-t06-add-end-to-end-classification-tests-with-fixture-messages)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** At least 5 examples per supported intent.; Negative/urgent sentiment examples pass.; Unknown/off-topic examples do not trigger actions.  
  **Done when:** Fixture accuracy target is documented.; Failures identify which intent regressed.

### P6 phase gate

- [ ] All P6 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P6 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P6 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.

---

## P7 — Privacy, audit, and dashboard
**Phase outcome:** Privacy redaction, audit events, and dashboard metrics are in place.
- [ ] **P7-T01 — Implement PII redaction**  
  **Type:** Code  
  **Goal:** Redact emails, phone numbers, likely passwords, and non-synthetic IDs before logging or case summaries.  
  **Primary files:** `internal/privacy/redact.go`, `internal/privacy/redact_test.go`  
  **Prompt:** [`P7-T01` in task-prompts.md](task-prompts.md#p7-t01-implement-pii-redaction)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Email redaction.; Phone redaction.; Password phrase redaction.; Real-looking student number redaction.; Synthetic S100001-style IDs are allowed for demo.  
  **Done when:** Redactor is used by logs, audit, session store, and CRM summaries.; Tests cover false positives and false negatives.
- [ ] **P7-T02 — Create audit event schema and store**  
  **Type:** Code  
  **Goal:** Record safe operational events for intent, sources, tool calls, workflow triggers, and escalations.  
  **Primary files:** `internal/audit/event.go`, `internal/audit/store.go`, `internal/audit/store_test.go`  
  **Prompt:** [`P7-T02` in task-prompts.md](task-prompts.md#p7-t02-create-audit-event-schema-and-store)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Record and query events by trace ID.; Event payload is redacted.; Concurrent writes are safe.; Retention timestamp is set.  
  **Done when:** Every orchestrator action can produce an audit event.; Audit store can run in memory for demo and later be backed by DB.
- [ ] **P7-T03 — Wire redacted structured logging across services**  
  **Type:** Code  
  **Goal:** Ensure logs are useful for debugging but safe for learner privacy.  
  **Primary files:** `internal/middleware/logging.go`, `internal/privacy/logging_test.go`  
  **Prompt:** [`P7-T03` in task-prompts.md](task-prompts.md#p7-t03-wire-redacted-structured-logging-across-services)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Chat message is not logged raw.; Trace ID is logged.; Tool status is logged without raw payload.; Secret-like values are redacted.  
  **Done when:** No log line contains real-looking PII from tests.; Logger uses `slog` fields, not string concatenation.
- [ ] **P7-T04 — Build admin metrics endpoint**  
  **Type:** Code  
  **Goal:** Summarize containment, escalation, top intents, confidence, workflow count, and low-confidence questions.  
  **Primary files:** `internal/handlers/admin.go`, `internal/handlers/admin_test.go`, `internal/audit/metrics.go`  
  **Prompt:** [`P7-T04` in task-prompts.md](task-prompts.md#p7-t04-build-admin-metrics-endpoint)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Metrics calculation from seeded audit events.; Endpoint requires admin/mock auth.; Empty store returns zeros safely.  
  **Done when:** Dashboard data can be generated after demo conversation.; Metrics are aggregate and redacted.
- [ ] **P7-T05 — Create minimal admin dashboard UI**  
  **Type:** Code  
  **Goal:** Give stakeholders an at-a-glance view of AI adoption and operational risk.  
  **Primary files:** `web/templates/admin.html`, `web/static/admin.js`, `internal/handlers/admin_ui.go`  
  **Prompt:** [`P7-T05` in task-prompts.md](task-prompts.md#p7-t05-create-minimal-admin-dashboard-ui)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Template renders metric labels.; Unauthorized request is rejected.; Review queue renders redacted question text only.  
  **Done when:** Dashboard shows top intents, escalations, workflows, low-confidence items, and stale-source warnings.; No raw PII appears.
- [ ] **P7-T06 — Add retention and export controls**  
  **Type:** Code  
  **Goal:** Document and enforce how demo audit data is kept, deleted, and exported.  
  **Primary files:** `internal/audit/retention.go`, `internal/audit/retention_test.go`, `docs/privacy-impact-lite.md`  
  **Prompt:** [`P7-T06` in task-prompts.md](task-prompts.md#p7-t06-add-retention-and-export-controls)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Expired events are purged.; Export excludes raw messages.; Retention config default is short for demo.  
  **Done when:** Privacy doc matches implementation behavior.; Demo can reset all stored data.

### P7 phase gate

- [ ] All P7 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P7 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P7 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.

---

## P8 — Workflow automation and Power Automate option
**Phase outcome:** Workflow simulator and optional Power Automate webhook path are ready.
- [ ] **P8-T01 — Build local workflow simulator service**  
  **Type:** Code  
  **Goal:** Simulate Power Automate-style HTTP-triggered flows for offline demos.  
  **Primary files:** `cmd/workflow-sim/main.go`, `internal/workflow/sim_handler.go`, `internal/workflow/sim_handler_test.go`  
  **Prompt:** [`P8-T01` in task-prompts.md](task-prompts.md#p8-t01-build-local-workflow-simulator-service)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Payment reminder payload accepted.; Missing idempotency key rejected.; Duplicate idempotency key returns same workflow ID.; Invalid payload returns 400.  
  **Done when:** Simulator returns workflow ID and status.; Workflow events are visible in dashboard.
- [ ] **P8-T02 — Implement workflow idempotency and retry policy**  
  **Type:** Code  
  **Goal:** Prevent duplicate payment reminders and handle transient failures safely.  
  **Primary files:** `internal/workflow/idempotency.go`, `internal/workflow/client_test.go`  
  **Prompt:** [`P8-T02` in task-prompts.md](task-prompts.md#p8-t02-implement-workflow-idempotency-and-retry-policy)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Same trace/student/action creates same idempotency key.; Transient 500 is retried within limit.; Permanent 400 is not retried.; Context cancellation stops retry.  
  **Done when:** No duplicate learner reminders in repeated calls.; Retry behavior is observable in audit logs.
- [ ] **P8-T03 — Add optional Power Automate webhook client**  
  **Type:** Code  
  **Goal:** Show RPA/workflow integration relevance while keeping local simulator as default.  
  **Primary files:** `internal/workflow/powerautomate.go`, `internal/workflow/powerautomate_test.go`  
  **Prompt:** [`P8-T03` in task-prompts.md](task-prompts.md#p8-t03-add-optional-power-automate-webhook-client)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Client sends expected JSON schema to `httptest.Server`.; Signature/header is included when configured.; Webhook URL missing falls back to simulator config.; Secrets are not logged.  
  **Done when:** Same interface supports simulator and Power Automate.; Docs explain secure webhook storage.
- [ ] **P8-T04 — Audit workflow outcomes and errors**  
  **Type:** Code  
  **Goal:** Make workflow automation measurable and debuggable.  
  **Primary files:** `internal/orchestrator/transcript.go`, `internal/audit/event.go`, `internal/workflow/client_test.go`  
  **Prompt:** [`P8-T04` in task-prompts.md](task-prompts.md#p8-t04-audit-workflow-outcomes-and-errors)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Success event recorded.; Failure event recorded with safe message.; Dashboard count increments.; No raw payload is stored.  
  **Done when:** Every workflow attempt has trace ID, action, status, and idempotency key hash.
- [ ] **P8-T05 — Document Power Automate flow schema and setup**  
  **Type:** Documentation  
  **Goal:** Provide enough detail for an interviewer to see how the Go webhook would connect to Power Automate.  
  **Primary files:** `docs/power-automate-flow.md`  
  **Prompt:** [`P8-T05` in task-prompts.md](task-prompts.md#p8-t05-document-power-automate-flow-schema-and-setup)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Review JSON examples against client tests.; Confirm docs do not include real webhook secrets.  
  **Done when:** Webhook request/response examples are present.; Security notes mention idempotency, secret storage, and replay protection.

### P8 phase gate

- [ ] All P8 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P8 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P8 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.

---

## P9 — Evaluation runner and quality gates
**Phase outcome:** Evaluation runner produces quality reports and fails critical regressions.
- [ ] **P9-T01 — Create JSONL evaluation dataset**  
  **Type:** Code  
  **Goal:** Define repeatable test cases for intents, sources, workflow actions, escalations, and safety behavior.  
  **Primary files:** `data/eval-questions.jsonl`, `internal/eval/dataset.go`, `internal/eval/dataset_test.go`  
  **Prompt:** [`P9-T01` in task-prompts.md](task-prompts.md#p9-t01-create-jsonl-evaluation-dataset)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Dataset parser handles valid JSONL.; Invalid rows fail with line number.; Required expected fields are validated.; Dataset includes critical safety cases.  
  **Done when:** At least 30 evaluation examples exist.; Each case has expected intent and expected action/source behavior.
- [ ] **P9-T02 — Build `cmd/eval` runner**  
  **Type:** Code  
  **Goal:** Run the assistant against the evaluation dataset and produce repeatable quality results.  
  **Primary files:** `cmd/eval/main.go`, `internal/eval/runner.go`, `internal/eval/runner_test.go`  
  **Prompt:** [`P9-T02` in task-prompts.md](task-prompts.md#p9-t02-build-cmd-eval-runner)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Runner calls fake chat client for each case.; Timeouts are captured as failures.; Results include latency.; CLI exits zero when all critical tests pass.  
  **Done when:** `go run ./cmd/eval` works locally.; Runner can target in-process fake or live local API.
- [ ] **P9-T03 — Implement scoring functions**  
  **Type:** Code  
  **Goal:** Score intent accuracy, source grounding, action correctness, escalation precision, safety, and latency.  
  **Primary files:** `internal/eval/score.go`, `internal/eval/score_test.go`  
  **Prompt:** [`P9-T03` in task-prompts.md](task-prompts.md#p9-t03-implement-scoring-functions)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Intent exact match scoring.; Expected source present scoring.; Expected action present/absent scoring.; Critical hallucination scoring.; Latency threshold scoring.  
  **Done when:** Scores are deterministic.; Critical policy errors are separated from minor misses.
- [ ] **P9-T04 — Generate JSON and Markdown evaluation reports**  
  **Type:** Code  
  **Goal:** Create portfolio-ready evidence of model/system quality.  
  **Primary files:** `internal/eval/report.go`, `internal/eval/report_test.go`, `reports/eval-summary.md`  
  **Prompt:** [`P9-T04` in task-prompts.md](task-prompts.md#p9-t04-generate-json-and-markdown-evaluation-reports)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** JSON report contains summary and per-case results.; Markdown report contains metrics table.; Report redacts user messages when configured.  
  **Done when:** Reports can be committed as sample outputs.; Failures are readable by non-technical stakeholders.
- [ ] **P9-T05 — Fail builds on critical evaluation failures**  
  **Type:** Code  
  **Goal:** Make responsible AI behavior part of the development gate, not a manual afterthought.  
  **Primary files:** `cmd/eval/main.go`, `internal/eval/gates.go`, `internal/eval/gates_test.go`, `Makefile`  
  **Prompt:** [`P9-T05` in task-prompts.md](task-prompts.md#p9-t05-fail-builds-on-critical-evaluation-failures)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Critical hallucination returns non-zero exit.; Missing required escalation returns non-zero exit.; Minor accuracy miss can be warning based on threshold.; Make target propagates exit code.  
  **Done when:** `make eval` fails for critical safety regressions.; README explains quality gate.
- [ ] **P9-T06 — Create review queue for failed/low-confidence answers**  
  **Type:** Code  
  **Goal:** Show how the system would improve over time through human review.  
  **Primary files:** `internal/eval/review_queue.go`, `internal/handlers/admin.go`, `internal/eval/review_queue_test.go`  
  **Prompt:** [`P9-T06` in task-prompts.md](task-prompts.md#p9-t06-create-review-queue-for-failed-low-confidence-answers)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Low-confidence answer added to queue.; Failed critical eval added to queue.; Duplicate review items collapse by normalized question.; Resolved item no longer appears as open.  
  **Done when:** Dashboard can display unresolved review items.; Review items contain sources/actions but no raw PII.

### P9 phase gate

- [ ] All P9 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P9 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P9 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.

---

## P10 — Docker, CI, and local developer experience
**Phase outcome:** Docker, CI, env safety, and smoke test support repeatable demos.
- [ ] **P10-T01 — Create Dockerfiles for Go services**  
  **Type:** Code  
  **Goal:** Package API and mock services as small reproducible containers.  
  **Primary files:** `Dockerfile`, `cmd/*`, `.dockerignore`  
  **Prompt:** [`P10-T01` in task-prompts.md](task-prompts.md#p10-t01-create-dockerfiles-for-go-services)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Add a shell or Makefile smoke test that builds the API image.; Unit tests still pass before image build.; Container starts and responds to `/healthz` in local smoke script.  
  **Done when:** Image does not include source secrets.; Non-root runtime user is used where practical.; Build is documented.
- [ ] **P10-T02 — Create Docker Compose local stack**  
  **Type:** Code  
  **Goal:** Start API, mock services, workflow simulator, and optional database with one command.  
  **Primary files:** `docker-compose.yml`, `.env.example`, `Makefile`  
  **Prompt:** [`P10-T02` in task-prompts.md](task-prompts.md#p10-t02-create-docker-compose-local-stack)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** `make compose-test` or smoke script waits for health endpoints.; Missing optional LLM config still allows deterministic demo mode.; Service URLs match README.  
  **Done when:** `make dev` or `docker compose up` starts the full demo.; Ports are documented.; Synthetic mode is default.
- [ ] **P10-T03 — Add CI workflow for Go tests and evaluation gates**  
  **Type:** Code  
  **Goal:** Demonstrate professional delivery discipline.  
  **Primary files:** `.github/workflows/ci.yml`, `Makefile`  
  **Prompt:** [`P10-T03` in task-prompts.md](task-prompts.md#p10-t03-add-ci-workflow-for-go-tests-and-evaluation-gates)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** CI runs `go test ./...`.; CI runs `go vet ./...`.; CI runs `make eval` with deterministic/fake providers.; CI fails on critical eval failure.  
  **Done when:** Pull requests cannot pass with failing tests.; CI avoids live model/API dependencies.
- [ ] **P10-T04 — Add environment sample and secret-safety checks**  
  **Type:** Code  
  **Goal:** Prevent accidental commit of real keys or private data.  
  **Primary files:** `.env.example`, `.gitignore`, `internal/config/config_test.go`  
  **Prompt:** [`P10-T04` in task-prompts.md](task-prompts.md#p10-t04-add-environment-sample-and-secret-safety-checks)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Config tests prove secrets are optional for deterministic mode.; Secret scanner or grep script rejects known key patterns in repo.; `.env` is ignored.  
  **Done when:** Only placeholders appear in docs.; No real webhook URL or API key is committed.
- [ ] **P10-T05 — Add one-command smoke test**  
  **Type:** Code  
  **Goal:** Make the project easy to verify during interview prep.  
  **Primary files:** `scripts/smoke.sh`, `Makefile`, `docs/test-plan.md`  
  **Prompt:** [`P10-T05` in task-prompts.md](task-prompts.md#p10-t05-add-one-command-smoke-test)  
  **Quality gate:** Strict TDD required. Use the prompt and tests before production code.  
  **Tests/review:** Smoke test starts stack or assumes stack is running.; Calls `/healthz`.; Calls transcript chat scenario.; Verifies workflow and CRM outputs for synthetic IDs.  
  **Done when:** `make smoke` gives a clear pass/fail.; Failures include actionable messages.

### P10 phase gate

- [ ] All P10 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P10 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P10 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.

---

## P11 — Portfolio polish and interview readiness
**Phase outcome:** Portfolio materials are polished for interview presentation.
- [ ] **P11-T01 — Polish README for applicant storytelling**  
  **Type:** Documentation  
  **Goal:** Make the project understandable to a hiring manager in under two minutes.  
  **Primary files:** `README.md`  
  **Prompt:** [`P11-T01` in task-prompts.md](task-prompts.md#p11-t01-polish-readme-for-applicant-storytelling)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Read README top-to-bottom and confirm it explains problem, solution, stack, demo, and privacy.  
  **Done when:** README includes screenshot placeholder, quickstart, architecture, success metrics, and TDD quality statement.
- [ ] **P11-T02 — Add architecture diagram and sequence diagram**  
  **Type:** Documentation  
  **Goal:** Show enterprise integration thinking visually.  
  **Primary files:** `docs/architecture.md`  
  **Prompt:** [`P11-T02` in task-prompts.md](task-prompts.md#p11-t02-add-architecture-diagram-and-sequence-diagram)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Trace golden path through diagram and verify every component exists in tasks or code.  
  **Done when:** Diagram includes chat UI, Go API, orchestrator, RAG, mock Banner/payment/CRM/LMS, workflow, audit, dashboard.
- [ ] **P11-T03 — Finalize 5-7 minute demo script**  
  **Type:** Documentation  
  **Goal:** Prepare a concise interview walkthrough that highlights the role requirements.  
  **Primary files:** `docs/demo-script.md`  
  **Prompt:** [`P11-T03` in task-prompts.md](task-prompts.md#p11-t03-finalize-5-7-minute-demo-script)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Dry-run the script and verify it fits the time box.; Confirm every spoken claim can be shown in UI, logs, docs, or tests.  
  **Done when:** Demo covers Tier 0 answer, Tier 1 workflow, urgent escalation, dashboard, and TDD/evaluation evidence.
- [ ] **P11-T04 — Prepare screenshots or short GIF placeholders**  
  **Type:** Documentation  
  **Goal:** Make the GitHub repo visually scannable.  
  **Primary files:** `README.md`, `docs/demo-script.md`  
  **Prompt:** [`P11-T04` in task-prompts.md](task-prompts.md#p11-t04-prepare-screenshots-or-short-gif-placeholders)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Check screenshots do not show real student data, real tokens, or private URLs.  
  **Done when:** README has image placeholders or links.; Captions explain what each screen proves.
- [ ] **P11-T05 — Run final release checklist**  
  **Type:** Documentation  
  **Goal:** Confirm the portfolio is coherent, safe, and demonstrable.  
  **Primary files:** `docs/phases-and-tasks.md`, `docs/test-plan.md`, `docs/privacy-impact-lite.md`  
  **Prompt:** [`P11-T05` in task-prompts.md](task-prompts.md#p11-t05-run-final-release-checklist)  
  **Quality gate:** Review-gated documentation task.  
  **Tests/review:** Run quickstart commands.; Run tests and eval.; Review docs for consistency.; Check links.  
  **Done when:** All must-have tasks are checked or clearly deferred.; ZIP/repo is ready to share.; Known limitations are honest and documented.

### P11 phase gate

- [ ] All P11 tasks above are complete or explicitly deferred with a reason.
- [ ] All code tasks in P11 have failing-test evidence before implementation.
- [ ] `go test ./...` passes after P11 code tasks.
- [ ] Relevant docs are updated with any changed behavior or assumptions.


---

## MVP release gate

Before calling the MVP interview-ready:

- [ ] `make test` passes.
- [ ] `make eval` passes with zero critical policy errors.
- [ ] `make smoke` passes against the local stack.
- [ ] Demo uses only synthetic learner records.
- [ ] Transcript request, unpaid payment reminder, financial-hold escalation, and urgent-sentiment escalation all work.
- [ ] Dashboard shows containment, escalation, workflow, low-confidence, and stale-source metrics.
- [ ] README, architecture, demo script, privacy notes, test plan, and evaluation report are consistent.
- [ ] Known limitations are documented honestly.
