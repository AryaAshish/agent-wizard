# agent-wizard

`agent-wizard` is a **local-first** CLI for managing reusable **agent skills** (folders containing `SKILL.md`) across projects.

It is modeled after workflows popularized by tools like [`reseed`](https://github.com/nattergabriel/reseed), with emphasis on deterministic resolution, reproducible pinning, drift checks, and safe sync semantics.

## ICP / defaults (frozen)

- Primary ICP: **teams** (solo + enterprise supported).
- Default install posture: **`manifest-only` + `sync`** (skills materialized locally for agents).
- Ambiguous ids across independent sources **fail unless** you qualify with `source/id`.

See also: [`docs/support-matrix.md`](docs/support-matrix.md), [`docs/privacy.md`](docs/privacy.md), [`docs/cli-contract.md`](docs/cli-contract.md).

## Commands (overview)

Core:

- `init`, `migrate`, `status` (supports `--json`, `--check-drifts`)
- `list` (`--installed`, `--source-name`, `-S`)
- `add`, `remove`, `sync` (`--dry-run`, `--prune`, `--strict-lock`)

Sources:

- `sources list|add|remove` registering `local`, `git`, `archive`

Reproducibility:

- `lock` writes [`agentskills.lock`](docs/spec/lockfile-schema.md)
- Drift tooling via `status --check-drifts` (exit **3** on drift — intended for CI)

Utility:

- `cache status|prune`
- `ci-check` evaluates optional env gates (`AGENT_WIZARD_ALLOWED_SOURCES`, `AGENT_WIZARD_MIN_SCHEMA_VERSION`)
- `catalog validate <file>`
- `import --from DIR --into LIBRARY`
- `pack add PACK`
- `browse` (minimal numeric picker via stdin), `watch` (poll-sync loop)

## Quickstart

```bash
go test ./...
go run . init

# register a skill library checkout (examples included in-repo)
go run . sources add --name lib --kind local --path ./examples/library

# opt into the example pack or add individual skills
go run . pack add android-starter
# or: go run . add pr-review

go run . lock
go run . sync
go run . status --json
```

## Source kinds

| kind | config fields | notes |
|------|-----------------|------|
| `local` | `path` | fastest loop |
| `git` | `gitUrl`, `gitRef`, `subdir` | uses on-disk cache under XDG cache home |
| `archive` | `archiveUrl` | zip download + **safe extract** guards |

Threat/risk framing: [`docs/security/threat-model.md`](docs/security/threat-model.md)

## Profiles (multi-output)

Profiles live in [`agentskills.yaml`](docs/spec/manifest-schema.md). If omitted, a synthetic profile targets `targetDir`.

## Contributing / release

See [`CONTRIBUTING.md`](CONTRIBUTING.md), [`SECURITY.md`](SECURITY.md), [`docs/release/release-checklist.md`](docs/release/release-checklist.md).
