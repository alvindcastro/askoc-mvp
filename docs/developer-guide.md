# Developer Guide

This guide is the day-to-day map for changing the AskOC Go MVP without losing the demo contract.

## Development Contract

- Keep the default path deterministic, offline, and synthetic-data-only.
- Use strict TDD for Go behavior changes: write a failing test, confirm the red state, implement the smallest passing change, then run the relevant package tests.
- Do not route real learner, Banner, payment, CRM, LMS, webhook, or provider data through the repo.
- Keep the portfolio demo focused on transcript guidance, transcript/payment status, payment reminder workflow, CRM handoff, audit, admin metrics, and evaluation evidence.
- Prefer existing package boundaries over new abstractions unless the current boundary is clearly insufficient.

Docs-only edits do not require a failing Go test first, but they should still match implemented behavior.

## Repository Map

| Path | Use it for |
|---|---|
| `cmd/api` | API server, chat UI, admin UI, route wiring, dependency construction. |
| `cmd/mock-banner` | Synthetic student profile and transcript status service. |
| `cmd/mock-payment` | Synthetic transcript payment status service. |
| `cmd/mock-crm` | Synthetic CRM case service. |
| `cmd/mock-lms` | Synthetic LMS access-status service. |
| `cmd/workflow-sim` | Local Power Automate-style payment reminder simulator. |
| `cmd/ingest` | Approved-source ingestion to local RAG chunks. |
| `cmd/eval` | JSONL evaluation runner and quality gate. |
| `internal/orchestrator` | Main transcript, grounding, workflow, escalation, and guardrail decisions. |
| `internal/handlers` | HTTP handlers, UI handlers, admin APIs, health checks, response helpers. |
| `internal/tools` | Typed HTTP clients for mock Banner, payment, CRM, and LMS services. |
| `internal/workflow` | In-process workflow client, simulator handler, idempotency, webhook client. |
| `internal/rag` | Source loading, chunking, freshness, retrieval, local retriever. |
| `internal/classifier` | Deterministic fallback and strict LLM classification parsing. |
| `internal/llm` | Provider-neutral LLM types and OpenAI-compatible REST client. |
| `internal/privacy` | Shared PII and secret redaction. |
| `internal/audit` | In-memory audit events, metrics, export, retention, reset, purge. |
| `internal/eval` | Dataset parsing, deterministic client, scoring, gates, review queue, reports. |
| `web/templates` | Server-rendered chat and admin HTML shells. |
| `web/static` | Small CSS and JavaScript assets for the local UI. |
| `data` | Synthetic students, source allowlist, RAG chunks, evaluation datasets. |
| `reports` | Generated evaluation evidence. |
| `docs` | Product, architecture, setup, testing, troubleshooting, and portfolio docs. |

## Common Change Paths

For API request or response shape changes:

1. Update or add tests in `internal/domain`, `internal/validation`, and `internal/handlers`.
2. Update `docs/api-spec.md`.
3. Run the handler package tests and `go test ./...`.
4. Add changelog notes when the external contract changes.

For orchestrator behavior changes:

1. Add or update decision-table tests in `internal/orchestrator`.
2. Keep low-confidence, stale-source, and unknown-intent cases safe.
3. Confirm audit actions remain redacted and traceable.
4. Run `go test ./internal/orchestrator ./internal/audit ./internal/eval`.
5. Run `make eval` if answer, action, handoff, source, or safety behavior changed.

For mock integration changes:

1. Update the relevant mock package tests under `internal/mock/*`.
2. Update typed client tests under `internal/tools`.
3. Keep fixture IDs synthetic and stable unless the test dataset is updated with the same change.
4. Run the focused mock/client tests, then `go test ./...`.

For workflow changes:

1. Update tests in `internal/workflow` and any transcript orchestration tests that depend on workflow actions.
2. Preserve idempotency-key hashing in audit metadata.
3. Keep the in-process workflow path as the default when `ASKOC_WORKFLOW_URL` is empty.
4. Run `go test ./internal/workflow ./internal/orchestrator ./cmd/workflow-sim`.

For RAG or approved-source changes:

1. Update source metadata in `data/seed-sources.json`.
2. Regenerate chunks with `go run ./cmd/ingest -sources data/seed-sources.json -out data/rag-chunks.json`.
3. Run `go test ./internal/rag ./internal/orchestrator`.
4. Run `make eval` if answer grounding or source expectations changed.

For UI changes:

1. Keep `/chat` and `/admin` usable as local demo tools, not marketing pages.
2. Keep text compact and operational.
3. Verify the UI manually with the Docker or manual local stack.
4. Run `go test ./internal/handlers ./cmd/api`.

For evaluation changes:

1. Update `data/eval-questions.jsonl` and the expected action/source/handoff fields together.
2. Run `make eval`.
3. Commit the regenerated `reports/eval-summary.json` and `reports/eval-summary.md` when the committed evidence is meant to change.

## Local Quality Gates

Use these before handing off code:

```bash
go test ./...
go vet ./...
make eval
make secret-check
git diff --check
```

Use these when the local stack or Docker path changed:

```bash
make docker-build
ASKOC_API_PORT=18080 make smoke
```

Use this when session concurrency or redaction-adjacent session behavior changed:

```bash
make test-race
```

## Documentation Sync

Update docs in the same change when implementation behavior changes:

| Behavior changed | Docs to check |
|---|---|
| API route, payload, auth, or error shape | `docs/api-spec.md`, `README.md`, `docs/how-to.md` |
| Setup command, env var, port, Docker service | `docs/setup.md`, `README.md`, `INDEX.md` |
| Test command or gate | `docs/testing.md`, `docs/test-plan.md`, `README.md` |
| Demo flow | `docs/demo-script.md`, `README.md` |
| Architecture boundary | `docs/architecture.md`, `docs/golang-implementation.md` |
| Privacy, audit, retention, redaction | `docs/privacy-impact-lite.md`, `docs/nice-to-knows.md` |
| Workflow/webhook behavior | `docs/power-automate-flow.md`, `docs/how-to.md` |

When in doubt, prefer updating one focused practical guide and linking to the deeper reference instead of duplicating long specs in several files.
