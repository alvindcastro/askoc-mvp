# AskOC Web App Revamp TDD Plan

This plan is for the `feat/ui-vamp` revamp branch. It turns `DESIGN.md` into a practical web-app execution plan for the existing **AskOC AI Concierge** MVP.

`DESIGN.md` names the visual system "VoiceBox." In this repo, treat that as visual language only: high contrast, editorial typography, square borders, flat surfaces, and restrained red accents. Do not turn AskOC into an editorial, blog, CMS, or marketing product.

## Source Inputs

- [ ] Review `DESIGN.md` before changing UI docs or code.
- [ ] Keep AskOC centered on transcript/payment support, synthetic enterprise checks, workflow reminders, CRM escalation, audit evidence, and evaluation gates.
- [ ] Use only synthetic learner, payment, workflow, CRM, LMS, and audit data.
- [ ] Preserve the current Go server-rendered app shape unless a task explicitly adds a tested route.
- [ ] Follow [Strict TDD Policy](tdd-policy.md) for every code task.

## Strict TDD Loop

Every code prompt in this plan must follow this sequence:

1. Write or update a failing test for the requested behavior.
2. Run the narrow test and confirm it fails for the expected missing behavior.
3. Implement the smallest production change.
4. Run the narrow test until green.
5. Run the broader relevant suite, normally `go test ./...`.
6. Refactor only while tests stay green.
7. Update related Markdown when behavior, routes, commands, fixtures, or UX assumptions change.

Do not mark a code checklist item done unless the red and green evidence exists.

## Web App Brainstorm

The revamp should make AskOC feel like a dense operational console for a learner-service automation demo. The first screen should be the actual chat/workflow surface, not a landing page.

The UI should help a reviewer scan:

- [ ] what the learner asked,
- [ ] what grounded answer was returned,
- [ ] which source citations, confidence, risk, and freshness states apply,
- [ ] which synthetic systems were checked,
- [ ] which workflow or CRM action happened,
- [ ] what trace ID ties the chat to audit/evaluation evidence,
- [ ] why no real learner data or live enterprise systems are involved.

Design direction:

- [ ] Use black `#0A0A0A`, white `#FAFAFA`, red `#EF4444`, and neutral surfaces.
- [ ] Use `Archivo Black` for major headings, `Work Sans` for UI/body, and `Space Mono` for trace IDs, synthetic IDs, route names, and action codes.
- [ ] Use 0px radius, 2px borders, no shadows, no gradients, no decorative blobs, and no rounded dashboard cards.
- [ ] Use red sparingly for active navigation, critical/error states, and hover emphasis.
- [ ] Make chat, source citations, action traces, escalation, admin metrics, and review rows compact and keyboard accessible.
- [ ] Prefer visible proof over explanation text: source chips, action rows, status chips, trace IDs, and redaction markers.

## Current Route And Screen Inventory

- [x] `GET /` and `GET /chat` render the learner chat shell.
- [x] `POST /api/v1/chat` accepts the learner chat JSON request from the web UI.
- [x] `GET /admin` renders the protected admin dashboard shell.
- [x] `GET /api/v1/admin/metrics` returns protected aggregate metrics.
- [x] `GET /api/v1/admin/audit/export` exports protected audit data.
- [x] `POST /api/v1/admin/audit/reset` resets demo audit data.
- [x] `POST /api/v1/admin/audit/purge` purges expired audit data.
- [x] `GET /api/v1/admin/review-items` returns protected review queue items.
- [x] `GET /static/*` serves `web/static`.
- [x] `GET /healthz` and `GET /readyz` expose service status.

Current screen surfaces:

- [ ] Learner chat: conversation stream, composer, synthetic student ID input, intent, sources, actions, escalation.
- [ ] Chat response details: currently too thin for the P11 proof story; source metadata and action details need richer rendering.
- [ ] Admin dashboard: token form, summary metrics, top intents, workflow counts, low-confidence review list, stale-source count, audit controls.
- [ ] Static styling: split between `app.css` and `admin.css`; no shared theme-token layer exists yet.

