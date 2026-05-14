# Pull request review

Produce a structured PR review: correctness, tests, rollout risk, security-sensitive surfaces, and concrete follow-ups—without bikeshedding style-only nits unless your team mandates them.

## When to use

- You have a diff (GitHub/GitLab PR or local branch) and need a reviewer checklist before merge.
- You want consistency across repos so junior reviewers don’t forget rollback / observability / auth boundaries.

## When not to use

- Hot incident response where latency beats thoroughness—use your runbook skill instead.
- Pure formatting-only changes where policy already says “approve if CI green.”

## Inputs

- Link or branch name and base branch (e.g. `feature/foo` vs `main`).
- Risk tier: user-facing API, auth/payments, migrations, infra—flag **high** when any apply.
- Test evidence: CI link or paste failing tests.

## Outputs

- Ordered findings: **blockers**, **should-fix**, **nice-to-have**.
- Explicit **ship / don’t ship** with rollback note if ship.

## Steps

1. Summarize intent in one sentence from PR title/body only—do not invent product goals.

```bash
git fetch origin && git log --oneline origin/main..HEAD | head -n 20
git diff origin/main...HEAD --stat
```

2. Trace correctness on critical paths touched (handlers, serializers, migrations). List assumptions.

```bash
git diff origin/main...HEAD -- '*.go' '*.ts' '*.tsx' '*.sql'
```

3. Map tests to behavior changed—note gaps where coverage is story-only.

```bash
go test ./... -count=1 2>&1 | tail -n 30
```

4. Security pass: secrets in code, injection, authz boundaries, dependency changes with audit impact.

```bash
git diff origin/main...HEAD -- go.mod go.sum package.json package-lock.json
```

5. Rollout: feature flags, migrations order, backwards compatibility, observability (logs/metrics) for new paths.

## Safety

- Do not post production credentials, customer PII, or internal URLs in review comments—redact tokens and hostnames.
- If the diff introduces `eval`, shell-outs with user input, or dynamic SQL without binds—**block** until justified.

## References

- Align tone with team norms; when unsure, prefer asking one clarifying question over approving ambiguity.
