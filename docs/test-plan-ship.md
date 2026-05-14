# Shipping test plan — phased execution

**Goal:** High confidence before public releases—**not** exhaustive coverage of every host, but **no known critical path untested** plus **explicit pre-tag gates** (failure, idempotency, clean replay).

## Principles

1. **Critical path first:** install → discover → **`add`** (creates manifest + sync when missing, v0.1.3+) → optional **`init`** for picker / older releases → second project → team sources → lock.
2. **Automate what’s stable:** `go test`, `scripts/*.sh`; avoid flaky network on every PR where possible.
3. **Separate release rehearsal:** curling npm / human demo optional on schedule, **mandatory before tag**.
4. **Evidence on tag:** Maintainer records SHA + date in release PR or [`docs/metrics.md`](./metrics.md).

---

## Execution order & dependencies

| Phase | Depends on | Typical runner |
|-------|-------------|----------------|
| P0 | — | Contributor / CI |
| P1–P3 | P0 locally | CI + local |
| P4 | P1–P3 | CI (`distribution-dry-run`) |
| P5–P5.2 | P0 (`go test` temp dirs imply clean HOME) | CI |
| P6.* | Automated green on candidate SHA | Maintainer |
| P7 | — | Nightly OK |

Docs-only PRs may skip **P4** locally; **shipping a tag still requires P4 + P6.3 + P6.4** below.

---

## Phase map

| Phase | Objective | Run | Pass gate |
|-------|-----------|-----|-----------|
| **P0 — Zero-state** | No stale `HOME` / cache masking bugs | Prefer fresh temp dirs (`mktemp -d`). For interactive repeats: unset or point `HOME` to a disposable dir; **`agent-wizard cache prune`** if you reused the same HOME | Doc’d env knobs: `HOME`, `INSTALL_DIR`, `PATH`, optional `AGENT_WIZARD_CACHE_DIR` reset |
| **P1 — Unit + integration** | Full Go suite | From repo root: `go test ./... -count=1` | Exit **0** |
| **P2 — Docs + coverage** | Matches CI helpers | `bash scripts/verify_docs.sh` && `bash scripts/check_coverage.sh` | Both exit **0** |
| **P3 — Vuln baseline** | Dependency scan | `govulncheck ./...` (or rely on [.github/workflows/ci.yml](../.github/workflows/ci.yml)) | Exit **0** or waiver in release notes |
| **P4 — Release-binary smoke** | Tarball + `install.sh` + npm shim + CLI | `bash scripts/distribution/smoke.sh` (optional `VERSION=vX.Y.Z-rc`) | **`distribution smoke passed`** |
| **P5 — Hermetic E2E** | Embedded `community` in `go test` | `go test ./internal/cli/... -count=1 -run 'TestEndUserFlow_'` | Exit **0** |
| **P5.1 — Failure + hints** | Bad skill/manifest/source UX | Implemented as `TestCLI_ErrorHints` in [`internal/cli/e2e_test.go`](../internal/cli/e2e_test.go): unknown skill sync, **corrupt YAML on `add` load**, malformed YAML manifest on `sync`, unresolved manifest source alongside `community` | Each case: **non‑zero exit** and error text contains **`hint:`** (preferred) or unambiguous actionable substring |
| **P5.2 — Idempotency** | Repeated `add` / `sync` | `TestCLI_IdempotentAddSync` + **`smoke.sh` runs `sync` twice** | Second `sync` exits **0**; first skill tree still valid |
| **P6 — Manual pre-tagrollup** | See subsections below | Checklists | Boxes + sign-off |

### Optional **P7 — Perf / fuzz**

`bash scripts/perf_smoke.sh` (CI nightly). Fuzz/golden snapshots: defer until justified.

---

## P6 Manual gates (maintainer checklist)

Complete before **tag**. Record **git SHA + date**.

### P6.1 Install paths (real installs)

- [ ] **curl \| sh** README path: binary on `PATH`, `agent-wizard --version` matches tag.
- [ ] **npm** `npx @aryaashish/agent-wizard --version` (Node 18+).
- [ ] (Optional) `go install github.com/aryaashish/agent-wizard@vX.Y.Z`.

### P6.2 Team / git skill library

- [ ] Temporary git repo (`file://…` or test org) with ≥1 skill.
- [ ] `sources add --kind git`; manifest lists that source.
- [ ] `list --source-name <name>` → `add … --source <name>` → `sync`.

### **P6.3 Canonical demo gate (blocks tag)**

- [ ] Repo **A:** follow README “Happy path” fenced block end-to-end on the **same major.minor as the tag** (v0.1.3+ one-line `add`; **v0.1.2** per README version note: `init` then `add` then `sync`).
- [ ] Repo **B:** repeat in a fresh directory — same commands, synced `SKILL.md` appears.
- [ ] **Do not tag** if P6.3 cannot be completed **the same calendar day** as the candidate SHA (fix or postpone release).

Depends on successful **P4** for that toolchain **or** the exact `go install` / tarball used in production README.

### P6.4 Clean environment replay

Run **either**:

- **`docker run` one-shot** — mount repo, install minimal deps (**bash**, **nodejs** for npm block in smoke), execute `scripts/distribution/smoke.sh`; **or**
- Equivalent **fresh VM** / teammate laptop with no prior `agent-wizard` install.

Proves artifacts do not rely on hidden local state beyond documented inputs.

Frequency: at minimum before **first GA of a minor** / after materializing install changes.

### P6.5 Security & leakage

- [ ] Skim [Threat model](./security/threat-model.md) vs current CLI/network behavior.
- [ ] Confirm no credential files committed (`.gitignore` + upstream secret scanning).

---

## Automated regression inventory

| Automated | Covers |
|-----------|--------|
| `TestEndUserFlow_EmbeddedCommunity_ListFilterAddSync`, `PackAddSync` | Embedded community listing, single skill, **pack bundle** (`android-starter` five skills) |
| `TestCLI_ErrorHints`, `TestCLI_IdempotentAddSync`, `TestEndUserFlow_EmbeddedCommunity_AddColdStartNoInit` | **P5.1 / P5.2**: corrupt manifest on `add`, malformed YAML on `sync`, unknown skill sync, bogus second manifest source alongside `community`, duplicate add + double `sync`, **cold `add` without `init`** |
| `bash scripts/distribution/smoke.sh` (cold `add`, then two `sync` calls) | **P4 / idempotency** on released-style binary |

---

## Green to ship

1. **P1–P3 + P4 + P5 + P5.1 + P5.2** pass on **the tag commit** on **ubuntu-latest** (+ **windows-latest**/**macos-latest** via CI matrix).
2. **P6.3** + **P6.4** completed and noted (SHA/date).
3. **Second reviewer** re-runs README happy path (no substitutions).
4. **CHANGELOG** entry for the version.

---

## What “production-grade confidence” means

We ship when automated paths (P1–P5.2), release-shaped binary smoke (P4), and explicit human gates (P6—including **failure UX**, **repeatability**, and **canonical demo**) are satisfied—not when every possible OS/agent combination has been tried. Patch releases absorb residual edge cases.
