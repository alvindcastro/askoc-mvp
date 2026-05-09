# AskOC UX Theme Brainstorm And Phased Task Board

This document applies `DESIGN.md` as the UX theme for the existing **AskOC AI Concierge** MVP. It does not change the product direction: AskOC remains a Go-based learner-service automation concierge for transcript/payment support, synthetic enterprise integrations, workflow automation, CRM escalation, audit evidence, and evaluation gates.

`DESIGN.md` names the visual system "VoiceBox"; in this repository, treat that as the **visual language only**: bold editorial typography, high contrast, square borders, no shadows, and a restrained red accent.

For the active revamp branch, [web-app-revamp-tdd-plan.md](web-app-revamp-tdd-plan.md) is the source-of-truth execution handoff. This file is retained as the theme brief and historical UX0-UX5 planning note; use the R0-R5 taxonomy for implementation status, test evidence, and changelog entries.

## 2026-05-09 Revamp Implementation Notes

- `/` remains an intentional chat alias because `cmd/api/main.go` routes both `/` and `/chat` to `ChatPageHandler`; no redirect or new route was added.
- Chat and admin now share the same visual language: high-contrast shell, active red route underline, square panels, 2px borders, visible synthetic-mode labels, and focus rings.
- Chat proof points are source chips, confidence/risk/freshness metadata, trace ID, action rows, workflow ID, CRM case ID, priority, and idempotency key labels when present.
- Admin proof points are redacted aggregate metrics, review queue filter, trace/queue/priority/status chips, stale-source count, evaluation-gate copy, and audit export/reset/purge controls.
- Verification is tracked through `go test ./internal/handlers -run 'Test(Chat|Admin).*Revamp|Test(Chat|Admin)StaticAssets'`, `go test ./internal/handlers`, `go test ./...`, `make eval`, `make secret-check`, and `git diff --check`.

## MVP UX Thesis

AskOC should feel like a crisp operational console for a learner-service automation demo, not a generic chatbot and not a separate editorial site. The UX should help an interviewer quickly see:

- what the learner asked,
- what source-grounded answer was returned,
- which synthetic systems were checked,
- which workflow or CRM action happened,
- what trace ID ties the conversation to audit and evaluation evidence,
- why no real learner data or real enterprise systems are involved.

## DESIGN.md Rules To Apply

- Use black `#0A0A0A`, white `#FAFAFA`, neutral surfaces, and restrained red `#EF4444`.
- Use `Archivo Black` for major headings, `Work Sans` for UI/body text, and `Space Mono` for trace IDs, synthetic IDs, route names, and action codes.
- Keep components square: `0px` radius except avatar-like thumbnails if ever used.
- Use borders and typographic hierarchy instead of shadows, gradients, or soft cards.
- Use red sparingly: active navigation underline, one featured status accent, destructive/error state, or primary hover.
- Use 2px borders on interactive elements.
- Use the `DESIGN.md` focus ring: `0 0 0 2px #FAFAFA, 0 0 0 4px #0A0A0A`.
- Do not add decorative hero imagery, rounded dashboard cards, gradient backgrounds, or marketing-page structure.

## Strict TDD Rule

Every code task must follow Red, Green, Refactor:

- [ ] Write or update the failing test first.
- [ ] Run the narrow test and confirm it fails for the expected reason.
- [ ] Implement the smallest production change.
- [ ] Run the narrow test until green.
- [ ] Run the broader suite for the touched area.
- [ ] Refactor only while tests stay green.
- [ ] Update docs when behavior, routes, data shape, or UX assumptions change.

## Phase Overview

| Phase | Outcome |
|---|---|
| UX0 - Theme mapping | `DESIGN.md` is mapped to AskOC MVP surfaces without changing product scope. |
| UX1 - App shell and navigation | The Go-rendered app shell uses the theme while preserving existing routes. |
| UX2 - Chat workflow surface | Learner chat makes sources, actions, synthetic data, and trace evidence easy to scan. |
| UX3 - Admin and audit evidence | Dashboard/review surfaces expose workflow, escalation, and evaluation proof clearly. |
| UX4 - Responsive and accessibility gates | Mobile/desktop layouts, keyboard use, focus states, and contrast are verified. |
| UX5 - Demo polish | README/demo docs explain the themed MVP and the tests that protect it. |

## UX0 - Theme Mapping

**Phase outcome:** `DESIGN.md` is treated as the AskOC visual theme, not as a new app concept.

- [ ] **UX0-T01 - Map DESIGN.md tokens to AskOC UI**
  - Type: Documentation
  - Goal: Define how colors, typography, borders, chips, buttons, inputs, and focus rings apply to chat, admin, eval, and demo pages.
  - Primary files: `docs/askoc-ux-theme-brainstorm.md`, `docs/askoc-ux-tdd-prompts.md`
  - Prompt: [`UX0-T01` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux0-t01-map-designmd-tokens-to-askoc-ui)
  - Review check: Confirm the MVP is still AskOC AI Concierge.
  - Done when: every theme decision maps to an existing AskOC surface.

