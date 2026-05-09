# How-To Playbook

Use these recipes for common local development, demo, and review tasks.

## Run The One-Command Demo Proof

```bash
make smoke
```

If port `8080` is occupied:

```bash
ASKOC_API_PORT=18080 make smoke
```

## Start The Full Stack For Manual Review

```bash
make compose-up
```

Visit:

```text
http://localhost:8080/chat
http://localhost:8080/admin
```

Stop the stack:

```bash
make compose-down
```

## Submit A Transcript Status Chat Request

```bash
curl -sS -H 'Content-Type: application/json' \
  -d '{"channel":"web","message":"I ordered my transcript but it has not been processed. My student ID is S100002.","student_id":"S100002"}' \
  http://localhost:8080/api/v1/chat
```

Expected action marker: `payment_reminder_triggered`.

## Submit A Financial-Hold Escalation Request

```bash
curl -sS -H 'Content-Type: application/json' \
  -d '{"channel":"web","message":"My transcript request has a financial hold and is not moving. My student ID is S100003.","student_id":"S100003"}' \
  http://localhost:8080/api/v1/chat
```

Expected action markers: `financial_hold_detected` and `crm_case_created`.

## Ask For A Grounded Transcript Answer

```bash
curl -sS -H 'Content-Type: application/json' \
  -d '{"channel":"web","message":"How do I order an official transcript?"}' \
  http://localhost:8080/api/v1/chat
```

Expected result: `transcript_request` intent with approved-source metadata. If local RAG chunks are missing or unusable, the assistant should use a safe fallback instead of inventing policy.

## Check Admin Metrics

```bash
curl -sS -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/metrics
```

The metrics are in-memory. Restarting the API clears them unless you exported the events first.

## Export, Purge, Or Reset Audit Events

Export redacted audit events:

```bash
curl -sS -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/audit/export
```

Purge expired events using the demo retention policy:

```bash
curl -sS -X POST -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/audit/purge
```

Reset the in-memory audit store:

```bash
curl -sS -X POST -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/audit/reset
```

## Run The Deterministic Evaluation Gate

```bash
make eval
```

This reads `data/eval-questions.jsonl`, writes `reports/eval-summary.json` and `reports/eval-summary.md`, and exits non-zero for critical safety failures.

To evaluate a running local API instead of the deterministic in-process client:

```bash
go run ./cmd/eval \
  -input data/eval-questions.jsonl \
  -output reports/eval-summary.json \
  -markdown-output reports/eval-summary.md \
  -base-url http://localhost:8080/api/v1 \
  -fail-on-critical
```

## Regenerate Local RAG Chunks

This command fetches allowlisted public URLs. It can fail when offline or when a source site is unavailable.

```bash
go run ./cmd/ingest \
  -sources data/seed-sources.json \
  -out data/rag-chunks.json
```

Run this after approved source metadata changes. Then run:

```bash
go test ./internal/rag ./internal/orchestrator
make eval
```

## Run The Workflow Simulator Manually

Start the simulator:

```bash
go run ./cmd/workflow-sim
```

Start the API against the simulator:

```bash
ASKOC_WORKFLOW_URL=http://localhost:8084/api/v1/automation/payment-reminder make dev
```

Leave `ASKOC_WORKFLOW_URL` empty when you want the in-process workflow client instead.

## Use Non-Default Compose Ports

```bash
ASKOC_API_PORT=18080 \
ASKOC_BANNER_PORT=18081 \
ASKOC_PAYMENT_PORT=18082 \
ASKOC_CRM_PORT=18083 \
ASKOC_WORKFLOW_PORT=18084 \
ASKOC_LMS_PORT=18085 \
make compose-up
```

Then use the overridden API URL, for example:

```text
http://localhost:18080/chat
```

## Run A Focused Go Test

```bash
go test ./internal/orchestrator -run TestTranscriptStatus
```

Use focused tests while implementing, then run `go test ./...` before handoff.

## Run Secret And Whitespace Checks

```bash
make secret-check
git diff --check
```

These are cheap checks to run before sharing a branch or recording demo evidence.
