# PR review (strict output)

Use when you need **severity-tagged, file-scoped findings** with **no generic praise**. For a broader checklist, use bundled `pr-review` first; use this skill when output must match the schema below exactly.

## When to use

- You have a PR diff (or branch range) and want machine-parseable review output.
- You must not restate code and must prioritize high-impact issues only.

## When not to use

- You lack any diff or file list (no `git` range, no patch, no PR URL with fetchable diff).
- Style-only nits are out of scope unless the team policy explicitly requires them.

## Inputs

- `YOUR_BASE_REF` and `YOUR_HEAD_REF` (or `YOUR_PR_URL` + instructions to fetch diff).
- Optional: `YOUR_RISK_HINT` (e.g. auth, payments, migrations).

## Outputs

Produce **only** this structure (no preamble, no closing pleasantries):

```
RISK_LEVEL: low|medium|high

ISSUES:
- [BLOCKER|SHOULD|NICE] YOUR_PATH:YOUR_LINE → one-line issue
  explanation: one sentence
  fix: one concrete fix (file + change intent)

MISSING_TESTS:
- bullet or "- none -"

SUMMARY_ONE_LINE:
```

If you cannot fill `ISSUES` truthfully from the diff, write `ISSUES: - none -` and lower `RISK_LEVEL` only when justified.

## Steps

1. Materialize the diff you are allowed to review.

```bash
git fetch origin
git diff YOUR_BASE_REF...YOUR_HEAD_REF --stat
git diff YOUR_BASE_REF...YOUR_HEAD_REF
```

2. Map changed surface area (entries, serializers, migrations, auth).

```bash
git diff YOUR_BASE_REF...YOUR_HEAD_REF --name-only | head -n 200
grep -RInE 'password|secret|token|api_key|BEGIN RSA|AWS_' -- YOUR_CHANGED_PATHS 2>/dev/null | head -n 50 || true
```

3. For each issue, cite **path:line** only if you verified the line in the diff; otherwise omit line.

## Stop and ask

Stop and request **base/head refs or PR diff** if `git diff` is empty and no patch was supplied.

## Reject if

- You would output generic advice without a cited change hunk.
- You would invent `YOUR_PATH:YOUR_LINE` not present in the supplied diff.

## Safety

- Do not execute destructive git commands; read-only inspection only.
- Treat logs and pasted secrets as sensitive; redact in explanations.
