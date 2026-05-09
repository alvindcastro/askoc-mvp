# LinkedIn Draft 04: Workflow Boundaries

## Audience

CTOs and engineering leads who care about system boundaries, failure modes, and automation design.

## Draft

I have been paying attention to the least flashy part of this repo: what happens after the message arrives.

For the transcript-status slice, the Go API does not jump straight from message to action. It classifies the request, checks synthetic Banner and payment services through typed clients, and triggers a local reminder workflow only when the record is safe for self-service. If the synthetic record has a financial hold, the system creates a mock CRM case instead.

That boundary is deliberate. Routine follow-up can be automated, but risky work should move to staff with context. The assistant can make the next step clearer, but it should not waive a hold, approve a transcript, or hide a decision inside a friendly response.

AI helped me draft and revise pieces of the implementation, but the automation rules stayed explicit in code and tests. That is the part I would want a CTO to inspect.

#automation #golang #systemsdesign #ai
