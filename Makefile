.PHONY: dev test test-race eval

dev:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./internal/session

eval:
	go run ./cmd/eval -input data/eval-questions.jsonl -output reports/eval-summary.json -markdown-output reports/eval-summary.md -fail-on-critical
