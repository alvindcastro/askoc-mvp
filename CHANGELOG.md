# Changelog

All notable MVP task changes are recorded here with what changed, where it changed, when it changed, why it changed, and how it was completed.

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
