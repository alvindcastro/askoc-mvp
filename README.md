# AskOC AI Concierge

**Go/Golang portfolio MVP for an AI/Automation Solutions Developer role**

AskOC AI Concierge is a privacy-aware learner-service automation MVP designed for a higher-education digital learner experience team. It demonstrates how conversational AI, retrieval-augmented generation, workflow automation, and mock enterprise integrations can reduce Tier 0 and Tier 1 learner-service volume while improving response quality and routing.

The implementation is intentionally **Go-first**: the API gateway, AI orchestrator, mock enterprise integrations, ingestion jobs, workflow simulator, audit service, and evaluation runner are all designed to be written in Go.

## One-sentence pitch

AskOC AI Concierge is a **Go-based AI learner-service agent** that uses RAG, workflow automation, mock Banner/CRM/LMS/payment integrations, and privacy-first synthetic data controls to answer common student questions and automate transcript/payment follow-up.

## Two-minute reviewer path

1. Read the problem, scope, and privacy boundary in this README.
2. Run the repeatable local proof:

```bash
make smoke
```

3. Keep the stack running when you want to inspect the UI:

```bash
make compose-up
```

4. Open `http://localhost:9080/chat` for the themed learner chat and `http://localhost:9080/admin` for the themed protected dashboard. The local admin token is `demo-admin-token`.
5. Skim [docs/developer-guide.md](docs/developer-guide.md), [docs/demo-script.md](docs/demo-script.md), [docs/architecture.md](docs/architecture.md), and [reports/eval-summary.md](reports/eval-summary.md) for local testing, the 5-7 minute walkthrough, integration diagrams, and quality evidence.

## Portfolio evidence at a glance

| Evidence | What it proves | Where to verify |
|---|---|---|
| One-command smoke test | Docker Compose can run the synthetic transcript/payment demo and CRM handoff | `make smoke`, `scripts/smoke.sh` |
| Developer testing guide | Local stack startup, alternate ports, manual checks, troubleshooting, and useful operating notes are documented in one place | [docs/developer-guide.md](docs/developer-guide.md) |
| Go service architecture | Chat UI, Go API, orchestrator, RAG, mock enterprise APIs, workflow, audit, dashboard, and eval are separated by typed boundaries | [docs/architecture.md](docs/architecture.md) |
| Responsible-AI gate | Intent, source, action, handoff, safety, and critical hallucination checks are repeatable | `make eval`, [reports/eval-summary.md](reports/eval-summary.md) |
| Privacy boundary | All learner, payment, LMS, workflow, and CRM data is visibly synthetic and redacted before audit/dashboard use | [docs/privacy-impact-lite.md](docs/privacy-impact-lite.md), `data/synthetic-students.json` |
| Themed UI proof | `DESIGN.md` is applied as AskOC visual polish: square borders, high contrast, active red route underline, source chips, action trace rows, trace IDs, admin review rows, and visible focus states | `/chat`, `/admin`, `go test ./internal/handlers` |
| TDD delivery | Code phases were implemented behind failing tests, package tests, full-suite tests, eval, smoke, and changelog evidence | [docs/phases-and-tasks.md](docs/phases-and-tasks.md), [CHANGELOG.md](CHANGELOG.md) |

## Screenshot and GIF placeholders

These placeholders are intentionally textual until final captures are generated from the local synthetic stack. See [docs/assets/README.md](docs/assets/README.md) for the capture manifest. Any future image or GIF must show only `S10000X` demo IDs, mock workflow IDs, mock CRM IDs, local URLs, and placeholder tokens.

| Placeholder | Intended capture | Caption/proof |
|---|---|---|
| Chat grounded answer | `/chat` after asking “How do I order my official transcript?” | Shows a Tier 0 answer with approved transcript source, confidence, risk, freshness, and safe action trace. |
| Transcript workflow | `/chat` after the `S100002` transcript-status prompt | Shows mock Banner/payment checks and `payment_reminder_triggered` without real payment data. |
| Escalation case | `/chat` after the `S100003` financial-hold or urgent-sentiment prompt | Shows privacy-aware CRM handoff, queue, priority, and synthetic case ID. |
| Admin dashboard | `/admin` after demo conversations | Shows containment, escalation, workflow, stale-source, low-confidence, and review-queue metrics using redacted data. |
| Eval report | `reports/eval-summary.md` after `make eval` | Shows zero critical safety failures and portfolio-readable quality evidence. |

## Why Go for this MVP

Go is a strong fit for an applicant demo because it lets the project look enterprise-ready without becoming heavy:

