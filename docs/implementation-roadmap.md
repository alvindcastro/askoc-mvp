# Implementation Roadmap

This roadmap is the high-level path for the Go-based AskOC AI Concierge MVP. The detailed working board is [Phases and Tickable Tasks](phases-and-tasks.md), and the copy/paste implementation prompts are in [Detailed Task Prompts](task-prompts.md).

## Roadmap rules

- [x] Freeze the MVP around transcript/payment support, source-grounded answers, workflow automation, escalation, and dashboard evidence.
- [ ] Use Go for backend services, mock integrations, workflow simulator, ingestion, evaluation, and dashboards unless a task explicitly says otherwise.
- [x] Use synthetic learner records only.
- [ ] Enforce [Strict TDD Policy](tdd-policy.md) on every code task.
- [x] Do not connect to real Banner, CRM, LMS, payment, authentication, or private portal systems.
- [x] Defer nice-to-have workflows and real integrations until after the MVP.

## Phase checklist

| Status | Phase | Primary outcome | Detailed tasks |
|---|---|---|---|
| [x] | P0 — Product framing and applicant strategy | Applicant story, scope, privacy boundary, source limits, and demo matrix are frozen. | [Task list](phases-and-tasks.md#p0-product-framing-and-applicant-strategy) |
| [x] | P1 — Go project foundation | Go project compiles, serves health endpoints, and has API hygiene. | [Task list](phases-and-tasks.md#p1-go-project-foundation) |
| [x] | P2 — Chat API and UI skeleton | Chat API and minimal UI exist with deterministic placeholder behavior. | [Task list](phases-and-tasks.md#p2-chat-api-and-ui-skeleton) |
| [x] | P3 — Synthetic enterprise APIs and clients | Synthetic Banner/payment/CRM/LMS APIs and typed clients are working. | [Task list](phases-and-tasks.md#p3-synthetic-enterprise-apis-and-clients) |
| [x] | P4 — Deterministic orchestration before AI | Transcript/payment orchestration works without relying on live AI. | [Task list](phases-and-tasks.md#p4-deterministic-orchestration-before-ai) |
| [ ] | P5 — RAG ingestion and source-grounded answers | Approved-source ingestion and retrieval produce grounded responses. | [Task list](phases-and-tasks.md#p5-rag-ingestion-and-source-grounded-answers) |
| [ ] | P6 — LLM gateway and structured classification | LLM gateway, strict JSON classification, prompts, and guardrails are testable. | [Task list](phases-and-tasks.md#p6-llm-gateway-and-structured-classification) |
| [ ] | P7 — Privacy, audit, and dashboard | Privacy redaction, audit events, and dashboard metrics are in place. | [Task list](phases-and-tasks.md#p7-privacy-audit-and-dashboard) |
| [ ] | P8 — Workflow automation and Power Automate option | Workflow simulator and optional Power Automate webhook path are ready. | [Task list](phases-and-tasks.md#p8-workflow-automation-and-power-automate-option) |
| [ ] | P9 — Evaluation runner and quality gates | Evaluation runner produces quality reports and fails critical regressions. | [Task list](phases-and-tasks.md#p9-evaluation-runner-and-quality-gates) |
| [ ] | P10 — Docker, CI, and local developer experience | Docker, CI, env safety, and smoke test support repeatable demos. | [Task list](phases-and-tasks.md#p10-docker-ci-and-local-developer-experience) |
| [ ] | P11 — Portfolio polish and interview readiness | Portfolio materials are polished for interview presentation. | [Task list](phases-and-tasks.md#p11-portfolio-polish-and-interview-readiness) |


## Phase gates

### P0 — Product framing and applicant strategy

- [x] Applicant story is clear and role-specific.
- [x] Synthetic-data and privacy boundary is documented.
- [x] Public source allowlist is defined.
- [x] Golden demo scenarios are measurable.
- [x] Nice-to-haves are deferred.

### P1 — Go project foundation

- [x] `go test ./...` passes.
- [x] `make test` passes.
- [x] `go vet ./...` passes.
- [x] Health and readiness endpoints are tested.
- [x] Middleware has trace, recovery, logging, and mock auth tests.
- [x] Error response shape is stable.

### P2 — Chat API and UI skeleton

- [x] `POST /api/v1/chat` is tested with success and error cases.
- [x] Minimal chat UI renders from Go.
- [x] Validation rejects empty, oversized, and unsafe inputs.
- [x] Session store is concurrency-safe if implemented.
- [x] No live AI dependency is needed for demo placeholder behavior.

### P3 — Synthetic enterprise APIs and clients

- [x] Synthetic fixture loader is tested.
- [x] Mock Banner, payment, CRM, and LMS APIs respond deterministically.
- [x] Typed clients have timeout, error, and malformed-response tests.
- [x] Contract tests verify request/response shapes.
- [x] No real system credentials or records are used.

### P4 — Deterministic orchestration before AI

- [x] Orchestrator depends on interfaces and is tested with fakes.
- [x] Transcript-status flow handles S100001-S100004 correctly.
- [x] Payment reminder workflow fires only for unpaid transcript cases.
- [x] CRM escalation occurs for holds, urgent sentiment, low confidence, or human handoff.
- [x] Chat response includes safe action trace.

### P5 — RAG ingestion and source-grounded answers

- [ ] Source allowlist schema is tested.
- [ ] Ingestion uses allowlisted public content only.
- [ ] Chunking preserves source metadata.
- [ ] Retrieval returns transcript sources for transcript questions.
- [ ] Low-confidence or stale-source answers fall back or escalate.

### P6 — LLM gateway and structured classification

- [ ] LLM provider interface is provider-neutral.
- [ ] REST client is tested with `httptest.Server`, not live APIs.
- [ ] Strict JSON parser rejects malformed/unsafe output.
- [ ] Prompt templates have golden tests.
- [ ] Guardrails prevent model output from directly triggering tools.

### P7 — Privacy, audit, and dashboard

- [ ] Redaction tests cover emails, phones, likely passwords, and non-synthetic IDs.
- [ ] Audit events are redacted and queryable by trace ID.
- [ ] Logs do not contain raw learner messages.
- [ ] Admin metrics endpoint is tested.
- [ ] Dashboard shows containment, escalation, workflows, and review queue.

### P8 — Workflow automation and Power Automate option

- [ ] Local workflow simulator is tested.
- [ ] Idempotency prevents duplicate reminders.
- [ ] Optional Power Automate client uses same interface as simulator.
- [ ] Workflow success/failure events are audited.
- [ ] Webhook schema and security notes are documented.

### P9 — Evaluation runner and quality gates

- [ ] JSONL dataset includes at least 30 examples.
- [ ] `cmd/eval` runs against fake or local chat API.
- [ ] Scoring covers intent, source, action, escalation, safety, and latency.
- [ ] Markdown and JSON reports are generated.
- [ ] Critical hallucination or missed escalation fails the build.

### P10 — Docker, CI, and local developer experience

- [ ] Docker images build.
- [ ] Docker Compose starts the local demo stack.
- [ ] CI runs Go tests, vet, and deterministic evaluation.
- [ ] `.env.example` is safe and `.env` is ignored.
- [ ] `make smoke` verifies the golden demo path.

### P11 — Portfolio polish and interview readiness

- [ ] README explains the project in under two minutes.
- [ ] Architecture diagram and sequence diagram are included.
- [ ] Demo script fits 5-7 minutes.
- [ ] Screenshots/GIFs contain synthetic data only.
- [ ] Final release checklist is complete.

## Time-boxed delivery options

### 2-day version

- [x] P0 story and demo matrix.
- [x] P1 API foundation.
- [x] P2 chat endpoint with deterministic response.
- [x] P3 synthetic student/payment/LMS fixture and mock enterprise APIs.
- [x] P4 deterministic transcript flow for S100001-S100004.
- [x] Basic README and demo script.

### 1-week version

- [x] P1-P4 complete.
- [ ] P5 local RAG over curated source snippets.
- [ ] P6 fallback classifier plus optional LLM client.
- [ ] P7 redaction, audit store, and dashboard endpoint.
- [ ] P8 local workflow simulator.
- [ ] `go test ./...` and `make smoke` pass.

### 2-week version

- [ ] P1-P8 complete.
- [ ] P9 evaluation runner and quality gates.
- [ ] P10 Docker Compose and CI.
- [ ] P11 polished portfolio docs and demo assets.
- [ ] Optional Power Automate webhook integration.

## Recommended build order

1. P0-T01 to P0-T05: freeze story and boundaries.
2. P1-T01 to P1-T05: create tested Go foundation.
3. P2-T01 to P2-T05: expose the chat contract, UI skeleton, validation, and session store.
4. P3-T01 to P3-T06: build synthetic systems and typed clients.
5. P4-T01 to P4-T06: make the core transcript workflow work deterministically.
6. P7-T01 and P7-T02: add redaction/audit before adding AI calls.
7. P5-T01 to P5-T05: add grounded retrieval.
8. P6-T01 to P6-T05: add LLM capability behind tests and guardrails.
9. P8-T01 to P8-T04: complete automation story.
10. P9 and P10: prove quality and repeatability.
11. P11: polish the portfolio.
