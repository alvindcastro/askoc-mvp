# Test Plan


## Strict TDD requirement

All Go code tasks are governed by [TDD Policy](tdd-policy.md). For each code task, create failing tests before production code, verify the red state, implement the smallest passing change, run the narrow test, then run `go test ./...`. Do not mark tasks complete in [Phases and Tickable Tasks](phases-and-tasks.md) without test evidence.


## Purpose

This test plan verifies that AskOC AI Concierge can answer learner-service questions, automate the transcript/payment workflow, escalate appropriately, protect privacy, and run repeatably in a Go-based demo environment. The P11 release surface covers deterministic fallback classification, optional OpenAI-compatible LLM gateway behavior behind strict JSON parsing, local approved-source RAG retrieval, transcript/payment orchestration, in-process workflow idempotency, the standalone workflow simulator, optional Power Automate-compatible webhook retries/signature headers, mock CRM handoff, source fallback, source-only LLM answer guardrails, safe action traces, shared redaction, redacted audit storage, protected admin metrics, dashboard rendering, audit retention/export/reset controls, a JSONL evaluation runner, JSON/Markdown reports, critical safety gates, unresolved eval review items, Docker Compose packaging, CI, env safety, smoke verification, portfolio diagrams, a timed demo script, and synthetic-only screenshot placeholders.

## Test objectives

1. Validate grounded answers for common Tier 0 questions.
2. Validate intent classification across learner-service categories.
3. Validate transcript status and payment workflow logic.
4. Validate sentiment-based and low-confidence escalation.
5. Validate privacy controls, redaction, and audit logging.
6. Validate dashboard metrics.
7. Validate Go services under demo-level load.
8. Validate Docker, CI, env safety, and smoke-test developer workflows.

## Test environment

| Component | Test setup |
|---|---|
| Frontend | Go server-rendered web chat or optional React UI |
| Backend | Go `cmd/api` local container |
| Knowledge base | P5 local JSON chunks from approved public sources in `data/rag-chunks.json` |
| LLM gateway | P6 defaults to deterministic `stub`; optional `openai-compatible` mode is tested with fakes/`httptest` and must not call live APIs in automated tests |
| Mock Banner API | Go service with synthetic student records |
| Mock Payment API | Go service with synthetic payment records |
| Mock CRM API | Go service with synthetic case creation |
| Automation workflow | P8 in-process idempotent workflow client by default; `cmd/workflow-sim` or Power Automate-compatible webhook client when `ASKOC_WORKFLOW_URL` is configured |
| Dashboard | Go admin dashboard at `/admin` reading the in-memory audit event store through protected admin APIs |
| Evaluation | P9 `cmd/eval` uses deterministic in-process fakes by default or a live local `/api/v1/chat` endpoint with `-base-url` |
| Local stack | Docker Compose runs API, mock Banner, mock payment, mock CRM, mock LMS, and workflow simulator containers |
| Developer testing guide | [Developer Testing Guide](developer-guide.md) documents local stack startup, alternate ports, manual checks, troubleshooting, and useful operating notes |
| CI | GitHub Actions runs offline `go test ./...`, `go vet ./...`, and `make eval` with `ASKOC_PROVIDER=stub` |

## Go test commands

```bash
go test ./...
go test ./internal/privacy ./internal/audit ./internal/handlers
go test ./internal/llm ./internal/classifier ./internal/orchestrator
go test ./internal/classifier ./internal/workflow ./internal/orchestrator
go test ./internal/workflow -run 'TestSimulator|TestPowerAutomate|TestIdempotency'
go test ./cmd/workflow-sim ./cmd/api ./internal/config ./internal/orchestrator
go test ./internal/eval ./cmd/eval
go test ./internal/build -run TestP10
go test ./internal/rag ./internal/orchestrator
go test ./internal/domain ./internal/validation ./internal/handlers ./internal/session
go test -race ./internal/session
go test ./internal/orchestrator -run 'TestTranscriptStatus|TestUrgent|TestLowConfidence'
go run ./cmd/eval -input data/eval-questions.jsonl -output reports/eval-summary.json -markdown-output reports/eval-summary.md -fail-on-critical
make secret-check
make eval
make docker-build
make smoke
ASKOC_API_PORT=9180 make smoke
```

P11 release verification uses package-specific privacy, audit, handler, LLM, classifier, RAG, workflow, config, API, simulator, orchestrator, eval, and build-artifact tests plus `go test ./...`. The deterministic evaluation gate is `make eval`; it must report zero critical failures before demo release. The local repeatability gate is `make smoke`; it builds and starts the Compose stack, checks `/healthz`, and verifies transcript workflow plus CRM action traces using synthetic IDs.

## Smoke checks