Known consistency risks handled during the revamp:

- [x] Chat UI copy no longer references a P2 placeholder; the shell describes the current synthetic transcript/payment MVP.
- [x] Chat rendering now covers source confidence, risk, freshness, workflow IDs, CRM case IDs, priority, idempotency keys, and trace proof.
- [x] Admin review rendering keeps a safe empty state and preserves useful metrics-derived review proof when the eval queue is empty.
- [x] `docs/api-spec.md` mentions routes such as feedback that are not wired in `cmd/api`; no UI control was added for missing endpoints.
- [x] `/` remains an intentional chat alias because `cmd/api/main.go` already maps `/` and `/chat` to `ChatPageHandler`.

## 2026-05-09 Implementation Evidence

- Red test: `go test ./internal/handlers -run 'Test(Chat|Admin).*Revamp|Test(Chat|Admin)StaticAssets'` failed first on missing nav landmarks, stale P2 copy, missing theme tokens, and missing evidence rendering contracts.
- Green tests: the same narrow command passes after template/static changes; `go test ./internal/handlers` also passes.
- Changed runtime files: `web/templates/chat.html`, `web/templates/admin.html`, `web/static/app.css`, `web/static/app.js`, `web/static/admin.css`, `web/static/admin.js`.
- Changed test files: `internal/handlers/ui_test.go`, `internal/handlers/admin_ui_test.go`.
- Changed documentation files: `README.md`, `docs/demo-script.md`, `docs/test-plan.md`, `docs/askoc-ux-theme-brainstorm.md`, `docs/web-app-revamp-tdd-plan.md`, `CHANGELOG.md`.
- Full verification: `go test ./...`, `make eval`, `make secret-check`, `git diff --check`, and an alternate `9xxx` `ASKOC_API_PORT` smoke run passed; default `make smoke` was blocked by an occupied API port.
- Manual visual checklist: inspect `/chat` and `/admin` at mobile and desktop widths for text overlap, hidden source/action rows, clipped controls, dominant red usage, focus visibility, and accidental raw learner data.

## Phase R0 - Baseline, Scope, And Evidence

**Outcome:** The revamp branch has a confirmed baseline, route inventory, and evidence contract before code changes.

- [x] **R0-T01 - Reconcile current web app status**
  - **Type:** Documentation
  - **Goal:** Confirm the existing route/screen inventory, stale UI copy, and proof gaps before implementation.
  - **Primary files:** `docs/web-app-revamp-tdd-plan.md`, `docs/askoc-ux-theme-brainstorm.md`
  - **Prompt:**

```text
You are implementing R0-T01 for the AskOC web app revamp.

Task goal: Reconcile the current web app status before code changes.

Instructions:
1. Inspect cmd/api/main.go, internal/handlers, web/templates, web/static, README.md, and docs/demo-script.md.
2. Confirm the route inventory and screen inventory.
3. List stale UI copy and proof gaps without editing code.
4. Keep the product scope AskOC AI Concierge; do not add new product routes.

Review checks:
- Existing routes are listed accurately.
- Placeholder/stale copy is called out.
- Missing UI proof points are tied to current API/demo expectations.
- No real learner data, secrets, or private URLs are added.

Expected response:
1. Files inspected.
2. Route and screen findings.
3. Consistency risks.
4. Docs changed.
```

  - **Completion checklist:**
    - [ ] Route inventory is confirmed.
    - [ ] Stale copy and proof gaps are listed.
    - [ ] No new route is implied without a code task.
    - [ ] No real data or secrets are added.

- [x] **R0-T02 - Freeze revamp task taxonomy**
  - **Type:** Documentation
  - **Goal:** Keep one revamp prompt/task structure across docs, code, verification, and demo updates.
  - **Primary files:** `docs/web-app-revamp-tdd-plan.md`, `docs/askoc-ux-tdd-prompts.md`
  - **Prompt:**

