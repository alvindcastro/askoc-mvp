# VoiceBox TDD Task Prompts

Use these prompts to implement the VoiceBox editorial web app from [VoiceBox Web App Brainstorm And Phased Task Board](voicebox-web-app-brainstorm.md). Each code prompt requires strict TDD.

## Global Context To Include With Every Prompt

```text
Project: VoiceBox editorial web app.
Purpose: Build a high-contrast opinion and cultural commentary web app based on DESIGN.md.
UX source: DESIGN.md.
Design system: black #0A0A0A, white #FAFAFA, red #EF4444, neutral surfaces, Archivo Black headings, Work Sans UI/body, Space Mono metadata, 0px radius, no shadows, no gradients, thick borders, restrained red accent.
Product rule: Build the usable editorial app first. Do not make a marketing landing page.
Data rule: Use synthetic editorial fixtures only. Do not include real unpublished articles, private notes, secrets, analytics tokens, or private URLs.
TDD rule: For code tasks, write a failing test first, verify the failure, implement the smallest production change, run the narrow test, then run the broader suite.
Accessibility rule: Every interactive control must have an accessible name, keyboard behavior, and the DESIGN.md focus ring.
Documentation rule: Update Markdown when routes, commands, fixtures, UX assumptions, or test evidence change.
```

## V0-T01 - Define The VoiceBox MVP Thesis

**Type:** Documentation
**Goal:** State the product purpose, target reader, target editor, first demo path, and non-goals.

### Copy/paste prompt

```text
You are implementing V0-T01 for the VoiceBox editorial web app.

Task title: Define the VoiceBox MVP thesis
Task goal: State the product purpose, target reader, target editor, first demo path, and non-goals.

Instructions:
1. Read DESIGN.md before editing.
2. Update only the Markdown needed for product framing.
3. Position the app as an editorial/commentary web app, not a generic landing page.
4. Keep all examples synthetic.
5. Make acceptance criteria observable.

Primary files:
- docs/voicebox-web-app-brainstorm.md
- README.md if the project README should mention VoiceBox

Review checks:
- The thesis names the target reader and editor.
- The first demo path covers home, archive, article detail, drafts, and preview.
- Non-goals reject real CMS publishing, authentication, comments, analytics tracking, and AI-generated articles.

Expected response:
1. List the review checks performed.
2. Summarize changed Markdown files.
3. Confirm no real unpublished content, secrets, or private URLs were added.
```

### Completion checklist

- [ ] Review checks were performed.
- [ ] Acceptance criteria are satisfied.
- [ ] Related docs are updated.
- [ ] No real unpublished content, secrets, or private URLs were added.

## V0-T02 - Freeze Routes And Information Architecture

**Type:** Documentation
**Goal:** Define the MVP routes and navigation labels before code.

### Copy/paste prompt

```text
You are implementing V0-T02 for the VoiceBox editorial web app.

Task title: Freeze routes and information architecture
Task goal: Define the MVP routes and navigation labels before code.

Instructions:
1. Read DESIGN.md before editing.
2. Define only routes needed for article discovery, reading, and local draft preview.
3. Keep navigation labels short and editorial.
4. Keep acceptance criteria testable.

Primary files:
- docs/voicebox-web-app-brainstorm.md

Required routes:
- /
- /archive
- /article/{slug}
- /drafts
- /drafts/{id}/preview

Review checks:
- Every route has a clear user job.
- Active navigation can be represented with a red underline.
- No route requires real authentication, real CMS data, or live analytics.

Expected response:
1. List the IA review checks performed.
2. Summarize route and navigation decisions.
3. Call out deferred routes.
```

### Completion checklist

- [ ] Route list is complete.
- [ ] Route jobs are testable.
- [ ] Deferred routes are explicit.
- [ ] No real integrations are required.

## V0-T03 - Define Editorial Fixture Schema

**Type:** Documentation
**Goal:** Specify synthetic article, author, category, and draft data before implementation.

### Copy/paste prompt

