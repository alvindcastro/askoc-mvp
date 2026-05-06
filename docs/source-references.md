# Source References

## Purpose

This document lists recommended public content and technical references for the AskOC AI Concierge MVP. The MVP should use only approved public sources and synthetic data.

## Source allowlist boundary

The learner-service knowledge base for the MVP is limited to `data/seed-sources.json`.

Allowlisted knowledge sources must be:

- public web pages already approved in this repo,
- directly relevant to learner-service support,
- represented with source URL, title, department, risk level, freshness rule, and knowledge-domain metadata,
- reviewed before demos when the source covers deadlines, fees, platform transitions, admissions, or policy.

The MVP must not ingest:

- private portal pages,
- pages requiring login or learner-specific authentication,
- myOkanagan account pages or personal dashboards,
- Banner, CRM, LMS, payment, or student-record system data,
- private files, email, chat exports, ticket queues, or copied learner records,
- URLs that are not in `data/seed-sources.json`.

If a question falls outside the allowlisted knowledge domains, the assistant should provide a safe fallback or create a mock staff handoff instead of guessing.

## Approved learner-service seed sources

The P0-T03 allowlist currently contains only the public Okanagan College URLs below. These URLs were selected from this repo's existing documentation and spot-checked as accessible public web pages on 2026-05-06.

## Okanagan College public content examples

### Office of the Registrar

Useful for Registrar scope and learner-service context.

URL:

```text
https://www.okanagancollege.ca/office-of-the-registrar
```

### Transcript Request Guidance

Useful for the primary transcript workflow.

URL:

```text
https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards
```

### Online Resources / myOkanagan / LMS Guidance

Useful for portal and LMS support questions.

URL:

```text
https://www.okanagancollege.ca/student-handbook/online-resources
```

### LMS Information

Useful for Moodle/Brightspace transition questions and digital learning support context.

URL:

```text
https://www.okanagancollege.ca/teaching-and-learning-innovations/lms
```

### Information Security Guide

Useful for privacy/security framing and responsible handling of personal information.

URL:

```text
https://www.okanagancollege.ca/it-services/it-security/information-security-guide
```

## Go technical references

The references in this section are implementation references only. They are not learner-service RAG sources and should not be used to answer learner policy/procedure questions.

### Go documentation

Useful for language, modules, testing, and standard tooling.

URL:

```text
https://go.dev/doc/
```

### Go web applications tutorial

Useful for simple server-rendered web chat using Go HTTP handlers and templates.

URL:

```text
https://go.dev/doc/articles/wiki/
```

### Go RESTful API with Gin tutorial

Useful if choosing Gin instead of standard `net/http` or `chi`.

URL:

```text
https://go.dev/doc/tutorial/web-service-gin
```

### Effective Go

Useful for idiomatic Go style.

URL:

```text
https://go.dev/doc/effective_go
```

## Azure and automation references

The references in this section are implementation references only. They are not learner-service RAG sources and should not be used to answer learner policy/procedure questions.

### Azure for Go developers

Useful for Go SDKs, Azure deployment, identity, storage, and cloud service integration.

URL:

```text
https://learn.microsoft.com/en-us/azure/developer/go/
```

### Azure AI for Go developers

Useful for building AI applications with Go and Azure AI services.

URL:

```text
https://learn.microsoft.com/en-us/azure/developer/go/azure-ai-for-go-developers
```

### Azure OpenAI REST API reference

Useful for implementing a Go REST client without depending on a specific SDK.

URL:

```text
https://learn.microsoft.com/en-us/azure/foundry/openai/reference
```

### Azure AI Search — Retrieval-Augmented Generation

Useful for explaining a cloud-based RAG architecture.

URL:

```text
https://learn.microsoft.com/en-us/azure/search/retrieval-augmented-generation-overview
```

### Azure AI Language

Useful for intent classification, sentiment analysis, and language understanding.

URL:

```text
https://learn.microsoft.com/en-us/azure/ai-services/language-service/
```

### Power Automate

Useful for workflow automation/RPA framing.

URL:

```text
https://learn.microsoft.com/en-us/power-automate/
```

### Power Automate HTTP request trigger and OAuth authentication

Useful for webhook-triggered workflows from Go services.

URL:

```text
https://learn.microsoft.com/en-us/power-automate/oauth-authentication
```

## Suggested knowledge base seed list

Start with 10–20 public pages covering:

- Transcript requests.
- Tuition and fee payment.
- Important dates.
- Refunds and withdrawals.
- Registration.
- myOkanagan.
- LMS access.
- Admissions/application status.
- International student support.
- Advising.
- IT support.
- Privacy/security guidance.

Do not add a suggested topic to `data/seed-sources.json` until a public source URL is approved and listed there.

## Ingestion metadata template

Each ingested document should store:

```json
{
  "source_url": "https://www.okanagancollege.ca/example-page",
  "title": "Example Page Title",
  "department": "Registrar",
  "content_type": "public_web_page",
  "retrieved_at": "2026-05-06",
  "last_reviewed_by": "applicant_demo",
  "risk_level": "medium"
}
```

For the P0 allowlist, seed records also include `allowlisted`, `private_portal`, `knowledge_domains`, `requires_freshness_check`, `stale_after_days`, and `verification_basis`.

## Source freshness rule

For demo purposes, any source involving fees, deadlines, platform transitions, admissions, or policy should be treated as time-sensitive.

Recommended rule:

- Re-ingest and review before each demo.
- Show retrieval timestamp in admin view.
- Escalate if a question depends on a missing or stale source.

If a source is stale, missing, or conflicts with another source, the assistant must not present the answer as authoritative. It should say that verified information is unavailable in the demo knowledge base and offer staff follow-up.
