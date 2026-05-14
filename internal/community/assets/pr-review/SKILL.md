# Pull request review

Structured PR review: correctness, tests, rollout, security-sensitive surfaces, concrete follow-ups. For **severity-tagged machine-parseable** output only, use bundled `pr-review-strict` instead.

## When to use

- You have a diff (`YOUR_BASE_REF` vs `YOUR_HEAD_REF`, or PR with fetchable patch).
- You need a merge checklist before ship.

## When not to use

- Incident response where latency beats thoroughness.
- No diff and no file list.

## Inputs

- `YOUR_BASE_REF`, `YOUR_HEAD_REF` (or equivalent).
- `YOUR_RISK_TIER`: low | medium | high (auth, payments, migrations, infra ⇒ high unless disproven).
- `YOUR_TEST_CMD` (e.g. `go test ./... -count=1` or `npm test`—read-only).

## Outputs

```
SHIP: ship|dont_ship

BLOCKERS:
- bullet or "- none -"

SHOULD_FIX:
- bullet or "- none -"

NICE_TO_HAVE:
- bullet or "- none -"

ROLLBACK_NOTE:
- one sentence or "- n/a -"

OPEN_QUESTIONS:
- bullet or "- none -"
```

## Steps

1. Intent from title/body only—do not invent product goals.

```bash
git fetch origin
git log --oneline YOUR_BASE_REF..YOUR_HEAD_REF | head -n 25
git diff YOUR_BASE_REF...YOUR_HEAD_REF --stat
```

2. Critical paths: handlers, serializers, migrations.

```bash
git diff YOUR_BASE_REF...YOUR_HEAD_REF --name-only | head -n 200
git diff YOUR_BASE_REF...YOUR_HEAD_REF -- '*.go' '*.ts' '*.tsx' '*.sql' | head -n 400
```

3. Tests and dependency deltas.

```bash
YOUR_TEST_CMD
git diff YOUR_BASE_REF...YOUR_HEAD_REF -- go.mod go.sum package.json package-lock.json pnpm-lock.yaml yarn.lock 2>/dev/null | head -n 120
```

4. Security pass: secrets, injection, authz; if `eval`/shell-outs with user input or dynamic SQL without binds ⇒ list under `BLOCKERS`.

## Stop and ask

Stop if `git diff YOUR_BASE_REF...YOUR_HEAD_REF` is empty and no patch was supplied.

## Reject if

- `SHIP: ship` while any `BLOCKER` remains unaddressed or unaccepted by an explicit risk note (not allowed—use `dont_ship` or move to `OPEN_QUESTIONS` with owner).
- Any finding lacks a file or hunk reference from the diff.

## Safety

- Redact tokens, hostnames with customer data, PII from pasted logs.
