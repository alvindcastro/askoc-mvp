# AskOC UX Theme TDD Prompts

Use these prompts to apply `DESIGN.md` as the UX theme for the existing AskOC AI Concierge MVP. The product remains the Go-based learner-service automation concierge; the theme work should improve the chat, workflow evidence, dashboard, and demo surfaces without changing the core scope.

## Global Context To Include With Every Prompt

```text
Project: AskOC AI Concierge MVP.
Purpose: Go-based learner-service automation concierge for transcript/payment questions, synthetic enterprise checks, workflow reminders, CRM escalation, redacted audit, dashboard evidence, and evaluation gates.
UX source: DESIGN.md.
Theme rule: Apply DESIGN.md as visual language only. Do not turn this repository into a separate VoiceBox/editorial product.
Design system: black #0A0A0A, white #FAFAFA, red #EF4444, neutral surfaces, Archivo Black headings, Work Sans UI/body, Space Mono metadata, 0px radius, no shadows, no gradients, thick borders, restrained red accent.
Data rule: Use synthetic learner and synthetic integration data only. Do not add real learner data, real payment data, private portal content, secrets, analytics tokens, or private URLs.
TDD rule: For code tasks, write a failing test first, verify the failure, implement the smallest production change, run the narrow test, then run the broader suite.
Accessibility rule: Every interactive control must have an accessible name, keyboard behavior, and the DESIGN.md focus ring.
Documentation rule: Update Markdown when routes, commands, fixtures, UX assumptions, or test evidence change.
```

## UX0-T01 - Map DESIGN.md Tokens To AskOC UI

**Type:** Documentation
**Goal:** Define how colors, typography, borders, chips, buttons, inputs, and focus rings apply to AskOC screens.

### Copy/paste prompt

```text
You are implementing UX0-T01 for AskOC AI Concierge.

Task title: Map DESIGN.md tokens to AskOC UI
Task goal: Define how DESIGN.md maps to chat, admin, eval, and demo surfaces.

Instructions:
1. Read DESIGN.md and the current AskOC docs before editing.
2. Keep AskOC as the MVP product.
3. Treat VoiceBox naming in DESIGN.md as visual-theme language only.
4. Map tokens to existing surfaces: chat, action trace, source citations, synthetic mode labels, admin metrics, review queue, evaluation evidence.
5. Do not introduce a new editorial/blog product scope.

Primary files:
- docs/askoc-ux-theme-brainstorm.md
- docs/askoc-ux-tdd-prompts.md

Review checks:
- The docs still identify AskOC AI Concierge as the product.
- Every UX rule maps to an existing MVP surface.
- No route or task implies CMS, blog publishing, comments, or real analytics.

Expected response:
1. List review checks performed.
2. Summarize docs changed.
3. Confirm no product-scope drift was introduced.
```

### Completion checklist

- [ ] Review checks were performed.
- [ ] AskOC remains the product.
- [ ] Theme decisions map to existing MVP surfaces.
- [ ] No new editorial/blog product scope was introduced.

## UX0-T02 - Freeze Themed Route Inventory

**Type:** Documentation
**Goal:** List existing MVP routes and their themed UX jobs.

### Copy/paste prompt

```text
You are implementing UX0-T02 for AskOC AI Concierge.

Task title: Freeze themed route inventory
Task goal: List the existing MVP routes and define what each themed screen must prove.

Instructions:
1. Inspect existing handlers/templates before editing docs.
2. List only routes that exist or are already part of the AskOC MVP plan.
3. Describe the UX job for chat, admin/dashboard, review/eval evidence, and static assets.
4. Do not add editorial article/archive/draft routes.

Primary files:
- docs/askoc-ux-theme-brainstorm.md

Review checks:
- Route inventory matches AskOC MVP behavior.
- No route implies a separate VoiceBox product.
- Synthetic mode and traceability are represented in route jobs.

Expected response:
1. List route review checks.
2. Summarize route inventory changes.
3. Call out any unknown routes that need source-code confirmation.
```

### Completion checklist

- [ ] Route inventory is AskOC-specific.
- [ ] Route jobs are testable.
- [ ] Synthetic mode is covered.
- [ ] No separate product routes were added.

## UX0-T03 - Define UX Acceptance Evidence

**Type:** Documentation
**Goal:** Define tests, screenshots, and manual checks that prove the theme works.

### Copy/paste prompt