- [ ] **UX0-T02 - Freeze themed route inventory**
  - Type: Documentation
  - Goal: List the existing MVP routes and their themed UX jobs.
  - Primary files: `docs/askoc-ux-theme-brainstorm.md`
  - Prompt: [`UX0-T02` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux0-t02-freeze-themed-route-inventory)
  - Review check: Confirm no new product routes imply a separate editorial app.
  - Done when: `/`, `/chat`, admin/dashboard routes, API evidence, and static assets are covered.

- [ ] **UX0-T03 - Define UX acceptance evidence**
  - Type: Documentation
  - Goal: Define what tests, screenshots, and manual checks prove the theme works.
  - Primary files: `docs/askoc-ux-theme-brainstorm.md`, `docs/test-plan.md`
  - Prompt: [`UX0-T03` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux0-t03-define-ux-acceptance-evidence)
  - Review check: Confirm evidence includes TDD, accessibility, responsive checks, and synthetic-data labeling.
  - Done when: each UX phase has observable proof.

## UX1 - App Shell And Navigation

**Phase outcome:** The existing Go-rendered app shell uses the `DESIGN.md` theme while preserving MVP routes and behavior.

- [ ] **UX1-T01 - Add theme tokens to static CSS**
  - Type: Code
  - Goal: Add reusable AskOC theme variables/classes for the `DESIGN.md` visual system.
  - Primary files: `web/static/app.css`, CSS contract tests or handler/template tests
  - Prompt: [`UX1-T01` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux1-t01-add-theme-tokens-to-static-css)
  - Tests first: failing test for required token names or themed class usage.
  - Done when: theme tokens exist and no shadows/gradients/rounded panels are introduced.

- [ ] **UX1-T02 - Theme app header and route navigation**
  - Type: Code
  - Goal: Make navigation concise, keyboard accessible, and consistent with the red active underline rule.
  - Primary files: `web/templates/chat.html`, static CSS/JS, handler/template tests
  - Prompt: [`UX1-T02` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux1-t02-theme-app-header-and-route-navigation)
  - Tests first: failing render test for nav landmark, active item, and synthetic-mode label.
  - Done when: app shell still routes to the existing AskOC chat/admin surfaces.

- [ ] **UX1-T03 - Standardize buttons, inputs, chips, and focus states**
  - Type: Code
  - Goal: Apply square, high-contrast component styling without changing API behavior.
  - Primary files: `web/static/app.css`, templates, UI tests
  - Prompt: [`UX1-T03` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux1-t03-standardize-buttons-inputs-chips-and-focus-states)
  - Tests first: failing test for accessible labels and focusable controls.
  - Done when: learner message input and action controls follow `DESIGN.md`.

## UX2 - Chat Workflow Surface

**Phase outcome:** The chat UI makes the AskOC workflow understandable at a glance.

- [ ] **UX2-T01 - Theme learner and assistant message layout**
  - Type: Code
  - Goal: Make conversation turns readable, compact, and visibly synthetic-demo-safe.
  - Primary files: `web/templates/chat.html`, `web/static/app.css`, `web/static/app.js`, UI tests
  - Prompt: [`UX2-T01` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux2-t01-theme-learner-and-assistant-message-layout)
  - Tests first: failing test for rendered message roles, safe empty state, and synthetic marker.
  - Done when: message layout supports transcript/payment demo without visual clutter.

- [ ] **UX2-T02 - Surface sources and confidence states**
  - Type: Code
  - Goal: Show source citations, stale/low-confidence warnings, and fallback states clearly.
  - Primary files: chat template/JS/CSS, handler or JS tests
  - Prompt: [`UX2-T02` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux2-t02-surface-sources-and-confidence-states)
  - Tests first: failing test for citation rendering and low-confidence fallback label.
  - Done when: grounded answers are visually distinct from fallback/escalation states.

- [ ] **UX2-T03 - Surface action trace and synthetic integrations**
  - Type: Code
  - Goal: Show classifier, mock Banner, mock payment, workflow, CRM, and trace ID actions as scan-friendly rows.
  - Primary files: chat template/JS/CSS, UI tests
  - Prompt: [`UX2-T03` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux2-t03-surface-action-trace-and-synthetic-integrations)
  - Tests first: failing test for action list rendering and trace ID display.
  - Done when: action evidence is visible without exposing raw internals.

## UX3 - Admin And Audit Evidence

**Phase outcome:** Admin/dashboard/review surfaces match the theme and make proof points easy to inspect.

