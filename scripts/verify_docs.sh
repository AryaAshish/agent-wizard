#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
go run . help >/dev/null
go run . catalog validate examples/curated-index.yaml >/dev/null
go run . status --json >/dev/null || true
