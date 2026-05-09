# Developer Testing Guide

## Purpose

This guide is the day-to-day local testing and troubleshooting companion for AskOC. It covers the Docker stack, repeatable test commands, manual smoke checks, and useful operating notes for the synthetic demo environment.

## Local Docker Stack

The Compose stack builds and runs these local images:

| Service | Image | Default host URL |
|---|---|---|
| API and web UI | `askoc-api:local` | `http://localhost:9080` |
| Mock Banner | `askoc-mock-banner:local` | `http://localhost:9081` |
| Mock payment | `askoc-mock-payment:local` | `http://localhost:9082` |
| Mock CRM | `askoc-mock-crm:local` | `http://localhost:9083` |
| Workflow simulator | `askoc-workflow-sim:local` | `http://localhost:9084` |
| Mock LMS | `askoc-mock-lms:local` | `http://localhost:9085` |

Use the default one-command proof when ports `9080`-`9085` are free:

```bash
make smoke
```

`make smoke` builds the images, starts Compose, verifies `/healthz`, posts the `S100002` unpaid transcript scenario, posts the `S100003` financial-hold scenario, and then tears the stack down.

Use this command when you want to keep the stack running for manual testing:

```bash
scripts/smoke.sh --compose --keep-stack
```

Use another contiguous `9xxx` host range when `9080` is already allocated:

```bash
ASKOC_API_PORT=9180 ASKOC_BANNER_PORT=9181 ASKOC_PAYMENT_PORT=9182 ASKOC_CRM_PORT=9183 ASKOC_WORKFLOW_PORT=9184 ASKOC_LMS_PORT=9185 scripts/smoke.sh --compose --keep-stack
```

When the default-port stack is running, open:

```text
Chat UI:     http://localhost:9080/chat
Admin UI:    http://localhost:9080/admin
Health:      http://localhost:9080/healthz
Readiness:   http://localhost:9080/readyz
Admin token: demo-admin-token
```

Stop the stack with:

```bash
docker compose down --remove-orphans
```

## Test Commands

Run the fast local suite first:

```bash
make test
```

Run focused package checks while developing:

```bash
go test ./internal/handlers
go test ./internal/classifier ./internal/workflow ./internal/orchestrator
go test ./internal/rag ./internal/orchestrator
go test ./internal/eval ./cmd/eval
go test ./internal/build -run TestP10
go test -race ./internal/session
```

Run the responsible-AI and safety gates:

```bash
make eval
make secret-check
```

Run the Docker proof and keep it running:

```bash
scripts/smoke.sh --compose --keep-stack
```

Run smoke checks against an already running default-port stack:

```bash
make compose-test
```

Run smoke checks against an already running stack on an alternate host port:

```bash
scripts/smoke.sh --base-url http://localhost:9180
```

## Manual API Checks

Health:

```bash
curl -fsS http://localhost:9080/healthz
```

Unpaid transcript workflow:

```bash
curl -fsS -H 'Content-Type: application/json' -X POST \
  -d '{"channel":"web","message":"I ordered my transcript but it has not been processed. My student ID is S100002.","student_id":"S100002"}' \
  http://localhost:9080/api/v1/chat
```

Financial-hold escalation:

```bash
curl -fsS -H 'Content-Type: application/json' -X POST \
  -d '{"channel":"web","message":"My transcript request has a financial hold and is not moving. My student ID is S100003.","student_id":"S100003"}' \
  http://localhost:9080/api/v1/chat
```

Admin metrics:

```bash
curl -fsS -H 'Authorization: Bearer demo-admin-token' http://localhost:9080/api/v1/admin/metrics
```

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| `permission denied while trying to connect to the docker API` | Current shell cannot access Docker | Start Docker Desktop or run from a shell with Docker permissions. |
| `Bind for 0.0.0.0:9080 failed: port is already allocated` | Another local process owns `9080` | Use an alternate `9xxx` range such as `9180`-`9185`. |
| Smoke health check times out | API did not become healthy or a dependency failed | Run `docker compose ps` and `docker compose logs --tail=200 api mock-banner mock-payment mock-crm workflow-sim`. |
| `make compose-test` fails while alternate ports are running | `make compose-test` is pinned to `http://localhost:9080` | Use `scripts/smoke.sh --base-url http://localhost:<api-port>`. |
| Chat response is missing `payment_reminder_triggered` | API cannot reach payment or workflow services, or fixture behavior changed | Check `docker compose ps`, API logs, and `data/synthetic-students.json` for `S100002`. |
| Chat response is missing `crm_case_created` | API cannot reach CRM, or hold routing changed | Check API and mock CRM logs, then verify `S100003` in `data/synthetic-students.json`. |
| Admin UI/API returns unauthorized | Missing mock bearer token | Use `demo-admin-token` for local admin routes. |
| Grounded transcript answers look stale or missing | Local RAG chunks were not generated or changed | Run `go run ./cmd/ingest -sources data/seed-sources.json -out data/rag-chunks.json` when refreshing approved chunks. |

## Nice To Knows

- The demo uses synthetic records only. Use `S100001` through `S100004` and avoid real learner data.
- `ASKOC_PROVIDER=stub` is the default. The local stack does not require live AI credentials.
- Compose service URLs inside the Docker network use service names such as `http://mock-banner:9081`; host URLs use `localhost`.
- `make smoke` is a pass/fail release proof and tears down the stack. Use `scripts/smoke.sh --compose --keep-stack` for exploratory testing.
- The Compose stack routes payment reminders through `workflow-sim` by default. Without `ASKOC_WORKFLOW_URL`, the API uses the in-process idempotent workflow client.
- PostgreSQL is not required for the current P11 local demo. The stack is deterministic without a database.
- `make eval` updates `reports/eval-summary.json` and `reports/eval-summary.md`.
- `.env`, local override files, real API keys, and real webhook URLs must stay out of git.