```text
You are implementing UX0-T03 for AskOC AI Concierge.

Task title: Define UX acceptance evidence
Task goal: Define what tests and review checks prove the themed AskOC MVP works.

Instructions:
1. Keep strict TDD language for every code task.
2. Include accessibility, responsive, redaction, synthetic-data labeling, and traceability checks.
3. Include documentation checks for README/demo updates.
4. Do not require live external services or real learner data.

Primary files:
- docs/askoc-ux-theme-brainstorm.md
- docs/test-plan.md if test-plan evidence changes

Review checks:
- Every UX phase has observable evidence.
- Code evidence starts with failing tests.
- Documentation evidence is explicit for non-code tasks.

Expected response:
1. List evidence checks added or reviewed.
2. Summarize docs changed.
3. Confirm no live-service dependency was introduced.
```

### Completion checklist

- [ ] UX evidence is observable.
- [ ] Code tasks require red-first tests.
- [ ] Documentation tasks have review checks.
- [ ] No live-service dependency was introduced.

## UX1-T01 - Add Theme Tokens To Static CSS

**Type:** Code
**Goal:** Add reusable AskOC theme variables/classes for the `DESIGN.md` visual system.

### Copy/paste prompt

```text
You are implementing UX1-T01 for AskOC AI Concierge.

Task title: Add theme tokens to static CSS
Task goal: Add reusable DESIGN.md theme variables/classes to the existing AskOC static CSS.

Strict TDD:
1. Write a failing test first for required theme tokens or themed class usage.
2. Run the narrow test and confirm it fails because tokens/classes are missing.
3. Implement the smallest CSS/template change.
4. Run the narrow test until green.
5. Run the broader relevant test suite.

Implementation requirements:
- Preserve existing AskOC routes and API behavior.
- Include black #0A0A0A, white #FAFAFA, red #EF4444, neutral surfaces, square radius, no shadows, no gradients, and focus ring tokens.
- Use Space Mono for trace IDs/action codes where applicable.
- Do not add unrelated frontend frameworks unless the repo already uses them.

Expected response:
1. Show red test name and expected failure.
2. Summarize CSS/template changes.
3. List test commands run.
4. Confirm no behavior or product-scope changes.
```

### Completion checklist

- [ ] Failing token/style test was written first.
- [ ] Theme tokens/classes exist.
- [ ] No shadows, gradients, or rounded panels were added.
- [ ] Broader tests pass.

## UX1-T02 - Theme App Header And Route Navigation

**Type:** Code
**Goal:** Make navigation concise, keyboard accessible, and consistent with the red active underline rule.

### Copy/paste prompt

```text
You are implementing UX1-T02 for AskOC AI Concierge.

Task title: Theme app header and route navigation
Task goal: Apply the DESIGN.md navigation treatment to existing AskOC app routes.

Strict TDD:
1. Write failing render/template tests for nav landmark, active item, synthetic-mode label, and keyboard-accessible links/buttons.
2. Verify red state.
3. Implement the smallest template/CSS/JS changes.
4. Run narrow and broader tests.

UX requirements:
- Active nav uses a restrained red underline.
- Header states the AskOC MVP product clearly.
- Synthetic mode is visible.
- No marketing hero or editorial site navigation is introduced.

Expected response:
1. Show red test evidence.
2. Summarize header/nav changes.
3. List tests run.
4. Confirm AskOC product scope remains unchanged.
```

### Completion checklist

- [ ] Header/nav tests were red first.
- [ ] Active route and synthetic mode render.
- [ ] Keyboard access is covered.
- [ ] AskOC routes are preserved.

## UX1-T03 - Standardize Buttons Inputs Chips And Focus States

**Type:** Code
**Goal:** Apply square, high-contrast component styling without changing API behavior.

### Copy/paste prompt

```text
You are implementing UX1-T03 for AskOC AI Concierge.

Task title: Standardize buttons, inputs, chips, and focus states
Task goal: Apply DESIGN.md component styling to existing learner/admin controls.

Strict TDD:
1. Write failing tests for accessible names, focusable controls, status chip rendering, and validation/error display.
2. Verify red state.
3. Implement the smallest template/CSS/JS changes.
4. Run narrow and broader tests.

UX requirements:
- Buttons use uppercase Work Sans, 2px borders, square corners.
- Inputs use 44px target height where practical, square corners, visible labels/help/error text.
- Status chips represent source/action/workflow/escalation states.
- Focus ring matches DESIGN.md.

Expected response:
1. Show red test evidence.
2. Summarize component changes.
3. List tests run.
4. Confirm API behavior did not change.
```

### Completion checklist