```text
You are implementing V0-T03 for the VoiceBox editorial web app.

Task title: Define editorial fixture schema
Task goal: Specify synthetic article, author, category, and draft data before implementation.

Instructions:
1. Define fixture fields before writing loader code.
2. Keep all fixture content synthetic.
3. Include validation rules that future code tests can enforce.
4. Avoid real authors, unpublished essays, private notes, secrets, and private URLs.

Primary files:
- docs/voicebox-web-app-brainstorm.md

Required schema coverage:
- article slug, status, title, deck, category, tags, byline, dateline, excerpt, body blocks, citations, pull quotes, related slugs
- author name, role, avatar placeholder, bio
- category slug, label, rubric description
- draft id, status, edited timestamp, validation state

Review checks:
- Required fields are enough to render home, archive, article detail, drafts, and preview.
- Duplicate slug/id behavior is specified.
- Draft and published states are visibly separate.

Expected response:
1. List schema review checks.
2. Summarize fixture decisions.
3. Confirm all examples are synthetic.
```

### Completion checklist

- [ ] Schema covers all MVP routes.
- [ ] Validation rules are documented.
- [ ] Draft and published states are distinct.
- [ ] Fixture examples are synthetic.

## V1-T01 - Initialize App Shell And Test Harness

**Type:** Code
**Goal:** Create the minimal runnable web app with a passing base test command.

### Copy/paste prompt

```text
You are implementing V1-T01 for the VoiceBox editorial web app.

Task title: Initialize app shell and test harness
Task goal: Create the minimal runnable web app with a passing base test command.

Strict TDD:
1. Write a failing test for rendering the root route first.
2. Run the narrow test and confirm it fails because the route/app shell is missing.
3. Implement the smallest app shell needed to pass.
4. Run the narrow test until green.
5. Run the broader test command for the app.

Implementation requirements:
- Root route renders a real app surface, not a marketing placeholder.
- The page includes a heading, navigation landmark, main landmark, and a visible synthetic editorial fixture marker.
- Use existing project conventions for framework, package manager, and test runner when present.
- Do not add real analytics, external CMS calls, secrets, or private URLs.

Expected response:
1. Show the failing test name and expected failure.
2. Summarize implementation changes.
3. List test commands run and results.
4. Confirm no real content or secrets were added.
```

### Completion checklist

- [ ] Failing test was written first.
- [ ] Narrow red state was verified.
- [ ] Root route passes.
- [ ] Broader test command passes.
- [ ] No real content or secrets were added.

## V1-T02 - Add VoiceBox Design Tokens

**Type:** Code
**Goal:** Implement colors, fonts, spacing, borders, focus rings, and square component defaults from `DESIGN.md`.

### Copy/paste prompt

```text
You are implementing V1-T02 for the VoiceBox editorial web app.

Task title: Add VoiceBox design tokens
Task goal: Implement the key DESIGN.md tokens in app styling.

Strict TDD:
1. Write a failing token contract test, style snapshot, or component test first.
2. Verify it fails because the VoiceBox tokens or classes are missing.
3. Implement the smallest styling layer needed to pass.
4. Run the narrow test and broader style/app tests.

Required UX constraints:
- black #0A0A0A, red #EF4444, white #FAFAFA, neutral surfaces
- Archivo Black headings, Work Sans UI/body, Space Mono metadata
- 0px radius except avatar thumbnails
- no shadows
- no gradients
- 2px borders for interactive elements
- focus ring: 0 0 0 2px #FAFAFA, 0 0 0 4px #0A0A0A

Expected response:
1. Show the red test evidence.
2. Summarize token/style changes.
3. List test commands run.
4. Call out any framework limitation or deferred visual check.
```

### Completion checklist

- [ ] Token contract test was red first.
- [ ] Required tokens exist.
- [ ] No shadows or gradients are introduced.
- [ ] Broader tests pass.

## V1-T03 - Load Synthetic Editorial Fixtures

**Type:** Code
**Goal:** Load deterministic article, author, category, and draft fixtures.

### Copy/paste prompt

```text
You are implementing V1-T03 for the VoiceBox editorial web app.

Task title: Load synthetic editorial fixtures
Task goal: Load deterministic article, author, category, and draft fixtures.

Strict TDD:
1. Write failing fixture loader tests first.
2. Include duplicate slug/id, missing field, invalid related slug, and draft/published separation cases.
3. Verify the tests fail for missing loader or validation.
4. Implement the smallest loader and fixture set.
5. Run narrow tests and the broader suite.

Implementation requirements:
- Fixtures are local and synthetic.
- Published articles and drafts are visibly separate.
- Related article slugs must resolve.
- Loader errors must be actionable and safe.

Expected response:
1. Show red test names and failures.
2. Summarize fixture and loader changes.
3. List test commands run.
4. Confirm no real unpublished work was added.
```

### Completion checklist

- [ ] Fixture tests were red first.
- [ ] Duplicate and missing-field cases are covered.
- [ ] Fixtures load deterministically.
- [ ] No real unpublished work was added.