`scripts/smoke.sh` can either start Docker Compose itself or test an already running local stack. Compose host ports default to `9080`-`9085`; set `ASKOC_API_PORT`, `ASKOC_BANNER_PORT`, `ASKOC_PAYMENT_PORT`, `ASKOC_CRM_PORT`, `ASKOC_WORKFLOW_PORT`, or `ASKOC_LMS_PORT` to another `9xxx` range when a default port is already in use.

```bash
make smoke
make compose-up
make compose-test
scripts/smoke.sh --compose --keep-stack
scripts/smoke.sh --base-url http://localhost:9080
```

`make smoke` is the release-style proof and tears the Compose stack down after the assertions. Use `scripts/smoke.sh --compose --keep-stack` for exploratory testing. `make compose-test` targets `http://localhost:9080`; when the stack is running on alternate ports, use `scripts/smoke.sh --base-url http://localhost:<api-port>`.

Expected smoke assertions:

- `/healthz` responds before the timeout,
- `S100002` transcript-status chat returns `payment_reminder_triggered`,
- `S100003` financial-hold chat returns `financial_hold_detected` and `crm_case_created`,
- failures print the endpoint or missing response marker that needs attention.

## R0-R5 web app revamp verification

The revamp keeps AskOC as the product and uses `DESIGN.md` only as visual theme guidance. The verified surface is still the existing Go-rendered route set: `/` and `/chat` for learner chat, `/admin` for protected dashboard review, `/static/*` for local assets, `/api/v1/chat` for chat, admin metrics/audit/review APIs, `/healthz`, and `/readyz`.

Automated gates:

```bash
go test ./internal/handlers -run 'Test(Chat|Admin).*Revamp|Test(Chat|Admin)StaticAssets'
go test ./internal/handlers
go test ./...
make eval
make secret-check
git diff --check
```

The revamp handler/static tests verify:

- stale `P2 placeholder` chat copy is absent,
- chat/admin navigation has landmarks, active route state, and synthetic-mode labeling,
- `#0A0A0A`, `#FAFAFA`, `#EF4444`, square radius, 2px borders, and the `DESIGN.md` focus ring exist in static CSS,
- gradients, shadows, and rounded dashboard panels are not introduced,
- chat rendering contracts include source confidence, risk, freshness, low-confidence/no-source fallback labels, trace IDs, workflow IDs, CRM case IDs, priority, and idempotency key evidence,
- admin rendering contracts include review trace IDs, queue, priority, status, redacted markers, safe empty state, and evaluation-gate copy.

Responsive and accessibility checks:

- CSS includes mobile breakpoints at `max-width: 820px` for chat, action trace, admin metrics, review rows, and controls.
- Key controls keep explicit labels: chat message, synthetic student ID, submit, admin token, audit export, purge, reset, and review filter.
- Keyboard focus uses `0 0 0 2px #FAFAFA, 0 0 0 4px #0A0A0A`.
- Manual screenshot review should inspect `/chat` and `/admin` at mobile and desktop widths for text overlap, clipped buttons, hidden source/action rows, dominant red usage, and accidental raw learner data.

## Test data

| Student ID | Transcript status | Payment status | Holds | Expected result |
|---|---|---|---|---|
| `S100001` | ready_for_processing | paid | none | Assistant says request is ready; no reminder |
| `S100002` | blocked_by_unpaid_fee | unpaid | mock_payment_hold | Payment reminder workflow triggered |
| `S100003` | needs_staff_review | review_required | mock_financial_hold | Registrar/Student Accounts CRM case created |
| `S100004` | not_found | not_applicable | none | Normal handoff for unresolved synthetic status |
| `S999999` | not found | n/a | n/a | Safe not-found response; no data leak |

## Intent test set

| Intent | Test prompt | Expected behavior |
|---|---|---|
| `transcript_request` | “How do I order my official transcript?” | Grounded answer with source |
| `transcript_status` | “I ordered my transcript and it has not arrived.” | Ask/check synthetic student status |
| `fee_payment` | “I paid my fee but my account still shows a balance.” | Route to payment status workflow or escalation |
| `myokanagan_login` | “I cannot log into myOkanagan.” | Provide portal support guidance or IT handoff |
| `lms_access` | “Where do I find my online course?” | Answer from LMS/online resources content |
| `registration_help` | “How do I register for classes?” | Grounded registration guidance |
| `refund_request` | “Can I get a refund after withdrawing?” | Grounded answer and caution about deadlines |
| `application_status` | “Has my application been received?” | Mock application-status flow or escalation |
| `key_dates` | “When is the withdrawal deadline?” | Source-grounded answer; avoid guessing dates |
| `human_handoff` | “I want to speak to a person.” | Create or offer escalation |
| `unknown` | “Can you write my essay?” | Safe refusal or redirect |

## Functional tests

### Test 1: Grounded transcript answer

**Prompt:**

```text
How do I order an official transcript?
```

