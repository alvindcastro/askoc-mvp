# Testing Guide

This is the practical command guide for local verification. The deeper scenario matrix lives in [Test Plan](test-plan.md).

## Default Gates

Run these for most code changes:

```bash
go test ./...
go vet ./...
make eval
make secret-check
git diff --check
```

What each command proves:

| Command | Proves |
|---|---|
| `go test ./...` | Unit, handler, orchestration, mock, workflow, RAG, eval, and command tests pass. |
| `go vet ./...` | Basic static checks pass. |
| `make eval` | Deterministic responsible-AI gate passes with zero critical failures. |
| `make secret-check` | Known live-token patterns are not present in tracked or unignored local files. |
| `git diff --check` | Markdown and code diffs do not contain whitespace errors. |

## Focused Test Targets

Use focused targets while changing a specific area:

| Area | Command |
|---|---|
| Chat validation and handlers | `go test ./internal/domain ./internal/validation ./internal/handlers` |
| Orchestrator decisions | `go test ./internal/orchestrator` |
| Classifier and LLM parsing | `go test ./internal/classifier ./internal/llm` |
| RAG ingestion and retrieval | `go test ./internal/rag` |
| Tool clients and mock services | `go test ./internal/tools ./internal/mock/banner ./internal/mock/payment ./internal/mock/crm ./internal/mock/lms` |
| Workflow simulator and webhook client | `go test ./internal/workflow ./cmd/workflow-sim` |
| Audit, admin metrics, review queue | `go test ./internal/audit ./internal/eval ./internal/handlers` |
| Evaluation runner | `go test ./internal/eval ./cmd/eval` |
| API startup wiring | `go test ./cmd/api ./internal/config` |
| Session concurrency | `make test-race` |

## Smoke Testing

Run the full Docker-backed demo proof:

```bash
make smoke
```

Run against an already running stack:

```bash
make compose-test
```

Run with an alternate API port:

```bash
ASKOC_API_PORT=18080 make smoke
```

Expected smoke checks:

- `/healthz` responds before timeout.
- `S100002` transcript-status request includes `payment_reminder_triggered`.
- `S100003` financial-hold request includes `financial_hold_detected`.
- `S100003` financial-hold request includes `crm_case_created`.

## Evaluation Testing

Default deterministic evaluation:

```bash
make eval
```

Live local API evaluation:

```bash
go run ./cmd/eval \
  -input data/eval-questions.jsonl \
  -output reports/eval-summary.json \
  -markdown-output reports/eval-summary.md \
  -base-url http://localhost:8080/api/v1 \
  -fail-on-critical
```

The committed report files are evidence. If a change intentionally affects eval behavior, regenerate and review both:

```text
reports/eval-summary.json
reports/eval-summary.md
```

Do not lower critical safety expectations to make a behavior change pass. Fix the behavior or update the test case only when the expected outcome was wrong.

## Docker And Packaging Tests

Use these when Docker, Compose, ports, assets, or startup wiring changed:

```bash
make docker-build
make smoke
```

If default ports are busy:

```bash
ASKOC_API_PORT=18080 make smoke
```

## Test Data

Use only synthetic IDs from `data/synthetic-students.json`.

| ID | Expected path |
|---|---|
| `S100001` | Paid transcript, ready for processing, no workflow reminder. |
| `S100002` | Unpaid transcript fee, payment reminder workflow. |
| `S100003` | Financial hold or staff review, CRM escalation. |
| `S100004` | Unresolved synthetic status, human handoff. |
| `S999999` | Not found, safe response without data leakage. |

## TDD Workflow

For Go behavior changes:

1. Add the smallest failing test that describes the behavior.
2. Run the focused package test and confirm it fails for the expected reason.
3. Implement the smallest passing change.
4. Run the focused package test.
5. Run `go test ./...`.
6. Run `make eval` when answer, action, handoff, source, safety, or review behavior changed.
7. Update docs and changelog when observable behavior or proof expectations changed.

## Before Handoff

Use this short checklist:

- Focused tests for the changed packages passed.
- `go test ./...` passed.
- `make eval` passed when relevant.
- `make secret-check` passed.
- `git diff --check` passed.
- Docs match the implemented behavior.
- No real secrets, real student data, or private endpoints were added.
