# Interview Demo Script

## Goal

Show a 5–7 minute Go-based MVP that maps directly to the AI/Automation Solutions Developer role: conversational AI, RAG, enterprise integrations, workflow automation, privacy, monitoring, and continuous improvement.

## Opening pitch

> “This is AskOC AI Concierge, a Go-based AI learner-service automation MVP. It answers common Tier 0 and Tier 1 learner questions using source-grounded RAG, checks synthetic Banner/payment records, triggers a payment reminder workflow, escalates complex cases into a mock CRM, and logs everything for dashboard review.”

## Demo setup

Run locally:

```bash
make dev
make seed
make ingest
```

Open:

```text
http://localhost:8080
```

Services:

```text
API/UI:       http://localhost:8080
Mock Banner:  http://localhost:8081
Mock Payment: http://localhost:8082
Mock CRM:     http://localhost:8083
Workflow sim: http://localhost:8084
Dashboard:    http://localhost:8080/admin
```

## Minute 1: Grounded Tier 0 answer

Ask:

```text
How do I order my official transcript?
```

Show:

- intent is `transcript_request`,
- response includes source link,
- response is concise,
- answer avoids unsupported claims,
- source confidence is visible.

Talking point:

> “The answer is not free-form guessing. The Go orchestrator retrieves approved content first, then asks the LLM to answer only from those sources.”

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

Open workflow simulator or Power Automate mock screen.

Show:

- idempotency key,
- workflow ID,
- audit event,
- safe notification summary.

Talking point:

> “The workflow layer is webhook-compatible with Power Automate, but I built a local Go simulator so the demo is reliable and testable.”

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

## Minute 5: Admin dashboard

Open:

```text
http://localhost:8080/admin
```

Show:

- total conversations,
- top intents,
- containment rate,
- escalation rate,
- workflow success,
- low-confidence answers,
- unresolved questions,
- source freshness.

Talking point:

> “The dashboard turns the assistant into an operational improvement tool. Staff can see what learners ask, where the model is uncertain, and what should be improved next.”

## Minute 6: Go architecture walkthrough

Show repository structure:

```text
cmd/api
cmd/mock-banner
cmd/mock-payment
cmd/mock-crm
cmd/workflow-sim
cmd/ingest
cmd/eval
internal/orchestrator
internal/rag
internal/tools
internal/privacy
internal/audit
```

Talking point:

> “I used Go because this role needs reliable integrations and automation services. The model is only one part of the system; the Go services enforce timeouts, typed tool calls, redaction, audit logging, and safe fallback.”

## Minute 7: Evaluation and tests

Run:

```bash
go test ./...
go run ./cmd/eval -input data/eval-questions.jsonl
```

Show:

- intent accuracy,
- source recall,
- workflow decision accuracy,
- critical hallucination count,
- privacy redaction tests.

Talking point:

> “I treat this as a maintained AI product. The evaluation runner checks model quality, RAG grounding, workflow decisions, and privacy controls after changes.”

## Expected demo path

All demo records are synthetic. Student IDs, payment states, holds, workflow IDs, and CRM cases must come from local Go services or fixtures only.

| Step | Scenario | Input | Expected output | Observable proof |
|---:|---|---|---|---|
| 1 | Transcript answer | “How do I order my official transcript?” | Grounded answer with source citation | Chat response shows `transcript_request`, source ID/link, confidence, and no unsupported claims |
| 2 | Unpaid payment workflow | “I ordered my transcript but it has not been processed. My student ID is S100002.” | Payment status is unpaid and reminder workflow is accepted | Chat response shows `transcript_status`, `payment_status_checked`, `payment_reminder_triggered`, workflow ID, and no CRM handoff |
| 3 | Financial-hold escalation | “My transcript still is not moving. My student ID is S100003.” | Financial hold is detected and staff handoff is created | Chat response shows `transcript_status`, `financial_hold_detected`, CRM case ID, and Registrar/Student Accounts handoff |
| 4 | Urgent sentiment escalation | “This is really frustrating. I need this transcript for a job application.” | Urgent/negative sentiment creates a priority CRM case using existing conversation context | Chat response shows urgent sentiment, `crm_case_created`, priority flag, case ID, and privacy-aware summary |
| 5 | Dashboard evidence | Open `/admin` after steps 1-4 | Metrics reflect the demo events | Dashboard shows conversation count, top intents, workflow success, escalation count, and audit events |

## Demo acceptance matrix

Each row must be verifiable from the Go API response, local mock service logs, workflow simulator, CRM simulator, audit events, or dashboard. Source references are approved public/curated learner-service content; synthetic records are the only data used for student/payment/hold state.

| ID | Scenario | Synthetic input | Expected intent | Expected source | Expected action | Expected handoff behavior | Pass evidence |
|---|---|---|---|---|---|---|---|
| D01 | Transcript answer | “How do I order my official transcript?” | `transcript_request` | Transcript ordering source chunk, such as `oc-transcript-request-2005-onwards` | Return concise grounded answer with citation | No handoff; keep learner in chat | Response includes cited source, answer confidence, and zero critical unsupported claims |
| D02 | Unpaid payment workflow | `S100002` transcript status prompt | `transcript_status` | Transcript/payment guidance source chunk plus synthetic payment record `S100002` | Check mock Banner, check mock Payment, trigger `payment_reminder_triggered` | No CRM handoff unless workflow fails or confidence is low | Response includes unpaid status, workflow ID, idempotency key, and audit event |
| D03 | Financial-hold escalation | `S100003` transcript status prompt | `transcript_status` | Transcript/hold guidance source chunk plus synthetic Banner record `S100003` | Check mock Banner, detect `financial_hold`, create CRM case | Handoff to Registrar/Student Accounts queue with minimal summary and case ID | Response includes hold-safe wording, CRM case ID, queue/priority, and no payment reminder |
| D04 | Urgent sentiment escalation | Frustrated urgent follow-up in an unresolved transcript conversation | `escalation_request` or transcript follow-up intent with urgent sentiment | Prior transcript source and active synthetic conversation context | Classify sentiment as urgent/negative and create priority CRM case | Priority handoff to staff; assistant does not promise a deadline or outcome | Response includes priority flag, CRM case ID, redacted summary, and safe expectation-setting |
| D05 | Low-confidence source fallback | Transcript-adjacent question with no approved source | `unknown` or low-confidence transcript intent | No acceptable source above threshold | Do not answer from model memory; ask clarifying question or route to staff | Low-confidence handoff if learner needs account-specific help | Response shows no citation used for unsupported claim and logs low-confidence review item |

Golden-path pass condition: D01-D04 pass in order during the 5-7 minute demo, and each expected action has at least one observable output. D05 is a safety check used when discussing fallback behavior.

## Backup plan

If LLM or internet access fails during demo:

- use pre-seeded local RAG chunks,
- use deterministic classifier fallback,
- show mock API responses,
- run Go tests and evaluation report,
- explain how cloud AI is abstracted behind `internal/llm`.

## Closing line

> “This MVP shows that I can build more than a chatbot. I can build a Go-based, privacy-aware AI automation service that integrates with enterprise systems, reduces email-driven work, and gives staff measurable insight into learner support demand.”