```text
You are implementing R0-T02 for the AskOC web app revamp.

Task goal: Freeze the revamp task taxonomy and make every phase prompt actionable.

Instructions:
1. Keep task IDs in the R0-R5 format for revamp work.
2. Ensure each task has type, goal, primary files, prompt, and tickable completion checklist.
3. For every code task, include the strict red-first TDD loop.
4. For documentation tasks, define review checks before editing.

Review checks:
- Every task maps to observable evidence.
- Code tasks require red test evidence before production changes.
- Documentation tasks have review-first checks.
- Existing AskOC scope is preserved.
```

  - **Completion checklist:**
    - [ ] Every revamp task has a prompt.
    - [ ] Every code task requires red-first TDD.
    - [ ] Every documentation task has review checks.
    - [ ] The task list remains tickable.

## Phase R1 - Design System Foundation

**Outcome:** `DESIGN.md` is mapped into reusable AskOC theme tokens and app-shell rules.

- [x] **R1-T01 - Map DESIGN.md tokens to AskOC surfaces**
  - **Type:** Documentation
  - **Goal:** Define token usage for chat, action trace, admin, review, and eval surfaces.
  - **Primary files:** `docs/askoc-ux-theme-brainstorm.md`, `docs/web-app-revamp-tdd-plan.md`
  - **Prompt:**

```text
You are implementing R1-T01 for the AskOC web app revamp.

Task goal: Map DESIGN.md tokens to AskOC surfaces.

Instructions:
1. Read DESIGN.md and existing UX docs.
2. Treat VoiceBox as visual language only.
3. Map colors, type, borders, chips, buttons, inputs, focus rings, and mono metadata to existing AskOC surfaces.
4. Explicitly reject gradients, shadows, rounded dashboard cards, decorative hero sections, and editorial/blog scope.

Review checks:
- Chat, source citations, action trace, escalation, admin metrics, review rows, and eval evidence are covered.
- AskOC remains the product.
- No unsupported route or product claim is added.
```

  - **Completion checklist:**
    - [ ] Token mapping covers all current surfaces.
    - [ ] Forbidden visual patterns are listed.
    - [ ] VoiceBox is framed as theme only.
    - [ ] AskOC scope is unchanged.

- [x] **R1-T02 - Add shared theme tokens to static CSS**
  - **Type:** Code
  - **Goal:** Add reusable CSS variables/classes for the `DESIGN.md` system.
  - **Primary files:** `web/static/app.css`, `web/static/admin.css`, UI/template tests
  - **Prompt:**

```text
You are implementing R1-T02 for the AskOC web app revamp.

Task goal: Add shared DESIGN.md theme tokens to the static CSS.

Strict TDD:
1. Write a failing test first for required token names, square radius, focus ring, or themed class usage.
2. Run the narrow test and confirm it fails because the token/class behavior is missing.
3. Implement the smallest CSS/template change.
4. Run the narrow test until green.
5. Run go test ./... before marking done.

Requirements:
- Include black #0A0A0A, white #FAFAFA, red #EF4444, neutral surfaces, 0px radius, 2px borders, no shadows, no gradients, and the DESIGN.md focus ring.
- Keep CSS lightweight; do not add a frontend framework.
- Preserve existing routes and API behavior.

Expected response:
1. Red test name and failure.
2. CSS/template changes.
3. Test commands run.
4. Confirmation that API behavior did not change.
```

  - **Completion checklist:**
    - [ ] Red test exists and failed first.
    - [ ] Theme tokens exist.
    - [ ] No shadows, gradients, or rounded panels are introduced.
    - [ ] `go test ./...` passes.

- [x] **R1-T03 - Theme app shell and navigation**
  - **Type:** Code
  - **Goal:** Make the app shell clear, compact, keyboard accessible, and consistent across chat/admin.
  - **Primary files:** `web/templates/chat.html`, `web/templates/admin.html`, static CSS/JS, handler/template tests
  - **Prompt:**

