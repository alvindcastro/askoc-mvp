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

| Step | Input | Expected output |
|---|---|---|
| 1 | “How do I order my official transcript?” | Grounded answer with source |
| 2 | `S100002` transcript status | Payment unpaid |
| 3 | Workflow trigger | Payment reminder accepted |
| 4 | Frustrated/urgent message | Priority CRM case |
| 5 | Dashboard | Metrics updated |

## Backup plan

If LLM or internet access fails during demo:

- use pre-seeded local RAG chunks,
- use deterministic classifier fallback,
- show mock API responses,
- run Go tests and evaluation report,
- explain how cloud AI is abstracted behind `internal/llm`.

## Closing line

> “This MVP shows that I can build more than a chatbot. I can build a Go-based, privacy-aware AI automation service that integrates with enterprise systems, reduces email-driven work, and gives staff measurable insight into learner support demand.”