**Expected:**

- intent is `transcript_request`,
- response includes approved source link,
- response does not invent unsupported policy,
- confidence is above threshold,
- no enterprise API call is made unless user asks for status.

### Test 2: Unpaid transcript workflow

**Prompt:**

```text
I ordered my transcript but it has not been processed. My student ID is S100002.
```

**Expected:**

- intent is `transcript_status`,
- mock Banner API is called,
- mock Payment API is called,
- payment status is `unpaid`,
- payment reminder workflow action is triggered with workflow ID and idempotency key,
- audit port records workflow attempted/completed events with a hashed idempotency key and retry attempt count when applicable,
- no CRM case is created unless learner is frustrated or workflow fails.

### Test 3: Financial hold escalation

**Prompt:**

```text
Can you check my transcript? My student ID is S100003.
```

**Expected:**

- transcript status is checked,
- payment status is review-required in the synthetic fixture,
- financial hold is detected,
- assistant avoids giving detailed financial judgment,
- CRM case is created for Registrar/Finance follow-up,
- learner receives case ID.

### Test 4: Negative sentiment escalation

**Prompt:**

```text
This is extremely frustrating. I need this transcript today for a job application.
```

**Expected:**

- sentiment is `negative` or `urgent`,
- priority escalation is triggered if conversation context involves unresolved transcript issue,
- case summary includes only necessary context,
- assistant acknowledges urgency without promising impossible outcomes.

### Test 5: Unknown answer fallback

**Prompt:**

```text
Can you guarantee my transfer credit will be approved?
```

**Expected:**

- assistant does not guarantee outcome,
- assistant provides general guidance only if grounded source exists,
- assistant escalates or recommends staff support,
- low-confidence event appears in dashboard.

## Go unit tests

| Package | Test focus |
|---|---|
| `internal/privacy` | PII redaction, likely password/token redaction, false positives, synthetic ID preservation, and safe log inputs |
| `internal/domain` | chat request/response JSON models, intent/source/action/escalation fields |
| `internal/validation` | empty, whitespace-only, oversized message, and synthetic student ID validation |
| `internal/fixtures` | synthetic fixture loading, duplicate rejection, required fields, synthetic ID enforcement |
| `internal/mock/banner` | known/unknown synthetic student profile and transcript/hold status handlers |
| `internal/mock/payment` | paid/unpaid transcript payment status handlers and safe unknown-payment errors |
| `internal/mock/crm` | case creation, priority routing, required summaries, and summary redaction |
| `internal/mock/lms` | synthetic LMS access-status lookup and unknown-course fallback |
| `internal/session` | create, append, read, expire, redaction, and concurrent access behavior |
| `internal/classifier` | valid structured output, invalid JSON fallback, confidence thresholds |
| `internal/llm` | provider-neutral payloads, safe provider errors, OpenAI-compatible request/response handling, rate-limit/retryable/timeout errors |
| `internal/rag` | chunking, metadata, retrieval top-k ranking, stale source flags |
| `internal/tools` | trace header forwarding, timeout handling, not-found/retryable/parse error mapping, response parsing, safe errors |
| `internal/workflow` | idempotency, retry behavior, duplicate prevention |
| `internal/orchestrator` | decision table for transcript/payment/escalation |
| `internal/audit` | redacted memory store, trace queries, metrics, export, reset, prune, and retention policy |
| `internal/handlers` | request validation, status codes, trace IDs, protected admin metrics, dashboard shell, audit export/reset/purge |
| `internal/eval` | JSONL parsing, deterministic runner behavior, scoring, reports, gates, and unresolved review queue items |

## P6 classification fixture gate

`data/classification-fixtures.jsonl` is the P6 synthetic intent/sentiment fixture set. The gate is 100% fixture intent accuracy for the supported demo intents: `transcript_request`, `transcript_status`, `fee_payment`, `human_handoff`, `escalation_request`, and `unknown`. The fixture test also requires at least five examples per intent, negative/urgent sentiment coverage, and no tool-trigger permission for unknown/off-topic fixtures.

## Example table-driven test cases

```go
func TestRedact(t *testing.T) {
    tests := []struct {
        name string
        in   string
        want string
    }{
        {"email", "Email me at student@example.com", "Email me at [REDACTED_EMAIL]"},
        {"password", "my password is Hunter2", "my password is [REDACTED_SECRET]"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := privacy.Redact(tt.in)
            if got != tt.want {
                t.Fatalf("got %q want %q", got, tt.want)
            }
        })
    }
}
```

## RAG quality tests

| Check | Expected result |
|---|---|
| Correct source retrieved | Relevant source appears in top 3 chunks |
| Source required | No policy answer without source |
| No unsupported claims | Answer facts map to retrieved chunks |
| Freshness metadata | Source date/indexed date available |
| Stale/high-risk source handling | Assistant asks for staff confirmation instead of presenting the answer as authoritative |
| Conflicting content | Assistant escalates or asks for clarification |

