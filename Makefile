.PHONY: dev test test-race

dev:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./internal/session
