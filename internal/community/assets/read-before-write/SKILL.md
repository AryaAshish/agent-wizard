# Read before write

Enumerate **what you read** and **why the change belongs there** before proposing edits.

## When to use

- Any non-trivial edit touching shared modules, APIs, or migrations.

## When not to use

- Single-line typo in an isolated doc with no coupling.

## Inputs

- `YOUR_TARGET_FILES` (explicit paths).
- `YOUR_REPO_ROOT`.

## Outputs

```
FILES_REVIEWED:
- path | reason read

PATTERNS:
- convention or invariant observed (bullet)

CHANGE_JUSTIFICATION:
- why this location and not an alternative

OPEN_QUESTIONS:
- bullet or "- none -"
```

## Steps

1. Confirm files exist and read them (or declare cannot access).

```bash
cd YOUR_REPO_ROOT
for f in YOUR_TARGET_FILES; do test -f "$f" && wc -l "$f" || echo "missing $f"; done
```

2. Pull local conventions (imports, layering, error handling style).

```bash
grep -RIn 'package ' -- YOUR_TARGET_FILES 2>/dev/null | head
grep -RInE 'func main|TODO|FIXME' -- YOUR_TARGET_FILES 2>/dev/null | head -n 40
```

3. If dependencies are unclear, expand search once, then list **stop** questions.

## Stop and ask

Stop if **no** `YOUR_TARGET_FILES` list was provided.

## Reject if

- `FILES_REVIEWED` contains a path you did not actually inspect.
- `CHANGE_JUSTIFICATION` references modules not listed in `FILES_REVIEWED`.

## Safety

- Read-only commands; no `rm`, no `git reset --hard`.
