#!/usr/bin/env sh
set -eu

if ! command -v git >/dev/null 2>&1; then
  printf '%s\n' "missing required command: git" >&2
  exit 1
fi

patterns='ghp_|xoxb-|AIza|sk-live-'

found=0
for file in $(git ls-files --cached --others --exclude-standard); do
  case "$file" in
    scripts/check-secrets.sh|internal/build/p10_artifacts_test.go)
      continue
      ;;
  esac

  if [ ! -f "$file" ]; then
    continue
  fi

  case "$file" in
    .env)
      printf '%s\n' ".env must stay local and ignored; remove it from git inputs" >&2
      found=1
      ;;
    .env.example)
      ;;
  esac

  if grep -E "$patterns" "$file" >/dev/null 2>&1; then
    printf 'secret-looking value found in %s\n' "$file" >&2
    found=1
  fi
done

if [ "$found" -ne 0 ]; then
  exit 1
fi

printf '%s\n' "secret check passed"
