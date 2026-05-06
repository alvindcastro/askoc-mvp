# AskOC AI Concierge

**Go/Golang portfolio MVP for an AI/Automation Solutions**

AskOC AI Concierge is a privacy-aware learner-service automation MVP designed for a higher-education digital learner experience team. It demonstrates how conversational AI, retrieval-augmented generation, workflow automation, and enterprise integrations can reduce Tier 0 and Tier 1 learner-service volume while improving response quality and routing.

The implementation is intentionally **Go-first**: the API gateway, AI orchestrator, mock enterprise integrations, ingestion jobs, workflow simulator, audit service, and evaluation runner are all designed to be written in Go.

## One-sentence pitch

AskOC AI Concierge is a **Go-based AI learner-service agent** that uses RAG, intent detection, sentiment analysis, mock Banner/CRM/LMS integrations, and Power Automate-style workflows to answer common student questions and automate transcript/payment follow-up.

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
| Intent recognition | Transcript, fees, portal, LMS, registration, refund, application status, key dates, escalation |
| Sentiment analysis | Flags frustrated or urgent learners for priority routing |
| Workflow automation | Transcript payment/status follow-up and CRM routing |
| Enterprise integration | Mock Banner, payment, CRM, LMS, and notification APIs |
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
| Web UI | Go server-rendered chat UI with `html/template` + HTMX, or optional React/Next.js |
| API gateway | Go `net/http` or `chi` router |
| AI orchestration | Go service using typed interfaces for LLM, retriever, classifier, and tools |
| LLM gateway | Azure OpenAI / OpenAI-compatible REST client written in Go |
| Retrieval | Azure AI Search REST client, PostgreSQL + pgvector, or local vector store |
| Intent/sentiment | LLM structured JSON output, Azure AI Language REST, or lightweight Go classifier |
| Automation | Power Automate cloud flow triggered by Go webhook client; local Go workflow simulator for demo |
| Mock systems | Go microservices for Banner, CRM, payment, LMS, and notifications |
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
  go.sum
  docker-compose.yml
  Makefile
  cmd/
    api/                  # public chat/API gateway
    mock-banner/          # synthetic Banner-style student API
    mock-crm/             # mock case management API
    mock-payment/         # synthetic payment-status API
    mock-lms/             # synthetic LMS-status API
    workflow-sim/         # local Power Automate-style simulator
    ingest/               # public content ingestion and chunking
    eval/                 # model/RAG evaluation runner
  internal/
    audit/
    classifier/
    config/
    handlers/
    llm/
    middleware/
    orchestrator/
    privacy/
    rag/
    tools/
    workflow/
  web/
    templates/
    static/
  data/
    synthetic-students.json
    eval-questions.jsonl
    seed-sources.json
  docs/
    architecture.md
    golang-implementation.md
    privacy-impact-lite.md
    api-spec.md
    test-plan.md
    model-evaluation.md
    demo-script.md
    implementation-roadmap.md
```

## Local demo commands

```bash
make dev
make seed
make ingest
make test
make eval
make smoke
```

Example service URLs:

```text
Chat UI:      http://localhost:8080
API:          http://localhost:8080/api/v1
Mock Banner:  http://localhost:8081
Mock Payment: http://localhost:8082
Mock CRM:     http://localhost:8083
Workflow sim: http://localhost:8084
Dashboard:    http://localhost:8080/admin
```

## Demo data policy

Do not use real student data. Use synthetic records only.

| Student ID | Name | Transcript payment | Hold | Expected result |
|---|---|---:|---|---|
| `S100001` | Maya Chen | Paid | None | Ready for processing |
| `S100002` | Jordan Patel | Unpaid | None | Payment reminder workflow |
| `S100003` | Alex Morgan | Paid | Financial hold | CRM escalation |
| `S100004` | Sam Rivera | Unknown | None | Human handoff |

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
| Audit coverage | 100% of tool/action calls logged |

## MVP boundaries

### In scope

- Go public API and orchestrator,
- public-content RAG,
- mock Go student/payment/CRM/LMS APIs,
- synthetic student records,
- transcript workflow,
- admin analytics dashboard,
- privacy and audit documentation.

### Out of scope

- real OC authentication,
- real Banner integration,
- real payment processing,
- real student records,
- scraping private portals,
- production fine-tuned LLM deployment.

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