```text
You are implementing R1-T03 for the AskOC web app revamp.

Task goal: Theme the app shell and route navigation.

Strict TDD:
1. Write failing render tests for nav landmark, active route, AskOC product title, synthetic-mode label, and keyboard-accessible links/buttons.
2. Verify the tests fail for the expected missing behavior.
3. Implement the smallest template/CSS/JS changes.
4. Run the narrow tests, then go test ./...

Requirements:
- Header must make AskOC AI Concierge the product.
- Active route uses restrained red underline.
- Synthetic mode is visible.
- Do not add marketing hero or editorial navigation.
- Decide whether `/` remains a chat alias or becomes explicit redirect/404 only in a tested route task.
```

  - **Completion checklist:**
    - [ ] Header/nav tests were red first.
    - [ ] Active route and synthetic mode render.
    - [ ] Keyboard access is covered.
    - [ ] Existing AskOC routes are preserved.

## Phase R2 - Chat Evidence Surface

**Outcome:** The learner chat proves source grounding, synthetic integrations, workflow action, escalation, and traceability at a glance.

- [x] **R2-T01 - Replace stale placeholder copy**
  - **Type:** Code
  - **Goal:** Remove P2 placeholder language and align empty states with the current orchestrated MVP.
  - **Primary files:** `web/templates/chat.html`, `internal/handlers/ui_test.go`, docs if copy changes demo expectations
  - **Prompt:**

```text
You are implementing R2-T01 for the AskOC web app revamp.

Task goal: Replace stale placeholder UI copy with current MVP copy.

Strict TDD:
1. Write a failing render test that rejects "P2 placeholder" and deterministic-placeholder copy.
2. Write or update a test that requires current MVP language: synthetic demo mode, transcript/payment support, source/action proof, and safe data boundary.
3. Verify red state.
4. Implement the smallest template copy change.
5. Run the narrow handler/template test, then go test ./...

Requirements:
- Do not change chat API behavior.
- Do not overclaim live AI, real system integration, or production deployment.
- Keep the first screen usable as the chat workflow surface.
```

  - **Completion checklist:**
    - [ ] Red test catches stale copy.
    - [ ] Empty state reflects current MVP behavior.
    - [ ] No API behavior changes.
    - [ ] Broader tests pass.

- [x] **R2-T02 - Theme learner and assistant message layout**
  - **Type:** Code
  - **Goal:** Make conversation turns compact, readable, responsive, and safe.
  - **Primary files:** `web/templates/chat.html`, `web/static/app.css`, `web/static/app.js`, UI tests
  - **Prompt:**

```text
You are implementing R2-T02 for the AskOC web app revamp.

Task goal: Theme learner and assistant messages.

Strict TDD:
1. Write failing UI/render tests for learner role, assistant role, empty state, synthetic marker, safe error state, and no text overlap hooks.
2. Verify red state.
3. Implement the smallest template/CSS/JS change.
4. Run narrow tests, then go test ./...

Requirements:
- Message layout stays compact and scan-friendly.
- Synthetic-demo safety remains visible.
- Raw internal errors are never displayed.
- Mobile and desktop layouts must not overlap.
```

  - **Completion checklist:**
    - [ ] Message tests were red first.
    - [ ] Role and empty/error states are covered.
    - [ ] Synthetic marker is visible.
    - [ ] Chat API contract is unchanged.

- [x] **R2-T03 - Surface sources, confidence, risk, and freshness**
  - **Type:** Code
  - **Goal:** Show grounded-answer evidence and fallback states clearly.
  - **Primary files:** `web/static/app.js`, `web/static/app.css`, `web/templates/chat.html`, handler/JS tests where available
  - **Prompt:**

```text
You are implementing R2-T03 for the AskOC web app revamp.

Task goal: Surface sources, confidence, risk, freshness, and fallback states.

Strict TDD:
1. Write failing tests for source citation rendering, confidence label, stale-source warning, high-risk marker, low-confidence fallback, and no-source escalation state.
2. Verify red state.
3. Implement the smallest rendering change.
4. Run narrow tests, then go test ./...

Requirements:
- Citations are compact and source-grounded.
- Warning states use semantic styling; red is not overused.
- No private URLs, real learner data, or unsupported source claims are displayed.
```

  - **Completion checklist:**
    - [ ] Source/confidence tests were red first.
    - [ ] Citation and fallback states render.
    - [ ] Private URLs are not introduced.
    - [ ] Broader tests pass.

