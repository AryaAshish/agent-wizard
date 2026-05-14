# Bundled community skills

Skills ship inside the CLI under [`internal/community/assets/`](../internal/community/assets/). On **v0.1.3+**, `agent-wizard add <id> --source community` creates the manifest and runs **sync** by default. On **v0.1.2** and earlier, run **`init`**, then **`add`**, then **`sync`**.

| Skill id | Path | Purpose |
|----------|------|---------|
| `pr-review` | `internal/community/assets/pr-review/` | Structured PR review checklist |
| `plan-review` | `internal/community/assets/plan-review/` | RFC/plan alignment before build |
| `launch-ready` | `internal/community/assets/launch-ready/` | Release GO/NO-GO gates |
| `go-ci-module` | `internal/community/assets/go-ci-module/` | Go modules CI hygiene (`go test`, caching) |
| `github-actions-matrix` | `internal/community/assets/github-actions-matrix/` | GitHub Actions workflow scaffold |
| `docker-image-hardening` | `internal/community/assets/docker-image-hardening/` | Safer multi-stage Dockerfiles |
| `supabase-migration-safety` | `internal/community/assets/supabase-migration-safety/` | Postgres / Supabase migration review |
| `cursor-rules-hooks` | `internal/community/assets/cursor-rules-hooks/` | Cursor rules, hooks, skill paths |
| `mcp-server-transports` | `internal/community/assets/mcp-server-transports/` | MCP stdio vs HTTP choices |
| `conventional-commits-release` | `internal/community/assets/conventional-commits-release/` | Commit discipline → changelog |

**Pack:** `android-starter` bundles `pr-review`, `plan-review`, `launch-ready`, `cursor-rules-hooks`, `github-actions-matrix` (`agent-wizard pack add android-starter`).
