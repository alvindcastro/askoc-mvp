# VoiceBox Web App Brainstorm And Phased Task Board

This document turns `DESIGN.md` into a buildable web app concept and a strict-TDD task board. It is intentionally separate from the AskOC implementation history so the VoiceBox app can be planned without rewriting completed MVP phases.

## Product Concept

VoiceBox is a high-contrast editorial web app for opinion blogs and cultural commentary. The app should feel like a sharp digital magazine: fast scanning, strong hierarchy, compact reading tools, and confident typography.

The strongest MVP direction is a **commentary publishing and reading workspace** with:

- a front page that promotes the latest featured argument,
- category and rubric navigation,
- article reading pages with pull quotes, citations, bylines, and related pieces,
- a lightweight editor preview for drafting opinion posts,
- search and filters for topics, authors, and status,
- strict accessibility and responsive checks for the bold visual system.

## UX Rules From `DESIGN.md`

- Use black `#0A0A0A`, white `#FAFAFA`, neutral surfaces, and one restrained red accent `#EF4444`.
- Use `Archivo Black` for display/headings, `Work Sans` for UI/body, and `Space Mono` for code-like metadata.
- Keep all UI square: `0px` border radius except full-radius author avatars.
- Use thick borders, typography, and contrast for hierarchy. Do not use shadows, gradients, soft panels, decorative blobs, or rounded cards.
- Use red sparingly: one prominent red element per viewport, active navigation underline, featured top border, destructive/error state, or primary-button hover.
- Keep pages typographic and editorial. Avoid decorative imagery unless it directly supports the story.
- Use visible focus rings: `0 0 0 2px #FAFAFA, 0 0 0 4px #0A0A0A`.

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
| V0 - Product framing | VoiceBox scope, audience, routes, data model, and TDD evidence rules are frozen. |
| V1 - App foundation | The web app boots with routing, fixture loading, tests, and the VoiceBox design tokens. |
| V2 - Editorial home and archive | Readers can scan featured, latest, category, and search/archive views. |
| V3 - Article reading experience | Article pages support rich editorial layout, pull quotes, related links, and accessible reading states. |
| V4 - Drafting and preview workflow | Editors can create draft content in a local-only workflow and preview it in the VoiceBox design system. |
| V5 - Quality, accessibility, and release polish | Tests, accessibility gates, responsive checks, and demo documentation are complete. |

## V0 - Product Framing

**Phase outcome:** VoiceBox is scoped as an editorial/commentary web app before implementation starts.

- [ ] **V0-T01 - Define the VoiceBox MVP thesis**
  - Type: Documentation
  - Goal: State the product purpose, target reader, target editor, and first demo path.
  - Primary files: `README.md`, `docs/voicebox-web-app-brainstorm.md`
  - Prompt: [`V0-T01` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v0-t01-define-the-voicebox-mvp-thesis)
  - Review check: Confirm the thesis says editorial/commentary web app, not a generic landing page.
  - Done when: audience, demo path, and non-goals are visible.

- [ ] **V0-T02 - Freeze routes and information architecture**
  - Type: Documentation
  - Goal: Define the MVP routes and navigation labels before code.
  - Primary files: `docs/voicebox-web-app-brainstorm.md`
  - Prompt: [`V0-T02` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v0-t02-freeze-routes-and-information-architecture)
  - Review check: Confirm each route supports article discovery, reading, or drafting.
  - Done when: `/`, `/archive`, `/article/{slug}`, `/drafts`, and `/drafts/{id}/preview` are specified.

- [ ] **V0-T03 - Define editorial fixture schema**
  - Type: Documentation
  - Goal: Specify synthetic article, author, category, and draft data before implementation.
  - Primary files: `docs/voicebox-web-app-brainstorm.md`
  - Prompt: [`V0-T03` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v0-t03-define-editorial-fixture-schema)
  - Review check: Confirm fixtures are synthetic and contain no real unpublished work.
  - Done when: schema includes article status, byline, dateline, category, tags, excerpt, body blocks, pull quotes, and related slugs.

## V1 - App Foundation

