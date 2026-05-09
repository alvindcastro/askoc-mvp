#!/usr/bin/env sh
set -eu

BASE_URL="${ASKOC_BASE_URL:-http://localhost:${ASKOC_API_PORT:-9080}}"
COMPOSE=0
KEEP_STACK=0
TIMEOUT_SECONDS="${ASKOC_SMOKE_TIMEOUT_SECONDS:-90}"

usage() {
  printf '%s\n' "usage: scripts/smoke.sh [--compose] [--keep-stack] [--base-url URL]"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --compose)
      COMPOSE=1
      ;;
    --keep-stack)
      KEEP_STACK=1
      ;;
    --base-url)
      shift
      if [ "$#" -eq 0 ]; then
        usage
        exit 2
      fi
      BASE_URL="$1"
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      exit 2
      ;;
  esac
  shift
done

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf 'missing required command: %s\n' "$1" >&2
    exit 1
  fi
}

cleanup() {
  if [ "$COMPOSE" -eq 1 ] && [ "$KEEP_STACK" -ne 1 ]; then
    docker compose down --remove-orphans >/dev/null 2>&1 || true
  fi
}

need curl
if [ "$COMPOSE" -eq 1 ]; then
  need docker
  trap cleanup EXIT INT TERM
  docker compose up --build -d
fi

wait_for_health() {
  deadline=$(( $(date +%s) + TIMEOUT_SECONDS ))
  while [ "$(date +%s)" -le "$deadline" ]; do
    if curl -fsS "$BASE_URL/healthz" >/dev/null 2>&1; then
      printf 'health check passed: %s/healthz\n' "$BASE_URL"
      return 0
    fi
    sleep 2
  done
  printf 'health check failed: %s/healthz did not respond within %s seconds\n' "$BASE_URL" "$TIMEOUT_SECONDS" >&2
  exit 1
}

post_chat() {
  message="$1"
  student_id="$2"
  output_file="$3"
  curl -fsS \
    -H 'Content-Type: application/json' \
    -X POST \
    -d "{\"channel\":\"web\",\"message\":\"$message\",\"student_id\":\"$student_id\"}" \
    "$BASE_URL/api/v1/chat" > "$output_file"
}

assert_contains() {
  file="$1"
  expected="$2"
  label="$3"
  if ! grep -q "$expected" "$file"; then
    printf 'smoke check failed: expected %s in chat response\n' "$label" >&2
    printf '%s\n' "response:" >&2
    sed -n '1,120p' "$file" >&2
    exit 1
  fi
}

UNPAID_RESPONSE="$(mktemp)"
HOLD_RESPONSE="$(mktemp)"
trap 'rm -f "$UNPAID_RESPONSE" "$HOLD_RESPONSE"; cleanup' EXIT INT TERM

wait_for_health

post_chat "I ordered my transcript but it has not been processed. My student ID is S100002." "S100002" "$UNPAID_RESPONSE"
assert_contains "$UNPAID_RESPONSE" "payment_reminder_triggered" "payment_reminder_triggered"

post_chat "My transcript request has a financial hold and is not moving. My student ID is S100003." "S100003" "$HOLD_RESPONSE"
assert_contains "$HOLD_RESPONSE" "financial_hold_detected" "financial_hold_detected"
assert_contains "$HOLD_RESPONSE" "crm_case_created" "crm_case_created"

printf '%s\n' "smoke checks passed"
