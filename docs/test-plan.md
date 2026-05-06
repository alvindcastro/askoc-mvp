# Test Plan


## Strict TDD requirement

All Go code tasks are governed by [TDD Policy](tdd-policy.md). For each code task, create failing tests before production code, verify the red state, implement the smallest passing change, run the narrow test, then run `go test ./...`. Do not mark tasks complete in [Phases and Tickable Tasks](phases-and-tasks.md) without test evidence.


## Purpose

This test plan verifies that AskOC AI Concierge can answer learner-service questions, automate the transcript/payment workflow, escalate appropriately, and protect privacy in a Go-based demo environment. P4 currently covers deterministic classifier, transcript/payment orchestration, in-process workflow idempotency, mock CRM handoff, and safe action traces; RAG, durable audit/dashboard, and the standalone workflow simulator remain later-phase coverage.

## Test objectives

1. Validate grounded answers for common Tier 0 questions.
2. Validate intent classification across learner-service categories.
3. Validate transcript status and payment workflow logic.
4. Validate sentiment-based and low-confidence escalation.
5. Validate privacy controls, redaction, and audit logging.
6. Validate dashboard metrics.
7. Validate Go services under demo-level load.

## Test environment

| Component | Test setup |
|---|---|
| Frontend | Go server-rendered web chat or optional React UI |
| Backend | Go `cmd/api` local container |
| Knowledge base | P4 transcript source placeholder; P5 adds approved public-page indexing |
| Mock Banner API | Go service with synthetic student records |
| Mock Payment API | Go service with synthetic payment records |
| Mock CRM API | Go service with synthetic case creation |
| Automation workflow | P4 in-process idempotent workflow port; P8 adds Go workflow simulator or Power Automate demo flow |
| Dashboard | Deferred to P7 Go dashboard endpoints reading audit/event store |

## Go test commands

```bash
go test ./...
go test ./internal/classifier ./internal/workflow ./internal/orchestrator
go test ./internal/domain ./internal/validation ./internal/handlers ./internal/session
go test -race ./internal/session
go test ./internal/orchestrator -run 'TestTranscriptStatus|TestUrgent|TestLowConfidence'
```

P4 verification uses the package-specific classifier, workflow, and orchestrator tests plus `go test ./...`. Privacy, evaluation, dashboard, and broad race coverage remain later-phase checks unless those packages exist.

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
- audit port records workflow attempted/completed events,
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
| `internal/privacy` | PII redaction, password warnings, safe summaries |
| `internal/domain` | chat request/response JSON models, intent/source/action/escalation fields |
| `internal/validation` | empty, whitespace-only, oversized message, and synthetic student ID validation |
| `internal/fixtures` | synthetic fixture loading, duplicate rejection, required fields, synthetic ID enforcement |
| `internal/mock/banner` | known/unknown synthetic student profile and transcript/hold status handlers |
| `internal/mock/payment` | paid/unpaid transcript payment status handlers and safe unknown-payment errors |
| `internal/mock/crm` | case creation, priority routing, required summaries, and summary redaction |
| `internal/mock/lms` | synthetic LMS access-status lookup and unknown-course fallback |
| `internal/session` | create, append, read, expire, redaction, and concurrent access behavior |
| `internal/classifier` | valid structured output, invalid JSON fallback, confidence thresholds |
| `internal/rag` | chunking, metadata, retrieval top-k ranking, stale source flags |
| `internal/tools` | trace header forwarding, timeout handling, not-found/retryable/parse error mapping, response parsing, safe errors |
| `internal/workflow` | idempotency, retry behavior, duplicate prevention |
| `internal/orchestrator` | decision table for transcript/payment/escalation |
| `internal/handlers` | request validation, status codes, trace IDs |

## Example table-driven test cases

```go
func TestRedact(t *testing.T) {
    tests := []struct {
        name string
        in   string
        want string
    }{
        {"email", "Email me at student@example.com", "Email me at [REDACTED_EMAIL]"},
        {"password", "my password is Hunter2", "my [REDACTED_SECRET]"},
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

After running test conversations, verify that the dashboard shows:

- total conversations,
- top intents,
- containment rate,
- escalation rate,
- automation success/failure,
- average response time,
- low-confidence answer count,
- unresolved questions,
- sentiment distribution.

## Exit criteria

The MVP is demo-ready when:

- at least 20 curated questions pass,
- all four synthetic student scenarios work,
- no answer is generated without source support for policy/procedure questions,
- low-confidence fallback works,
- escalation works,
- audit logs are redacted,
- dashboard metrics update after interactions,
- `go test ./...` passes,
- README and demo script are complete.
