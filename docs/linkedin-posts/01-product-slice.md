# LinkedIn Draft 01: Product Slice

## Audience

Hiring managers, CTOs, and engineering leads who want to see practical product judgment.

## Draft

I have been playing with a small AI automation idea in Go: what if a student-services assistant behaved less like a demo chatbot and more like a real workflow tool? I kept the scope narrow on purpose, because narrow scope is where the engineering decisions become visible.

The repo focuses on transcript support. A learner asks how to order an official transcript, then asks why the request has not moved. The system answers from approved source chunks, checks synthetic Banner and payment records, triggers a local reminder workflow for unpaid payment, and routes financial holds to a mock CRM handoff.

That slice is small, but it touches the parts that matter: retrieval, classification, typed integrations, workflow actions, audit evidence, and fallback behavior. It also keeps human judgment where it belongs. The system does not approve requests, waive holds, or promise outcomes.

AI helped me move faster while building, mostly as a drafting and review assistant. The product shape still came from the boring questions: what data is allowed, what action is safe, what gets logged, and what must go to staff.

That is the kind of AI work I enjoy. Useful, bounded, and testable.

#golang #ai #automation #softwareengineering