## V2-T01 - Build Editorial Home Page

**Type:** Code
**Goal:** Render a first-screen editorial front page using VoiceBox hierarchy.

### Copy/paste prompt

```text
You are implementing V2-T01 for the VoiceBox editorial web app.

Task title: Build editorial home page
Task goal: Render a first-screen editorial front page using VoiceBox hierarchy.

Strict TDD:
1. Write failing home page tests for featured article, latest list, category labels, and active nav.
2. Verify failure before production changes.
3. Implement the smallest page/component/template changes.
4. Run the narrow tests and broader suite.

UX requirements:
- Use a massive editorial headline.
- Use compact latest/rubric lists with bottom borders.
- Use only one prominent red accent in the viewport.
- Do not use decorative gradients, rounded cards, or shadow panels.
- The first screen must be the usable editorial app, not a marketing page.

Expected response:
1. Show red test evidence.
2. Summarize UI changes.
3. List tests run.
4. Confirm DESIGN.md constraints checked.
```

### Completion checklist

- [ ] Home page tests were red first.
- [ ] Featured and latest content render.
- [ ] Active navigation is visible.
- [ ] DESIGN.md constraints are respected.

## V2-T02 - Build Archive With Filters And Search

**Type:** Code
**Goal:** Let readers filter by category, author, tag, and text query.

### Copy/paste prompt

```text
You are implementing V2-T02 for the VoiceBox editorial web app.

Task title: Build archive with filters and search
Task goal: Let readers filter by category, author, tag, and text query.

Strict TDD:
1. Write failing tests for category filter, author filter, tag filter, text query, empty state, and selected chip state.
2. Confirm the expected red state.
3. Implement the smallest filtering and rendering changes.
4. Run narrow and broader tests.

UX requirements:
- Use square filter chips with uppercase labels.
- Selected chips invert to black.
- Empty state is useful and compact.
- Keyboard users can reach and operate filters.

Expected response:
1. Show red test evidence.
2. Summarize filter/search behavior.
3. List tests run.
4. Call out any deferred full-text search work.
```

### Completion checklist

- [ ] Filter tests were red first.
- [ ] Query and chip states are covered.
- [ ] Empty state is covered.
- [ ] Broader tests pass.

## V2-T03 - Add Responsive Navigation

**Type:** Code
**Goal:** Provide category navigation that works on mobile and desktop.

### Copy/paste prompt

```text
You are implementing V2-T03 for the VoiceBox editorial web app.

Task title: Add responsive navigation
Task goal: Provide category navigation that works on mobile and desktop.

Strict TDD:
1. Write failing tests for active route, keyboard operation, mobile menu open/close, and accessible names.
2. Verify red state.
3. Implement the smallest responsive navigation changes.
4. Run component/route tests and any available e2e tests.

UX requirements:
- Active nav uses the red underline.
- Text must not overlap at mobile or desktop widths.
- Controls use visible focus rings.
- Avoid rounded menu panels and shadows.

Expected response:
1. Show red test evidence.
2. Summarize navigation behavior.
3. List tests run.
4. Note viewport checks performed.
```

### Completion checklist

- [ ] Navigation tests were red first.
- [ ] Mobile and desktop states are covered.
- [ ] Keyboard behavior is covered.
- [ ] No overlap issues are observed.

## V3-T01 - Build Article Detail Route

**Type:** Code
**Goal:** Render article body blocks, headline, byline, dateline, category, and citations.

### Copy/paste prompt

```text
You are implementing V3-T01 for the VoiceBox editorial web app.

Task title: Build article detail route
Task goal: Render article body blocks, headline, byline, dateline, category, and citations.

Strict TDD:
1. Write failing tests for known slug rendering, missing slug not-found behavior, metadata, and citation rendering.
2. Verify the red state.
3. Implement the smallest route and rendering logic.
4. Run narrow and broader tests.

UX requirements:
- Headline uses display typography.
- Body copy remains readable and compact.
- Category overline is uppercase.
- Not-found state is safe and editorial, not a stack trace.

Expected response:
1. Show red test evidence.
2. Summarize route and rendering changes.
3. List tests run.
4. Confirm unknown articles are handled safely.
```

### Completion checklist

- [ ] Article tests were red first.
- [ ] Known and unknown slug cases are covered.
- [ ] Metadata and citations render.
- [ ] Broader tests pass.

## V3-T02 - Add Pull Quotes And Related Links