- fast HTTP APIs and webhook services,
- simple deployable binaries,
- clean concurrency with `context` and goroutines,
- strong typing for tool calls and workflow events,
- easy Docker deployment,
- excellent fit for microservices, integrations, and automation middleware.

## Problem statement

Learners often contact college service teams with repeatable questions about transcripts, fees, registration, portals, LMS access, refunds, key dates, and application status. Many inquiries are handled through manual email triage, which can delay answers and consume administrative time.

This MVP shows how a college could:

- answer common learner questions instantly,
- ground answers in approved public content,
- automate low-risk transactions,
- escalate complex or sensitive issues to staff,
- provide analytics for continuous improvement,
- protect learner privacy through synthetic data, redaction, and audit logging.

## MVP scope

The MVP focuses on one high-value use case:

> A learner asks how to order an official transcript and why a request has not been processed. The assistant retrieves official guidance, checks a mock student/payment record, triggers a payment reminder workflow if needed, and escalates frustrated or unresolved cases to a mock CRM queue.

## Key capabilities

| Capability | Included in MVP |
|---|---|
| Conversational AI | Web chatbot with grounded answers and safe fallback |
| RAG | Retrieval over curated public learner-service content |
| Intent recognition | Transcript/payment workflow plus safe fallback for other learner-service questions |
| Sentiment analysis | Flags frustrated or urgent learners for priority routing |
| Workflow automation | Transcript payment/status follow-up and CRM routing |
| Enterprise integration | Mock Banner, payment, CRM, LMS, and workflow APIs |
| Admin dashboard | Containment, escalation, top intents, confidence, unresolved questions |
| Privacy by design | Synthetic records, PII redaction, audit logs, retention rules |
| Deployment | Dockerized Go services with optional Azure deployment path |


## Execution plan and strict TDD

The build plan is now organized around phase checklists and task-level prompts:

- [docs/brainstorm.md](docs/brainstorm.md) — MVP brainstorming notes and implementation decisions.
- [docs/phases-and-tasks.md](docs/phases-and-tasks.md) — phase-by-phase tickable task board.
- [docs/task-prompts.md](docs/task-prompts.md) — detailed prompts for each task.
- [docs/tdd-policy.md](docs/tdd-policy.md) — strict TDD policy for every Go code task.

Every code task follows the same non-negotiable loop:

```text
write failing test → verify failure → implement minimal Go code → run tests → refactor while green
```

A task is not done until the relevant package test and `go test ./...` pass. AI, RAG, workflow, privacy, and orchestration changes also require safe fallback or evaluation tests.

## Suggested Go-first tech stack

| Layer | Suggested implementation |
|---|---|
| Web UI | Go server-rendered chat/admin UI with `html/template`, small JavaScript, and `DESIGN.md` theme tokens |
| API gateway | Go `net/http` or `chi` router |
| AI orchestration | Go service using typed interfaces for LLM, retriever, classifier, and tools |
| LLM gateway | Azure OpenAI / OpenAI-compatible REST client written in Go |
| Retrieval | Azure AI Search REST client, PostgreSQL + pgvector, or local vector store |
| Intent/sentiment | LLM structured JSON output, Azure AI Language REST, or lightweight Go classifier |
| Automation | Power Automate cloud flow triggered by Go webhook client; local Go workflow simulator for demo |
| Mock systems | Go microservices for Banner, CRM, payment, LMS, and workflow simulation |
| Database | PostgreSQL for app/audit data; SQLite acceptable for local MVP |
| Observability | `log/slog`, trace IDs, OpenTelemetry-ready middleware |
| Testing | Go table-driven tests, `httptest`, contract tests, and JSONL evaluation runner |
| Deployment | Docker Compose locally; Azure Container Apps or Kubernetes optional |

## Target user journeys

### Journey 1: Tier 0 knowledge answer

1. Learner asks: “How do I order my official transcript?”
2. Go API accepts `POST /api/v1/chat`.
3. Orchestrator identifies `transcript_request` intent.
4. Retriever returns approved source chunks.
5. LLM gateway produces a grounded answer with source links.
6. Audit service logs trace ID, intent, confidence, and source IDs.

### Journey 2: Tier 1 transaction support

1. Learner asks: “I ordered my transcript but it hasn’t been processed.”
2. Assistant asks for a synthetic student ID for the demo.
3. Go mock Banner API confirms student status.
4. Go mock payment API checks transcript payment status.
5. If unpaid, Go orchestration service triggers a Power Automate webhook or local Go workflow simulator.
6. If paid but blocked, Go mock CRM API creates a case for Registrar follow-up.

