# Interview Demo Script

## Goal

Show a 5–7 minute Go-based MVP that maps directly to the AI/Automation Solutions Developer role: conversational AI, RAG, enterprise integrations, workflow automation, privacy, monitoring, and continuous improvement.

## Opening pitch

> “This is AskOC AI Concierge, a Go-based learner-service automation MVP for higher education. It answers transcript questions from approved source chunks, checks synthetic Banner and payment records, triggers an idempotent workflow reminder, escalates urgent or blocked cases into a mock CRM, shows redacted dashboard evidence, and backs the demo with Go tests, an evaluation gate, Docker Compose, and a one-command smoke test. It uses only synthetic learner data and does not require live AI by default.”

## Time-boxed run sheet

| Minute | Focus | Screen or command | Proof point |
|---:|---|---|---|
| 0:00-0:45 | Problem, Go architecture, privacy boundary | README and chat UI | Hiring manager can understand the project quickly and see synthetic-data-only scope. |
| 0:45-1:45 | Tier 0 grounded answer | `/chat` transcript-order question | Point at source chips, confidence, risk, freshness, and trace ID proof. |
| 1:45-3:15 | Tier 1 transcript/payment workflow | `/chat` with `S100002` | Point at synthetic action rows, workflow ID, idempotency key, and mock Banner/payment proof. |
| 3:15-4:30 | Urgent or blocked escalation | `/chat` urgent message or `S100003` | Point at CRM case ID, priority chip, queue, and safe synthetic handoff copy. |
| 4:30-5:30 | Dashboard evidence | `/admin` | Point at redacted aggregate metrics, review queue filter, status chips, stale-source count, and audit controls. |
| 5:30-7:00 | TDD/evaluation/release proof | terminal and reports | `make test`, `make eval`, `make smoke`, and reports show repeatable quality gates. |

## Demo setup

Run locally:

```bash
make smoke
make compose-up
# Optional refresh if public pages are reachable:
go run ./cmd/ingest -sources data/seed-sources.json -out data/rag-chunks.json
```

For manual non-container development, start `cmd/mock-banner`, `cmd/mock-payment`, `cmd/mock-crm`, and `cmd/workflow-sim` in separate terminals before `make dev`.

If port `8080` is already in use during prep, run `ASKOC_API_PORT=18080 make smoke` and open `http://localhost:18080`.

Open:

```text
http://localhost:8080
```

Services:

```text
API/UI:       http://localhost:8080
Admin UI:     http://localhost:8080/admin
Mock Banner:  http://localhost:8081
Mock Payment: http://localhost:8082
Mock CRM:     http://localhost:8083
Workflow:     http://localhost:8084/api/v1/automation/payment-reminder or in-process fallback
Mock LMS:     http://localhost:8085
Dashboard:    protected admin metrics and redacted review queue
```

## Minute 1: Grounded Tier 0 answer

Ask:

```text
How do I order my official transcript?
```

Show:

- intent is `transcript_request`,
- response includes source link, source confidence, risk level, and freshness status,
- response is concise,
- answer avoids unsupported claims,
- deterministic action trace is visible,
- the themed proof panel shows the active route underline, `Trace ID`, `Sources`, and `Action trace` sections without exposing raw learner data.

Talking point:

> “This answer is grounded by the P5 local retriever. The API only cites chunks from `data/rag-chunks.json`, which is generated from the approved source allowlist, and it falls back or asks for staff confirmation when confidence is low or the source is stale/high-risk.”

## Minute 2: Tier 1 transaction support

Ask:

```text
I ordered my transcript but it has not been processed. My student ID is S100002.
```

Show:

- Go API receives chat request,
- orchestrator classifies `transcript_status`,
- mock Banner API is called,
- mock payment API is called,
- payment status is `unpaid`,
- assistant explains the next step.

Talking point:

> “This mirrors an enterprise integration pattern without touching real systems. Banner, payment, and CRM are represented by typed Go clients and synthetic APIs.”

## Minute 3: Workflow automation

Show workflow event:

```text
payment_reminder_triggered
```

Show:

- idempotency key,
- workflow ID,
- tested audit-port event,
- safe reminder summary,
- the chat action trace row uses synthetic labels and mono IDs instead of raw payloads.

Talking point:

> “The workflow boundary is idempotent and testable without external services. P8 adds a standalone Go simulator plus an optional Power Automate-compatible webhook client behind the same interface, including retry handling and redacted audit metadata.”

## Minute 4: Sentiment and escalation

Ask:

```text
This is really frustrating. I need this transcript for a job application.
```