**Type:** Code
**Goal:** Render pull quotes with the red left border and related links with compact list styling.

### Copy/paste prompt

```text
You are implementing V3-T02 for the VoiceBox editorial web app.

Task title: Add pull quotes and related links
Task goal: Render pull quotes and related links on article pages.

Strict TDD:
1. Write failing tests for pull quote block rendering, missing quote content, valid related slugs, and invalid related slug validation.
2. Verify red state.
3. Implement the smallest block and related-link rendering changes.
4. Run narrow and broader tests.

UX requirements:
- Pull quote uses a 4px red left border.
- Related links use compact list rows and bottom borders.
- Red remains restrained.
- No cards inside cards.

Expected response:
1. Show red test evidence.
2. Summarize block rendering changes.
3. List tests run.
4. Confirm invalid related slugs are rejected or safely ignored.
```

### Completion checklist

- [ ] Pull quote tests were red first.
- [ ] Related links are validated.
- [ ] VoiceBox red accent rule is respected.
- [ ] Broader tests pass.

## V3-T03 - Add Reading Controls

**Type:** Code
**Goal:** Add accessible share, copy link, and reading progress controls without clutter.

### Copy/paste prompt

```text
You are implementing V3-T03 for the VoiceBox editorial web app.

Task title: Add reading controls
Task goal: Add accessible share, copy link, and reading progress controls without clutter.

Strict TDD:
1. Write failing tests for accessible names, keyboard activation, copy feedback state, and progress state.
2. Verify red state.
3. Implement the smallest controls.
4. Run narrow and broader tests.

UX requirements:
- Prefer familiar icons with tooltips where available.
- Do not add visible instructional text for obvious controls.
- Keep controls compact and square.
- Provide clear focus and feedback states.

Expected response:
1. Show red test evidence.
2. Summarize reading controls.
3. List tests run.
4. Call out browser API mocking used for copy/share behavior.
```

### Completion checklist

- [ ] Control tests were red first.
- [ ] Keyboard behavior is covered.
- [ ] Copy/share behavior is safely mocked.
- [ ] Broader tests pass.

## V4-T01 - Build Draft List And Draft States

**Type:** Code
**Goal:** Show drafts by status with clear editorial metadata.

### Copy/paste prompt

```text
You are implementing V4-T01 for the VoiceBox editorial web app.

Task title: Build draft list and draft states
Task goal: Show drafts by status with clear editorial metadata.

Strict TDD:
1. Write failing tests for draft sorting, status chips, empty state, and draft metadata.
2. Verify red state.
3. Implement the smallest draft list.
4. Run narrow and broader tests.

UX requirements:
- Status chips use DESIGN.md chip styling.
- Drafts are clearly separate from published articles.
- No route implies real production publishing.

Expected response:
1. Show red test evidence.
2. Summarize draft list behavior.
3. List tests run.
4. Confirm draft-only local workflow language.
```

### Completion checklist

- [ ] Draft list tests were red first.
- [ ] Status chips are covered.
- [ ] Draft/published separation is visible.
- [ ] Broader tests pass.

## V4-T02 - Build Local Draft Editor Form

**Type:** Code
**Goal:** Add local-only draft editing fields for title, rubric, excerpt, body, tags, and pull quote.

### Copy/paste prompt

```text
You are implementing V4-T02 for the VoiceBox editorial web app.

Task title: Build local draft editor form
Task goal: Add local-only draft editing fields for title, rubric, excerpt, body, tags, and pull quote.

Strict TDD:
1. Write failing tests for required title, required category/rubric, body length, tag parsing, field errors, and no accidental publish.
2. Verify red state.
3. Implement the smallest form state and validation logic.
4. Run narrow and broader tests.

UX requirements:
- Inputs are 44px high where appropriate, square, and bordered.
- Labels are uppercase Work Sans.
- Helper and error text follow DESIGN.md.
- Invalid drafts cannot be treated as publish-ready.

Expected response:
1. Show red test evidence.
2. Summarize form and validation behavior.
3. List tests run.
4. Confirm local-only draft scope.
```

### Completion checklist

- [ ] Validation tests were red first.
- [ ] Field errors are covered.
- [ ] No publish action is introduced.
- [ ] Broader tests pass.

## V4-T03 - Build Editorial Preview

**Type:** Code
**Goal:** Preview draft content using the same article layout as published articles.

### Copy/paste prompt

