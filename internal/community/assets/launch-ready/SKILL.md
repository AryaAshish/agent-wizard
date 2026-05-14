# Launch readiness checklist

Verify a release candidate is safe to ship to production users: CI, migrations, observability, rollback, comms—without skipping “boring” operational checks.

## When to use

- Tagging a release that touches persistence, auth, billing, or infra configs.
- First launch of a feature flag path that widens audience beyond internal testers.

## When not to use

- Docs-only or cosmetic releases—run CI only and skip infra sections below.

## Inputs

- Release notes draft or changelog subset for this deploy.
- Target environment order (staging → prod).
- Owner on-call handle for rollback decisions.

## Outputs

- **GO / NO-GO** with explicit reasons for NO-GO items (tracked tickets acceptable only if risk accepted by owner).

## Steps

1. Confirm CI green on exact SHA being deployed.

```bash
git rev-parse HEAD
git status --short
```

2. Database migrations: applied order safe on live volume? Expand-only before contract changes?

```bash
# Example — adapt to your migration runner
ls -la YOUR_MIGRATION_DIR 2>/dev/null || echo "No migration dir — skip or locate tools."
```

3. Feature flags / kill switches documented—who can flip off without redeploy?

4. Observability: dashboards/alerts updated for new failure modes—not only happy-path logs.

```bash
# Quick grep for obviously missing HTTP handler wiring (adapt languages)
rg -n "TODO\(metrics\)|FIXME\(alert\)" YOUR_SERVICE_SRC || true
```

5. Rollback: prior artifact tag verified deployable; destructive migrations avoided or compensated.

```bash
git tag -l 'v*' | tail -n 5
```

6. Comms: status page / customer notice needed? Who posts?

## Safety

- Never paste prod credentials into tickets or chat logs—use secret references only.
- If rollback requires manual data repair—**NO-GO** until runbook exists and is rehearsed on staging.

## References

- Pair this with `pr-review` for the final merge and `plan-review` for pre-build scope alignment.
