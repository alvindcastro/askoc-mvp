# Setup Guide

This guide gets the local AskOC demo running from a fresh checkout. It assumes a Linux, macOS, or WSL shell from the repository root.

## Prerequisites

- Go 1.22 or newer.
- Docker with Docker Compose v2 for the full local stack.
- `make`, `curl`, and a POSIX shell.
- No live student data, live Banner data, or real payment data.

The default runtime uses deterministic `stub` mode and synthetic fixtures, so no LLM key or webhook secret is required.

## First Check

Run the fastest repository health check first:

```bash
go test ./...
```

Then run the repeatable demo proof:

```bash
make smoke
```

`make smoke` builds the Docker Compose stack, waits for the API health check, verifies the unpaid transcript workflow for `S100002`, verifies the financial-hold CRM escalation for `S100003`, and tears the stack down when it finishes.

If port `8080` is already in use, override the API host port:

```bash
ASKOC_API_PORT=18080 make smoke
```

## Full Docker Demo

Use this when you want the UI and mock services running together:

```bash
make compose-up
```

Then open:

```text
Chat UI:  http://localhost:8080/chat
Admin UI: http://localhost:8080/admin
Health:   http://localhost:8080/healthz
Ready:    http://localhost:8080/readyz
```

Stop the stack from another terminal:

```bash
make compose-down
```

To run the smoke path and keep the stack alive afterward:

```bash
scripts/smoke.sh --compose --keep-stack
```

## Manual Local Services

Use this when you are actively debugging code without rebuilding containers. Start each service in its own terminal:

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

If `ASKOC_WORKFLOW_URL` is empty, the API uses the in-process deterministic workflow client. The standalone workflow simulator is useful when you want to inspect webhook-like behavior separately.

## Environment Files

`.env.example` documents the local defaults and contains placeholders only. Copy it to `.env` for local experimentation, but do not commit `.env` or any real secrets.

Important defaults:

| Variable | Default | Notes |
|---|---|---|
| `ASKOC_HTTP_ADDR` | `:8080` | API listen address. |
| `ASKOC_PROVIDER` | `stub` | Deterministic mode with no live model calls. |
| `ASKOC_RAG_CHUNKS_PATH` | `data/rag-chunks.json` | Local approved-source chunks. |
| `ASKOC_AUTH_ENABLED` | `false` | If set to `true`, mock bearer auth protects all routes. |
| `ASKOC_AUTH_TOKEN` | empty | Admin APIs use `demo-admin-token` when this is empty. |
| `ASKOC_WORKFLOW_URL` | empty | Empty uses in-process workflow; Compose points to `workflow-sim`. |

Run the secret scanner before committing environment or docs changes:

```bash
make secret-check
```

## Admin Access

The admin HTML shell is available at `/admin`. The admin data APIs require a bearer token:

```bash
curl -sS -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/metrics
```

If `ASKOC_AUTH_TOKEN` is set, use that value instead of `demo-admin-token` for admin API calls.

## Optional OpenAI-Compatible Mode

Automated tests and the default demo must stay in `stub` mode. To manually test the OpenAI-compatible provider path, set all required provider values:

```bash
ASKOC_PROVIDER=openai-compatible \
ASKOC_PROVIDER_MODEL=demo-model-name \
ASKOC_PROVIDER_ENDPOINT=https://example.invalid/v1/chat/completions \
ASKOC_PROVIDER_API_KEY=placeholder-only \
make dev
```

Do not commit real provider endpoints, API keys, webhook secrets, or captured responses containing secrets.
