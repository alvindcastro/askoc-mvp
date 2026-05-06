.PHONY: dev test test-race eval secret-check docker-build compose-up compose-down compose-test smoke

dev:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./internal/session

eval:
	go run ./cmd/eval -input data/eval-questions.jsonl -output reports/eval-summary.json -markdown-output reports/eval-summary.md -fail-on-critical

secret-check:
	scripts/check-secrets.sh

docker-build:
	docker build --build-arg APP=api -t askoc-api:local .

compose-up:
	docker compose up --build

compose-down:
	docker compose down --remove-orphans

compose-test:
	scripts/smoke.sh --base-url http://localhost:8080

smoke:
	scripts/smoke.sh --compose
