.PHONY: dev test

dev:
	go run ./cmd/api

test:
	go test ./...