- [ ] Component tests were red first.
- [ ] Buttons, inputs, chips, and focus states are covered.
- [ ] Error states remain safe.
- [ ] API behavior is unchanged.

## UX2-T01 - Theme Learner And Assistant Message Layout

**Type:** Code
**Goal:** Make conversation turns readable, compact, and visibly synthetic-demo-safe.

### Copy/paste prompt

```text
You are implementing UX2-T01 for AskOC AI Concierge.

Task title: Theme learner and assistant message layout
Task goal: Improve chat message readability while preserving the existing chat API contract.

Strict TDD:
1. Write failing UI/render tests for learner role, assistant role, empty state, synthetic marker, and safe error state.
2. Verify red state.
3. Implement the smallest template/CSS/JS changes.
4. Run narrow and broader tests.

UX requirements:
- Conversation stays compact and scannable.
- Synthetic-demo status is visible.
- Message text does not overlap on mobile or desktop.
- Raw internal errors are never shown.

Expected response:
1. Show red test evidence.
2. Summarize chat layout changes.
3. List tests run.
4. Confirm chat API contract is unchanged.
```

### Completion checklist

- [ ] Message layout tests were red first.
- [ ] Role and empty/error states are covered.
- [ ] Synthetic marker is visible.
- [ ] Chat API contract is unchanged.

## UX2-T02 - Surface Sources And Confidence States

**Type:** Code
**Goal:** Show source citations, stale/low-confidence warnings, and fallback states clearly.

### Copy/paste prompt

```text
You are implementing UX2-T02 for AskOC AI Concierge.

Task title: Surface sources and confidence states
Task goal: Make RAG/source confidence visible in the chat response surface.

Strict TDD:
1. Write failing tests for citation rendering, stale-source label, low-confidence fallback, and no-source escalation state.
2. Verify red state.
3. Implement the smallest rendering change.
4. Run narrow and broader tests.

UX requirements:
- Citations are compact and source-grounded.
- Warning states use semantic styling, not excessive red.
- Fallback/escalation states are visibly different from grounded answers.
- No private URLs or real learner data are displayed.

Expected response:
1. Show red test evidence.
2. Summarize source/confidence rendering.
3. List tests run.
4. Confirm no private source data was added.
```

### Completion checklist

- [ ] Source/confidence tests were red first.
- [ ] Citation and fallback states are covered.
- [ ] Private URLs are not introduced.
- [ ] Broader tests pass.

## UX2-T03 - Surface Action Trace And Synthetic Integrations

**Type:** Code
**Goal:** Show classifier, mock Banner, mock payment, workflow, CRM, and trace ID actions as scan-friendly rows.

### Copy/paste prompt

```text
You are implementing UX2-T03 for AskOC AI Concierge.

Task title: Surface action trace and synthetic integrations
Task goal: Make synthetic integration and workflow evidence visible in the chat UI.

Strict TDD:
1. Write failing tests for action row rendering, action status chips, trace ID display, idempotency key display if present, and safe hidden internals.
2. Verify red state.
3. Implement the smallest rendering change.
4. Run narrow and broader tests.

UX requirements:
- Trace IDs and action codes use Space Mono.
- Rows use border-led list styling.
- Mock Banner/payment/workflow/CRM actions are clearly marked synthetic.
- Internal errors and raw payloads are not exposed.

Expected response:
1. Show red test evidence.
2. Summarize action trace UI changes.
3. List tests run.
4. Confirm no raw internal payloads are displayed.
```

### Completion checklist

- [ ] Action trace tests were red first.
- [ ] Trace ID and statuses render.
- [ ] Synthetic integrations are labeled.
- [ ] Raw internals are not exposed.

## UX3-T01 - Theme Admin Metrics Dashboard

**Type:** Code
**Goal:** Display containment, escalation, workflow, review queue, and evaluation metrics in a dense operational layout.

### Copy/paste prompt

```text
You are implementing UX3-T01 for AskOC AI Concierge.

Task title: Theme admin metrics dashboard
Task goal: Apply DESIGN.md visual language to existing admin metrics/dashboard surfaces.

Strict TDD:
1. Write failing tests for metric labels, synthetic-mode marker, redacted values, and empty/loading/error states.
2. Verify red state.
3. Implement the smallest template/CSS/JS changes.
4. Run narrow and broader tests.

UX requirements:
- Use dense operational bands/lists, not rounded cards.
- No raw learner messages are shown.
- Metrics remain readable on desktop and mobile.
- Red is reserved for critical/error emphasis.

Expected response:
1. Show red test evidence.
2. Summarize dashboard changes.
3. List tests run.
4. Confirm redaction is preserved.
```

