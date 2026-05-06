# Brainstorm: Go-Based AskOC AI Concierge MVP

This brainstorm captures the best implementation direction for a portfolio MVP targeting an AI/Automation Solutions Developer role. The recommended path is intentionally practical: build a reliable deterministic workflow first, then add AI/RAG behind tested interfaces.

## Recommended MVP direction

Build a **Go-based learner-service automation concierge** that demonstrates:

- source-grounded transcript and student-service answers,
- synthetic Banner/payment/CRM/LMS integrations,
- transcript payment/status automation,
- urgent-sentiment escalation,
- privacy-aware audit logs,
- operational dashboard metrics,
- strict TDD and evaluation gates.

The MVP should not be positioned as “just a chatbot.” It should be positioned as a small learner-service automation platform.

## Best demo workflow

Use transcript support as the core demo because it naturally touches knowledge retrieval, transaction status, payment confirmation, workflow automation, escalation, and auditability.

Golden flow:

1. Learner asks how to order a transcript.
2. Assistant answers from approved public/curated sources.
3. Learner asks why the transcript has not been processed and provides synthetic ID `S100002`.
4. Go orchestrator checks mock Banner and mock payment APIs.
5. Payment status is unpaid.
6. Workflow simulator or Power Automate webhook sends a payment reminder.
7. Learner expresses frustration and urgency.
8. Sentiment/urgency classification triggers priority CRM case creation.
9. Dashboard shows the intent, workflow action, escalation, and trace ID.

## Strong architectural decisions

- [ ] Use Go for API gateway, orchestrator, mock services, workflow simulator, ingestion, and evaluation runner.
- [ ] Keep UI simple: Go templates + small JavaScript/HTMX is enough for the interview demo.
- [ ] Build deterministic transcript workflow before live LLM integration.
- [ ] Put all external dependencies behind interfaces.
- [ ] Use `httptest.Server` for client tests.
- [ ] Use fakes for orchestrator tests.
- [ ] Make RAG retrieval swappable: local retriever first, Azure AI Search or pgvector later.
- [ ] Use structured JSON for LLM classification and validate it before actions.
- [ ] Require source citations or safe fallback for policy/procedure answers.
- [ ] Audit every tool call and workflow action with redacted payloads.

## MVP feature options

| Feature idea | Include now? | Reason |
|---|---:|---|
| Transcript request answer | Yes | Clear Tier 0 use case |
| Transcript payment/status workflow | Yes | Strong Tier 1 automation story |
| Mock Banner API | Yes | Shows enterprise integration thinking |
| Mock payment API | Yes | Supports automation workflow |
| Mock CRM API | Yes | Shows human handoff and routing |
| Mock LMS API | Yes, small | Shows broader learner digital experience |
| Power Automate webhook | Optional after simulator | Good role fit, but local simulator keeps demo reliable |
| Voice assistant | No | Nice-to-have, not needed for MVP |
| Fine-tuned model | No | RAG + guardrails is more appropriate for MVP |
| Real OC login | No | Out of scope and privacy risk |
| Real payment processing | No | Out of scope and unnecessary |
| Private portal scraping | No | Privacy/security risk |
| Kubernetes | Stretch | Good polish but not needed for MVP proof |

## Phase strategy

The safest implementation order is:

- [ ] Freeze product story and privacy boundaries.
- [ ] Create tested Go foundation.
- [ ] Build chat API and simple UI.
- [ ] Build synthetic enterprise APIs.
- [ ] Build deterministic orchestrator workflow.
- [ ] Add privacy redaction and audit before AI calls.
- [ ] Add RAG ingestion/retrieval.
- [ ] Add LLM gateway and strict JSON classification.
- [ ] Add workflow simulator and optional Power Automate client.
- [ ] Add evaluation runner.
- [ ] Add Docker/CI/smoke test.
- [ ] Polish portfolio docs and demo script.

The detailed task board is in [phases-and-tasks.md](phases-and-tasks.md).

## TDD brainstorming notes

Strict TDD is especially valuable in this MVP because AI systems can become vague and hard to trust. The code should prove behavior through deterministic tests wherever possible.

TDD priorities:

- [ ] Test orchestration decisions before adding model calls.
- [ ] Test fallback behavior more heavily than happy paths.
- [ ] Test redaction before logging or storing conversation text.
- [ ] Test tool-trigger thresholds before workflow automation.
- [ ] Test malformed LLM output.
- [ ] Test stale/low-confidence retrieval.
- [ ] Test duplicate workflow prevention.
- [ ] Test critical evaluation failures as build failures.

## Interview positioning

Use language like:

> “I designed the MVP as an automation product, not just a chatbot. The Go backend uses interfaces, strict TDD, deterministic mock integrations, RAG with citations, workflow idempotency, audit logs, and evaluation gates so the system is maintainable and safe.”

## Risks to call out honestly

- [ ] Public-source content must be refreshed before a real deployment.
- [ ] Real Banner/CRM/LMS integrations would require institutional API access, security review, and privacy impact assessment.
- [ ] LLM outputs must be validated and constrained before any action is taken.
- [ ] The MVP uses synthetic data and does not represent production-grade authentication or authorization.
- [ ] Evaluation metrics are demo-scale and would need expansion with real stakeholder-reviewed test cases.

## Recommended first coding task

Start with **P1-T03 — Implement health and readiness endpoints** or **P2-T02 — Implement `POST /api/v1/chat` handler** because they make TDD visible quickly.

For example, the first visible TDD cycle could be:

1. Write `TestHealthzReturnsOKJSON` using `httptest`.
2. Run it and see it fail because the handler does not exist.
3. Add the smallest handler.
4. Run the test until green.
5. Add method-not-allowed test.
6. Refactor response helper.
7. Run `go test ./...`.

This creates a simple but credible story: the project is built with disciplined engineering from the first commit.