### Journey 3: Escalation with sentiment signal

1. Learner says: “This is really frustrating. I need this for a job application.”
2. Sentiment is classified as negative/urgent.
3. Assistant summarizes the case.
4. Mock CRM case is created with transcript context, payment status, conversation summary, and priority flag.

## Repository structure

```text
askoc-ai-concierge/
  README.md
  go.mod
  Makefile
  Dockerfile
  docker-compose.yml
  .dockerignore
  .env.example
  .github/workflows/ci.yml
  cmd/
    api/                  # API server with health/readiness, chat API, and chat UI routes
    eval/                 # JSONL evaluation runner and quality gate reports
    ingest/               # Approved-source ingestion to local RAG chunks
    mock-banner/          # Synthetic student profile, transcript status, and hold API
    mock-payment/         # Synthetic transcript payment status API
    mock-crm/             # Synthetic CRM case creation API
    mock-lms/             # Synthetic LMS access-status API
    workflow-sim/         # Local Power Automate-style workflow simulator
  internal/
    audit/
    classifier/
    config/
    domain/
    eval/
    fixtures/
    handlers/
    llm/
    middleware/
    mock/
    orchestrator/
    privacy/
    rag/
    session/
    tools/
    validation/
    build/
  web/
    templates/
    static/
  data/
    eval-questions.jsonl
    synthetic-students.json
    seed-sources.json
    rag-chunks.json
  reports/
    eval-summary.json
    eval-summary.md
  scripts/
    check-secrets.sh
    smoke.sh
  docs/
    ...
```

The P11-ready chat API uses guarded orchestration and evaluation: deterministic fallback intent/sentiment classification remains the default, while optional `openai-compatible` provider mode adds a tested REST LLM gateway, strict JSON classification parsing, versioned prompts, source-only answer guardrails, local RAG retrieval over approved public chunks, typed mock Banner/payment/CRM clients, idempotent payment-reminder workflow clients, a standalone local workflow simulator, an optional Power Automate-compatible webhook client with retry and signature headers, a safe action trace, CRM handoff routing for holds, urgent sentiment, low confidence, or explicit human handoff, shared PII redaction, an in-memory audit event store, protected admin metrics, a minimal dashboard, audit export/reset/purge controls, a redacted eval review queue endpoint, a deterministic JSONL evaluation runner, Docker Compose local stack, offline CI gate, safe env sample, secret check, one-command smoke test, architecture diagrams, demo script, screenshot placeholders, and final release checklist.

## Quickstart and commands

```bash
make dev
make test
make test-race
make eval
make secret-check
make docker-build
make compose-up
make compose-test
make smoke
go test ./...
go vet ./...
go test ./internal/build -run TestP10
go test ./internal/eval ./cmd/eval
go run ./cmd/ingest -sources data/seed-sources.json -out data/rag-chunks.json
go run ./cmd/eval -input data/eval-questions.jsonl -output reports/eval-summary.json -markdown-output reports/eval-summary.md
go run ./cmd/workflow-sim
go run ./cmd/mock-banner
go run ./cmd/mock-payment
go run ./cmd/mock-crm
go run ./cmd/mock-lms
```

For the repeatable Docker demo, run `make smoke`. It builds the API and mock-service images with Docker Compose, waits for `/healthz`, posts the unpaid `S100002` transcript-status scenario, and posts the `S100003` financial-hold scenario that creates a mock CRM case. Use `make compose-up` when you want to keep the stack running, then `make compose-test` to smoke-test an already running default-port stack. The Compose stack uses `ASKOC_PROVIDER=stub`, synthetic fixtures, service-DNS URLs such as `http://mock-banner:9081`, and the local workflow simulator URL by default. Host ports default to `9080`-`9085`; if that range is already in use, choose another `9xxx` range with the override commands in [docs/developer-guide.md](docs/developer-guide.md).

For manual local development without containers, start the mock Banner, payment, CRM, and optionally workflow simulator services in separate terminals before `make dev`. The API loads local RAG chunks from `data/rag-chunks.json` at startup and talks to typed mock services through configurable local URLs. If `ASKOC_WORKFLOW_URL` is empty, the API uses the in-process idempotent workflow client; set `ASKOC_WORKFLOW_URL=http://localhost:9084/api/v1/automation/payment-reminder` to route reminders through `cmd/workflow-sim`, or point it at a Power Automate HTTP trigger for the optional webhook path. Auth is disabled by default for learner chat. Admin metrics, unresolved eval review items, audit export, purge, and reset routes require a bearer token; by default use `demo-admin-token`, or set `ASKOC_AUTH_TOKEN=<demo-token>` to reuse the configured mock token.

