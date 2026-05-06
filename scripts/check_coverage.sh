#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

MIN_TOTAL="${MIN_COVERAGE_TOTAL:-25.0}"

go test ./... -coverprofile=coverage.out -covermode=atomic >/dev/null
TOTAL="$(go tool cover -func=coverage.out | awk '/^total:/ {print $3}' | tr -d '%')"

awk -v got="$TOTAL" -v min="$MIN_TOTAL" 'BEGIN { if (got+0 < min+0) { exit 1 } }'

echo "coverage_total=${TOTAL}% (min=${MIN_TOTAL}%)"