- [x] **R2-T04 - Surface action trace, workflow, CRM, and trace IDs**
  - **Type:** Code
  - **Goal:** Render synthetic integration proof as scan-friendly action rows.
  - **Primary files:** `web/static/app.js`, `web/static/app.css`, `web/templates/chat.html`, UI tests
  - **Prompt:**

```text
You are implementing R2-T04 for the AskOC web app revamp.

Task goal: Surface action trace, workflow, CRM, and trace IDs.

Strict TDD:
1. Write failing tests for action row rendering, status chips, trace ID display, workflow ID, CRM case ID, priority queue, idempotency key if present, and safe hidden internals.
2. Verify red state.
3. Implement the smallest rendering change.
4. Run narrow tests, then go test ./...

Requirements:
- Trace IDs, synthetic IDs, and action codes use Space Mono.
- Mock Banner/payment/workflow/CRM actions are clearly marked synthetic.
- Rows use border-led list styling.
- Raw payloads, secrets, and internal errors are not exposed.
```

  - **Completion checklist:**
    - [ ] Action trace tests were red first.
    - [ ] Trace ID and action statuses render.
    - [ ] Workflow and CRM proof is visible when present.
    - [ ] Raw internals are not exposed.

## Phase R3 - Admin, Audit, And Evaluation Evidence

**Outcome:** The admin surface exposes aggregate proof, review queues, and evaluation status without leaking raw learner content.

- [x] **R3-T01 - Theme admin metrics dashboard**
  - **Type:** Code
  - **Goal:** Apply the `DESIGN.md` operational layout to metrics and audit controls.
  - **Primary files:** `web/templates/admin.html`, `web/static/admin.css`, `web/static/admin.js`, `internal/handlers/admin_ui_test.go`
  - **Prompt:**

```text
You are implementing R3-T01 for the AskOC web app revamp.

Task goal: Theme the admin metrics dashboard.

Strict TDD:
1. Write failing tests for metric labels, synthetic-mode marker, token form labels, redacted values, audit controls, loading state, empty state, and error state.
2. Verify red state.
3. Implement the smallest template/CSS/JS changes.
4. Run narrow tests, then go test ./...

Requirements:
- Use dense operational bands/lists, not rounded cards.
- No raw learner messages are shown.
- Metrics remain readable on mobile and desktop.
- Red is reserved for critical/error emphasis.
```

  - **Completion checklist:**
    - [ ] Dashboard tests were red first.
    - [ ] Metrics and empty/error states are covered.
    - [ ] Redaction is preserved.
    - [ ] Broader tests pass.

- [x] **R3-T02 - Theme audit and review queue rows**
  - **Type:** Code
  - **Goal:** Make trace-linked audit/review items easy to scan and filter.
  - **Primary files:** `web/static/admin.js`, `web/static/admin.css`, admin tests
  - **Prompt:**

```text
You are implementing R3-T02 for the AskOC web app revamp.

Task goal: Theme audit and review queue rows.

Strict TDD:
1. Write failing tests for trace ID, queue, priority, status chip, redaction marker, filter control, safe empty state, and no overwrite of useful audit-derived review evidence by an empty eval queue.
2. Verify red state.
3. Implement the smallest rendering/styling changes.
4. Run narrow tests, then go test ./...

Requirements:
- Rows use bottom borders and compact metadata.
- Priority/error states use semantic styling.
- Trace IDs use Space Mono.
- Redacted data markers remain visible.
- Review behavior is unchanged unless a separate tested behavior task says otherwise.
```

  - **Completion checklist:**
    - [ ] Audit/review tests were red first.
    - [ ] Trace, queue, priority, and status render.
    - [ ] Redaction markers are visible.
    - [ ] Review semantics are unchanged.

