# Technical plan review

Turn a written plan into **decisions**, **dependencies**, **cut-lines**, and **testable milestones** before coding spreads.

## When to use

- Multi-surface initiative (backend, frontend, infra, data).
- Schema or cross-repo fork risk.

## When not to use

- Single-file bugfix with obvious scope.

## Inputs

- `YOUR_REPO_ROOT` (for migration discovery commands).
- `YOUR_PLAN_PATH` or pasted outline.
- `YOUR_SCOPE_WINNER`: fixed_date | flexible_scope (pick one).
- Constraints: compliance, offline, SLA, migration downtime budget (bullets).

## Outputs

```
APPROVED_SCOPE:
- bullet

DEFERRED:
- id or description | reason

DEPENDENCIES:
- system or team | blocking yes/no

RISKS:
- risk | owner | mitigation

CUT_LINES:
- milestone | minimum proof

DEFINITION_OF_DONE:
- QA-executable bullet list

GAPS:
- bullet or "- none -"
```

## Steps

1. Goals vs non-goals; replace passive “will be handled” with owner nouns.

```bash
test -f YOUR_PLAN_PATH && grep -Eo '[A-Z]+-[0-9]+' YOUR_PLAN_PATH | sort -u || printf '%s\n' "No YOUR_PLAN_PATH"
```

2. Dataflow prose: inputs → systems → persistence → outputs; mark sync vs async.

3. Migrations / rollout phases + rollback per phase.

```bash
find YOUR_REPO_ROOT -path '*migration*' -name '*.sql' 2>/dev/null | head -n 30
```

4. Test strategy: what fails first if wrong (unit vs integration vs contract).

## Stop and ask

Stop if neither `YOUR_PLAN_PATH` nor a pasted plan body was provided.

## Reject if

- Any `RISKS` row lacks `owner` or `mitigation`.
- `DEFINITION_OF_DONE` contains subjective adjectives without observable checks.

## Safety

- Flag missing RPO/RTO/downtime language for high blast-radius plans; do not invent numbers.
