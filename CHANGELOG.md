# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## Unreleased

- **v0.1.3 (planned patch):** `add` bootstraps a missing manifest (headless init) and runs **sync** by default (`--no-sync` to skip); README / `docs/SKILLS.md` call out **v0.1.2** vs newer one-liner happy path; Go **1.25.10** in CI + `toolchain`; `config.DefaultPath` honors **`HOME`** on Windows so tests and scripted installs isolate global config.
- README adoption funnel (five-line opener, differentiation table, single ≤20-line happy path).
- Ten bundled community skills with structured playbooks + [`docs/SKILLS.md`](docs/SKILLS.md).
- CLI: `list --filter`, tighter `init` next steps, actionable `sync`/`add` hints when manifests break.
- Contributor ergonomics: `.env.example`, expanded `.gitignore`, issue templates, ROADMAP, `docs/show-hn.md`, `docs/metrics.md`.
- Shipping: phased [`docs/test-plan-ship.md`](docs/test-plan-ship.md); e2e for embedded pack sync, CLI error hints, idempotent sync; malformed-manifest hint on `sync`; distribution smoke runs `sync` twice.

## [v0.1.2]

Maintenance tag line referenced prior to changelog alignment—install via GitHub Releases or `go install github.com/aryaashish/agent-wizard@v0.1.2`.

## Earlier

- Initial OSS-oriented CLI scaffolding: manifest/lockfile workflows, git/archive sources, drift checks, CI helpers, docs, and examples.
