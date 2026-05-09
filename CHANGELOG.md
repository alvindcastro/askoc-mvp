# Changelog

All notable MVP task changes are recorded here with what changed, where it changed, when it changed, why it changed, and how it was completed.

## 2026-05-09 - Developer Testing Guide And Docker Stack Notes

### Local testing workflow documentation

- What: added a developer testing guide with Docker stack startup, alternate port commands, test commands, manual API checks, troubleshooting, and useful local operating notes.
- Where: `docs/developer-guide.md`, `README.md`, `docs/test-plan.md`, `docs/golang-implementation.md`, `INDEX.md`.
- When: 2026-05-09.
- Why: make local testing and handoff clearer for the running Compose stack, especially when default port `8080` is occupied.
- How: built the local Docker images, attempted the default smoke path, confirmed the `8080` bind conflict, restarted the stack on `18080`-`18085` with `scripts/smoke.sh --compose --keep-stack`, and documented the verified workflow.

## 2026-05-09 - R0-R5 Web App Revamp Implementation

### R0-T01 - Reconcile current web app status

- What: confirmed `/`, `/chat`, `/admin`, static assets, chat API, admin metrics/audit/review APIs, health, and readiness routes; recorded stale P2 copy and evidence gaps.
- Where: `docs/web-app-revamp-tdd-plan.md`, `docs/askoc-ux-theme-brainstorm.md`.
- When: 2026-05-09.
- Why: establish the actual server-rendered AskOC surface before changing UI code.
- How: inspected `cmd/api/main.go`, `internal/handlers`, `web/templates`, `web/static`, README, and demo docs; kept `/` as an intentional chat alias and added no unsupported routes.

### R0-T02 - Freeze revamp task taxonomy

- What: kept the active revamp status in R0-R5 terms and retained UX0-UX5 as historical theme planning only.
- Where: `docs/web-app-revamp-tdd-plan.md`, `docs/askoc-ux-theme-brainstorm.md`.
- When: 2026-05-09.
- Why: avoid two competing task taxonomies for the same AskOC UI work.
- How: marked R0-R5 task headings complete after implementation evidence existed and documented the R0-R5 source-of-truth rule.

### R1-T01 - Map DESIGN.md tokens to AskOC surfaces

- What: mapped VoiceBox visual language to AskOC chat, source citations, action trace, escalation, admin metrics, review rows, eval evidence, inputs, buttons, chips, focus rings, and route navigation.
- Where: `docs/askoc-ux-theme-brainstorm.md`, `docs/test-plan.md`.
- When: 2026-05-09.
- Why: apply `DESIGN.md` as visual polish without changing the product into an editorial app.
- How: documented black/white/red tokens, square borders, no gradients/shadows/rounded panels, mono metadata, and restrained red usage.

### R1-T02 and R1-T03 - Theme tokens, shell, and navigation

- What: added shared theme contracts and high-contrast route shells with active route state, synthetic-mode labels, keyboard focus, and AskOC-first product naming.
- Where: `web/static/app.css`, `web/static/admin.css`, `web/templates/chat.html`, `web/templates/admin.html`, `internal/handlers/ui_test.go`, `internal/handlers/admin_ui_test.go`.
- When: 2026-05-09.
- Why: make the existing Go-rendered app feel like a compact operational console while preserving current routes and APIs.
- How: wrote failing handler/static tests for nav landmarks, active routes, token names, square radius, 2px borders, focus ring, and forbidden visual drift; then implemented the smallest template/CSS changes.

### R2-T01 through R2-T04 - Chat evidence surface

- What: replaced stale P2 placeholder copy; added compact learner/assistant messaging, source evidence, confidence/risk/freshness chips, low-confidence/no-source fallback labels, trace ID display, workflow/CRM/action rows, priority, and idempotency-key rendering.
- Where: `web/templates/chat.html`, `web/static/app.css`, `web/static/app.js`, `internal/handlers/ui_test.go`.
- When: 2026-05-09.
- Why: let a reviewer scan what the learner asked, what answer was grounded, which synthetic systems/actions ran, and what trace ties the chat to audit/eval proof.
- How: extended red handler/static asset tests for stale copy and rendering contracts, verified failure, then updated the template and JavaScript without changing the chat API contract.

### R3-T01 through R3-T03 - Admin, audit, and eval evidence

- What: themed the admin dashboard, metric strip, audit controls, review filter, review rows, trace/queue/priority/status/redacted chips, safe empty state, and evaluation-gate copy.
- Where: `web/templates/admin.html`, `web/static/admin.css`, `web/static/admin.js`, `internal/handlers/admin_ui_test.go`, `docs/test-plan.md`.
- When: 2026-05-09.
- Why: expose aggregate workflow/escalation/review evidence without raw learner messages or unsupported production claims.
- How: added failing admin shell/static tests first, preserved existing protected API IDs/data attributes, kept review semantics unchanged, and documented eval evidence in the test plan.

### R4-T01 through R4-T03 - Accessibility, responsive, and theme drift gates

- What: added automated contracts for labels, landmarks, focus ring, mobile breakpoints, forbidden gradients/rounded panels, restrained red usage, and source/action/admin control coverage.
- Where: `internal/handlers/ui_test.go`, `internal/handlers/admin_ui_test.go`, `web/static/app.css`, `web/static/admin.css`, `docs/test-plan.md`.
- When: 2026-05-09.
- Why: prevent the themed UI from drifting into low-contrast, overlapping, marketing-style, or inaccessible patterns.
- How: verified red/green static tests, added `max-width: 820px` responsive rules, kept visible labels for chat/admin controls, and recorded a manual desktop/mobile screenshot checklist.

### R5-T01 through R5-T03 - Reviewer path and final evidence

- What: updated reviewer/demo/test documentation and recorded exact revamp evidence, commands, limitations, and changed files.
- Where: `README.md`, `docs/demo-script.md`, `docs/test-plan.md`, `docs/web-app-revamp-tdd-plan.md`, `docs/askoc-ux-theme-brainstorm.md`, `CHANGELOG.md`.
- When: 2026-05-09.
- Why: make the revamp understandable as AskOC UX polish, not a separate product, and keep the final evidence auditable.
- How: added README themed-UI proof, demo “point at” callouts, R0-R5 test evidence, manual visual checklist, and this changelog entry after implementation.

### R0-R5 verification evidence

- What: verified the revamp with red/green handler/static contracts and broader handler tests.
- Where: `internal/handlers`, `web/templates`, `web/static`, docs listed above.
- When: 2026-05-09.
- Why: ensure UI changes are test-backed and do not alter API behavior.
- How: red command failed first: `go test ./internal/handlers -run 'Test(Chat|Admin).*Revamp|Test(Chat|Admin)StaticAssets'`; green commands passed: same narrow command, `go test ./internal/handlers`, `go test ./...`, `make eval`, `make secret-check`, `git diff --check`, and `ASKOC_API_PORT=18080 make smoke`. The default `make smoke` image build succeeded but could not bind occupied port `8080`, so the documented alternate port was used.

## 2026-05-09 - AskOC DESIGN.md UX Theme Prompt Plan

### AskOC UX theme TDD prompt pack