`make eval` is the responsible-AI quality gate. It runs `cmd/eval` against `data/eval-questions.jsonl` using the deterministic in-process evaluator by default, writes `reports/eval-summary.json` and `reports/eval-summary.md`, and exits non-zero for critical safety regressions such as unsupported critical claims or missing required escalation.

`make secret-check` scans tracked and unignored local files for known live-token patterns while allowing `.env.example` placeholders. `.env` and local override files are ignored by default.

Current environment settings:

| Variable | Default | Purpose |
|---|---|---|
| `ASKOC_HTTP_ADDR` | `:9080` | API listen address |
| `ASKOC_AUTH_ENABLED` | `false` | Enables mock bearer-token auth |
| `ASKOC_AUTH_TOKEN` | empty | Demo bearer token when auth is enabled |
| `ASKOC_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error` |
| `ASKOC_WORKFLOW_URL` | empty | Optional workflow webhook URL; empty uses the in-process client; redacted from config output |
| `ASKOC_WORKFLOW_TIMEOUT_SECONDS` | `5` | Tool and workflow webhook client timeout |
| `ASKOC_WORKFLOW_SIGNATURE` | empty | Optional workflow webhook signature/header value; redacted from config output |
| `ASKOC_WORKFLOW_SIGNATURE_HEADER` | `X-AskOC-Workflow-Signature` | Header name used when `ASKOC_WORKFLOW_SIGNATURE` is set |
| `ASKOC_WORKFLOW_MAX_RETRIES` | `1` | Retry count for transient workflow webhook `5xx` responses |
| `ASKOC_BANNER_URL` | `http://localhost:9081` | Mock Banner base URL used by P4 orchestration |
| `ASKOC_PAYMENT_URL` | `http://localhost:9082` | Mock payment base URL used by P4 orchestration |
| `ASKOC_CRM_URL` | `http://localhost:9083` | Mock CRM base URL used by P4 orchestration |
| `ASKOC_RAG_CHUNKS_PATH` | `data/rag-chunks.json` | Local approved-source chunks used by P5 retrieval |
| `ASKOC_PROVIDER` | `stub` | `stub` keeps deterministic mode; `openai-compatible` enables the tested REST LLM gateway |
| `ASKOC_PROVIDER_MODEL` | `demo-placeholder` | Provider model/deployment name when LLM mode is enabled |
| `ASKOC_PROVIDER_ENDPOINT` | empty | OpenAI-compatible or Azure chat completions endpoint; required only for `openai-compatible` |
| `ASKOC_PROVIDER_TIMEOUT_SECONDS` | `5` | LLM request timeout |
| `ASKOC_PROVIDER_API_KEY` | empty | Provider API key; required only for `openai-compatible` and redacted from config output |
| `ASKOC_API_PORT` | `9080` | Optional Docker Compose host port override for the API |
| `ASKOC_BANNER_PORT` | `9081` | Optional Docker Compose host port override for mock Banner |
| `ASKOC_PAYMENT_PORT` | `9082` | Optional Docker Compose host port override for mock payment |
| `ASKOC_CRM_PORT` | `9083` | Optional Docker Compose host port override for mock CRM |
| `ASKOC_WORKFLOW_PORT` | `9084` | Optional Docker Compose host port override for workflow simulator |
| `ASKOC_LMS_PORT` | `9085` | Optional Docker Compose host port override for mock LMS |

Current service URLs:

```text
Chat UI:   http://localhost:9080/chat
Chat API:  http://localhost:9080/api/v1/chat
Admin UI:  http://localhost:9080/admin
Admin API: http://localhost:9080/api/v1/admin/metrics
Review API: http://localhost:9080/api/v1/admin/review-items
Health:    http://localhost:9080/healthz
Readiness: http://localhost:9080/readyz

Mock Banner:  http://localhost:9081/api/v1/students/S100002
Mock Payment: http://localhost:9082/api/v1/students/S100002/payment-status
Mock CRM:     http://localhost:9083/api/v1/crm/cases
Workflow Sim: http://localhost:9084/api/v1/automation/payment-reminder
Mock LMS:     http://localhost:9085/api/v1/students/S100001/lms-access?course_id=DEMO-LMS-101
```

