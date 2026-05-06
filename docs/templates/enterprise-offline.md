# Template: enterprise / offline-first

1. Mirror skills as a local directory or Git repository accessible without public internet.
2. Use `sources add` with `--kind local` (or internal Git hosting with `--kind git`).
3. Prefer `installMode: manifest-only` and run `sync` in CI with `--strict-lock` once `agentskills.lock` exists.
4. Periodically run `agent-wizard status --check-drifts --strict-digest` in CI (exit code 3 signals drift).
