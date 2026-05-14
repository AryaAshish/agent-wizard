# Bundled community skills

Skills ship inside the CLI under [`internal/community/assets/`](../internal/community/assets/). On **v0.1.3+**, `agent-wizard add <id> --source community` creates the manifest and runs **sync** by default. On **v0.1.2** and earlier, run **`init`**, then **`add`**, then **`sync`**.

Browse with **`agent-wizard list --source-name community`**: skills are sorted by id in two aligned columns (summary from the first paragraph under the title in each `SKILL.md`; missing summary shows `-`). Try **`agent-wizard create-skill <id>`** to scaffold a new skill locally before contributing (see [CONTRIBUTING.md](../CONTRIBUTING.md)).

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
| `pr-review-strict` | `internal/community/assets/pr-review-strict/` | PR findings with fixed severity schema |
| `root-cause-debugger` | `internal/community/assets/root-cause-debugger/` | Hypothesis triage from logs/errors |
| `read-before-write` | `internal/community/assets/read-before-write/` | Enumerate reads before proposing edits |
| `minimal-diff-enforcer` | `internal/community/assets/minimal-diff-enforcer/` | Smallest change plan per concern |
| `test-first-repair` | `internal/community/assets/test-first-repair/` | Regression test shape before fix |
| `api-contract-guardrails` | `internal/community/assets/api-contract-guardrails/` | Breaking API / schema change triage |
| `dependency-upgrade-advisor` | `internal/community/assets/dependency-upgrade-advisor/` | Bump risk + verify commands |
| `migration-planner` | `internal/community/assets/migration-planner/` | Ordered DB migrations + rollback |
| `security-review` | `internal/community/assets/security-review/` | Security triage from diff/grep |
| `release-notes-from-commits` | `internal/community/assets/release-notes-from-commits/` | User-facing notes from commit range |

**Pack:** `android-starter` bundles `pr-review`, `plan-review`, `launch-ready`, `cursor-rules-hooks`, `github-actions-matrix` (`agent-wizard pack add android-starter`).