- What: added an AskOC MVP UX theme brainstorm, tickable phase/task board, and copy/paste strict-TDD prompts grounded in `DESIGN.md`.
- Where: `docs/askoc-ux-theme-brainstorm.md`, `docs/askoc-ux-tdd-prompts.md`, `INDEX.md`, `CHANGELOG.md`.
- When: 2026-05-09.
- Why: apply the `DESIGN.md` visual system to the existing AskOC AI Concierge MVP without changing the transcript/payment automation product scope.
- How: reviewed the design system, kept the AskOC phase board stable, separated UX theme planning into dedicated Markdown files, and made all implementation tasks tickable with red/green/refactor evidence.

## 2026-05-06 - P11 Portfolio Polish And Interview Readiness

### P11-T01 - Polish README for applicant storytelling

- What: added a two-minute reviewer path, portfolio evidence table, quickstart framing, screenshot/GIF placeholder section, and P11-ready capability summary.
- Where: `README.md`.
- When: 2026-05-06.
- Why: make the project understandable to a hiring manager quickly while preserving the Go-first, synthetic-data-only story.
- How: reviewed the README against the P11 prompt for problem, solution, stack, demo, privacy, success metrics, architecture, and TDD evidence, then moved the most important reviewer actions to the top.

### P11-T02 - Add architecture diagram and sequence diagram

- What: refreshed the Mermaid architecture diagram and added a P11 interview sequence diagram covering Tier 0 answer, Tier 1 workflow, urgent escalation, audit, dashboard, and evaluation evidence.
- Where: `docs/architecture.md`.
- When: 2026-05-06.
- Why: show enterprise integration thinking visually without claiming unsupported services.
- How: traced the golden path through existing commands/packages and removed stale notification-service overclaims from architecture-facing docs.

### P11-T03 - Finalize 5-7 minute demo script

- What: tightened the opening pitch, added a timed run sheet, updated P10-era wording to P11 release wording, and moved full test commands into prep/release evidence.
- Where: `docs/demo-script.md`.
- When: 2026-05-06.
- Why: make the interview walkthrough concise enough for a 5-7 minute slot while keeping every spoken claim observable.
- How: dry-run reviewed the script against the required Tier 0 answer, Tier 1 workflow, urgent escalation, dashboard, and TDD/evaluation proof points.

### P11-T04 - Prepare screenshots or short GIF placeholders

- What: added textual screenshot/GIF placeholders, captions, a capture manifest, and a portfolio screenshot privacy checklist.
- Where: `README.md`, `docs/demo-script.md`, `docs/assets/README.md`, `docs/privacy-impact-lite.md`.
- When: 2026-05-06.
- Why: make the GitHub repo visually scannable without committing binary captures that could accidentally expose real data, tokens, or private URLs.
- How: defined the future capture set for chat grounded answer, transcript workflow, CRM escalation, admin dashboard, and eval report, then documented allowed and disallowed visible values.

### P11-T05 - Run final release checklist