Show:

- sentiment classified as negative/urgent,
- assistant does not promise impossible turnaround,
- mock CRM case created,
- case summary is minimal and privacy-aware,
- learner receives case ID.

Talking point:

> “Sentiment does not make final decisions alone. It increases routing priority when combined with unresolved context and safe business rules.”

## Minute 5: P11 audit dashboard, redaction, and evaluation

Talking point:

> “The response body still shows the safe decision trace directly, and the system records redacted audit events for orchestrator actions, workflow outcomes, guardrails, and escalations. Workflow audit metadata stores a hashed idempotency key rather than the raw key. The admin dashboard summarizes top intents, containment, escalations, workflows, low-confidence review items, stale-source warnings, and unresolved eval review items without displaying raw PII.”

Open:

```text
http://localhost:8080/admin
```

Use the default local admin token:

```text
demo-admin-token
```

Show:

- total conversations and containment rate,
- top intents,
- workflow success/failure counts,
- low-confidence review items with redacted question text,
- stale-source warning count,
- audit export/reset/purge controls,
- the review queue filter, trace IDs, priority/status chips, and `redacted` marker.

## Minute 6: Go architecture walkthrough

Show repository structure:

```text
cmd/api
cmd/mock-banner
cmd/mock-payment
cmd/mock-crm
internal/orchestrator
internal/rag
internal/classifier
internal/workflow
internal/tools
internal/audit
```

Talking point:

> “I used Go because this role needs reliable integrations and automation services. The model is only one part of the system; the Go services enforce timeouts, typed tool calls, redaction, audit logging, and safe fallback.”

## Minute 7: Evaluation and tests

Run:

```bash
go test ./...
go test ./internal/privacy ./internal/audit ./internal/handlers
go test ./internal/llm ./internal/classifier ./internal/orchestrator
go test ./internal/rag ./internal/orchestrator
go test ./internal/classifier ./internal/workflow ./internal/orchestrator
go test ./internal/build -run TestP10
make eval
make secret-check
make smoke
```

Show:

- intent accuracy,
- source retrieval and stale/high-risk fallback,
- workflow decision accuracy,
- safe action traces,
- CRM summary redaction and audit/dashboard redaction tests.

Talking point:

> “I treat this as a maintained automation product. The test suite has deterministic unit coverage for RAG allowlisting, ingestion, chunking, retrieval, source fallback, classification, LLM gateway behavior, prompt guardrails, workflow idempotency, simulator contracts, webhook retries, transcript decisions, action traces, CRM escalation, shared redaction, audit storage, admin metrics, dashboard rendering, retention controls, repo-level Docker and CI contracts, and a broader JSONL evaluation runner that fails critical regressions.”
> “P6 adds the guarded LLM layer: OpenAI-compatible calls are optional, strict JSON is validated before use, prompt templates are versioned, and low-confidence or ungrounded model output falls back instead of triggering tools.”

## Expected demo path

All demo records are synthetic. Student IDs, payment states, holds, workflow IDs, and CRM cases must come from local Go services or fixtures only.

| Step | Scenario | Input | Expected output | Observable proof |
|---:|---|---|---|---|
| 1 | Transcript answer | “How do I order my official transcript?” | Grounded answer with approved transcript source chunk | Chat response shows `transcript_request`, source ID/link, confidence, risk/freshness metadata, and no unsupported claims |
| 2 | Unpaid payment workflow | “I ordered my transcript but it has not been processed. My student ID is S100002.” | Payment status is unpaid and reminder workflow is accepted | Chat response shows `transcript_status`, `payment_status_checked`, `payment_reminder_triggered`, workflow ID, and no CRM handoff |
| 3 | Financial-hold escalation | “My transcript still is not moving. My student ID is S100003.” | Financial hold is detected and staff handoff is created | Chat response shows `transcript_status`, `financial_hold_detected`, CRM case ID, and Registrar/Student Accounts handoff |
| 4 | Urgent sentiment escalation | “This is really frustrating. I need this transcript for a job application.” | Urgent/negative sentiment creates a priority CRM case | Chat response shows urgent sentiment, `crm_case_created`, priority flag, case ID, and privacy-aware summary |
| 5 | P11 TDD/eval/smoke evidence | Run package tests, eval gate, secret check, and smoke gate | Tests prove LLM gateway behavior, strict classification, prompt guardrails, source grounding, deterministic decisions, workflow simulator/webhook behavior, redaction, audit metrics, dashboard controls, Docker/CI contracts, env safety, eval quality gates, and release readiness | `go test ./internal/eval ./cmd/eval`, `go test ./internal/build -run TestP10`, `go test ./...`, `make eval`, `make secret-check`, and `make smoke` pass |