### Completion checklist

- [ ] Dashboard tests were red first.
- [ ] Metrics and empty/error states are covered.
- [ ] Redaction is preserved.
- [ ] Broader tests pass.

## UX3-T02 - Theme Audit And Review Queue Rows

**Type:** Code
**Goal:** Make trace-linked audit/review items easy to scan and filter.

### Copy/paste prompt

```text
You are implementing UX3-T02 for AskOC AI Concierge.

Task title: Theme audit and review queue rows
Task goal: Improve audit/review queue readability without changing review semantics.

Strict TDD:
1. Write failing tests for trace ID, queue, priority, status chip, redaction marker, filter control, and safe empty state.
2. Verify red state.
3. Implement the smallest rendering and styling changes.
4. Run narrow and broader tests.

UX requirements:
- Rows use bottom borders and compact metadata.
- Priority/error states use semantic styling.
- Trace IDs use Space Mono.
- Redacted data markers remain visible.

Expected response:
1. Show red test evidence.
2. Summarize audit/review UI changes.
3. List tests run.
4. Confirm review behavior is unchanged.
```

### Completion checklist

- [ ] Audit/review tests were red first.
- [ ] Trace, queue, priority, and status render.
- [ ] Redaction markers are visible.
- [ ] Review semantics are unchanged.

## UX3-T03 - Theme Evaluation Report Surface

**Type:** Code or Documentation
**Goal:** Make evaluation pass/fail evidence visually consistent in UI or docs.

### Copy/paste prompt

```text
You are implementing UX3-T03 for AskOC AI Concierge.

Task title: Theme evaluation report surface
Task goal: Make evaluation gate evidence easy to inspect.

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
- Do not hide known limitations.

Expected response:
1. Show red test or review-first evidence.
2. Summarize evaluation surface changes.
3. List tests/review checks run.
4. Confirm no evaluation claims were overstated.
```

### Completion checklist

- [ ] Red test or review-first check happened.
- [ ] Pass/fail evidence is visible.
- [ ] Warnings are not hidden.
- [ ] Claims are not overstated.

## UX4-T01 - Add Accessibility Checks For Themed Controls

**Type:** Code
**Goal:** Verify landmarks, labels, focus order, and keyboard behavior for chat/admin controls.

### Copy/paste prompt

```text
You are implementing UX4-T01 for AskOC AI Concierge.

Task title: Add accessibility checks for themed controls
Task goal: Protect themed chat/admin controls with accessibility tests.

Strict TDD:
1. Write a failing accessibility or keyboard test first.
2. Verify it fails for a real missing label, landmark, focus state, or keyboard behavior.
3. Implement the smallest fix.
4. Run accessibility/keyboard tests and broader tests.

Required coverage:
- chat input and submit
- action/source controls
- admin dashboard controls
- review queue controls
- navigation

Expected response:
1. Show red accessibility evidence.
2. Summarize fixes.
3. List tests run.
4. Call out residual manual checks.
```

### Completion checklist

- [ ] Accessibility test was red first.
- [ ] Key controls are covered.
- [ ] Keyboard behavior is covered.
- [ ] Broader tests pass.

## UX4-T02 - Add Responsive Layout Checks

**Type:** Code
**Goal:** Verify chat, action trace, and dashboard content do not overlap on mobile or desktop.

### Copy/paste prompt

```text
You are implementing UX4-T02 for AskOC AI Concierge.

Task title: Add responsive layout checks
Task goal: Verify themed MVP screens fit mobile and desktop viewports.

Strict TDD:
1. Add a failing viewport, screenshot, or layout assertion test first where tooling exists.
2. Verify red state.
3. Implement the smallest responsive CSS/layout change.
4. Run viewport checks and broader tests.

Required checks:
- mobile chat
- mobile action trace
- desktop chat
- desktop admin dashboard
- desktop review/eval evidence

Expected response:
1. Show red visual/layout evidence.
2. Summarize responsive changes.
3. List viewport/test commands run.
4. Call out manual screenshot review if needed.
```

### Completion checklist

- [ ] Responsive test/check was red first where tooling exists.
- [ ] Mobile and desktop states are covered.
- [ ] No overlap or overflow is observed.
- [ ] Broader tests pass.

## UX4-T03 - Add Contrast And Red-Accent Review Gate

**Type:** Code or Documentation
**Goal:** Check high contrast and restrained red usage across MVP screens.