## Automation tests

| Scenario | Expected workflow result |
|---|---|
| unpaid transcript | payment reminder sent |
| paid transcript, no hold | no reminder; status response only |
| financial hold | CRM case created |
| unknown student | safe not-found response; no case unless requested |
| duplicate unpaid request | no duplicate reminder within configured window |
| automation timeout | learner informed; CRM case created or retry logged |

## Privacy and security tests

| Test | Expected result |
|---|---|
| User enters email address | Email is redacted in logs |
| User enters password | Assistant tells user not to share passwords; password redacted |
| User asks for another student's record | Refusal; no API data returned |
| Prompt injection attempt | System ignores unsafe instruction |
| Logs reviewed | No raw payment card, password, or government ID values |
| Repo scan | No secrets, tokens, or real records |

## Performance tests

| Metric | Target |
|---|---:|
| Chat response, RAG only | Under 5 seconds |
| Chat response with two mock API calls | Under 7 seconds |
| Workflow trigger response | Under 3 seconds after decision |
| Dashboard load | Under 3 seconds for demo dataset |
| Concurrent demo users | 10 users locally |
| Go API unit test suite | Under 10 seconds locally |

## Contract tests

Run contract tests between `cmd/api` and mock services:

- API expects Banner response shape.
- API expects payment response shape.
- API expects CRM response shape.
- API expects workflow simulator response shape.
- Invalid payloads return structured errors.

## Dashboard validation

After running test conversations, verify that the dashboard at `/admin` shows:

- total conversations,
- top intents,
- containment rate,
- escalation rate,
- automation success/failure,
- low-confidence answer count,
- stale-source warning count,
- redacted review queue items.

## Exit criteria

The MVP is demo-ready when:

- at least 30 curated synthetic evaluation questions pass,
- all four synthetic student scenarios work,
- no answer is generated without source support for policy/procedure questions,
- low-confidence fallback works,
- escalation works,
- audit logs are redacted,
- dashboard metrics update after interactions,
- README, architecture diagrams, demo script, privacy notes, test plan, and evaluation report are consistent,
- known limitations are honest,
- `make test` passes,
- `make eval` passes with zero critical policy errors,
- `make smoke` passes against the local stack.

## P11 release checklist

| Check | Evidence |
|---|---|
| README explains the project in under two minutes | README has a two-minute reviewer path, quickstart, architecture pointer, success metrics, TDD quality statement, and screenshot placeholders. |
| Architecture and sequence diagrams are present | `docs/architecture.md` includes Mermaid high-level and interview sequence diagrams covering chat UI, Go API, orchestrator, RAG, mock Banner/payment/CRM/LMS, workflow, audit, and dashboard. |
| Demo script fits 5-7 minutes | `docs/demo-script.md` has a minute-by-minute run sheet and keeps full test commands as prep/release evidence rather than live narration. |
| Screenshot/GIF placeholders are privacy-reviewed | README and demo script list placeholders with captions; `docs/privacy-impact-lite.md` defines the capture review checks. |
| Release commands pass | `make test`, `make eval`, `make smoke`, and `make secret-check` are the final local evidence set. |
| Synthetic-only boundary holds | `data/synthetic-students.json`, redaction tests, privacy notes, and screenshot review rules prohibit real learner data, real payment data, private URLs, and secrets. |
| Limitations are honest | README and privacy notes keep real authentication, real Banner/payment/CRM/LMS integrations, private portal scraping, and production deployment out of scope. |

## P11 local release evidence

Latest local verification on 2026-05-06:

| Command or review | Result |
|---|---|
| README top-to-bottom review | Pass: problem, solution, stack, quickstart, privacy, success metrics, architecture pointer, screenshot placeholders, and TDD quality statement are present. |
| Diagram golden-path trace | Pass: chat UI, Go API, orchestrator, RAG, mock Banner/payment/CRM/LMS, workflow, audit, dashboard, and eval/report components are represented. |
| Demo-script dry run | Pass: the run sheet fits the 5-7 minute interview window when full test commands are treated as prep/release evidence rather than live narration. |
| Screenshot/GIF privacy review | Pass: placeholders only; future captures must follow the privacy checklist before binary assets are added. |
| `make test` | Pass: `go test ./...` completed successfully. |
| `make eval` | Pass: 34/34 cases passed with zero critical failures and refreshed `reports/eval-summary.md`. |
| `make secret-check` | Pass: no known live-token patterns detected. |
| `git diff --check` | Pass: no whitespace errors. |
| `make smoke` | Default port `9080` was occupied locally after images built; rerun with a documented alternate `9xxx` API port override passed health, transcript workflow, and CRM smoke checks. |
