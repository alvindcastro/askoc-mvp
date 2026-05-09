# Nice To Knows

These are practical details that are easy to miss when reviewing or extending the repo.

## The Default Demo Is Offline

`ASKOC_PROVIDER=stub` is the default. The normal demo path does not call a live LLM, SIS, CRM, LMS, payment processor, or Power Automate endpoint. It uses deterministic Go logic, synthetic fixtures, local RAG chunks, and mock services.

## Synthetic IDs Are Intentional

The `S100001` style IDs are synthetic and safe to preserve in logs, audit action references, screenshots, and docs. Real-looking numeric IDs, email addresses, phone numbers, passwords, tokens, and API-key-like values should be redacted by shared privacy code before storage or display.

## The Admin Page And Admin APIs Are Different

`/admin` serves the local dashboard shell. The admin data APIs require a bearer token:

```text
Authorization: Bearer demo-admin-token
```

If `ASKOC_AUTH_ENABLED=true`, the mock auth middleware protects all routes, including the chat UI and health checks.

## Workflow Has Two Local Modes

When `ASKOC_WORKFLOW_URL` is empty, `cmd/api` uses the in-process idempotent workflow client. In Docker Compose, `ASKOC_WORKFLOW_URL` points to `workflow-sim`, which exercises an HTTP webhook-shaped path.

Both modes should keep learner-facing action traces safe and should avoid storing raw idempotency keys in audit metadata.

## Evaluation Reports Are Generated Evidence

`make eval` writes:

```text
reports/eval-summary.json
reports/eval-summary.md
```

Those files are portfolio evidence. If they change, review whether the behavior change was intentional before committing them.

## RAG Is Local JSON For This MVP

The implemented retrieval path reads local chunks from `data/rag-chunks.json`. Embeddings, pgvector, Azure AI Search, and production content pipelines are intentionally out of scope for the default demo.

If retrieval confidence is low, or a source is stale/high-risk, the assistant should fall back or add caution instead of inventing policy.

## Tool Calls Are Go Decisions

The LLM does not directly call arbitrary tools. The Go orchestrator decides when to check Banner, payment, workflow, CRM, or handoff paths based on validated intent, confidence, synthetic student ID handling, source guardrails, and sentiment.

## Compose Uses Service DNS

Inside Docker Compose, the API reaches mock services with names such as:

```text
http://mock-banner:8081
http://mock-payment:8082
http://mock-crm:8083
http://workflow-sim:8084
```

Manual local runs should use localhost URLs instead.

## The Audit Store Is In Memory

Audit events, admin metrics, and review queue state are process-local demo state. Restarting the API clears them. Use export before reset when you need to inspect evidence:

```bash
curl -sS -H 'Authorization: Bearer demo-admin-token' \
  http://localhost:8080/api/v1/admin/audit/export
```

## Port Overrides Are First-Class

The Compose stack supports host port overrides with environment variables such as `ASKOC_API_PORT`. This is the preferred fix when a default local port is already occupied.

## Docs And Code Should Move Together

When changing visible behavior, check the practical docs and the deeper specs. Most behavior changes touch at least one of:

- `README.md`
- `docs/api-spec.md`
- `docs/setup.md`
- `docs/how-to.md`
- `docs/testing.md`
- `docs/troubleshooting.md`
- `docs/demo-script.md`
- `CHANGELOG.md`
