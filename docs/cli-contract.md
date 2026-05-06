# CLI contract (v1)

## Grammar

- Prefer subcommands: `init`, `list`, `add`, `sync`, `status`, `sources`, `lock`, `migrate`, `cache`, `ci-check`, `browse`, `watch`, `import`, `pack`, `catalog`, `icp`.
- Flags complement subcommands; keep stable long flags over time.

## Machine output

- `status --json` emits a single JSON object suitable for CI parsing.
- Exit codes:
  - `0` success
  - `1` general failure
  - `3` lockfile drift detected (`status --check-drifts`)
