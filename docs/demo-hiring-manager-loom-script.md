# Hiring Manager Loom Script

Use this script for a concise 5-7 minute Loom aimed at a hiring manager. Keep the screen on the working product: README, chat UI, admin dashboard, and terminal proof. Use only synthetic demo data.

## Script

Hi, I am Alvin. Thanks for taking the time to review this.

I built AskOC AI Concierge to show how I approach AI automation work: find a repeated service workflow, make the data boundary clear, automate the low-risk steps, and leave judgment with staff.

The problem is simple. Learners ask many of the same questions about transcripts, payments, holds, portals, and deadlines. Staff answer those questions through email and manual triage. That creates delay for learners and repeated work for the service team.

This MVP focuses on one narrow path: transcript support. A learner asks how to order an official transcript, then asks why a request has not moved. The system answers from approved source chunks, checks synthetic Banner and payment records, triggers a payment reminder when the record is unpaid, and creates a mock CRM case when the request needs staff review.

The scope is narrow on purpose. A narrow slice is easier to test, easier to explain, and easier to improve.

Here is the flow.

First, the learner asks, "How do I order my official transcript?"

The chat API classifies the question as a transcript request. The retriever searches local approved source chunks. The answer includes source evidence, confidence, risk, and freshness metadata. The assistant does not answer from memory when it lacks support. It either cites the approved source or falls back safely.

This matters because a useful AI service must show why the answer can be trusted. The source is part of the product, not an afterthought.

Second, the learner asks for request status with a synthetic student ID, such as `S100002`.

The Go API calls typed mock Banner and payment services. The record shows an unpaid transcript payment. The orchestrator triggers an idempotent payment-reminder workflow. The response shows the action trace: Banner checked, payment checked, reminder triggered, and workflow ID returned.

No real student system runs here. Banner, payment, CRM, LMS, and workflow services are local mocks with synthetic records. That keeps the demo safe while still showing the integration pattern.

Third, the learner asks about `S100003`, a synthetic record with a financial hold.

The system detects the hold and does not trigger a self-service payment reminder. It creates a mock Registrar and Student Accounts CRM case. The summary is redacted. The learner receives a clear handoff message and a synthetic case ID.

This is the main product judgment in the demo: automate the routine step, escalate the risky step. The system supports staff; it does not approve, deny, waive, or promise an outcome.

Fourth, the learner says, "This is really frustrating. I need this transcript for a job application."

The system classifies urgent negative sentiment and creates a priority staff handoff. Sentiment alone does not make a final decision. It changes routing priority when the context justifies it.

Now I will open the admin dashboard.

The dashboard shows aggregate evidence: top intents, containment, escalations, workflow outcomes, low-confidence items, stale-source warnings, and review queue items. It does not show raw private messages. The audit path stores redacted metadata and hashed workflow identifiers where appropriate.

That is the privacy boundary I wanted to make visible. The demo should prove useful automation without normalizing careless data handling.

I made several engineering choices to keep the project practical.

I used Go because this role needs reliable services, integrations, and automation boundaries. The code separates the API, orchestrator, classifier, retriever, tool clients, workflow client, audit service, and mock systems. Each piece has a clear job.

I kept deterministic mode as the default. The demo can run without a live LLM, without real Banner, without real CRM, without a live payment processor, and without Power Automate. Optional provider and webhook paths sit behind typed interfaces, so they can be tested without controlling the whole product.

The test strategy follows the same idea: prove the behavior that carries risk. The suite covers classification, source grounding, workflow idempotency, CRM handoff, redaction, audit behavior, dashboard metrics, Docker artifacts, smoke checks, and JSONL evaluation cases.

At the end of the walkthrough I can run `make smoke` to prove the local Docker stack. I can run `make eval` to prove the responsible-AI gate. I can run `go test ./...` to prove the Go packages.

If I continued this project, I would improve three areas.

First, I would expand the review queue so staff can resolve low-confidence items and feed better examples back into evaluation.

Second, I would add stronger role-based access around admin and review screens before any real deployment.

Third, I would broaden the workflow catalog only after each new workflow has clear acceptance criteria, source boundaries, audit rules, and escalation rules.

This project reflects how I like to work. I define the smallest useful product slice. I make system boundaries explicit. I test the behavior that can harm users if it fails. I avoid cleverness when plain code will do.

The result is more than a chatbot. It is a Go-based, privacy-aware automation service that answers common questions, coordinates workflow actions, and gives staff measurable evidence.

Thanks for watching.

## Recording Cues

| Time | Screen | Show |
|---:|---|---|
| 0:00-0:45 | README | Problem, scope, synthetic-data boundary |
| 0:45-1:45 | `/chat` | Grounded transcript answer with source metadata |
| 1:45-3:00 | `/chat` | `S100002` unpaid transcript workflow |
| 3:00-4:10 | `/chat` | `S100003` financial-hold CRM handoff |
| 4:10-5:10 | `/admin` | Redacted metrics, audit evidence, review queue |
| 5:10-6:30 | Terminal | `make smoke`, `make eval`, or `go test ./...` evidence |
