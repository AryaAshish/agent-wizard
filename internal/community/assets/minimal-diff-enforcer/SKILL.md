# Minimal diff enforcer

Produce the **smallest** change plan: fewest files, fewest hunks, one concern per commit message intent.

## When to use

- You have `YOUR_REQUESTED_CHANGE` and existing code context.

## When not to use

- You do not have the current file contents or diff baseline.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_FILES_IN_SCOPE` (may be one file)
- `YOUR_REQUESTED_CHANGE` (one sentence)
- Optional: `YOUR_SYMBOL` for `git grep` narrowing

## Outputs

```
DIFF_PLAN:
- file: YOUR_PATH
  lines/regions: concise description (no pasted full file)

JUSTIFICATION_PER_HUNK:
- why this hunk is necessary

NON_GOALS:
- what you will not touch

ROLLBACK_NOTE:
- one sentence
```

## Steps

1. Establish baseline.

```bash
cd YOUR_REPO_ROOT
git status --porcelain
git diff -- YOUR_FILES_IN_SCOPE
```

2. Identify smallest edit surface.

```bash
git grep -n 'YOUR_SYMBOL' -- YOUR_FILES_IN_SCOPE || true
```

3. If the request implies multiple concerns, split into **multiple outputs** (user must re-invoke per concern) — do not bundle in one plan.

## Stop and ask

Stop if `git diff` cannot be produced and no file contents were attached.

## Reject if

- `DIFF_PLAN` includes a file outside `YOUR_FILES_IN_SCOPE` without explicit expansion request.
- The plan contains a wide refactor when the request is local.

## Safety

- No destructive commands; prefer `git diff` / `grep`.