## Demo acceptance matrix

Each P11 row must be verifiable from the Go API response, local RAG chunks, local mock service logs, the workflow simulator or in-process workflow response, CRM simulator output, protected admin metrics, redacted audit-store tests, Docker/CI artifact tests, smoke-script output, or `reports/eval-summary.md`. Source references are approved public/curated learner-service content; synthetic records are the only data used for student/payment/hold state.

| ID | Scenario | Synthetic input | Expected intent | Expected source | Expected action | Expected handoff behavior | Pass evidence |
|---|---|---|---|---|---|---|---|
| D01 | Transcript answer | “How do I order my official transcript?” | `transcript_request` | P5 local chunk `oc-transcript-request-2005-onwards-seed-001` or refreshed chunk from the same allowlisted source | Return concise grounded answer with source link | No handoff; keep learner in chat | Response includes source, retrieval confidence, risk/freshness metadata, and zero critical unsupported claims |
| D02 | Unpaid payment workflow | `S100002` transcript status prompt | `transcript_status` | Transcript/payment guidance source chunk plus synthetic payment record `S100002` | Check mock Banner, check mock Payment, trigger `payment_reminder_triggered` | No CRM handoff unless workflow fails or confidence is low | Response includes unpaid status, workflow ID, idempotency key, and redacted workflow audit events |
| D03 | Financial-hold escalation | `S100003` transcript status prompt | `transcript_status` | Transcript/hold guidance source chunk plus synthetic Banner record `S100003` | Check mock Banner, detect `financial_hold`, create CRM case | Handoff to Registrar/Student Accounts queue with minimal summary and case ID | Response includes hold-safe wording, CRM case ID, queue/priority, and no payment reminder |
| D04 | Urgent sentiment escalation | Frustrated urgent transcript message | `escalation_request` or human handoff intent with urgent sentiment | Transcript source chunk when transcript context is present | Classify sentiment as urgent/negative and create priority CRM case | Priority handoff to staff; assistant does not promise a deadline or outcome | Response includes priority flag, CRM case ID, redacted summary, and safe expectation-setting |
| D05 | Low-confidence source fallback | Transcript-adjacent question with no approved source | `unknown` or low-confidence transcript intent | No acceptable source above threshold | Do not answer from model memory; ask clarifying question or route to staff | Low-confidence handoff if learner needs account-specific help | Response shows no citation used for unsupported claim and logs low-confidence review item |

Golden-path pass condition: D01-D04 pass in order during the 5-7 minute demo, and each expected action has at least one observable output. D05 is a safety check used when discussing fallback behavior.

## Screenshot and GIF placeholders

Final captures should be generated only from the local synthetic stack. The placeholder manifest lives in [docs/assets/README.md](assets/README.md). Before adding any image or GIF, inspect it for real student data, real tokens, private URLs, browser profile details, and unredacted raw messages.

| Placeholder | Capture moment | Caption |
|---|---|---|
| `docs/assets/chat-grounded-answer-placeholder` | Transcript-order Tier 0 answer in `/chat` | Proves source-grounded answer with approved transcript citation, confidence, risk, and freshness metadata. |
| `docs/assets/transcript-workflow-placeholder` | `S100002` transcript-status flow in `/chat` | Proves mock Banner/payment checks and idempotent payment reminder workflow with synthetic IDs only. |
| `docs/assets/crm-escalation-placeholder` | `S100003` financial hold or urgent-sentiment escalation | Proves mock CRM handoff, priority routing, safe expectation-setting, and minimal redacted summary. |
| `docs/assets/admin-dashboard-placeholder` | `/admin` after demo scenarios | Proves aggregate containment, escalation, workflow, stale-source, low-confidence, and review-queue metrics without raw PII. |
| `docs/assets/eval-report-placeholder` | `reports/eval-summary.md` after `make eval` | Proves repeatable responsible-AI evidence and zero critical safety failures. |

## Backup plan

If LLM or internet access fails during demo:

- use pre-seeded local RAG chunks,
- use deterministic classifier fallback,
- show mock API responses,
- run Go tests and evaluation report,
- explain how cloud AI is abstracted behind `internal/llm`.

## Closing line

> “This MVP shows that I can build more than a chatbot. I can build a Go-based, privacy-aware AI automation service that integrates with enterprise systems, reduces email-driven work, and gives staff measurable insight into learner support demand.”
