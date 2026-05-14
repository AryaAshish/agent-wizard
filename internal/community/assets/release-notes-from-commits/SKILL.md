# Release notes from commits

Turn a **commit range** into user-facing notes: features, fixes, breaking changes, migrations.

## When to use

- Preparing `YOUR_TAG` notes from `YOUR_FROM_TAG..YOUR_TO_REF`.

## When not to use

- Tags/refs missing and user cannot supply a commit range.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_FROM_TAG` (or SHA)
- `YOUR_TO_REF` (tag/branch/SHA)

## Outputs

```
HIGHLIGHTS:
- user-facing bullet (max 7)

FIXES:
- bullet or "- none -"

BREAKING:
- bullet or "- none -"

MIGRATIONS:
- bullet or "- none -"

UPGRADE_STEPS:
1. ordered step
2. ...

COMMITS_CONSIDERED:
- YOUR_FROM_TAG..YOUR_TO_REF
```

## Steps

1. Collect commits and conventional markers.

```bash
cd YOUR_REPO_ROOT
git fetch --tags 2>/dev/null || true
git log --oneline YOUR_FROM_TAG..YOUR_TO_REF | head -n 200
git log YOUR_FROM_TAG..YOUR_TO_REF --grep='BREAK' --oneline | head -n 50
```

2. Map modules/dirs to user impact.

```bash
git diff --name-status YOUR_FROM_TAG..YOUR_TO_REF | head -n 200
```

3. Drop internal-only churn (“chore”, “ci”) unless it affects consumers.

## Stop and ask

Stop if `git rev-parse YOUR_FROM_TAG` or `YOUR_TO_REF` fails.

## Reject if

- `HIGHLIGHTS` includes a claim not supported by any listed commit subject/body in range.

## Safety

- Read-only git; no force-push instructions.
