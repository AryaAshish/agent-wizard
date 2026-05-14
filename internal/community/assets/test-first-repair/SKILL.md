# Test-first repair

Define the failing scenario and **regression test shape** before describing code changes.

## When to use

- A bug with observable failure (test name, stack, assertion).

## When not to use

- No failing signal and no reproduction path.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_FAILING_COMMAND` (e.g. `go test ./... -run YOUR_TEST -count=1`)
- Optional: `YOUR_TEST_FILE`

## Outputs

```
FAILING_SCENARIO:
- expected vs actual (one bullet each)

REGRESSION_TEST:
- file: YOUR_TEST_PATH
- case name: YOUR_TEST_NAME
- assertion shape: bullet

FIX:
- minimal code change intent (no patch dump)

WHY_IT_WORKS:
- one sentence tied to root cause

SIDE_EFFECT_CHECK:
- bullet list
```

## Steps

1. Reproduce or extract the failure boundary.

```bash
cd YOUR_REPO_ROOT
YOUR_FAILING_COMMAND
```

2. Locate existing tests around the failure.

```bash
grep -RIn 'YOUR_TEST_NAME' -- . 2>/dev/null | head -n 40
```

3. If reproduction fails, stop — do not invent a fix.

## Stop and ask

Stop if `YOUR_FAILING_COMMAND` is missing or cannot be adapted to the repo’s language/tooling.

## Reject if

- `FIX` appears before `REGRESSION_TEST` is defined.
- `REGRESSION_TEST` names a file that does not exist and was not proposed as new.

## Safety

- Tests may run locally; do not require network unless user supplied credentials for a sandbox.
