# Contributing

Thank you for helping improve `agent-wizard`.

## Workflow

1. Open an issue for behavior or schema contract changes—use templates under `.github/ISSUE_TEMPLATE/` (`bug_report`, `feature_request`, `skill_request`).
2. Prefer **small commits** with [Conventional Commits](https://www.conventionalcommits.org/) prefixes (`feat:`, `fix:`, `docs:`, `chore:`).
3. Ensure `go test ./...` and `bash scripts/verify_docs.sh` pass locally before pushing.
4. Release candidates: execute **phases P1–P3 + P4 + P5 (incl. P5.1/P5.2 via `go test`)** then complete **manual P6** in [docs/test-plan-ship.md](docs/test-plan-ship.md)—**do not tag** without **P6.3** (canonical README demo ×2 repos **using the same install path as users**: e.g. `VERSION=vX.Y.Z curl …/install.sh | sh` or the published release asset, not only `@main`) and **P6.4** (Docker or fresh-environment smoke). Link your signed checklist in the release PR.

Maintainers triage **best-effort ~weekly**; critical regressions jump the queue.

## Adding or updating a bundled (“community”) skill

Skills ship from `internal/community/assets/<skill-id>/SKILL.md` — picked up via `go:embed`. Each skill id must match its directory name.

### Scaffold with the CLI (optional)

```bash
agent-wizard create-skill <skill-id>
```

Creates `./<skill-id>/SKILL.md` with the standard section skeleton. The **first paragraph under the `#` title** becomes the short summary shown in **`agent-wizard list`** (truncate ~80 characters). Fill that paragraph before your first `##` heading.

To contribute to the bundled library: fork the repo, copy the folder to `internal/community/assets/<skill-id>/`, open a PR—see the post-create hints and this section’s checklist.

**Launch-ready checklist**

1. **One workflow** per skill—no generic catch-alls.
2. Sections: *When to use*, *When not to use*, *Inputs*, *Outputs*, *Steps*, *Safety*.
3. At least **two** runnable shell/command blocks OR explicit `YOUR_*` placeholders with guidance.
4. State OS/tool assumptions (e.g. Linux/macOS, Docker required).
5. No filler prose—every paragraph should change behavior.

After edits: run `go test ./...`, refresh [`docs/SKILLS.md`](docs/SKILLS.md), and if packs reference your skill update `.agent-wizard-pack.yaml`.

Request ideas without a PR: open **Request a bundled skill**.

## PR checklist

- [ ] `go fmt ./...` on touched Go
- [ ] `go test ./...`
- [ ] `bash scripts/verify_docs.sh`
- [ ] README / CHANGELOG updated when user-visible behavior changes
- [ ] No secrets or machine-specific paths committed

## Code style

- Deterministic CLI behavior; backwards-compatible defaults unless bumping documented schema versions.
- Avoid introducing alternate CLI frameworks—stay consistent with existing stdlib `flag` routing unless agreed otherwise.

See also [SECURITY.md](SECURITY.md).
