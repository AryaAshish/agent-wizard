#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p .artifacts

# Performance smoke: keep bounded benchmark runtime but stable enough for regressions.
go test ./internal/engine -run '^$' -bench BenchmarkSyncDryRunMediumLibrary -benchmem -count=3 | tee .artifacts/perf-smoke.txt
