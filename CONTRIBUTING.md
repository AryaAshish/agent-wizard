# Contributing

Thank you for helping improve `agent-wizard`.

## Workflow

1. Discuss larger changes via an issue when behavior or schema contracts change.
2. Prefer small commits and clear messages (optional Conventional Commits encouraged).
3. Ensure `go test ./...` and `scripts/verify_docs.sh` pass locally.

## Code style

- Format with `gofmt`/`go fmt`.
- Keep CLI behavior deterministic and backwards compatible unless bumping documented schema versions.
