# Privacy Impact Lite

## Purpose

This document captures privacy-by-design and responsible-AI controls for the AskOC AI Concierge MVP. It is not a formal institutional Privacy Impact Assessment, but it shows the applicant understands privacy, security, and ethical AI expectations in a learner-service environment.

## MVP privacy position

The MVP must use **synthetic data only**.

The demo data boundary is:

- **Synthetic learners only**: learner names, IDs, programs, statuses, holds, and support history are invented for the demo.
- **Synthetic IDs only**: demo student identifiers use the `S10000X` pattern and do not represent OC student numbers.
- **Mock payments only**: payment status, balances, transaction IDs, and reminder workflow results are fake and must not be treated as financial records.
- **Mock CRM cases only**: escalation IDs, case summaries, owners, and statuses are generated for the demo and do not represent real institutional cases.
- **Public source grounding only**: policy/procedure answers may use approved public web pages or manually curated public snippets, never private portal content.

Do not use:

- real student records,
- real payment information,
- real authentication tokens,
- real CRM cases,
- private portal content,
- scraped content behind login.

Use only public content and synthetic records.

Every demo fixture should be visibly fake through fields such as `synthetic: true`, `demo_record: true`, `data_notice`, or a `SYNTH-`/`MOCK-` identifier prefix. Screenshots, logs, dashboard rows, and case examples must preserve those fake-data markers.

## Personal information inventory

| Data element | MVP use | Real production equivalent | Risk |
|---|---|---|---|
| Synthetic student ID | Demo transaction lookup | Student number | Medium |
| Preferred name | Demo personalization | Student profile | Medium |
| Transcript status | Mock status check | Registrar record | Medium/high |
| Payment status | Mock payment flag only | Financial record | High |
| Conversation text | Learner request | Support conversation | Medium/high |
| Intent/sentiment | Routing | Analytics/routing metadata | Medium |
| CRM summary | Mock handoff | Case management record | Medium/high |
| Source links | Grounding | Public web content | Low |

## Go implementation controls

| Control | Go package or pattern |
|---|---|
| Redaction before logging | `internal/privacy` |
| Trace IDs | `internal/middleware` |
| Minimal audit events | `internal/audit` |
| Typed tool calls | `internal/tools` |
| Context timeouts | `context.WithTimeout` for downstream calls |
| Safe errors | Handler-level error mapping |
| Secrets outside repo | Environment variables and secret store |
| Auth stub | Bearer token middleware for local demo |
| Rate limiting | Middleware or reverse proxy |
| Prompt-injection handling | Treat retrieved text as untrusted content |

## Redaction requirements

Before writing logs, audit events, or CRM summaries, redact:

- email addresses,
- phone numbers,
- passwords or password-like phrases,
- payment card-like numbers,
- government ID-like patterns,
- access tokens,
- long free-form sensitive text where not needed.

Example redacted audit message:

```json
{
  "trace_id": "trace_01JABC456",
  "conversation_id": "conv_01JABC123",
  "role": "learner",
  "content_redacted": "I paid yesterday and my email is [REDACTED_EMAIL].",
  "created_at": "2026-05-06T12:00:00Z"
}
```

## Data minimization

The MVP should store:

- conversation ID,
- trace ID,
- redacted message text,
- predicted intent,
- confidence,
- source IDs,
- workflow/case IDs,
- status summaries.

The MVP should not store:

- raw passwords,
- payment card data,
- real student records,
- full unredacted conversation transcripts,
- private files or portal data.

## Demo fixture rules

`data/synthetic-students.json` is the only approved learner fixture for the demo boundary. It must contain invented records only.

Allowed fixture fields:

- synthetic student ID,
- fake preferred name,
- broad fake program label,
- transcript workflow status,
- mock payment status and synthetic transaction ID,
- mock CRM case status and synthetic case ID,
- mock LMS account status and demo course-access status,
- non-sensitive demo notes.

Not allowed in fixtures:

- real learner names,
- real student numbers or government IDs,
- real email addresses or phone numbers,
- real payment card, bank, tax, or account data,
- real Banner, CRM, LMS, or portal identifiers,
- LMS course content, grades, submissions, or activity records,
- private portal URLs or scraped content.

If a future demo needs a new scenario, add a new clearly fake record rather than copying or anonymizing a real learner.

## Logging rules

| Log type | Allowed | Not allowed |
|---|---|---|
| Application logs | Trace ID, status, duration, service name | Raw user secrets |
| Audit events | Redacted content, source IDs, tool summaries | Full private records |
| Tool-call logs | Tool name, synthetic ID, status, duration | Full downstream payloads with sensitive values |
| Dashboard metrics | Aggregates and review queue | Direct personal identifiers unless needed for demo |

## Prompt and retrieval safety

Policy/procedure answers must be source-grounded.

Knowledge sources are limited to the public allowlist in `data/seed-sources.json`. Private portal scraping, authenticated pages, personal account pages, and unofficial cached copies are out of scope for the MVP.

The assistant should not answer with certainty when:

- retrieved source confidence is low,
- source is stale,
- source conflicts with another source,
- question involves deadlines, fees, eligibility, holds, immigration, disability accommodation, legal matters, or financial decisions.

P6 enforces this boundary in `internal/rag`, `internal/classifier`, `internal/llm`, and `internal/orchestrator`: ingestion rejects unallowlisted URLs, local retrieval returns approved chunk citations only, strict JSON classification is validated before use, low-confidence classification or retrieval falls back, ungrounded LLM answers are rejected, and stale or high-risk sources require staff confirmation instead of an authoritative answer.

Safe fallback:

```text
I do not have enough verified information to answer that confidently. I can create a case for staff follow-up or point you to the relevant office.
```

## Tool safety

The LLM should not directly execute arbitrary tools. The Go orchestrator should decide tool calls using typed logic.

Allowed MVP tool actions:

- check synthetic student profile,
- check synthetic transcript status,
- check synthetic payment status,
- trigger synthetic payment reminder,
- create mock CRM case,
- write audit event.

Disallowed MVP tool actions:

- change real student record,
- process real payment,
- send real official transcript,
- modify real CRM/Banner/LMS data,
- send sensitive details to unapproved systems.

## Responsible AI controls

| Risk | Control |
|---|---|
| Hallucinated policy answer | RAG required for policy/procedure answers |
| Overconfident answer | Confidence threshold and fallback |
| Bias in sentiment routing | Review false positives/negatives; do not use sentiment alone for adverse action |
| Privacy leak | Redaction and minimal summaries |
| Prompt injection | Do not follow instructions inside retrieved pages or user messages that override system policy |
| Stale content | Show indexed date and review stale sources |
| Unclear source | Require source links for policy answer |

## Retention recommendation for MVP

| Data | Demo retention |
|---|---:|
| Synthetic records | Keep in repo |
| Redacted conversation logs | 7–30 days for demo |
| Audit events | 30–90 days for demo |
| Evaluation outputs | Keep as non-sensitive artifacts |
| Secrets/tokens | Never commit |

## Production notes

A production version should include:

- formal privacy review,
- institutional security review,
- role-based access control,
- approved data residency and retention rules,
- SSO integration,
- consent and notice language,
- human-in-the-loop review for sensitive workflows,
- data processing agreements for AI vendors,
- monitoring for model drift and unsafe outputs.

## Applicant talking point

> “I designed the Go services so the LLM cannot directly mutate records. The model classifies and drafts grounded responses, while the Go orchestrator enforces tool allowlists, source requirements, timeouts, redaction, and audit logging.”
