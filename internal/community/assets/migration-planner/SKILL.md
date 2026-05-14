# Migration planner

Plan **ordered** schema/data migrations with **rollback** and **downtime** classification.

## When to use

- DB schema change, backfill, index add, or column type change.

## When not to use

- No target database kind (postgres/mysql/sqlite) and no migration tool detected.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_DB_KIND`
- `YOUR_MIGRATION_TOOL` (e.g. goose, atlas, flyway, prisma) or “unknown”

## Outputs

```
MIGRATION_STEPS:
1. forward step (idempotent where possible)
2. ...
3. ...

BACKWARD_PLAN:
- ordered rollback or "- destructive -"

LOCKS_AND_INDEXES:
- expected lock risk: low|medium|high

DATA_BACKFILL:
- yes|no | strategy

VALIDATION_QUERIES:
- sql snippet or "- n/a -"

STOP_IF:
- bullet conditions (e.g. row count > YOUR_THRESHOLD)
```

## Steps

1. Discover migration layout.

```bash
cd YOUR_REPO_ROOT
find . -maxdepth 5 -type d \( -name 'migrations' -o -name 'migrate' -o -name 'db' \) 2>/dev/null
grep -RInE 'CREATE TABLE|ALTER TABLE|CREATE INDEX' -- YOUR_MIGRATION_DIR 2>/dev/null | head -n 40
```

2. Inspect recent migrations for patterns (transactions, locks).

```bash
ls -lt YOUR_MIGRATION_DIR 2>/dev/null | head
```

3. If destructive, require explicit user acknowledgment in `STOP_IF` (do not proceed silently).

## Stop and ask

Stop if `YOUR_MIGRATION_DIR` cannot be identified and user gave no file paths.

## Reject if

- Forward plan adds a **non-null column without default** without a staged rollout.

## Safety

- Do not execute DDL against production; planning only unless user runs commands themselves.