### Copy/paste prompt

```text
You are implementing UX4-T03 for AskOC AI Concierge.

Task title: Add contrast and red-accent review gate
Task goal: Prevent theme drift into low contrast, too much red, gradients, shadows, or rounded UI.

Strict TDD or review-first:
1. If changing code, write a failing style/visual assertion first.
2. If changing docs, define a manual review checklist first.
3. Verify red/review state.
4. Implement the smallest change.
5. Run tests or review checks.

Review requirements:
- Red is not dominant.
- Text contrast is high.
- No shadows, gradients, rounded panels, or decorative blobs are introduced.
- Focus states remain visible.

Expected response:
1. Show red test or review-first evidence.
2. Summarize changes.
3. List tests/review checks run.
4. Confirm theme constraints.
```

### Completion checklist

- [ ] Red test or review-first check happened.
- [ ] Red accent usage is restrained.
- [ ] Contrast and focus states are checked.
- [ ] No forbidden visual patterns were introduced.

## UX5-T01 - Update README Reviewer Path For Themed MVP

**Type:** Documentation
**Goal:** Explain the themed UI in the existing AskOC reviewer path.

### Copy/paste prompt

```text
You are implementing UX5-T01 for AskOC AI Concierge.

Task title: Update README reviewer path for themed MVP
Task goal: Explain DESIGN.md theming as UX polish for the existing AskOC MVP.

Instructions:
1. Keep the README led by AskOC AI Concierge.
2. Mention DESIGN.md only as the visual theme.
3. Keep transcript/payment workflow, synthetic integrations, TDD, privacy, and evaluation as the main story.
4. Do not describe VoiceBox as a separate app.

Primary files:
- README.md

Review checks:
- README still explains the AskOC MVP in under two minutes.
- Theme language does not replace product language.
- No unsupported production claims are added.

Expected response:
1. List review checks performed.
2. Summarize README changes.
3. Confirm no product-scope drift.
```

### Completion checklist

- [ ] README remains AskOC-first.
- [ ] Theme is framed as visual polish.
- [ ] No unsupported claims are added.
- [ ] No real data or secrets are added.

## UX5-T02 - Update Demo Script For Themed Screens

**Type:** Documentation
**Goal:** Add where to point during the demo: chat, source chips, action trace, dashboard, eval evidence.

### Copy/paste prompt

```text
You are implementing UX5-T02 for AskOC AI Concierge.

Task title: Update demo script for themed screens
Task goal: Show how the themed UI supports the existing 5-7 minute AskOC demo.

Instructions:
1. Keep the demo flow centered on transcript answer, payment workflow, CRM escalation, dashboard evidence, and TDD/eval proof.
2. Add brief visual callouts for themed chat, source citations, action trace, synthetic labels, and admin evidence.
3. Do not add a separate editorial/blog demo.

Primary files:
- docs/demo-script.md

Review checks:
- Every spoken UI claim maps to a visible screen or test.
- Demo remains 5-7 minutes.
- AskOC product story remains unchanged.

Expected response:
1. List review checks performed.
2. Summarize demo script changes.
3. Confirm no separate product story was introduced.
```

### Completion checklist

- [ ] Demo remains AskOC-centered.
- [ ] UI callouts map to visible evidence.
- [ ] Demo remains time-boxed.
- [ ] No separate product story is introduced.

## UX5-T03 - Record Final UX Quality Evidence

**Type:** Documentation
**Goal:** Summarize tests, visual checks, limitations, and out-of-scope UI work.

### Copy/paste prompt

```text
You are implementing UX5-T03 for AskOC AI Concierge.

Task title: Record final UX quality evidence
Task goal: Record what UX tests/checks passed and what remains out of scope.

Instructions:
1. Update changelog/test-plan docs after implementation evidence exists.
2. Include TDD, accessibility, responsive, redaction, synthetic-data, and no-secret checks.
3. Keep limitations honest.
4. Do not claim production deployment or real system integration.

Primary files:
- CHANGELOG.md
- docs/test-plan.md

Review checks:
- Evidence lists exact commands or manual checks.
- Known gaps are not hidden.
- No real data, secrets, or private URLs are added.

Expected response:
1. List review checks performed.
2. Summarize docs changed.
3. List final test commands and results.
4. Confirm known gaps and out-of-scope items.
```

### Completion checklist

- [ ] Evidence lists commands/checks.
- [ ] Known gaps are documented.
- [ ] No unsupported production claims are made.
- [ ] No real data, secrets, or private URLs are added.
