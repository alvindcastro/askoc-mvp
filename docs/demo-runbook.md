# End-To-End Demo Runbook

Use this runbook when preparing or presenting the full AskOC local demo. It covers what to turn on, how to verify the stack, what to click or send, what evidence to show, and how to shut everything down.

The default demo is fully local and deterministic. It uses synthetic learner records, local RAG chunks, mock Banner/payment/CRM/LMS services, and the local workflow simulator. It does not require a live LLM, live SIS, live CRM, live LMS, live payment processor, Power Automate, or real student data.

## Stack

| Layer | Service | Default URL | Started by Compose | Demo proof |
|---|---|---|---:|---|
| API, chat UI, admin UI | `cmd/api` | `http://localhost:8080` | yes | Chat, admin metrics, health, readiness |
| Mock Banner | `cmd/mock-banner` | `http://localhost:8081` | yes | Synthetic student, transcript, and hold checks |
| Mock payment | `cmd/mock-payment` | `http://localhost:8082` | yes | Synthetic transcript payment status |
| Mock CRM | `cmd/mock-crm` | `http://localhost:8083` | yes | Synthetic staff handoff case creation |
| Workflow simulator | `cmd/workflow-sim` | `http://localhost:8084` | yes | Payment reminder workflow acceptance |
| Mock LMS | `cmd/mock-lms` | `http://localhost:8085` | yes | Optional synthetic LMS access status |
| Local RAG data | `data/rag-chunks.json` | loaded by API | yes | Approved-source transcript answer |
| Eval runner | `cmd/eval` | local command | no | Responsible-AI quality report |

Compose wires the API to service-DNS URLs such as `http://mock-banner:8081` and `http://workflow-sim:8084/api/v1/automation/payment-reminder`. Host URLs are only for your browser, curl, and inspection.

## Prerequisites

- Run commands from the repository root.
- Install Go 1.22 or newer.
- Install Docker with Docker Compose v2.
- Have `make`, `curl`, and a POSIX shell available.
- Use only synthetic demo IDs such as `S100002` and `S100003`.

No secrets are needed for the default path. Keep `ASKOC_PROVIDER=stub` for interview demos unless you are explicitly testing the optional OpenAI-compatible provider path.

## Fast Proof

Run this before a demo to prove the Docker stack and golden API paths work:

```bash
make smoke
```

`make smoke` builds and starts the Compose stack, waits for `/healthz`, verifies the unpaid transcript workflow for `S100002`, verifies the financial-hold CRM escalation for `S100003`, and tears the stack down when it finishes.

If host port `8080` is occupied:

```bash
ASKOC_API_PORT=18080 make smoke
```

## Start The Full Stack

For a live UI walkthrough, start the stack and keep it running:

```bash
make compose-up
```

This runs in the foreground. Leave that terminal open. In a second terminal, verify the running stack:

```bash
docker compose ps
curl -sS http://localhost:8080/healthz
curl -sS http://localhost:8080/readyz
```

Alternative prep path: start the stack, run the smoke checks, and leave the containers running:

```bash
scripts/smoke.sh --compose --keep-stack
```

Open these screens:

```text
Chat UI:  http://localhost:8080/chat
Admin UI: http://localhost:8080/admin
Health:   http://localhost:8080/healthz
Ready:    http://localhost:8080/readyz
```

Use `demo-admin-token` when the admin UI or admin API asks for a token.

## Optional Clean Slate

If you ran smoke checks before presenting and want a clean dashboard, reset the API's in-memory audit state:

```bash
curl -sS -X POST \
  -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/audit/reset
```

The demo state is process-local. Restarting the API also clears admin metrics and audit events.

## Demo Flow

Run the spoken walkthrough from [docs/demo-script.md](demo-script.md). Use this operational sequence for the end-to-end path.

