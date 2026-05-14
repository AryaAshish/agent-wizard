# Shipping test plan (agent-wizard)

Goal: **Production-grade confidence** on the paths users hit daily—happy path, release binary, **failures**, **idempotency**, and **clean-environment** runs—without new CI systems (extend `go test`, existing scripts, optional one-shot Docker).

## Principles

1. Run phases **in order**; later phases assume earlier ones passed.
2. **Tag blocker:** P4 + **P6 canonical gate** + green-to-ship footer must pass for the **same commit** as the tag.
3. Prefer **hermetic tests** (`t.TempDir`, temp `HOME`); avoid relying on the maintainer’s laptop cache.

---

## Phase overview

| Phase | What | Pass |
|-------|------|------|
| **P0** | Zero-state / cache hygiene | Commands exit 0 |
| **P1** | Unit + integration (`go test`) | Exit 0 |
| **P2** | Docs + coverage scripts | Exit 0 |
| **P3** | Vulnerability scan | Exit 0 (or documented waiver) |
| **P4** | Release-binary smoke + **idempotency** | `smoke.sh` + repeat commands OK |
| **P5** | Embedded community E2E (`go test`) | Matching tests pass |
| **P5.1** | **Failure & error UX** | Assertions below |
| **P6** | **Manual gate** (install paths, team git, security) + **canonical demo blocks tag** | Checklist complete |

---

## P0 — Zero-state and repeatability prep

**Why:** Hidden `HOME` config, stale `go test` cache, or leftover `agentskills.yaml` in the repo root can make “green locally” a false positive.

**Run (pick all that apply before P1–P4 locally):**

```bash
go clean -testcache
# Ensure no accidental manifest in repo root used by verify_docs:
test ! -f agentskills.yaml || echo "WARN: agentskills.yaml in repo root may affect verify_docs"
```

**CI:** Already uses clean checkouts—no extra step.

**Optional clean shell (no hidden local deps):** from repo root, one-shot:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.22-bookworm bash -lc \
  'apt-get update -qq && apt-get install -y -qq git ca-certificates >/dev/null && go test ./... -count=1'
```

Use this when validating a release candidate if your host has exotic env vars or global git hooks (Docker needs network once for module download).

---

## P1 — Go tests

```bash
go test ./... -count=1
```

**Pass:** exit 0 on ubuntu + macos (+ windows if you ship Windows support).

---

## P2 — Docs + coverage

```bash
bash scripts/verify_docs.sh
bash scripts/check_coverage.sh
```

**Pass:** exit 0.

---

## P3 — govulncheck

```bash
govulncheck ./...
```

**Pass:** exit 0, or CI-equivalent job green.

---

## P4 — Release binary smoke + idempotency

**Smoke (real install shape):**

```bash
bash scripts/distribution/smoke.sh
```

**Pass:** ends with `distribution smoke passed` (release tarball + `install.sh` + npm wrapper + **init → list --filter → add → sync**).

**Idempotency (same project dir, no errors on repeat):** after `smoke.sh` succeeds, re-use its temp project is awkward because the script deletes it—instead run **manually** in a fresh dir (copy-paste):

```bash
proj="$(mktemp -d)" home="$(mktemp -d)"
export HOME="$home" PATH="$HOME/go/bin:$PATH"  # adjust if install dir differs
cd "$proj" && git init -q
agent-wizard init </dev/null
agent-wizard add pr-review --source community
agent-wizard sync
agent-wizard add pr-review --source community   # second add: no-op / stable
agent-wizard sync
agent-wizard sync                               # second sync: stable
test -f .agents/skills/pr-review/SKILL.md
```

**Pass:** second `add` / `sync` do not corrupt tree; `SKILL.md` still present.

**Backlog (optional):** fold the idempotency block into `scripts/distribution/smoke.sh` so CI enforces it every run.

---

## P5 — Embedded community E2E (hermetic)

```bash
go test ./internal/cli/... -count=1 -run 'TestEndUserFlow_EmbeddedCommunity'
```

**Pass:** all tests matching the run pass.

**Backlog:** add `TestEndUserFlow_EmbeddedCommunity_PackAddSync` (`pack add android-starter` + `sync` + assert pack skills on disk).

---

## P5.1 — Failure scenarios and error messages

Automated today (run every PR):

```bash
go test ./internal/cli/... -count=1 -run 'TestNegative_'
```

Covers **ambiguous skill** across sources and **strict-lock / drift** failure paths.

**Manual spot-checks before tag** (script until covered by tests—each must show a **clear error** and, where implemented, a **hint** such as `Try:` or `agent-wizard init`):

| Scenario | Suggested repro | Expect |
|----------|-----------------|--------|
| **Invalid / unknown skill id** | `add does-not-exist --source community` then `sync` | `sync` fails resolving skill; stderr/stdout mentions missing skill or similar (no silent success) |
| **Broken manifest** | Temp project: corrupt `agentskills.yaml` to invalid YAML, run `sync` | Non-zero exit; message points at manifest/config, not stack-only |
| **Missing / wrong source** | Manifest lists `sources: [community]` but global config has no such source (fresh `HOME` + hand-edited manifest) or `list --source-name nope` | Clear “source not found” style message |
| **add without manifest** | Empty dir, no `init`, run `add` | Hint path per current CLI (`init` guidance) |

**Pass:** each row behaves as expected; if output regresses, **block tag** until fixed or doc updated with known limitation.

---

## P6 — Manual gate + canonical demo (blocks tag)

Complete **once per release SHA**; paste evidence (checkboxes + date + commit) in the release PR or [`docs/metrics.md`](metrics.md).

### P6.0 — Canonical demo (required)

Same as launch plan: **Repo A** full path, **Repo B** repeat with identical commands; non-interactive `init` OK.

**Pass:** `pr-review` (or chosen demo skill) on disk in both trees within the timing bar you use for demos.

**Rule:** If P6.0 fails, **do not tag**—fix or document a known blocker explicitly.

### P6.1 — Install surfaces

- [ ] curl + `install.sh` → `agent-wizard --version` matches tag
- [ ] npm `npx` or global install path (optional but recommended once per major)

### P6.2 — Team git library

- [ ] `sources add --kind git` against `file://` or real test repo; `list` / `add` / `sync`

### P6.3 — Security

- [ ] Skim [`docs/security/threat-model.md`](security/threat-model.md); `git grep` for obvious secret patterns before tag

---

## Green to ship (summary)

1. **P0–P5 + P5.1** green on the tag commit (CI + any manual rows in P5.1).
2. **P4** `smoke.sh` green on that commit.
3. **P6** checklist done; **P6.0 canonical demo** explicitly signed off.
4. Second reviewer runs README happy path on that tag (not the author).

---

## Optional (Tier D / backlog)

- Golden CLI snapshots, fuzz YAML, load `sync`—defer until pain appears.
- Encode P5.1 manual rows as `go test` cases when stable.
- `workflow_dispatch` job that echoes P0–P4 commands for maintainers (no new infra—reuse scripts only).

---

## What “100% sure” means here

All **documented real-world classes** (happy, release, **failure**, **repeat**, **clean**) are either **automated** or **explicitly run and signed off** before the tag. Residual risk is unknown host quirks—address with patch releases and issue templates.
