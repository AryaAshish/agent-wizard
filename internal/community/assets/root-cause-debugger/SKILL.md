# Root cause debugger

Turn **signals** (errors, flaky tests, metrics) into **at most three falsifiable hypotheses**, then eliminate until one root cause remains.

## When to use

- A failure is reproducible or has a captured stack/log snippet.
- You need a minimal fix hypothesis before editing code.

## When not to use

- No error text, no log excerpt, no failing command output, and no repro steps.

## Inputs

- `YOUR_ERROR_TEXT` (or CI log tail path `YOUR_LOG_FILE`).
- `YOUR_REPO_ROOT` (working tree with code).

## Outputs

```
PROBLEM_SUMMARY:
  one sentence

HYPOTHESES:
1. statement | validate with: one command or observation
2. ...
3. ... (use "- none -" for unused slots)

ELIMINATED:
- hypothesis N: reason

ROOT_CAUSE:
  one sentence tied to evidence

FIX:
  minimal change description (no large rewrite)

VALIDATION:
- ordered checks after fix
```

## Steps

1. Capture the first failure signature (not noise later in the log).

```bash
test -f YOUR_LOG_FILE && tail -n 200 YOUR_LOG_FILE || printf '%s\n' "YOUR_ERROR_TEXT"
```

2. Locate call sites / owners in repo.

```bash
cd YOUR_REPO_ROOT
grep -RInF 'YOUR_UNIQUE_TOKEN_FROM_ERROR' -- . 2>/dev/null | head -n 40
find . -maxdepth 4 -type f \( -name '*.go' -o -name '*.ts' -o -name '*.py' \) 2>/dev/null | head
```

3. For each hypothesis, define one **cheap** falsification (single command or reading one file).

## Stop and ask

Stop if **both** are missing: (a) error text or log tail, (b) repo root path.

## Reject if

- More than three hypotheses appear (merge weaker ones).
- A hypothesis has no falsification step.

## Safety

- Do not paste production secrets back into chat; reference paths only.