- [ ] **UX3-T01 - Theme admin metrics dashboard**
  - Type: Code
  - Goal: Display containment, escalation, workflow, review queue, and evaluation metrics in a dense operational layout.
  - Primary files: admin template/static assets/tests
  - Prompt: [`UX3-T01` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux3-t01-theme-admin-metrics-dashboard)
  - Tests first: failing test for metric labels, synthetic-mode marker, and no raw learner text.
  - Done when: metrics are readable and redacted.

- [ ] **UX3-T02 - Theme audit and review queue rows**
  - Type: Code
  - Goal: Make trace-linked audit/review items easy to scan and filter.
  - Primary files: admin template/static assets/tests
  - Prompt: [`UX3-T02` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux3-t02-theme-audit-and-review-queue-rows)
  - Tests first: failing test for trace ID, queue, priority, status chip, and redaction marker.
  - Done when: review rows use border-led list styling and safe data only.

- [ ] **UX3-T03 - Theme evaluation report surface**
  - Type: Code or Documentation
  - Goal: Make evaluation pass/fail evidence visually consistent in UI or docs.
  - Primary files: `reports/eval-summary.md`, admin UI, docs
  - Prompt: [`UX3-T03` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux3-t03-theme-evaluation-report-surface)
  - Tests/review first: failing or manual check for missing critical gate/status evidence.
  - Done when: evaluation evidence is visible and honest about warnings.

## UX4 - Responsive And Accessibility Gates

**Phase outcome:** The themed MVP is usable on desktop and mobile and supports keyboard/screen-reader basics.

- [ ] **UX4-T01 - Add accessibility checks for themed controls**
  - Type: Code
  - Goal: Verify landmarks, labels, focus order, and keyboard behavior for chat/admin controls.
  - Primary files: UI tests, templates, static JS/CSS
  - Prompt: [`UX4-T01` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux4-t01-add-accessibility-checks-for-themed-controls)
  - Tests first: failing accessibility test for a missing label, landmark, or focus state.
  - Done when: key controls are keyboard reachable with visible focus.

- [ ] **UX4-T02 - Add responsive layout checks**
  - Type: Code
  - Goal: Verify chat, action trace, and dashboard content do not overlap on mobile or desktop.
  - Primary files: responsive CSS, visual/e2e tests if available
  - Prompt: [`UX4-T02` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux4-t02-add-responsive-layout-checks)
  - Tests first: failing viewport or layout assertion where tooling exists.
  - Done when: no text overlap, button overflow, or hidden action evidence is observed.

- [ ] **UX4-T03 - Add contrast and red-accent review gate**
  - Type: Code or Documentation
  - Goal: Check high contrast and restrained red usage across MVP screens.
  - Primary files: UI tests, CSS, `docs/test-plan.md`
  - Prompt: [`UX4-T03` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux4-t03-add-contrast-and-red-accent-review-gate)
  - Tests/review first: failing or manual check for overused red or low-contrast state.
  - Done when: red is used as an accent, not a dominant palette.

## UX5 - Demo Polish

**Phase outcome:** The docs and demo script explain that `DESIGN.md` themes the AskOC MVP without changing the product.

- [ ] **UX5-T01 - Update README reviewer path for themed MVP**
  - Type: Documentation
  - Goal: Explain the themed UI in the existing AskOC reviewer path.
  - Primary files: `README.md`
  - Prompt: [`UX5-T01` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux5-t01-update-readme-reviewer-path-for-themed-mvp)
  - Review check: Confirm README still leads with AskOC AI Concierge.
  - Done when: theme is described as UX polish, not a new product.

- [ ] **UX5-T02 - Update demo script for themed screens**
  - Type: Documentation
  - Goal: Add where to point during the demo: chat, source chips, action trace, dashboard, eval evidence.
  - Primary files: `docs/demo-script.md`
  - Prompt: [`UX5-T02` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux5-t02-update-demo-script-for-themed-screens)
  - Review check: Every spoken UI claim maps to a visible screen or test.
  - Done when: demo remains 5-7 minutes and product story remains AskOC.

- [ ] **UX5-T03 - Record final UX quality evidence**
  - Type: Documentation
  - Goal: Summarize tests, visual checks, limitations, and out-of-scope UI work.
  - Primary files: `CHANGELOG.md`, `docs/test-plan.md`
  - Prompt: [`UX5-T03` in askoc-ux-tdd-prompts.md](askoc-ux-tdd-prompts.md#ux5-t03-record-final-ux-quality-evidence)
  - Review check: Evidence includes TDD, accessibility, responsive, redaction, and synthetic-data checks.
  - Done when: future sessions know exactly what was verified.

## Out Of Scope

- [ ] Repositioning the repo as an editorial/blog product.
- [ ] Replacing the AskOC transcript/payment workflow.
- [ ] Adding real authentication, real Banner, real payment, real CRM, real LMS, or real learner records.
- [ ] Adding marketing landing pages.
- [ ] Adding decorative gradients, rounded dashboard cards, soft shadows, or unrelated imagery.
- [ ] Changing API behavior just to support visual polish.