The chat API validates JSON requests, rejects empty or oversized messages, accepts synthetic student IDs in the `S` plus six digits format, includes trace IDs in responses and action results, routes transcript/payment decisions through the orchestrator, and uses P5 retrieval plus P6 source guardrails for transcript-request answers.
P3 tool clients forward `X-Trace-ID` headers and map not-found, retryable, parse, timeout, and external-service failures into typed errors. P4 adds deterministic classifier/orchestrator tests and an in-process workflow port that returns idempotent synthetic workflow IDs. P5 adds allowlist parsing, deterministic ingestion, chunking, local retrieval, and stale/high-risk source fallback tests. P6 adds the optional tested LLM gateway, strict JSON parser, prompt golden tests, classification fixtures, and low-confidence/source guardrails. P7 adds shared redaction for logs, sessions, audit payloads, and CRM summaries; audit events for orchestrator actions, workflow outcomes, guardrails, and escalations; protected aggregate admin metrics; redacted review queue items; and demo audit retention/export/reset controls. P8 adds `cmd/workflow-sim`, a Power Automate-compatible HTTP client, idempotency-key hashing in workflow audit metadata, and retry attempt counts for webhook responses. P9 adds `data/eval-questions.jsonl`, `cmd/eval`, `internal/eval`, JSON/Markdown reports, critical gate failures, and unresolved eval review queue support. P10 adds multi-service Docker packaging, local Compose orchestration, CI, env safety, and smoke verification. P11 adds portfolio-facing README polish, diagram review, a timed demo script, synthetic screenshot placeholders, and a final release checklist.

## Demo data policy

Do not use real student data. Use synthetic records only.

| Student ID | Name | Transcript payment | Hold | Expected result |
|---|---|---:|---|---|
| `S100001` | Demo Learner One | Paid | None | Ready for processing |
| `S100002` | Demo Learner Two | Unpaid | Mock payment hold | Payment reminder workflow |
| `S100003` | Demo Learner Three | Review required | Mock financial hold | CRM escalation |
| `S100004` | Demo Learner Four | Not applicable | None | Human handoff |

The same fixture now includes synthetic LMS account status and demo course-access records. It does not include LMS course content.

## Success metrics

| Metric | Demo target |
|---|---:|
| Grounded answer accuracy | 90%+ on curated test questions |
| Intent classification accuracy | 85%+ |
| Tier 0/Tier 1 containment | 50–70% in simulation |
| Human escalation precision | 90% for low-confidence or sensitive cases |
| Average response time | Under 5 seconds |
| Critical hallucination rate | 0 critical policy errors in test set |
| Automation completion | 100% for happy-path transcript workflow |
| Audit coverage | Orchestrator actions, workflow outcomes, guardrails, and escalations emit redacted audit events |

## MVP boundaries

### In scope

- Go public API and orchestrator,
- public-content RAG,
- mock Go student/payment/CRM/LMS APIs,
- synthetic student records,
- transcript/payment workflow,
- local workflow simulator and optional webhook path,
- admin analytics dashboard,
- privacy and audit documentation.

### Out of scope

- real OC authentication,
- real Banner integration,
- real payment processing,
- real student records,
- scraping private portals,
- production fine-tuned LLM deployment,
- nice-to-have workflows beyond the transcript/payment demo slice.

## References for content grounding

Use public, approved pages only during the MVP build. Examples:

- Okanagan College Office of the Registrar: https://www.okanagancollege.ca/office-of-the-registrar
- Okanagan College transcript request guidance: https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards
- Okanagan College online resources: https://www.okanagancollege.ca/student-handbook/online-resources
- Okanagan College LMS information: https://www.okanagancollege.ca/teaching-and-learning-innovations/lms
- Okanagan College information security guide: https://www.okanagancollege.ca/it-services/it-security/information-security-guide
- Go documentation: https://go.dev/doc/
- Go web applications tutorial: https://go.dev/doc/articles/wiki/
- Go RESTful API with Gin tutorial: https://go.dev/doc/tutorial/web-service-gin
- Azure for Go developers: https://learn.microsoft.com/en-us/azure/developer/go/
- Azure AI for Go developers: https://learn.microsoft.com/en-us/azure/developer/go/azure-ai-for-go-developers
- Azure AI Search RAG overview: https://learn.microsoft.com/en-us/azure/search/retrieval-augmented-generation-overview
- Azure AI Language documentation: https://learn.microsoft.com/en-us/azure/ai-services/language-service/
- Azure OpenAI REST API reference: https://learn.microsoft.com/en-us/azure/foundry/openai/reference
- Microsoft Power Automate documentation: https://learn.microsoft.com/en-us/power-automate/
