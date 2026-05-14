# Supabase / Postgres migration safety

Review SQL migrations targeting Postgres (including Supabase-hosted): ordering, locks, rollback realism, RLS implications, and destructive ops.

## When to use

- Pull requests touching `*.sql`, Supabase CLI migrations, or schema dumps promoted to prod.

## When not to use

- Local-only scratch databases—still avoid destructive defaults but urgency differs.

## Inputs

- Migration filenames ordering convention (`YYYYMMDDHHMMSS_description.sql`).
- Expected traffic pattern during deploy (online vs maintenance window).

## Outputs

- **SAFE / NEEDS REVISION** with bullet rationale referencing statements.

## Steps

1. Classify DDL: additive vs reordering columns vs drops/renames—drops **always** justify backups.

```bash
find supabase/migrations prisma/migrations db/migrations -name '*.sql' 2>/dev/null | tail -n 20
```

2. Lock risk: long `ALTER TABLE` on hot paths—prefer additive nullable columns then backfill jobs.

3. Rollback story: reversible migrations or paired down migrations—not “restore from backup” as only plan unless accepted.

4. RLS/policies: new tables exposed via PostgREST must ship policies aligned with tenant model—grep `ENABLE ROW LEVEL SECURITY`.

```bash
rg -n "ROW LEVEL SECURITY|CREATE POLICY" YOUR_SQL_DIR || true
```

5. Secrets: connection strings belong in CI/env—not committed `.env`.

## Safety

- Never paste prod DB URLs into GitHub issues—rotate if leaked.
- `DROP DATABASE`, `TRUNCATE ... CASCADE`, unchecked `DELETE` without `WHERE`—treat as incident-class unless staged.

## References

- Pair with `launch-ready` before tagging releases touching schema.
