# Launch readiness checklist

Verify a release candidate is safe to ship: CI, migrations, observability, rollback, comms.

## When to use

- Tagging a release touching persistence, auth, billing, or infra configs.
- Widening a feature flag beyond internal testers.

## When not to use

- Docs-only or cosmetic release (still output `GO` only if CI scope is explicitly docs-only).

## Inputs

- `YOUR_SHA` or tag to deploy.
- `YOUR_ENV_ORDER` (e.g. staging then prod).
- `YOUR_ONCALL` rollback decision handle.

## Outputs

```
VERDICT: GO|NO_GO

NO_GO_ITEMS:
- reason | ticket or owner | "- none -"

CHECKS:
- [ ] CI green on YOUR_SHA
- [ ] migrations safe / ordered
- [ ] kill switch documented
- [ ] observability for new failure modes
- [ ] rollback artifact verified
- [ ] customer comms if needed

NOTES:
- bullet or "- none -"
```

## Steps

1. Exact revision and working tree cleanliness.

```bash
cd YOUR_REPO_ROOT
git rev-parse YOUR_SHA
git status --short
```

2. Migrations on live volume risk.

```bash
ls -lt YOUR_MIGRATION_DIR 2>/dev/null | head -n 15 || echo "Locate YOUR_MIGRATION_DIR"
```

3. Observability gaps.

```bash
grep -RInE 'TODO\(metrics\)|FIXME\(alert\)|stub metrics' -- YOUR_SERVICE_SRC 2>/dev/null | head -n 40 || true
```

4. Rollback tag.

```bash
git tag -l 'v*' | tail -n 8
```

## Stop and ask

Stop if `YOUR_SHA` cannot be resolved with `git rev-parse`.

## Reject if

- `VERDICT: GO` while `NO_GO_ITEMS` is non-empty.
- Destructive migration without rehearsed runbook ⇒ `NO_GO` (do not waive silently).

## Safety

- No prod credentials in output; secret references only.

## References

- Pair with `pr-review` or `pr-review-strict` for merge review; `plan-review` for pre-build scope.