**Phase outcome:** The app boots with tested routes, fixtures, and reusable VoiceBox styling primitives.

- [ ] **V1-T01 - Initialize app shell and test harness**
  - Type: Code
  - Goal: Create the minimal runnable web app with a passing base test command.
  - Primary files: app entrypoint, route tests, package config, test config
  - Prompt: [`V1-T01` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v1-t01-initialize-app-shell-and-test-harness)
  - Tests first: failing smoke test for rendering the root route.
  - Done when: the root route renders and the base test command passes.

- [ ] **V1-T02 - Add VoiceBox design tokens**
  - Type: Code
  - Goal: Implement colors, fonts, spacing, borders, focus rings, and square component defaults from `DESIGN.md`.
  - Primary files: global CSS/theme file, token tests or visual contract tests
  - Prompt: [`V1-T02` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v1-t02-add-voicebox-design-tokens)
  - Tests first: failing test or snapshot proving token names and key values exist.
  - Done when: no shadows, gradients, or non-avatar radius are introduced.

- [ ] **V1-T03 - Load synthetic editorial fixtures**
  - Type: Code
  - Goal: Load deterministic article, author, category, and draft fixtures.
  - Primary files: fixture data, fixture loader, fixture loader tests
  - Prompt: [`V1-T03` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v1-t03-load-synthetic-editorial-fixtures)
  - Tests first: invalid fixture and duplicate slug tests.
  - Done when: published and draft content loads deterministically.

## V2 - Editorial Home And Archive

**Phase outcome:** Readers can discover commentary through featured, latest, category, and search/archive views.

- [ ] **V2-T01 - Build editorial home page**
  - Type: Code
  - Goal: Render a first-screen editorial front page using VoiceBox hierarchy.
  - Primary files: home route/component/template, home tests
  - Prompt: [`V2-T01` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v2-t01-build-editorial-home-page)
  - Tests first: failing test for featured article, latest list, category labels, and active nav.
  - Done when: the page uses massive headline typography without decorative hero art.

- [ ] **V2-T02 - Build archive with filters and search**
  - Type: Code
  - Goal: Let readers filter by category, author, tag, and text query.
  - Primary files: archive route/component/template, search/filter tests
  - Prompt: [`V2-T02` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v2-t02-build-archive-with-filters-and-search)
  - Tests first: failing tests for query filtering, empty states, and selected chip styling.
  - Done when: filters are deterministic and accessible.

- [ ] **V2-T03 - Add responsive navigation**
  - Type: Code
  - Goal: Provide category navigation that works on mobile and desktop.
  - Primary files: nav component/template, navigation tests
  - Prompt: [`V2-T03` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v2-t03-add-responsive-navigation)
  - Tests first: failing tests for active route, keyboard access, and mobile menu state.
  - Done when: active nav uses the red underline and no layout text overlaps.

## V3 - Article Reading Experience

**Phase outcome:** Article pages deliver the bold editorial reading experience promised by `DESIGN.md`.

- [ ] **V3-T01 - Build article detail route**
  - Type: Code
  - Goal: Render article body blocks, headline, byline, dateline, category, and citations.
  - Primary files: article route/component/template, article tests
  - Prompt: [`V3-T01` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v3-t01-build-article-detail-route)
  - Tests first: failing tests for known slug, missing slug, and metadata rendering.
  - Done when: unknown articles return a safe not-found state.

- [ ] **V3-T02 - Add pull quotes and related links**
  - Type: Code
  - Goal: Render pull quotes with the red left border and related links with compact list styling.
  - Primary files: article blocks, related-links component, tests
  - Prompt: [`V3-T02` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v3-t02-add-pull-quotes-and-related-links)
  - Tests first: failing tests for block rendering and related slug validation.
  - Done when: pull quotes use one restrained red accent and related links remain scan-friendly.

- [ ] **V3-T03 - Add reading controls**
  - Type: Code
  - Goal: Add accessible share, copy link, and reading progress controls without clutter.
  - Primary files: reading controls, tests
  - Prompt: [`V3-T03` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v3-t03-add-reading-controls)
  - Tests first: failing tests for button labels, keyboard behavior, and copy feedback state.
  - Done when: controls use icon or compact button patterns and visible focus rings.

