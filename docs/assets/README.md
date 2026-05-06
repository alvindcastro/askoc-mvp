# Portfolio Capture Placeholders

P11 uses text placeholders instead of committed binary screenshots until final captures are generated from the local synthetic stack. This keeps the portfolio package reviewable and avoids accidentally committing real learner data, real tokens, private URLs, or browser-profile details.

Before adding any `.png`, `.jpg`, `.gif`, or `.webm` file in this directory, apply the screenshot/GIF review rules in [Privacy Impact Lite](../privacy-impact-lite.md#portfolio-screenshot-and-gif-review).

| Placeholder | Future file | Capture source | Caption |
|---|---|---|---|
| Chat grounded answer | `chat-grounded-answer.png` or `chat-grounded-answer.gif` | `/chat` after asking “How do I order my official transcript?” | Tier 0 answer grounded in approved transcript source chunks with confidence, risk, freshness, and safe action trace. |
| Transcript workflow | `transcript-workflow.png` or `transcript-workflow.gif` | `/chat` after the `S100002` transcript-status prompt | Mock Banner/payment checks and idempotent payment reminder workflow using synthetic IDs only. |
| CRM escalation | `crm-escalation.png` or `crm-escalation.gif` | `/chat` after the `S100003` financial-hold prompt or urgent-sentiment prompt | Mock CRM handoff with queue, priority, synthetic case ID, and privacy-aware summary. |
| Admin dashboard | `admin-dashboard.png` or `admin-dashboard.gif` | `/admin` after demo scenarios | Aggregate containment, escalation, workflow, stale-source, low-confidence, and review-queue metrics without raw PII. |
| Eval report | `eval-report.png` | `reports/eval-summary.md` after `make eval` | Repeatable responsible-AI quality evidence with zero critical safety failures. |

Allowed visible values:

- synthetic student IDs such as `S100002`,
- mock workflow or CRM identifiers with `SYNTH-` or `MOCK-` prefixes,
- local URLs such as `localhost:8080`,
- placeholder admin token text such as `demo-admin-token`.

Disallowed visible values:

- real learner names, emails, phone numbers, passwords, payment data, government-style IDs, or student records,
- real API keys, provider endpoints, webhook URLs, `.env` contents, or shell history,
- private portal pages or authenticated institutional systems.
