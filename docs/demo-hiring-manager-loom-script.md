# Hiring Manager Loom Script

Use this for a 5-7 minute Loom aimed at a hiring manager or engineering lead. Keep the screen on the working product: README, chat UI, admin dashboard, and terminal proof. Use only synthetic demo data.

## Setup

Start the local stack before recording:

```bash
make compose-up
```

Open:

```text
http://localhost:8080/chat
http://localhost:8080/admin
```

Use this local admin token:

```text
demo-admin-token
```

Keep a terminal ready with one proof command or a recent output window:

```bash
make smoke
make eval
go test ./...
```

## Spoken Script

Hi, I am Alvin. Thanks for taking a few minutes to review this.

This is AskOC AI Concierge, a Go-based MVP for a learner-service automation workflow. I built it to show how I approach AI automation work: start with one repeated service problem, keep the data boundary clear, automate the low-risk steps, and leave judgment with staff.

The problem is familiar in student services. Learners ask the same questions about transcripts, payments, holds, portals, and deadlines. Staff end up doing a lot of manual triage through email. That creates delay for learners and repeated work for the team.

For this MVP, I kept the slice narrow on purpose. The system focuses on transcript support. A learner asks how to order an official transcript, then asks why a request has not moved. The assistant answers from approved source chunks, checks synthetic Banner and payment records, triggers a local reminder workflow for unpaid payment, and routes financial holds to a mock CRM handoff.

The scope is small, but it touches the parts that matter: retrieval, classification, typed integrations, workflow actions, audit evidence, privacy boundaries, and fallback behavior.

First, I will show the Tier 0 answer.

In the chat, I ask: "How do I order my official transcript?"

The API classifies this as a transcript request. The retriever searches local approved source chunks, and the answer comes back with source metadata, confidence, risk, and freshness signals. If the source support is weak, the assistant is designed to fall back instead of sounding confident without evidence.

That is an important product choice. For a staff-facing workflow, the answer is not enough. The system also has to show why the answer can be trusted.

Next, I will show the transcript-status workflow with a synthetic student ID.

I ask: "I ordered my transcript but it has not been processed. My student ID is S100002."

This is not a real student record. The demo uses local synthetic data. The Go API calls typed mock Banner and payment services. Banner confirms the transcript state, payment shows the transcript fee is unpaid, and the orchestrator triggers a local payment-reminder workflow.

The response shows the action trace: intent classified, Banner checked, payment checked, reminder triggered, and a synthetic workflow ID returned.

This is the kind of automation boundary I wanted to show. The assistant is not just chatting. It is coordinating a workflow through typed services, and the action is limited to a safe self-service reminder.

Now I will show the handoff path.

I ask: "My transcript request has a financial hold and is not moving. My student ID is S100003."

This synthetic record has a financial hold, so the system does not push the learner through the same self-service path. It creates a mock Registrar and Student Accounts CRM case with a minimal redacted summary. The learner gets a clear handoff message and a synthetic case ID.

That is the main product judgment in the demo: automate routine follow-up, escalate risky work with context. The system does not approve transcripts, waive holds, promise outcomes, or hide a decision inside a friendly response.

I also included sentiment routing. If the learner says, "This is really frustrating. I need this transcript for a job application," the system can classify the message as urgent or negative and route the case with higher priority. Sentiment does not make the final decision. It only helps route the work when the context supports escalation.

Now I will open the admin dashboard.

The dashboard shows aggregate and redacted evidence: total conversations, containment, escalations, workflow outcomes, top intents, low-confidence items, stale-source warnings, and review queue items. It is designed to help staff see what the assistant did without exposing raw private details.

That privacy boundary is deliberate. Banner, payment, CRM, LMS, workflow IDs, and student IDs are all synthetic. Audit events are redacted. The CRM handoff summary is minimal. The demo proves the integration pattern without normalizing careless handling of student records.

On the engineering side, this is intentionally Go-first. The API, orchestrator, classifier, retriever, mock enterprise services, workflow simulator, audit service, dashboard, and evaluation runner are all separated behind clear boundaries.

The demo also runs without a live LLM by default. Deterministic mode makes the core behavior repeatable. Optional provider and webhook paths sit behind typed interfaces, so the product shape can be tested without depending on a live enterprise system.

To prove the behavior, I keep three checks close to the demo.

`make smoke` starts the local Docker stack and checks the unpaid transcript path plus the financial-hold CRM handoff.

`make eval` runs JSONL cases for intent, source grounding, workflow actions, escalation, and safety behavior.

`go test ./...` proves the Go packages.

Those commands are not decoration. They force the MVP to behave like software. If source grounding breaks, if a workflow action changes, or if a safety case regresses, the project should show it.

If I continued this project, I would expand the staff review queue, add stronger role-based access for admin views, and only add new workflows after each one had clear source rules, audit rules, and escalation rules.

The point of this MVP is not that transcripts are the only problem worth solving. The point is the engineering pattern: useful AI automation should be bounded, testable, privacy-aware, and honest about where human judgment belongs.

Thanks for watching.

## Recording Cues

| Time | Screen | Show |
|---:|---|---|
| 0:00-0:45 | README | Problem, narrow transcript scope, synthetic-data boundary |
| 0:45-1:45 | `/chat` | Transcript-order answer with source metadata |
| 1:45-3:00 | `/chat` | `S100002` unpaid transcript workflow and action trace |
| 3:00-4:15 | `/chat` | `S100003` financial-hold CRM handoff |
| 4:15-5:20 | `/admin` | Aggregate metrics, redacted audit evidence, review queue |
| 5:20-6:30 | Terminal | `make smoke`, `make eval`, or `go test ./...` proof |
| 6:30-7:00 | README or terminal | Close with engineering pattern and next steps |

## Demo Prompts

```text
How do I order my official transcript?
```

```text
I ordered my transcript but it has not been processed. My student ID is S100002.
```

```text
My transcript request has a financial hold and is not moving. My student ID is S100003.
```

```text
This is really frustrating. I need this transcript for a job application.
```

## Keep Visible

- Synthetic demo mode.
- Source metadata for the grounded transcript answer.
- `payment_reminder_triggered` for `S100002`.
- `crm_case_created` for `S100003`.
- Redacted dashboard metrics.
- One terminal proof command or recent passing output.

## Avoid Saying

- That the MVP uses real student data.
- That it connects to live Banner, CRM, LMS, payment, or production Power Automate by default.
- That the assistant approves transcripts, waives holds, or promises outcomes.
- That the demo is production-ready.