## V4 - Drafting And Preview Workflow

**Phase outcome:** Editors can create local draft content and preview it in the same editorial system.

- [ ] **V4-T01 - Build draft list and draft states**
  - Type: Code
  - Goal: Show drafts by status with clear editorial metadata.
  - Primary files: drafts route/component/template, draft tests
  - Prompt: [`V4-T01` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v4-t01-build-draft-list-and-draft-states)
  - Tests first: failing tests for draft status chips, empty state, and sorting.
  - Done when: status chips follow `DESIGN.md` semantics.

- [ ] **V4-T02 - Build local draft editor form**
  - Type: Code
  - Goal: Add local-only draft editing fields for title, rubric, excerpt, body, tags, and pull quote.
  - Primary files: draft editor, form validation, tests
  - Prompt: [`V4-T02` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v4-t02-build-local-draft-editor-form)
  - Tests first: failing tests for required fields, field errors, and no accidental publish.
  - Done when: invalid drafts cannot be previewed as publish-ready.

- [ ] **V4-T03 - Build editorial preview**
  - Type: Code
  - Goal: Preview draft content using the same article layout as published articles.
  - Primary files: preview route/component/template, preview tests
  - Prompt: [`V4-T03` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v4-t03-build-editorial-preview)
  - Tests first: failing tests for preview rendering, unsaved changes, and safe draft-only labels.
  - Done when: preview clearly marks draft status and never implies production publishing.

## V5 - Quality, Accessibility, And Release Polish

**Phase outcome:** The app is demo-ready, accessible, responsive, and documented.

- [ ] **V5-T01 - Add accessibility and keyboard tests**
  - Type: Code
  - Goal: Check labels, landmarks, focus order, contrast-sensitive states, and keyboard navigation.
  - Primary files: accessibility tests, affected components
  - Prompt: [`V5-T01` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v5-t01-add-accessibility-and-keyboard-tests)
  - Tests first: failing accessibility test for at least one missing landmark or label.
  - Done when: keyboard navigation works across home, archive, article, and draft routes.

- [ ] **V5-T02 - Add responsive visual checks**
  - Type: Code
  - Goal: Verify mobile and desktop layouts do not overlap and keep the VoiceBox typographic rules.
  - Primary files: visual/e2e tests, responsive CSS
  - Prompt: [`V5-T02` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v5-t02-add-responsive-visual-checks)
  - Tests first: failing viewport test or screenshot diff for a known layout issue.
  - Done when: home, archive, article, and draft preview pass mobile and desktop checks.

- [ ] **V5-T03 - Finalize demo script and release checklist**
  - Type: Documentation
  - Goal: Document how to demo the app, what tests prove, and what remains out of scope.
  - Primary files: README/demo docs, release checklist
  - Prompt: [`V5-T03` in voicebox-tdd-prompts.md](voicebox-tdd-prompts.md#v5-t03-finalize-demo-script-and-release-checklist)
  - Review check: Run through the demo script and confirm every claim has visible app or test evidence.
  - Done when: test commands, route walkthrough, design constraints, and non-goals are documented.

## Suggested Demo Path

- [ ] Open the home page and show the featured argument, category navigation, latest list, and red active underline.
- [ ] Use archive filters to narrow by category and search query.
- [ ] Open an article and show headline, byline, dateline, pull quote, citations, and related links.
- [ ] Open drafts, edit a synthetic draft, trigger validation, and preview the draft.
- [ ] Run the test command and show TDD evidence in the prompt/task docs.

## Out Of Scope For MVP

- [ ] Real authentication or publishing permissions.
- [ ] Real CMS integrations.
- [ ] User comments or moderation.
- [ ] Payment, subscriptions, or paywalls.
- [ ] AI-generated articles.
- [ ] Analytics that track real readers.
- [ ] Decorative hero imagery, gradients, rounded visual language, or shadow-based layout.