- What: marked P11 and MVP release gates complete, aligned roadmap/test-plan/scope/brainstorm docs, corrected stale notification overclaims, and recorded final release evidence.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `docs/test-plan.md`, `docs/mvp-scope.md`, `docs/brainstorm.md`, `docs/golang-implementation.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: confirm the portfolio is coherent, safe, demonstrable, and honest about limitations.
- How: verified `make test`, `make eval`, `make secret-check`, and `git diff --check`; default `make smoke` built images but found local port `8080` occupied, then the documented `ASKOC_API_PORT=18080 make smoke` path passed health, transcript workflow, and CRM smoke checks.

## 2026-05-06 - P10 Docker, CI, And Local Developer Experience

### P10-T01 - Create Dockerfiles for Go services

- What: added a multi-stage Dockerfile that builds any Go service under `cmd/*` with an `APP` build argument, copies only runtime assets needed for the demo, and runs as a non-root user, plus Docker ignore rules for local secrets and build outputs.
- Where: `Dockerfile`, `.dockerignore`, `internal/build/p10_artifacts_test.go`.
- When: 2026-05-06.
- Why: make the API and mock services reproducible as small local containers without copying `.env` or source-control metadata into images.
- How: added failing P10 artifact tests first for Dockerfile and `.dockerignore` requirements, confirmed the missing-file red state, then implemented the minimal image build contract and ignore list.

### P10-T02 - Create Docker Compose local stack

- What: added a Docker Compose stack for `api`, `mock-banner`, `mock-payment`, `mock-crm`, `mock-lms`, and `workflow-sim` with deterministic stub-provider defaults, service-DNS integration URLs, exposed local ports, and health checks.
- Where: `docker-compose.yml`, `Makefile`, `README.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/demo-script.md`, `docs/test-plan.md`, `internal/build/p10_artifacts_test.go`.
- When: 2026-05-06.
- Why: replace multi-terminal demo startup with a repeatable synthetic local stack.
- How: started from failing artifact tests for compose services, provider defaults, service URLs, health checks, and Makefile targets, then added the Compose stack and documented the port/service contract.

### P10-T03 - Add CI workflow for Go tests and evaluation gates

- What: added a GitHub Actions workflow that checks out the repo, installs Go 1.22, runs `go test ./...`, runs `go vet ./...`, and runs `make eval` with `ASKOC_PROVIDER=stub`.
- Where: `.github/workflows/ci.yml`, `internal/build/p10_artifacts_test.go`, `README.md`, `docs/test-plan.md`, `docs/implementation-roadmap.md`.
- When: 2026-05-06.
- Why: make the deterministic responsible-AI quality gate part of pull-request delivery without requiring live model credentials.
- How: wrote the failing CI artifact contract first, confirmed the workflow file was missing, then added the offline CI job and docs.

### P10-T04 - Add environment sample and secret-safety checks

- What: added a placeholder-only `.env.example`, a `.gitignore` that keeps `.env` and local override files out of git, and a `make secret-check` target backed by a grep-style scanner for known live-token patterns.
- Where: `.env.example`, `.gitignore`, `scripts/check-secrets.sh`, `Makefile`, `internal/build/p10_artifacts_test.go`, `README.md`, `docs/test-plan.md`.
- When: 2026-05-06.
- Why: keep deterministic demo mode usable without secrets and reduce the risk of committing real webhook URLs or API keys.
- How: extended the failing P10 environment test to require secret-check coverage, confirmed the missing target failure, then added the sample env file, ignore rules, scanner script, and Make target.

### P10-T05 - Add one-command smoke test

- What: added `scripts/smoke.sh` and Makefile targets for `make smoke`, `make compose-up`, `make compose-test`, and `make compose-down`; the smoke script can start Compose, wait for `/healthz`, and verify `S100002` payment workflow plus `S100003` CRM handoff response markers.
- Where: `scripts/smoke.sh`, `Makefile`, `docs/test-plan.md`, `docs/demo-script.md`, `README.md`, `internal/build/p10_artifacts_test.go`.
- When: 2026-05-06.
- Why: give interview prep and local review one clear pass/fail command for the golden synthetic demo path.
- How: wrote the failing smoke-script and Makefile artifact tests first, then implemented actionable shell checks for health, chat workflow, and CRM action traces.

### P10 review evidence

- What: completed P10 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/demo-script.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, README, architecture notes, implementation guide, test plan, demo script, and changelog aligned with the implemented Docker, CI, env-safety, and smoke-test workflow.
- How: confirmed the red state with `go test ./internal/build -run TestP10`, implemented the smallest passing P10 artifacts, then verified with `go test ./internal/build -run TestP10`, `make secret-check`, `go test ./...`, `go vet ./...`, `make eval`, `make docker-build`, `ASKOC_API_PORT=18080 make smoke`, and `git diff --check`.

## 2026-05-06 - P9 Evaluation Runner And Quality Gates

### P9-T01 - Create JSONL evaluation dataset

- What: added a 34-case JSONL evaluation dataset covering transcript answers, transcript status, payment workflow decisions, financial-hold escalation, human handoff, urgent sentiment, prompt injection, password redaction, unauthorized record-access refusal, and off-topic fallback.
- Where: `data/eval-questions.jsonl`, `internal/eval/dataset.go`, `internal/eval/dataset_test.go`.
- When: 2026-05-06.
- Why: make the demo quality bar repeatable and synthetic-data-only instead of relying on manual spot checks.
- How: wrote failing parser and fixture validation tests first for valid JSONL, line-numbered invalid rows, required expected fields, at least 30 cases, and critical safety coverage, then implemented strict JSONL parsing and case validation.

### P9-T02 - Build `cmd/eval` runner

- What: added `cmd/eval` plus a reusable runner that can call either a deterministic in-process chat client or a live local `/api/v1/chat` endpoint, records latency, captures timeouts, and writes reports.
- Where: `cmd/eval/main.go`, `cmd/eval/main_test.go`, `internal/eval/runner.go`, `internal/eval/runner_test.go`, `internal/eval/deterministic_client.go`.
- When: 2026-05-06.
- Why: let the quality gate run offline by default while still supporting a local API contract check during demos.
- How: wrote failing runner and CLI tests first for per-case fake-client calls, timeout capture, latency recording, and zero exit on a passing deterministic dataset, then implemented the runner, HTTP client, deterministic fixture-backed client, and command flags.

### P9-T03 - Implement scoring functions

- What: added deterministic scoring for intent, sentiment, source match, expected/forbidden actions, handoff queue, safety checks, forbidden critical claims, and latency warnings.
- Where: `internal/eval/score.go`, `internal/eval/score_test.go`.
- When: 2026-05-06.
- Why: separate critical policy/safety regressions from minor misses and make action/source/handoff expectations machine-checkable.
- How: wrote failing score tests first for exact intent/source/action/handoff matches, forbidden critical answer substrings, and latency thresholds, then implemented action aliasing for evaluation language such as `grounded_answer_returned` and `sentiment_classified`.

### P9-T04 - Generate JSON and Markdown evaluation reports

- What: added JSON and Markdown report generation with summary metrics, per-case results, quality-gate status, prompt redaction support, and committed sample outputs.
- Where: `internal/eval/report.go`, `internal/eval/report_test.go`, `reports/eval-summary.json`, `reports/eval-summary.md`.
- When: 2026-05-06.
- Why: provide portfolio-ready evidence that a reviewer can read without replaying every test case.
- How: wrote failing report tests first for JSON summary/per-case output, Markdown metrics tables, and prompt redaction, then generated fresh reports with `go run ./cmd/eval`.

### P9-T05 - Fail builds on critical evaluation failures

- What: added gate evaluation and `make eval`; critical safety, hallucination, or required-escalation failures return a non-zero exit while minor accuracy misses are warnings.
- Where: `internal/eval/gates.go`, `internal/eval/gates_test.go`, `cmd/eval/main.go`, `Makefile`, `README.md`.
- When: 2026-05-06.
- Why: make responsible AI behavior part of the developer workflow rather than a manual demo checklist.
- How: wrote failing gate tests first for critical hallucination, missing escalation, and warning-only minor misses, then wired gate results into `cmd/eval` and the Makefile target.

### P9-T06 - Create review queue for failed/low-confidence answers

- What: added an in-memory eval review queue with duplicate collapse, resolution, redacted questions, source/action context, critical failure reasons, and a protected admin endpoint consumed by the dashboard.
- Where: `internal/eval/review_queue.go`, `internal/eval/review_queue_test.go`, `internal/handlers/admin.go`, `internal/handlers/admin_test.go`, `internal/handlers/admin_ui.go`, `cmd/api/main.go`, `web/templates/admin.html`, `web/static/admin.js`, `docs/api-spec.md`.
- When: 2026-05-06.
- Why: show how failed or low-confidence answers can become a staff review workflow without exposing raw PII.
- How: wrote failing review queue and admin endpoint tests first for low-confidence additions, critical failures, duplicate normalized questions, resolution, sources/actions, and redaction, then implemented the queue and dashboard wiring.

### P9 review evidence

- What: completed P9 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/model-evaluation.md`, `docs/demo-script.md`, `docs/brainstorm.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, README, API contract, architecture notes, implementation guide, test plan, evaluation spec, demo script, and changelog aligned with the implemented P9 quality gate.
- How: confirmed the red test state for new P9 packages, implemented the smallest passing code, generated reports with `go run ./cmd/eval -input data/eval-questions.jsonl -output reports/eval-summary.json -markdown-output reports/eval-summary.md -fail-on-critical`, then verified `go test ./internal/eval ./cmd/eval ./internal/handlers`, `make eval`, `go test ./...`, `make test`, `go vet ./...`, and `git diff --check`.

## 2026-05-06 - P8 Workflow Automation And Power Automate Option

### P8-T01 - Build local workflow simulator service

- What: added a standalone local workflow simulator with a Power Automate-style payment reminder endpoint, deterministic workflow IDs, duplicate idempotency handling, health/readiness routes, and protected admin metrics/export over simulator audit events.
- Where: `cmd/workflow-sim/main.go`, `internal/workflow/sim_handler.go`, `internal/workflow/sim_handler_test.go`.
- When: 2026-05-06.
- Why: make the workflow automation demo runnable offline while preserving the same request/response shape a Power Automate HTTP trigger would use.
- How: wrote failing simulator handler tests first for accepted payloads, missing idempotency keys, duplicate keys, invalid JSON, and hashed audit metadata, then implemented the handler and standalone command using the existing in-memory workflow client and audit store.

### P8-T02 - Implement workflow idempotency and retry policy

- What: added redaction-safe idempotency-key hashing, exposed workflow attempt counts, and added transient webhook retry behavior with cancellation handling.
- Where: `internal/workflow/idempotency.go`, `internal/workflow/idempotency_test.go`, `internal/workflow/client.go`, `internal/workflow/powerautomate.go`, `internal/workflow/powerautomate_test.go`.
- When: 2026-05-06.
- Why: prevent duplicate learner reminders and make retry behavior visible without storing raw idempotency keys in audit data.
- How: started from failing workflow tests for deterministic hashes, transient `5xx` retry, permanent `400` no-retry, and context cancellation, then implemented the smallest client changes to satisfy those cases.

### P8-T03 - Add optional Power Automate webhook client

- What: added a `PaymentReminderSender`-compatible HTTP webhook client that posts the tested JSON schema, forwards trace and idempotency headers, optionally sends a configured signature header, redacts secrets from errors/config output, and is selected by `cmd/api` when `ASKOC_WORKFLOW_URL` is set.
- Where: `internal/workflow/powerautomate.go`, `internal/workflow/powerautomate_test.go`, `internal/config/config.go`, `internal/config/config_test.go`, `cmd/api/main.go`, `cmd/api/main_test.go`.
- When: 2026-05-06.
- Why: show how the Go orchestration boundary can connect to Power Automate or the local simulator without making external webhooks mandatory for the demo.
- How: wrote failing `httptest.Server`, config, and API wiring tests first, then implemented webhook client config, optional signature headers, retry limits, and default in-process fallback wiring.

### P8-T04 - Audit workflow outcomes and errors

- What: updated workflow audit records to include trace ID, action, status, hashed idempotency metadata, safe failure messages, workflow references, and retry attempt counts when available; simulator audit records no longer store raw idempotency keys.
- Where: `internal/orchestrator/transcript.go`, `internal/orchestrator/orchestrator.go`, `internal/orchestrator/orchestrator_test.go`, `internal/workflow/sim_handler.go`, `internal/workflow/sim_handler_test.go`.
- When: 2026-05-06.
- Why: make workflow automation measurable and debuggable in the admin dashboard/export path without leaking learner identifiers, webhook secrets, or raw idempotency keys.
- How: added failing orchestrator and simulator audit tests first, confirmed missing hash metadata, then reused `workflow.IdempotencyKeyHash` across orchestrator action audits and simulator events.

### P8-T05 - Document Power Automate flow schema and setup

- What: updated the workflow schema, setup notes, endpoint contract, runtime env vars, simulator behavior, security guidance, and P8 status across project docs.
- Where: `docs/power-automate-flow.md`, `docs/api-spec.md`, `README.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/demo-script.md`, `docs/implementation-roadmap.md`, `docs/phases-and-tasks.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the project documentation aligned with the implemented local simulator, webhook client, retry policy, and audit behavior instead of describing those paths as future work.
- How: reviewed the P8 task prompts and stale P7/P8 claims, then synchronized request/response examples, secure webhook storage guidance, command lists, environment variables, and task checkboxes.

### P8 review evidence

- What: completed P8 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/power-automate-flow.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/demo-script.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: provide PR-ready evidence that workflow automation, optional webhook behavior, and docs moved together.
- How: verified focused workflow/orchestrator/config/API tests during implementation, then ran `go test ./...`, `make test`, `go vet ./...`, `make test-race`, `go test ./... -coverprofile=/tmp/askoc-p8-coverage.out`, `go tool cover -func=/tmp/askoc-p8-coverage.out`, and `git diff --check`.

## 2026-05-06 - P7 Privacy, Audit, And Dashboard

### P7-T01 - Implement PII redaction

- What: added shared redaction for email addresses, separated/compact phone numbers, likely password/passcode/token/API-key assignments, and real-looking numeric IDs while preserving synthetic `S100001`-style demo IDs.
- Where: `internal/privacy/redact.go`, `internal/privacy/redact_test.go`, `internal/session/store.go`, `internal/middleware/logging.go`, `internal/orchestrator/orchestrator.go`, `internal/mock/crm/crm_test.go`.
- When: 2026-05-06.
- Why: enforce one privacy boundary for logs, sessions, audit payloads, and CRM summaries without erasing synthetic demo traceability.
- How: wrote failing redaction and logging tests first, confirmed missing `privacy.Redact`, then wired the shared redactor through session persistence, request logging, orchestrator CRM summaries, and mock CRM expectations.

### P7-T02 - Create audit event schema and store

- What: added timestamped audit event constants, JSON tags, an in-memory concurrency-safe audit store, trace queries, redacted storage, export with message content omitted, reset, and prune support.
- Where: `internal/audit/event.go`, `internal/audit/store.go`, `internal/audit/store_test.go`.
- When: 2026-05-06.
- Why: give the demo a safe operational event stream for intent, tool, workflow, guardrail, and escalation evidence before a database-backed audit store exists.
- How: implemented the audit package behind failing store tests for trace lookup, redaction, concurrent writes, retention timestamps, export, reset, and pruning, then aligned the store with the shared privacy redactor.

### P7-T03 - Wire redacted structured logging across services

- What: updated request logging to use the shared redactor, added privacy logging coverage for query secrets and PII, and kept request logs structured with method, path, status, and trace ID only.
- Where: `internal/middleware/logging.go`, `internal/middleware/logging_test.go`, `internal/privacy/logging_test.go`.
- When: 2026-05-06.
- Why: preserve useful local troubleshooting logs without storing raw learner messages or secret-like values.
- How: added a failing privacy logging test proving raw body and query PII are absent, then routed `BasicRedactor` through `privacy.Redact` while keeping `slog` field-based logging.

### P7-T04 - Build admin metrics endpoint

- What: added aggregate audit summaries for total conversations, containment/escalation rates, top intents, workflow counts/failures, low-confidence items, stale-source warnings, and action/status/type buckets behind a protected admin endpoint.
- Where: `internal/audit/metrics.go`, `internal/handlers/admin.go`, `internal/handlers/admin_test.go`, `cmd/api/main.go`, `cmd/api/main_test.go`.
- When: 2026-05-06.
- Why: make staff-facing demo metrics observable from actual redacted audit events instead of response-body-only traces.
- How: wrote failing seeded-audit and empty-store handler tests first, then implemented `SummaryFromEvents`, bearer-token admin checks, `GET /api/v1/admin/metrics`, and default local admin token handling.

### P7-T05 - Create minimal admin dashboard UI

- What: added a server-rendered admin dashboard shell with JavaScript for protected metrics refresh plus top intents, escalation, workflow, low-confidence review, stale-source, export, purge, and reset controls.
- Where: `internal/handlers/admin_ui.go`, `internal/handlers/admin_ui_test.go`, `web/templates/admin.html`, `web/static/admin.js`, `web/static/admin.css`, `cmd/api/main.go`.
- When: 2026-05-06.
- Why: give stakeholders a quick operational view of adoption and risk while keeping sensitive review text redacted.
- How: wrote failing dashboard-render tests first, then implemented `GET /admin`, the dashboard template, admin JavaScript, and restrained dashboard styling.

### P7-T06 - Add retention and export controls

- What: added a seven-day demo retention policy, purge helper, protected audit export endpoint, protected purge endpoint, protected reset endpoint, and docs for export/reset/purge behavior.
- Where: `internal/audit/retention.go`, `internal/audit/retention_test.go`, `internal/handlers/admin.go`, `internal/handlers/admin_test.go`, `docs/privacy-impact-lite.md`, `docs/api-spec.md`.
- When: 2026-05-06.
- Why: make demo audit data short-lived, resettable, and exportable without exposing raw messages.
- How: wrote failing retention and admin control tests first, then implemented default retention pruning, export with omitted message content, reset, purge, and matching privacy/API documentation.

### P7 review evidence

- What: completed P7 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/demo-script.md`, `docs/test-plan.md`, `docs/privacy-impact-lite.md`, `docs/api-spec.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, README, demo script, test plan, privacy notes, and API contract aligned with implemented privacy, audit, dashboard, and retention behavior.
- How: marked P7 tasks and gates complete after red-to-green package evidence, then verified focused P7 packages and the full Go suite with `go test ./...`.

## 2026-05-06 - P6 LLM Gateway And Structured Classification

### P6-T01 - Define LLM provider interface and request/response types

- What: added provider-neutral LLM request/response types, role-tagged messages, grounding source and citation models, token usage metadata, timeout validation, and typed provider errors.
- Where: `internal/llm/types.go`, `internal/llm/types_test.go`.
- When: 2026-05-06.
- Why: hide provider-specific details behind a replaceable Go interface before adding OpenAI-compatible runtime calls.
- How: wrote failing JSON schema, provider error, unwrap, and timeout validation tests first, confirmed missing `internal/llm` symbols, then implemented the minimal provider-neutral types and safe error helpers.

### P6-T02 - Implement OpenAI/Azure-compatible REST client

- What: added a testable OpenAI/Azure-compatible chat completions client with configured endpoint, API key, model, timeout, context cancellation, safe status mapping, and prompt/body-safe errors.
- Where: `internal/llm/openai_client.go`, `internal/llm/openai_client_test.go`.
- When: 2026-05-06.
- Why: allow optional live-provider demo wiring without making automated tests or the default local path depend on a real LLM service.
- How: wrote `httptest.Server` payload, success parsing, 429, 500, timeout, and pre-canceled-context tests first, confirmed the missing client failures, then implemented the REST client without any live API calls.

### P6-T03 - Parse strict JSON classification output

- What: added strict JSON classification parsing for `intent`, `intent_confidence`, `sentiment`, `urgency`, `needs_handoff`, and `reason`, with safe fallback for malformed, non-strict, unknown, or low-confidence output.
- Where: `internal/classifier/llm_parser.go`, `internal/classifier/llm_parser_test.go`.
- When: 2026-05-06.
- Why: ensure model output is typed and validated before it can influence transcript/payment tool decisions.
- How: wrote parser tests for valid JSON, malformed JSON, unknown fields, trailing payloads, missing required fields, unknown intents, out-of-range confidence, and below-threshold tool disabling before implementing strict decoding and safe fallback results.

### P6-T04 - Create prompt templates for classification and grounded answers

- What: added versioned prompt templates for strict JSON classification and source-only grounded answers, plus prompt drift, privacy-boundary, and source-rule tests.
- Where: `internal/orchestrator/prompts.go`, `internal/orchestrator/prompts_test.go`.
- When: 2026-05-06.
- Why: make LLM instructions reviewable, testable, and aligned with synthetic-data and no-direct-tool-calling rules.
- How: wrote failing prompt tests first for strict JSON wording, allowed labels, source-only answer rules, synthetic-data privacy, and a golden prompt version, then implemented plain text prompt builders with `PromptVersion`.

### P6-T05 - Add LLM fallback and guardrail behavior

- What: added a guarded LLM-backed classifier, deterministic fallback on model/parser failure, source-only answer validation, prompt-version audit metadata, low-confidence classification blocking before sensitive tool calls, and optional LLM answer generation only when approved sources are present.
- Where: `internal/orchestrator/ai_guardrails.go`, `internal/orchestrator/ai_guardrails_test.go`, `internal/orchestrator/orchestrator.go`, `internal/orchestrator/grounded_answer.go`, `internal/audit/event.go`, `cmd/api/main.go`, `cmd/api/main_test.go`, `internal/config/config.go`, `internal/config/config_test.go`.
- When: 2026-05-06.
- Why: keep the core demo reliable and safe if a model times out, returns malformed output, or suggests an unsupported answer, while preserving deterministic local behavior by default.
- How: wrote failing orchestrator guardrail tests for timeout fallback, unsourced answer rejection, low-confidence no-tool behavior, and source-guarded LLM answers before implementing audit metadata, provider wiring, config validation for `ASKOC_PROVIDER_*`, and default stub mode.

### P6-T06 - Add end-to-end classification tests with fixture messages

- What: added a JSONL classification fixture set with 31 synthetic examples, at least five examples per supported intent, urgent/negative sentiment coverage, unknown/off-topic no-tool assertions, and regression messages that identify the fixture and expected intent.
- Where: `data/classification-fixtures.jsonl`, `internal/classifier/e2e_test.go`, `internal/classifier/fallback.go`, `docs/model-evaluation.md`, `docs/test-plan.md`.
- When: 2026-05-06.
- Why: measure deterministic intent and sentiment behavior before using structured LLM classification in the demo.
- How: wrote fixture coverage and accuracy tests first, confirmed classifier fixture regressions, then expanded fallback keyword coverage and documented the 100% synthetic fixture accuracy target.

### P6 review evidence

- What: completed P6 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/model-evaluation.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, runtime configuration, architecture notes, implementation guide, test plan, and evaluation expectations aligned with the optional guarded LLM gateway and strict classification behavior.
- How: marked P6 tasks and gates complete after red-to-green package evidence, then verified the P6 package targets, full Go suite, vet, make test, and whitespace checks.

## 2026-05-06 - P5 RAG Ingestion And Source-Grounded Answers

### P5-T01 - Define source allowlist schema

- What: added the typed source allowlist schema with URL, title, department, risk, retrieved date, freshness, allowlist, and private-portal validation.
- Where: `internal/rag/source.go`, `internal/rag/source_test.go`, `data/seed-sources.json`, `docs/source-references.md`.
- When: 2026-05-06.
- Why: ensure only approved public learner-service sources can be used for retrieval and source-grounded policy/procedure answers.
- How: wrote failing allowlist parsing tests first, confirmed missing `internal/rag` schema failures, then implemented context-aware JSON loading, HTTPS/private-marker validation, source lookup maps, and per-source freshness metadata in the seed fixture.

### P5-T02 - Implement ingestion fetcher with allowlist and cleaning boundaries

- What: added deterministic approved-source fetching, safe HTML cleaning, content hashing, typed RAG errors, and the `cmd/ingest` local chunk writer.
- Where: `internal/rag/ingest.go`, `internal/rag/ingest_test.go`, `cmd/ingest/main.go`.
- When: 2026-05-06.
- Why: fetch only allowlisted public content and prevent private or unlisted pages from entering the demo knowledge base.
- How: wrote `httptest.Server` tests for unallowlisted URL rejection, script/style/nav removal, content hash creation, and network failure typing before implementing the fetcher, cleaner, and CLI.

### P5-T03 - Implement chunking with metadata preservation

- What: added configurable word chunking with stable chunk IDs, overlap support, no-empty-chunk behavior, and copied source metadata.
- Where: `internal/rag/chunk.go`, `internal/rag/chunk_test.go`.
- When: 2026-05-06.
- Why: make retrieved public content citeable by chunk while preserving the source URL, title, risk level, and freshness status.
- How: wrote short-content, long-content, stable-ID, and metadata-copy tests first, confirmed missing chunking symbols, then implemented deterministic chunk generation and ID hashing.

### P5-T04 - Create local retrieval implementation

- What: added a swappable local retriever, chunk loading from JSON, confidence scoring, top-k limiting, source metadata projection, and pre-seeded local demo chunks.
- Where: `internal/rag/retrieve.go`, `internal/rag/local_retriever.go`, `internal/rag/retrieve_test.go`, `data/rag-chunks.json`, `internal/config/config.go`, `internal/config/config_test.go`, `cmd/api/main.go`, `README.md`.
- When: 2026-05-06.
- Why: provide offline demo-safe retrieval before Azure AI Search, pgvector, or live embedding services exist.
- How: wrote retrieval ranking, unrelated-query, limit, and metadata tests first, then implemented token scoring, confidence thresholds, `ASKOC_RAG_CHUNKS_PATH`, and API startup wiring with a safe disabled fallback if chunks are unavailable.

### P5-T05 - Add grounded answer source packaging

- What: added retrieval-backed transcript-request answers and transcript-status source attachment with de-duplicated citations, source confidence/risk/freshness metadata, low-confidence fallback, and source-grounding action traces.
- Where: `internal/orchestrator/grounded_answer.go`, `internal/orchestrator/grounded_answer_test.go`, `internal/orchestrator/orchestrator.go`, `internal/domain/chat.go`, `web/static/app.js`, `docs/api-spec.md`, `docs/demo-script.md`.
- When: 2026-05-06.
- Why: replace the P4 hard-coded transcript source placeholder with approved-source retrieval and prevent unsupported policy answers.
- How: wrote orchestrator tests for sufficient sources, duplicate citations, low-confidence fallback, transcript-status citation attachment, and source metadata before implementing grounded source packaging and deterministic safe answer text.

### P5-T06 - Flag stale or high-risk sources

- What: added stale-source and high-risk source confidence decisions, plus staff-confirmation behavior for risky source matches.
- Where: `internal/rag/freshness.go`, `internal/rag/freshness_test.go`, `internal/orchestrator/grounded_answer.go`, `docs/privacy-impact-lite.md`, `docs/model-evaluation.md`, `docs/test-plan.md`.
- When: 2026-05-06.
- Why: avoid authoritative answers from outdated or sensitive public content, especially for policy, fee, deadline, or eligibility-adjacent questions.
- How: wrote stale, high-risk confidence, and fresh-source tests first, then implemented threshold decisions and orchestrator caution metadata/action handling.

### P5 review evidence

- What: completed P5 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/model-evaluation.md`, `docs/demo-script.md`, `docs/source-references.md`, `docs/privacy-impact-lite.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API contract, architecture, demo script, source rules, privacy notes, and verification expectations aligned with implemented local RAG behavior.
- How: marked P5 tasks and gates complete after confirming the red test state, then verified `go test ./internal/rag ./internal/orchestrator ./internal/config ./cmd/api ./cmd/ingest`, `go test ./...`, `go vet ./...`, `make test`, and `git diff --check` pass.

## 2026-05-06 - P4 Deterministic Orchestration Before AI

### P4-T01 - Define orchestrator ports and dependency injection

- What: added the deterministic orchestrator package with small ports for classifier, retriever, LLM, Banner, payment, workflow, CRM, and audit dependencies, then wired `cmd/api` to the orchestrator.
- Where: `internal/orchestrator/orchestrator.go`, `internal/orchestrator/orchestrator_test.go`, `internal/audit/event.go`, `cmd/api/main.go`, `internal/config/config.go`, `internal/config/config_test.go`.
- When: 2026-05-06.
- Why: make chat orchestration testable with fakes and keep the API independent of live AI, retrieval, and raw HTTP tool details.
- How: wrote constructor and compile-time interface tests first, confirmed missing package/dependency failures, then implemented dependency validation, disabled retrieval/LLM ports, mock-service URL config, API wiring, and a no-op audit recorder.

### P4-T02 - Implement fallback intent and sentiment classifier

- What: added deterministic fallback classification for transcript request, transcript status, fee payment, human handoff, escalation, unknown intent, and neutral/negative/urgent sentiment.
- Where: `internal/classifier/fallback.go`, `internal/classifier/fallback_test.go`, `internal/domain/chat.go`.
- When: 2026-05-06.
- Why: provide reliable demo behavior before live LLM classification exists and prevent low-confidence results from triggering sensitive tool calls.
- How: wrote table-driven classifier and low-confidence threshold tests first, then implemented keyword-based intent/sentiment mapping and typed confidence gating.

### P4-T03 - Implement transcript-status decision flow

- What: added transcript-status orchestration for `S100001` ready/no reminder, `S100002` unpaid/reminder, `S100003` financial-hold escalation, `S100004` unresolved handoff, and missing synthetic ID prompts.
- Where: `internal/orchestrator/transcript.go`, `internal/orchestrator/orchestrator_test.go`, `docs/test-plan.md`, `docs/model-evaluation.md`.
- When: 2026-05-06.
- Why: make the core Tier 1 transcript/payment demo work from deterministic synthetic Banner and payment states.
- How: wrote fake-tool orchestrator tests for each synthetic record before implementation, then added Banner/payment decision branches, safe answers, source packaging, and no-tool behavior when the student ID is missing.

### P4-T04 - Trigger payment reminder workflow from orchestrator

- What: added an in-process idempotent payment-reminder workflow port and orchestrator workflow attempt/completion/failure action handling.
- Where: `internal/workflow/client.go`, `internal/workflow/client_test.go`, `internal/orchestrator/transcript.go`, `internal/orchestrator/orchestrator_test.go`, `docs/api-spec.md`, `docs/power-automate-flow.md`.
- When: 2026-05-06.
- Why: connect unpaid transcript decisions to workflow automation without waiting for the P8 standalone simulator or a real Power Automate webhook.
- How: wrote idempotency and workflow-failure tests first, then implemented stable idempotency keys, local synthetic workflow IDs, safe workflow-failure messages, and audit-port events for attempted/completed/failed workflow attempts.

### P4-T05 - Create CRM escalation summary and priority routing

- What: added CRM handoff creation for financial holds, unresolved synthetic transcript status, urgent/negative sentiment, human handoff, and low-confidence fallback.
- Where: `internal/orchestrator/escalation.go`, `internal/orchestrator/orchestrator_test.go`, `docs/demo-script.md`, `docs/architecture.md`.
- When: 2026-05-06.
- Why: turn unresolved or sensitive deterministic conversations into staff-ready mock CRM cases without automating approvals, denials, or financial judgments.
- How: wrote CRM routing and summary-redaction tests first, then implemented Registrar/Student Accounts routing, priority staff routing, normal learner-support handoff, safe learner messages, and redacted CRM summaries.

### P4-T06 - Return action trace in chat responses

- What: extended learner-facing chat actions with `trace_id` and `idempotency_key`, and returned deterministic action names/statuses for classifier, Banner, payment, workflow, and CRM steps.
- Where: `internal/domain/chat.go`, `internal/orchestrator/orchestrator.go`, `internal/orchestrator/transcript.go`, `docs/api-spec.md`, `README.md`.
- When: 2026-05-06.
- Why: make the interview demo explainable from the response body without exposing logs or raw internal errors.
- How: wrote action-trace assertions in orchestrator tests first, then attached trace IDs to every action, included workflow IDs/idempotency keys where relevant, and kept internal workflow/CRM errors out of answers and action messages.

### P4 review evidence

- What: completed P4 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/model-evaluation.md`, `docs/demo-script.md`, `docs/power-automate-flow.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API contract, demo script, workflow notes, and verification expectations aligned with deterministic P4 behavior and explicitly leave RAG, durable audit/dashboard, and standalone workflow simulation to later phases.
- How: marked P4 tasks and gates complete after confirming the red test state, then verified `go test ./internal/classifier ./internal/workflow ./internal/orchestrator`, `go test ./internal/config ./cmd/api ./internal/classifier ./internal/workflow ./internal/orchestrator`, `go test ./...`, `go vet ./...`, and `make test` pass.

## 2026-05-06 - P3 Synthetic Enterprise APIs And Clients

### P3-T01 - Create synthetic fixture loader

- What: added a validated synthetic fixture loader for learner, transcript, payment, CRM, and LMS demo records.
- Where: `internal/fixtures/loader.go`, `internal/fixtures/loader_test.go`, `data/synthetic-students.json`, `docs/privacy-impact-lite.md`.
- When: 2026-05-06.
- Why: let mock services and tests reuse one deterministic synthetic data source without touching real learner, payment, CRM, or LMS records.
- How: wrote failing loader tests for expected students, duplicate IDs, missing required fields, and synthetic-only validation, then implemented context-aware JSON loading and fixture validation.

### P3-T02 - Build mock Banner-style student API

- What: added the synthetic Banner-style HTTP service for student profile and transcript status/hold lookups.
- Where: `cmd/mock-banner/main.go`, `internal/mock/banner`, `internal/tools/banner_client.go`, `internal/tools/banner_client_test.go`, `docs/api-spec.md`.
- When: 2026-05-06.
- Why: provide deterministic student and transcript state for later orchestration without any real Banner integration or credentials.
- How: wrote handler, contract, client, trace-header, not-found, and malformed-response tests first, then implemented `net/http` handlers and the typed Banner client.

### P3-T03 - Build mock payment API

- What: added the synthetic transcript payment-status HTTP service and typed payment client.
- Where: `cmd/mock-payment/main.go`, `internal/mock/payment`, `internal/tools/payment_client.go`, `internal/tools/payment_client_test.go`, `docs/api-spec.md`, `docs/golang-implementation.md`.
- When: 2026-05-06.
- Why: simulate paid and unpaid transcript scenarios without accepting or storing real payment details.
- How: wrote tests for paid, unpaid, unknown-payment, response contract, and canceled-context timeout behavior, then implemented the handler and context-aware client.

### P3-T04 - Build mock CRM API for case creation and queue routing

- What: added the synthetic CRM case-creation service with queue, priority, case ID, and redacted summary output.
- Where: `cmd/mock-crm/main.go`, `internal/mock/crm`, `internal/tools/crm_client.go`, `internal/tools/crm_client_test.go`, `docs/api-spec.md`.
- When: 2026-05-06.
- Why: show human handoff and staff-facing routing without creating real cases or storing raw conversation PII.
- How: wrote tests for normal case creation, priority case creation, empty-summary rejection, redaction, and 5xx retryable client errors before implementing the in-memory mock CRM handler and typed client.

### P3-T05 - Build mock LMS API

- What: added the synthetic LMS access-status service and typed LMS client for demo account/course access checks only.
- Where: `cmd/mock-lms/main.go`, `internal/mock/lms`, `internal/tools/lms_client.go`, `internal/tools/lms_client_test.go`, `data/synthetic-students.json`, `docs/api-spec.md`.
- When: 2026-05-06.
- Why: support basic LMS access questions in the enterprise integration surface while keeping LMS course content, grades, submissions, and activity records out of scope.
- How: wrote tests for known synthetic access, unknown-course fallback, unknown student, and canceled-context timeout behavior, then implemented the mock LMS handler and client on local port 8085.

### P3-T06 - Add typed enterprise clients with shared error model

- What: added shared tool error kinds and HTTP client helpers for Banner, payment, CRM, and LMS integrations.
- Where: `internal/tools/errors.go`, `internal/tools/http.go`, `internal/tools/*_client.go`, `internal/tools/*_test.go`, `README.md`, `docs/test-plan.md`.
- When: 2026-05-06.
- Why: keep later orchestrator code independent from raw HTTP details while preserving branchable not-found, retryable, external-service, parse, and timeout errors.
- How: wrote `httptest.Server` client tests for trace ID forwarding, context cancellation, 404 mapping, 5xx mapping, and malformed JSON, then implemented context-aware clients that emit `X-Trace-ID`.

### P3 review evidence

- What: completed P3 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `docs/privacy-impact-lite.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API contracts, local commands, privacy assumptions, and verification expectations aligned with the implemented synthetic enterprise APIs and clients.
- How: marked P3 tasks and gates complete after confirming the red test state, reconciled the README demo table with the `mock_payment_hold` fixture/API contract, then verified `go test ./internal/fixtures ./internal/mock/banner ./internal/mock/payment ./internal/mock/crm ./internal/mock/lms ./internal/tools`, `go test ./cmd/mock-banner ./cmd/mock-payment ./cmd/mock-crm ./cmd/mock-lms`, `go test ./...`, and `go vet ./...` pass.

## 2026-05-06 - P2 Chat API And UI Skeleton

### P2-T01 - Define chat domain models

- What: added provider-neutral chat request/response domain models for intents, sentiment, source citations, action traces, and escalation metadata.
- Where: `internal/domain/chat.go`, `internal/domain/chat_test.go`, `docs/api-spec.md`, `docs/architecture.md`.
- When: 2026-05-06.
- Why: establish the stable API contract before real retrieval, AI orchestration, and enterprise integrations are added.
- How: wrote JSON round-trip and validation-facing tests first, confirmed the missing package failure, then implemented typed Go structs and constants.

### P2-T02 - Implement `POST /api/v1/chat` handler with deterministic placeholder response

- What: added the chat HTTP handler, chat service interface, deterministic placeholder service, trace ID propagation, strict JSON decoding, and safe service-error handling.
- Where: `internal/handlers/chat.go`, `internal/handlers/chat_test.go`, `cmd/api/main.go`, `docs/api-spec.md`, `docs/golang-implementation.md`.
- When: 2026-05-06.
- Why: expose the public chat contract without relying on live AI, RAG, workflow, or mock enterprise services.
- How: wrote `httptest` coverage for success, invalid JSON, missing message, unsupported method, validation errors, and safe service failure before implementing the handler and placeholder.

### P2-T03 - Serve minimal Go web chat UI

- What: added a Go-rendered chat page, static JavaScript, static CSS, and routes for `/chat`, `/`, and `/static/`.
- Where: `web/templates/chat.html`, `web/static/app.js`, `web/static/app.css`, `internal/handlers/ui.go`, `internal/handlers/ui_test.go`, `cmd/api/main.go`, `README.md`.
- When: 2026-05-06.
- Why: provide an interview-friendly usable chat surface while keeping frontend complexity low.
- How: wrote rendering and static-content tests first, confirmed missing handler/static asset failures, then added the template, static assets, and route registration.

### P2-T04 - Add request validation and safe client-facing errors

- What: added chat request validation for missing, whitespace-only, oversized messages and synthetic-only student ID shape.
- Where: `internal/validation/chat.go`, `internal/validation/chat_test.go`, `internal/handlers/chat.go`, `internal/handlers/chat_test.go`, `docs/api-spec.md`, `README.md`.
- When: 2026-05-06.
- Why: keep malformed requests and accidental non-demo identifiers from destabilizing the local demo or leaking raw input.
- How: added table-driven validation tests and handler assertions that error responses do not echo request bodies, then implemented safe validation codes and messages.

### P2-T05 - Add in-memory conversation session store

- What: added a concurrency-safe in-memory session store with configurable TTL, create/append/read/expire behavior, and redaction before message persistence.
- Where: `internal/session/store.go`, `internal/session/store_test.go`, `internal/handlers/chat.go`, `Makefile`, `docs/test-plan.md`, `README.md`.
- When: 2026-05-06.
- Why: track short synthetic demo conversations so follow-up behavior can be layered in later phases without storing raw PII.
- How: wrote create/append/read, expiration, redaction, and concurrent access tests first, confirmed missing package failures, then implemented the mutex-protected store and `make test-race` target.

### P2 review evidence

- What: completed P2 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/architecture.md`, `docs/golang-implementation.md`, `docs/test-plan.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API surface, local commands, and verification expectations aligned with the implemented chat skeleton.
- How: marked P2 tasks and gates complete after the red/green cycle and verified `go test ./internal/domain ./internal/validation ./internal/handlers ./internal/session`, `go test -race ./internal/session`, `go test ./...`, `go test -race ./...`, `make test`, `make test-race`, and `go vet ./...` pass.

## 2026-05-06 - P1 Go Project Foundation

### P1-T01 - Initialize Go module, repository layout, and developer commands

- What: created the initial Go module, API skeleton, build smoke test, and developer commands.
- Where: `go.mod`, `Makefile`, `cmd/api/main.go`, `internal/build/smoke_test.go`.
- When: 2026-05-06.
- Why: establish a compilable Go workspace and predictable local commands before feature work.
- How: added failing compile/Makefile smoke tests first, confirmed `go test ./...` failed for missing implementation, then added the minimal module, `make test`, and `make dev` targets.

### P1-T02 - Build typed configuration loader

- What: added typed environment configuration with defaults, overrides, validation, and redacted string output.
- Where: `internal/config/config.go`, `internal/config/config_test.go`, `README.md`.
- When: 2026-05-06.
- Why: make HTTP, auth, logging, workflow, and provider settings explicit without leaking secrets.
- How: added table-driven tests for defaults, overrides, invalid booleans/timeouts/log levels, and secret redaction, then implemented `Load`, `LoadFromEnv`, and safe formatting.

### P1-T03 - Implement health and readiness endpoints

- What: added `GET /healthz` and `GET /readyz` handlers with JSON responses, method checks, trace IDs, and safe dependency status reporting.
- Where: `internal/handlers/health.go`, `internal/handlers/health_test.go`, `docs/api-spec.md`, `README.md`.
- When: 2026-05-06.
- Why: provide operational endpoints for local demo startup checks and later Docker/deployment health checks.
- How: added `httptest` coverage for healthy responses, readiness dependency success/failure, 405 responses, trace propagation, and non-leaky dependency failures before implementing handlers.

### P1-T04 - Add middleware for trace ID, panic recovery, request logging, and mock auth

- What: added middleware chaining, trace ID propagation, panic recovery, mock bearer auth, request logging, and basic redaction.
- Where: `internal/middleware/chain.go`, `internal/middleware/trace.go`, `internal/middleware/recover.go`, `internal/middleware/auth.go`, `internal/middleware/logging.go`, `internal/middleware/*_test.go`, `docs/task-prompts.md`, `docs/phases-and-tasks.md`.
- When: 2026-05-06.
- Why: establish API hygiene before chat, orchestration, and mock integration endpoints are added.
- How: wrote failing tests for trace preservation/generation, safe panic conversion, auth enabled/disabled behavior, and logging redaction hooks, then implemented standard-library middleware.

### P1-T05 - Create JSON response and error helpers

- What: added response helpers for successful JSON payloads and stable safe API errors.
- Where: `internal/handlers/respond.go`, `internal/handlers/respond_test.go`, `internal/handlers/health.go`.
- When: 2026-05-06.
- Why: standardize response shape for later handlers and avoid leaking raw Go errors to clients.
- How: added tests for success headers/body, `{error:{code,message,trace_id}}` error shape, and unsupported-value fallback before implementing helper functions.

### P1 review evidence

- What: completed P1 status and documentation sync.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/api-spec.md`, `docs/task-prompts.md`, `CHANGELOG.md`.
- When: 2026-05-06.
- Why: keep the task board, roadmap, API surface, and local commands aligned with the implemented Go foundation.
- How: marked P1 tasks and gates complete after the red/green cycle and verified `go test ./...`, `make test`, and `go vet ./...` pass.

## 2026-05-06 - P0 Product Framing And Applicant Strategy

### P0-T01 - Write the applicant story and MVP thesis

- What: tightened the applicant story and one-sentence pitch around a Go-based learner-service agent using RAG, workflow automation, mock enterprise integrations, and privacy-first synthetic data controls.
- Where: `README.md`, `docs/mvp-scope.md`.
- When: 2026-05-06.
- Why: make the project purpose, learner pain point, and AI/Automation Solutions Developer role mapping explicit before implementation starts.
- How: updated the README pitch and scope language, froze the primary transcript/payment workflow, and removed wording that could imply access to real OC systems.

### P0-T02 - Define synthetic data and privacy boundary

- What: documented the synthetic-data-only boundary and created the synthetic learner fixture.
- Where: `docs/privacy-impact-lite.md`, `data/synthetic-students.json`, `README.md`.
- When: 2026-05-06.
- Why: make it clear that learner records, IDs, payments, transcript states, and CRM cases are invented demo artifacts only.
- How: added fake-data markers, fixture rules, synthetic ID patterns, and four visibly fake demo records using `S10000X`, `SYNTH-*`, and `MOCK-*` identifiers.

### P0-T03 - Confirm source allowlist and knowledge-domain limits

- What: created the public source allowlist and documented retrieval boundaries.
- Where: `data/seed-sources.json`, `docs/source-references.md`, `docs/privacy-impact-lite.md`.
- When: 2026-05-06.
- Why: prevent private portal scraping, unapproved source ingestion, stale-source overconfidence, and learner-specific data leakage.
- How: listed approved public Okanagan College URLs already present in repo docs, spot-checked them as accessible public pages on 2026-05-06, separated implementation references from learner-service RAG sources, added freshness metadata, and defined fallback behavior for stale or missing sources.

### P0-T04 - Create demo acceptance matrix

- What: turned the interview demo into measurable acceptance scenarios.
- Where: `docs/demo-script.md`, `docs/model-evaluation.md`.
- When: 2026-05-06.
- Why: ensure the transcript answer, unpaid payment workflow, financial-hold escalation, and urgent sentiment escalation have observable pass criteria before code implementation.
- How: added D01-D05 demo cases with expected intent, source, action, handoff behavior, and pass evidence; aligned source checks with the source allowlist fixture.

### P0-T05 - Freeze MVP scope and defer nice-to-haves

- What: froze the MVP around transcript/payment support and deferred nonessential workflows and real integrations.
- Where: `docs/mvp-scope.md`, `docs/implementation-roadmap.md`, `README.md`, `docs/phases-and-tasks.md`.
- When: 2026-05-06.
- Why: keep the MVP narrow enough for strict TDD delivery and avoid overbuilding beyond the applicant demo.
- How: marked P0 complete, updated phase gates, clarified that non-transcript learner-service topics use fallback or handoff, and documented that real authentication, Banner, payment, CRM, LMS, and private portal integrations are out of scope.

### Review evidence

- What: completed documentation review checks for all P0 tasks.
- Where: `docs/phases-and-tasks.md`, `docs/implementation-roadmap.md`, `data/synthetic-students.json`, `data/seed-sources.json`, and the updated P0 Markdown files.
- When: 2026-05-06.
- Why: P0 contains documentation tasks only, so Go failing-test evidence and `go test ./...` are not applicable.
- How: used JSON validation, public URL spot-checks, source/fixture consistency checks, targeted content searches for required and prohibited terms, and Markdown whitespace checks.
