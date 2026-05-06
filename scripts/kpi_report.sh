#!/usr/bin/env bash
set -euo pipefail

# Minimal KPI helper for pilot teams: captures repository-local smoke timing.
root="$(cd "$(dirname "$0")/.." && pwd)"

start="$(date +%s)"
(cd "$root" && go run . help >/dev/null)

end="$(date +%s)"
echo "cli_help_latency_sec=$((end-start))"