- [x] **R3-T03 - Surface evaluation pass/fail evidence**
  - **Type:** Code or Documentation
  - **Goal:** Make evaluation gate evidence visible and honest.
  - **Primary files:** `reports/eval-summary.md`, `docs/test-plan.md`, admin UI if surfaced there
  - **Prompt:**

```text
You are implementing R3-T03 for the AskOC web app revamp.

Task goal: Surface evaluation pass/fail evidence.

Strict TDD or review-first:
1. If changing code, write a failing render/report test first.
2. If changing only docs/reports, define the review check before editing.
3. Verify the red/review state.
4. Implement the smallest change.
5. Run relevant tests or review checks.

Requirements:
- Critical pass/fail status is prominent.
- Warnings remain honest and visible.
- Report content remains synthetic and redacted.
- Do not hide known limitations or overstate production readiness.
```

  - **Completion checklist:**
    - [ ] Red test or review-first check happened.
    - [ ] Pass/fail evidence is visible.
    - [ ] Warnings are not hidden.
    - [ ] Claims are not overstated.

## Phase R4 - Accessibility, Responsive, And Visual Gates

**Outcome:** The themed MVP is usable on desktop/mobile and has explicit visual-quality checks.

- [x] **R4-T01 - Add accessibility checks for themed controls**
  - **Type:** Code
  - **Goal:** Verify landmarks, labels, focus order, and keyboard behavior.
  - **Primary files:** UI tests, templates, static JS/CSS
  - **Prompt:**

```text
You are implementing R4-T01 for the AskOC web app revamp.

Task goal: Add accessibility checks for themed controls.

Strict TDD:
1. Write a failing accessibility or keyboard test first.
2. Verify it fails for a real missing label, landmark, focus state, or keyboard behavior.
3. Implement the smallest fix.
4. Run accessibility/keyboard tests, then go test ./...

Required coverage:
- chat input and submit
- source/action controls
- admin token and audit controls
- review queue controls
- navigation
- visible focus ring
```

  - **Completion checklist:**
    - [ ] Accessibility test was red first.
    - [ ] Key controls are covered.
    - [ ] Keyboard behavior is covered.
    - [ ] Broader tests pass.

- [x] **R4-T02 - Add responsive layout checks**
  - **Type:** Code
  - **Goal:** Verify chat, action trace, dashboard, and review rows do not overlap.
  - **Primary files:** responsive CSS, visual/e2e tests if available, manual screenshot notes
  - **Prompt:**

```text
You are implementing R4-T02 for the AskOC web app revamp.

Task goal: Add responsive layout checks.

Strict TDD:
1. Add a failing viewport, screenshot, or layout assertion first where tooling exists.
2. Verify red state.
3. Implement the smallest responsive CSS/layout change.
4. Run viewport checks and go test ./...

Required checks:
- mobile chat
- mobile action trace
- desktop chat
- desktop admin dashboard
- desktop review/eval evidence

If no automated browser tooling exists, document the manual screenshot checklist before editing CSS.
```

  - **Completion checklist:**
    - [ ] Responsive test/check was red first where tooling exists.
    - [ ] Mobile and desktop states are covered.
    - [ ] No overlap or overflow is observed.
    - [ ] Broader tests pass.

- [x] **R4-T03 - Add contrast and theme drift gate**
  - **Type:** Code or Documentation
  - **Goal:** Prevent low contrast, too much red, gradients, shadows, and rounded UI drift.
  - **Primary files:** CSS tests, `docs/test-plan.md`, manual review notes
  - **Prompt:**

```text
You are implementing R4-T03 for the AskOC web app revamp.

Task goal: Add contrast and theme drift gates.

Strict TDD or review-first:
1. If changing code, write a failing style/visual assertion first.
2. If changing docs, define a manual review checklist first.
3. Verify red/review state.
4. Implement the smallest change.
5. Run tests or review checks.

Review requirements:
- Red is not dominant.
- Text contrast is high.
- Focus states remain visible.
- No shadows, gradients, rounded panels, decorative blobs, or marketing hero patterns are introduced.
```

  - **Completion checklist:**
    - [ ] Red test or review-first check happened.
    - [ ] Red accent usage is restrained.
    - [ ] Contrast and focus states are checked.
    - [ ] Forbidden visual patterns are absent.