```text
You are implementing V4-T03 for the VoiceBox editorial web app.

Task title: Build editorial preview
Task goal: Preview draft content using the same article layout as published articles.

Strict TDD:
1. Write failing tests for preview rendering, draft status label, invalid draft fallback, and unsaved-change behavior if supported.
2. Verify red state.
3. Implement the smallest preview route and shared rendering reuse.
4. Run narrow and broader tests.

UX requirements:
- Preview reuses article typography and block rendering.
- Draft status is clearly marked.
- Preview does not imply content was published.

Expected response:
1. Show red test evidence.
2. Summarize preview behavior.
3. List tests run.
4. Call out any shared article rendering reuse.
```

### Completion checklist

- [ ] Preview tests were red first.
- [ ] Draft status is visible.
- [ ] Invalid draft fallback is covered.
- [ ] Broader tests pass.

## V5-T01 - Add Accessibility And Keyboard Tests

**Type:** Code
**Goal:** Check labels, landmarks, focus order, contrast-sensitive states, and keyboard navigation.

### Copy/paste prompt

```text
You are implementing V5-T01 for the VoiceBox editorial web app.

Task title: Add accessibility and keyboard tests
Task goal: Check labels, landmarks, focus order, contrast-sensitive states, and keyboard navigation.

Strict TDD:
1. Write a failing accessibility or keyboard test first.
2. Verify it fails for a real missing label, landmark, focus, or keyboard behavior.
3. Implement the smallest accessibility fix.
4. Run accessibility, keyboard, and broader tests.

Required coverage:
- home
- archive
- article detail
- drafts
- preview
- buttons, filter chips, navigation, form inputs

Expected response:
1. Show red accessibility evidence.
2. Summarize fixes.
3. List tests run.
4. Call out any residual manual accessibility checks.
```

### Completion checklist

- [ ] Accessibility test was red first.
- [ ] Key routes are covered.
- [ ] Keyboard behavior is covered.
- [ ] Broader tests pass.

## V5-T02 - Add Responsive Visual Checks

**Type:** Code
**Goal:** Verify mobile and desktop layouts do not overlap and keep the VoiceBox typographic rules.

### Copy/paste prompt

```text
You are implementing V5-T02 for the VoiceBox editorial web app.

Task title: Add responsive visual checks
Task goal: Verify mobile and desktop layouts do not overlap and keep the VoiceBox typographic rules.

Strict TDD:
1. Add a failing viewport, screenshot, or layout assertion test first.
2. Verify the red state.
3. Implement the smallest responsive CSS/layout change.
4. Run viewport checks and broader tests.

Required viewport checks:
- mobile home
- mobile article
- mobile draft editor or preview
- desktop home
- desktop archive
- desktop article

UX requirements:
- No text overlaps.
- No button text overflows.
- Headlines remain editorial but fit their containers.
- Red accent remains restrained per viewport.

Expected response:
1. Show red visual/layout evidence.
2. Summarize responsive fixes.
3. List viewport/test commands run.
4. Call out any manual screenshot review.
```

### Completion checklist

- [ ] Visual/layout test was red first.
- [ ] Mobile and desktop routes are covered.
- [ ] No overlap or overflow is observed.
- [ ] Broader tests pass.

## V5-T03 - Finalize Demo Script And Release Checklist

**Type:** Documentation
**Goal:** Document how to demo the app, what tests prove, and what remains out of scope.

### Copy/paste prompt

```text
You are implementing V5-T03 for the VoiceBox editorial web app.

Task title: Finalize demo script and release checklist
Task goal: Document how to demo the app, what tests prove, and what remains out of scope.

Instructions:
1. Update only Markdown needed for demo and release evidence.
2. Include the route walkthrough and test commands.
3. Include DESIGN.md constraints that were verified.
4. Keep limitations honest.
5. Do not claim production CMS, real publishing, analytics, or authentication.

Primary files:
- README.md
- docs/voicebox-web-app-brainstorm.md
- any dedicated demo/release doc created for VoiceBox

Review checks:
- Every demo claim maps to a route, UI state, fixture, or test.
- Out-of-scope items are explicit.
- TDD evidence is summarized.
- No real content, secrets, or private URLs were added.

Expected response:
1. List review checks performed.
2. Summarize docs changed.
3. List final test commands and results.
4. Confirm remaining known gaps.
```

### Completion checklist

- [ ] Demo route walkthrough is documented.
- [ ] Test evidence is documented.
- [ ] Non-goals are explicit.
- [ ] No unsupported production claims are made.
