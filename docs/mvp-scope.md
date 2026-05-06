# MVP Scope

## Product goal

Build a Go-based applicant portfolio MVP that proves the ability to design, develop, deploy, and maintain AI-powered learner-service automation.

The MVP should demonstrate:

- conversational AI for Tier 0 and Tier 1 learner support,
- RAG over approved public sources,
- real-time decision workflows,
- mock enterprise integrations,
- workflow automation/RPA integration,
- privacy-aware logging and escalation,
- measurable operations dashboard.

## Primary use case

The first use case is transcript request support:

> A learner asks how to order a transcript and why their request has not been processed. The assistant answers from approved sources, checks synthetic student/payment records, triggers a payment reminder workflow when appropriate, and escalates complex or frustrated cases.

## Why this use case

Transcript support is strong for an MVP because it combines:

- a clear public policy/procedure answer,
- a common student-service inquiry,
- a transaction status check,
- a payment requirement,
- a natural escalation path,
- a realistic integration pattern with Registrar, Finance, CRM, and notification systems.

## MVP personas

| Persona | Need |
|---|---|
| Learner | Fast, accurate answer without waiting for email |
| Registrar staff | Fewer repetitive inquiries and better case summaries |
| Finance staff | Clear routing for payment or hold-related issues |
| IT/digital learner experience team | Safe integration pattern and measurable automation |
| Manager | Dashboard showing adoption, containment, and risk |

## In-scope capabilities

| Area | In scope |
|---|---|
| Go backend | `net/http` or `chi` API service with typed request/response models |
| Chat UI | Basic web chat served by Go templates or separate frontend |
| RAG | Go ingestion command, chunk metadata, retrieval, source-grounded responses |
| Intent detection | Transcript, fees, portal, LMS, registration, refunds, application status, dates, escalation |
| Sentiment detection | Negative/urgent signal for priority escalation |
| Mock integrations | Banner-style student API, payment API, CRM API, LMS API, notification API |
| Workflow automation | Power Automate webhook client and local Go workflow simulator |
| Audit | Redacted action logs, tool-call logs, trace IDs |
| Dashboard | Operational metrics and review queue |
| Evaluation | JSONL test set and Go evaluation runner |

## Out-of-scope items

| Area | Out of scope for MVP |
|---|---|
| Real student data | Use synthetic records only |
| Real Banner access | Mock APIs only |
| Real payments | Simulated status and reminders only |
| Real OC authentication | Mock bearer token or local demo auth |
| Private portal scraping | Public content only |
| Production model fine-tuning | Use RAG and structured output first |
| Legal/policy guarantee | Escalate uncertain policy decisions |

## Supported intents

| Intent | Example prompt | MVP behavior |
|---|---|---|
| `transcript_request` | “How do I order my transcript?” | Grounded answer with source |
| `transcript_status` | “Why has my transcript not been processed?” | Check synthetic status and payment |
| `fee_payment` | “I paid but still see a balance.” | Route to payment workflow/escalation |
| `myokanagan_login` | “I cannot log into myOkanagan.” | Provide source-grounded support guidance |
| `lms_access` | “Where is my online course?” | Answer from LMS/online resources source |
| `registration_help` | “How do I register?” | Provide source-grounded next steps |
| `refund_request` | “Can I get a refund?” | Answer carefully and escalate deadline-sensitive cases |
| `application_status` | “Has my application been received?” | Mock status flow or handoff |
| `key_dates` | “When is the withdrawal deadline?” | Retrieve source; avoid guessing |
| `human_handoff` | “I need a person.” | Create mock CRM case |
| `unknown` | “Can you write my essay?” | Redirect or decline unrelated request |

## Golden demo path

1. Ask: “How do I order my official transcript?”
2. Assistant answers with source-grounded steps.
3. Ask: “I ordered it but it has not been processed. My student ID is S100002.”
4. Go orchestrator calls mock Banner and payment APIs.
5. Payment API returns `unpaid`.
6. Go workflow client triggers Power Automate or local workflow simulator.
7. Assistant says a payment reminder was sent.
8. Learner says: “This is really frustrating. I need it for a job application.”
9. Sentiment is urgent/negative.
10. Go orchestrator creates priority mock CRM case.
11. Dashboard shows the interaction, intent, workflow event, and escalation.

## Demo success criteria

| Criterion | Target |
|---|---:|
| Curated Tier 0 questions answered correctly | 20+ |
| Synthetic transcript scenarios working | 4/4 |
| Source-grounded policy answers | 100% |
| Critical unsupported claims | 0 |
| Workflow happy path | 100% |
| Redacted logs | 100% for configured PII patterns |
| Local demo startup | One command with Docker Compose |

## Applicant portfolio angle

The project should show more than chatbot skills. It should show that the applicant can build an operational automation product:

- clean Go service architecture,
- reliable API integrations,
- responsible AI design,
- workflow automation,
- privacy awareness,
- measurable continuous improvement.