## Phase R5 - Reviewer Path And Final Evidence

**Outcome:** The README/demo/test docs explain the revamp as UX polish for AskOC and record exact evidence.

- [x] **R5-T01 - Update README reviewer path for revamp**
  - **Type:** Documentation
  - **Goal:** Explain themed UI as polish for the existing AskOC MVP.
  - **Primary files:** `README.md`
  - **Prompt:**

```text
You are implementing R5-T01 for the AskOC web app revamp.

Task goal: Update the README reviewer path for the themed MVP.

Instructions:
1. Keep README led by AskOC AI Concierge.
2. Mention DESIGN.md only as visual theme guidance.
3. Keep transcript/payment workflow, synthetic integrations, TDD, privacy, and evaluation as the main story.
4. Do not describe VoiceBox as a separate app.

Review checks:
- README still explains the AskOC MVP in under two minutes.
- Theme language does not replace product language.
- No unsupported production claims are added.
- No real data, secrets, or private URLs are added.
```

  - **Completion checklist:**
    - [ ] README remains AskOC-first.
    - [ ] Theme is framed as visual polish.
    - [ ] No unsupported claims are added.
    - [ ] No real data or secrets are added.

- [x] **R5-T02 - Update demo script for themed screens**
  - **Type:** Documentation
  - **Goal:** Add concise visual callouts for chat, source chips, action trace, dashboard, and eval evidence.
  - **Primary files:** `docs/demo-script.md`
  - **Prompt:**

```text
You are implementing R5-T02 for the AskOC web app revamp.

Task goal: Update the demo script for themed screens.

Instructions:
1. Keep the demo centered on transcript answer, payment workflow, CRM escalation, dashboard evidence, and TDD/eval proof.
2. Add where to point: themed chat, source citations, action trace, synthetic labels, admin metrics, review rows, and eval evidence.
3. Keep the demo 5-7 minutes.
4. Do not add a separate editorial/blog demo.

Review checks:
- Every spoken UI claim maps to a visible screen, test, report, or fixture.
- Demo remains time-boxed.
- AskOC product story remains unchanged.
```

  - **Completion checklist:**
    - [ ] Demo remains AskOC-centered.
    - [ ] UI callouts map to visible evidence.
    - [ ] Demo remains time-boxed.
    - [ ] No separate product story is introduced.

- [x] **R5-T03 - Record final UX quality evidence**
  - **Type:** Documentation
  - **Goal:** Summarize tests, screenshots/manual checks, limitations, and out-of-scope UI work.
  - **Primary files:** `CHANGELOG.md`, `docs/test-plan.md`
  - **Prompt:**

```text
You are implementing R5-T03 for the AskOC web app revamp.

Task goal: Record final UX quality evidence.

Instructions:
1. Update changelog/test-plan docs after implementation evidence exists.
2. Include TDD, accessibility, responsive, redaction, synthetic-data, no-secret, eval, and smoke evidence.
3. Keep limitations honest.
4. Do not claim production deployment or real system integration.

Review checks:
- Evidence lists exact commands or manual checks.
- Known gaps are not hidden.
- No real data, secrets, or private URLs are added.
```

  - **Completion checklist:**
    - [ ] Evidence lists commands/checks.
    - [ ] Known gaps are documented.
    - [ ] No unsupported production claims are made.
    - [ ] No real data, secrets, or private URLs are added.

## Out Of Scope For This Revamp

- [ ] Replacing the AskOC transcript/payment workflow.
- [ ] Adding editorial/blog/CMS behavior.
- [ ] Adding marketing landing pages.
- [ ] Adding real authentication, real Banner, real payment, real CRM, real LMS, or real learner records.
- [ ] Adding UI controls for API routes that are not implemented and tested.
- [ ] Changing API behavior solely for visual polish.
- [ ] Hiding known evaluation, redaction, or source-grounding limitations.