| Step | Screen | Input | Expected proof |
|---:|---|---|---|
| 1 | `/chat` | `How do I order my official transcript?` | `transcript_request` answer grounded in approved transcript source metadata |
| 2 | `/chat` | `I ordered my transcript but it has not been processed. My student ID is S100002.` | Mock Banner and payment checks, unpaid status, `payment_reminder_triggered`, workflow ID |
| 3 | `/chat` | `My transcript request has a financial hold and is not moving. My student ID is S100003.` | Financial hold detected, no payment reminder, `crm_case_created`, synthetic CRM case ID |
| 4 | `/chat` | `This is really frustrating. I need this transcript for a job application.` | Urgent or negative sentiment signal, priority staff handoff, safe expectation-setting |
| 5 | `/admin` | token `demo-admin-token` | Containment, escalation, workflow, stale-source, low-confidence, and review metrics without raw PII |
| 6 | terminal | `make eval` | Updated `reports/eval-summary.md` with responsible-AI gate evidence |

You can also send the key API requests from a terminal.

Grounded transcript answer:

```bash
curl -sS -H 'Content-Type: application/json' \
  -d '{"channel":"web","message":"How do I order an official transcript?"}' \
  http://localhost:8080/api/v1/chat
```

Unpaid transcript workflow:

```bash
curl -sS -H 'Content-Type: application/json' \
  -d '{"channel":"web","message":"I ordered my transcript but it has not been processed. My student ID is S100002.","student_id":"S100002"}' \
  http://localhost:8080/api/v1/chat
```

Financial-hold escalation:

```bash
curl -sS -H 'Content-Type: application/json' \
  -d '{"channel":"web","message":"My transcript request has a financial hold and is not moving. My student ID is S100003.","student_id":"S100003"}' \
  http://localhost:8080/api/v1/chat
```

Admin metrics:

```bash
curl -sS -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/metrics
```

## Quality Evidence

Use these commands during prep or at the end of the walkthrough:

```bash
go test ./...
make eval
make secret-check
make compose-test
```

Use `make compose-test` when the stack is already running. Use `make smoke` when you want the script to start and stop Compose for you.

## Port Overrides

If one or more default host ports are occupied, override them when starting Compose:

```bash
ASKOC_API_PORT=18080 \
ASKOC_BANNER_PORT=18081 \
ASKOC_PAYMENT_PORT=18082 \
ASKOC_CRM_PORT=18083 \
ASKOC_WORKFLOW_PORT=18084 \
ASKOC_LMS_PORT=18085 \
make compose-up
```

Then open the API on the overridden host port:

```text
http://localhost:18080/chat
http://localhost:18080/admin
```

The API still talks to the mock services through internal Compose DNS. You normally only need the overridden API port for browser and curl requests.

## Manual Startup Without Containers

Use this path when debugging Go code and you do not want to rebuild containers. Start each command in its own terminal:

```bash
go run ./cmd/mock-banner
go run ./cmd/mock-payment
go run ./cmd/mock-crm
go run ./cmd/workflow-sim
go run ./cmd/mock-lms
```

Then start the API:

```bash
ASKOC_WORKFLOW_URL=http://localhost:8084/api/v1/automation/payment-reminder make dev
```

If `ASKOC_WORKFLOW_URL` is empty, the API uses the in-process deterministic workflow client instead of the standalone workflow simulator.

## Optional RAG Refresh

The repository already includes local chunks in `data/rag-chunks.json`. Refresh only when approved source metadata changes or you intentionally want to test ingestion:

```bash
go run ./cmd/ingest \
  -sources data/seed-sources.json \
  -out data/rag-chunks.json
```

This command can fail when offline or when a source site is unavailable. A refresh is not required for the default demo.

## Shutdown

Stop the Compose stack:

```bash
make compose-down
```

If you used a detached smoke run with `--keep-stack`, the same shutdown command applies.

## Troubleshooting Checks

| Symptom | Check | Fix |
|---|---|---|
| Port conflict | `docker compose ps` or the Compose startup error | Use `ASKOC_API_PORT=18080` or the full port override block |
| API not ready | `curl -sS http://localhost:8080/healthz` | Wait for builds to finish, then inspect `docker compose ps` |
| Admin API returns unauthorized | Confirm bearer token | Use `Authorization: Bearer demo-admin-token` unless `ASKOC_AUTH_TOKEN` is set |
| Workflow marker missing | Confirm workflow simulator is healthy in Compose | Run `make compose-test` or restart with `make compose-down` then `make compose-up` |
| RAG answer falls back | Confirm `data/rag-chunks.json` exists | Keep the seeded file, or refresh only when public sources are reachable |

