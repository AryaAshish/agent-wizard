# Shipping test plan (agent-wizard)

Goal: **high confidence** before public releases—not “100%” in the mathematical sense (impossible for all hosts and future agents), but **no known critical path untested** and **CI + manual gates** documented.

## Principles

1. **Critical path first:** install → init → discover → add/sync → second project → lock (teams).
2. **Automate what is stable:** `go test`, scripts without network where possible.
3. **Quarantine flakiness:** network/curl/npm/real GitHub in **nightly** or **pre-tag** jobs, not necessarily every PR if flaky.
4. **Evidence on tag:** maintainer runs **release rehearsal** once per minor/major and pastes checklist into release notes or internal log.

---

## Tier A — Must pass on every PR (`go test ./...`)

| Area | Today | Gap to close |
|------|--------|----------------|
| Pure helpers | `internal/skills`, `internal/model`, `internal/engine` helpers | Keep table-driven tests when adding logic |
| CLI smoke | `internal/cli/run_test.go` | Expand only with stable assertions |
| **Embedded community** | `TestEndUserFlow_EmbeddedCommunity_ListFilterAddSync` | Add sibling: `pack add android-starter` + `sync` + assert 5 skill dirs |
| Local library + pack + lock | `TestEndUserFlow_LocalSource_Pack_Lock_Sync_StatusJSON` | Keep in sync when pack schema changes |
| Ambiguity / strict-lock | `TestNegative_*` in `e2e_test.go` | Add one test per new resolution rule |
| Community legacy config | `TestInitMigratesLegacyCommunityGitInGlobalConfig` | Extend if migration rules grow |

**Action:** Add **`TestEndUserFlow_EmbeddedCommunity_PackAddSync`** (init → `pack add android-starter` → `sync` → assert expected skill dirs from pack YAML).

---

## Tier B — Must pass before tagging a release (automated locally / CI job)

Run on **linux + macOS** at minimum (Windows for paths and git behavior where different).

| Check | Command / script | Owner |
|-------|------------------|--------|
| Full unit + integration | `go test ./... -count=1` | CI |
| Coverage floor | `bash scripts/check_coverage.sh` | CI |
| Docs contract | `bash scripts/verify_docs.sh` | CI |
| Vuln scan | `govulncheck ./...` | CI (already) |
| **Distribution smoke** | `bash scripts/distribution/smoke.sh` (if present; extend) | CI or `workflow_dispatch` |
| Perf smoke | `bash scripts/perf_smoke.sh` | CI nightly OK |

**Gap:** Ensure `distribution/smoke.sh` exercises **released binary** or **go install** path matching README—not only `go run .`.

---

## Tier C — Manual / pre-ship rehearsal (checklist)

Do once per **release candidate** (RC tag or `main` freeze before tag). Record date + git SHA in `docs/metrics.md` or release notes.

### C1 — Install paths (real world)

- [ ] **curl + install.sh** on clean shell (macOS or Linux): binary on `PATH`, `agent-wizard --version` matches expected tag.
- [ ] **npm** `npx @aryaashish/agent-wizard --version` (or global install) on machine with Node 18+.
- [ ] **go install** `github.com/aryaashish/agent-wizard@vX.Y.Z` from clean module cache optional.

### C2 — Canonical demo (launch plan)

- [ ] Repo A: non-interactive `init` → `list --source-name community --filter pr` → `add pr-review --source community` → `sync` → open `SKILL.md`.
- [ ] Repo B: repeat; confirm **no** manual re-copy from external doc.
- [ ] Optional: interactive `init` + picker path once (human).

### C3 — Team / git source (sharing)

- [ ] Create temp **bare** or regular git repo with 1–2 skills + optional pack on `file://` or private test org.
- [ ] `sources add --name team --kind git --git-url …`
- [ ] `agentskills.yaml` lists `team` in `sources` (or use `sources` + manifest edit as today).
- [ ] `list --source-name team`, `add … --source team`, `sync` — files present.

### C4 — Failure UX

- [ ] `sync` with bad skill id → error + **hint** still readable.
- [ ] `add` without manifest → hint mentions `init`.
- [ ] Offline or blocked cache: behavior acceptable (clear message).

### C5 — Security posture

- [ ] Skim `docs/security/threat-model.md` vs current CLI (no new network calls without doc).
- [ ] Confirm **no** secrets in repo: `.gitignore` + secret scan (GitHub secret scanning / `git grep` for keys).

---

## Tier D — Optional hardening (post v0.1 or if time)

| Item | Why optional |
|------|----------------|
| Golden CLI output snapshots | Brittle on Windows paths / ordering |
| Fuzz on YAML parsers | Diminishing returns vs manifest size |
| Matrix: old Go patch | Policy: support only documented Go version |
| Load test `sync` huge tree | Until users report scale |

---

## Definition of “green to ship”

1. **Tier A** green on **ubuntu + macos** (and **windows** if project commits to Windows support—already in CI matrix).
2. **Tier B** green for the same commit as the tag.
3. **Tier C** checklist completed for that tag (can be same day as tag; document in release PR).
4. **CHANGELOG** + **README** happy path re-run by someone who did **not** write the release (second pair of eyes).

---

## Implementation order (suggested)

1. **Add** `TestEndUserFlow_EmbeddedCommunity_PackAddSync` (highest ROI, still hermetic).
2. **Extend** `scripts/distribution/smoke.sh` (or add `scripts/e2e_install_smoke.sh`) to: build release binary with ldflags → run Tier A subset against binary.
3. **Add** GitHub **workflow_dispatch** workflow “Release rehearsal” calling install smoke + Tier C checklist job (markdown summary artifact).
4. **Document** in `CONTRIBUTING.md` PR section: “Release PRs must link completed Tier C checklist.”

---

## What “100% sure” means here

We ship when **known critical flows** are **automated** and **release-specific risks** are **explicitly signed off**—not when every theoretical host configuration has been tried. Unknowns after that are handled by **patch releases**, issue templates, and rollback docs.
