# LinkedIn Draft 03: Source Grounding

## Audience

CTOs, hiring managers, and product-minded engineering leads evaluating AI reliability.

## Draft

I have been playing with the RAG part of this repo, and the main rule is simple: the model's confidence is not the product.

For a transcript answer, I wanted the system to show where the answer came from. The demo uses local source chunks from an allowlist. When a learner asks how to order an official transcript, the API returns source metadata, confidence, risk, and freshness signals. If the source support is weak, the system falls back instead of pretending.

That design makes the assistant less flashy, but more useful. A student-services workflow needs evidence that a reviewer can inspect. The source is part of the product surface, not an implementation detail hidden behind the chat response.

AI helped me draft and revise the language, but the application has to decide what counts as support. That boundary is where trust starts.

#ai #rag #softwareengineering #golang
