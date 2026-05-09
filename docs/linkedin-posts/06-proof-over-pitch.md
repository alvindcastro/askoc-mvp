# LinkedIn Draft 06: Proof Over Pitch

## Audience

Hiring managers and CTOs who prefer working evidence over broad claims.

## Draft

I like portfolio projects that can be tested, not only explained. For AskOC AI Concierge, I wanted the proof to be easy for someone else to run.

`make smoke` starts the local stack and checks the unpaid transcript path plus the financial-hold CRM handoff. `make eval` runs JSONL cases for intent, source grounding, workflow actions, escalation, and safety behavior. `go test ./...` proves the Go packages.

Those commands are not decoration. They force the demo to behave like software. If source grounding breaks, if a workflow action changes, or if a safety case regresses, the project should show it.

AI helped me draft and inspect parts of the work, but the engineering loop stayed the same: write the failing test, make the behavior pass, keep the scope narrow, and show the evidence.

That is still the job: useful behavior, bounded scope, and proof someone else can run.

#tdd #golang #ai #automation #softwareengineering
