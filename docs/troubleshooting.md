# Troubleshooting

Use this guide when local setup, smoke tests, admin APIs, or demo workflows do not behave as expected.

## Quick Diagnostics

Check the API:

```bash
curl -i --max-time 5 http://localhost:8080/healthz
curl -i --max-time 5 http://localhost:8080/readyz
```

Check Compose services:

```bash
docker compose ps
docker compose logs --tail=200 api mock-banner mock-payment mock-crm workflow-sim
```

Check the common gates:

```bash
go test ./...
make eval
make secret-check
git diff --check
```

## Common Issues

| Symptom | Likely cause | Fix |
|---|---|---|
| `make smoke` fails because port `8080` is in use | Another local service is using the API port | Run `ASKOC_API_PORT=18080 make smoke` or stop the service using port `8080`. |
| `make compose-up` fails with bind errors on mock ports | One of ports `8081` through `8085` is in use | Override `ASKOC_BANNER_PORT`, `ASKOC_PAYMENT_PORT`, `ASKOC_CRM_PORT`, `ASKOC_WORKFLOW_PORT`, or `ASKOC_LMS_PORT`. |
| Docker smoke does not reach `/healthz` | Docker daemon is unavailable, image build failed, or API exited during startup | Run `docker compose ps` and `docker compose logs --tail=200 api`. |
| `make dev` starts, but transcript-status chat fails | The mock services are not running on localhost | Start `cmd/mock-banner`, `cmd/mock-payment`, and `cmd/mock-crm`, or use `make compose-up`. |
| Compose API cannot reach mock services | Localhost URLs were used inside Compose | Compose must use service DNS names such as `http://mock-banner:8081`, which are already set in `docker-compose.yml`. |
| Manual API cannot reach mock services | Compose service DNS names were exported into local shell | For manual local runs, use localhost URLs such as `http://localhost:8081`. |
| Admin API returns `401` | Missing or wrong bearer token | Use `Authorization: Bearer demo-admin-token`, unless `ASKOC_AUTH_TOKEN` is set. |
| `/admin` loads but metrics are empty | No chat actions have been recorded in this API process | Submit a transcript request first, then reload metrics. |
| Every route returns `401` | `ASKOC_AUTH_ENABLED=true` is set | Include the configured bearer token on every request or unset `ASKOC_AUTH_ENABLED` for local demo mode. |
| Startup fails with provider config error | `ASKOC_PROVIDER=openai-compatible` is set without endpoint or API key | Use default `ASKOC_PROVIDER=stub`, or set endpoint and API key for manual provider testing only. |
| Grounded transcript answer falls back | `data/rag-chunks.json` is missing, empty, or not readable | Regenerate chunks with `go run ./cmd/ingest -sources data/seed-sources.json -out data/rag-chunks.json` when public network access is available. |
| `make eval` changes files | Evaluation reports are generated artifacts | Review `reports/eval-summary.json` and `reports/eval-summary.md`; commit them only when the evidence intentionally changed. |
| `make secret-check` fails on a placeholder | Placeholder text matched a live-token pattern | Replace it with a safer fake value or update the scanner only if the scanner is clearly wrong. |
| Smoke response is missing `payment_reminder_triggered` | The unpaid transcript path changed or payment mock data is not aligned | Check `data/synthetic-students.json`, `internal/orchestrator`, and mock payment service tests. |
| Smoke response is missing `crm_case_created` | The financial-hold escalation path changed or CRM mock is unavailable | Check `internal/orchestrator`, `internal/mock/crm`, and API logs. |

## Port Override Example

Run the whole stack on alternate host ports:

```bash
ASKOC_API_PORT=18080 \
ASKOC_BANNER_PORT=18081 \
ASKOC_PAYMENT_PORT=18082 \
ASKOC_CRM_PORT=18083 \
ASKOC_WORKFLOW_PORT=18084 \
ASKOC_LMS_PORT=18085 \
make compose-up
```

Then open:

```text
http://localhost:18080/chat
```

## Reset Local State

The API audit store is in memory. Restarting the API clears it. You can also reset it without restarting:

```bash
curl -sS -X POST -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/audit/reset
```

For Compose, stop containers and remove orphans:

```bash
make compose-down
```

## When A Test Fails

Start with the narrow package that failed, not the full suite:

```bash
go test ./internal/orchestrator -run TestName
```

Then run:

```bash
go test ./...
```

If the failure involves answer quality, source grounding, actions, handoff, or safety gates, run:

```bash
make eval
```

## Safe Debugging Rules

- Do not paste real student IDs, real payment details, real tokens, or private webhook URLs into fixtures or docs.
- Keep `ASKOC_PROVIDER=stub` unless you are explicitly testing provider integration by hand.
- Prefer synthetic IDs such as `S100002` and `S100003` when reproducing workflow behavior.
- If logs include unexpected sensitive text, fix redaction before sharing output.
