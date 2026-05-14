# Supabase / Postgres migration safety

Review SQL migrations: ordering, locks, rollback realism, RLS, destructive ops. For **generic DB rollout plans** (non-Supabase-specific), use bundled `migration-planner`.

## When to use

- PR touches `*.sql`, `supabase/migrations`, or schema promoted to prod.

## When not to use

- No migration files path and no SQL diff.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_SQL_DIR` (e.g. `supabase/migrations`)
- Traffic: `YOUR_DEPLOY_WINDOW` online | maintenance

## Outputs

```
VERDICT: SAFE|NEEDS_REVISION

ISSUES:
- [lock|data_loss|rls|rollback|secret] YOUR_FILE → one line

ROLLBACK_REALISM:
- sufficient | weak | destructive_accepted

RLS_NOTES:
- bullet or "- n/a -"

FOLLOW_UPS:
- bullet or "- none -"
```

## Steps

1. List tail of migration chain.

```bash
cd YOUR_REPO_ROOT
find YOUR_SQL_DIR -name '*.sql' 2>/dev/null | sort | tail -n 25
```

2. Scan destructive / policy patterns.

```bash
grep -RInE 'DROP TABLE|DROP COLUMN|TRUNCATE|DELETE FROM|ALTER TYPE|RENAME' -- YOUR_SQL_DIR 2>/dev/null | head -n 50
grep -RInE 'ROW LEVEL SECURITY|CREATE POLICY|ENABLE ROW LEVEL SECURITY' -- YOUR_SQL_DIR 2>/dev/null | head -n 40
```

3. Latest files deep read (first ~80 lines each).

```bash
ls -lt YOUR_SQL_DIR 2>/dev/null | head -n 8
```

## Stop and ask

Stop if `YOUR_SQL_DIR` does not exist and user gave no alternative glob.

## Reject if

- `VERDICT: SAFE` while `ISSUES` lists unresolved `data_loss` or missing rollback for destructive DDL without explicit `destructive_accepted` in `ROLLBACK_REALISM`.

## Safety

- Never paste prod DB URLs; rotate if leaked in chat.

## References

- Pair with `launch-ready` before tagging schema releases.
