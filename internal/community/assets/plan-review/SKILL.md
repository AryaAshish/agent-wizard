# Technical plan review

Turn a written plan (RFC, tech spec, ticket epic) into **decisions**, **dependencies**, **cut-lines**, and **testable milestones**—before coding spreads across branches.

## When to use

- A multi-week initiative needs alignment across backend, frontend, infra, or data.
- You’re about to fork repos or schema—the cost of fixing direction later is high.

## When not to use

- Single-file bugfixes with obvious scope—ship with a normal PR description instead.

## Inputs

- Plan doc link or pasted outline.
- Fixed ship date vs flexible scope (say which wins).
- Known constraints: compliance, offline clients, SLA, migration downtime budget.

## Outputs

- **Approved scope** vs deferred backlog references (IDs).
- Risks with owners and mitigations—not a laundry list without owners.
- **Definition of done** that QA can execute without interpreting intent.

## Steps

1. Extract goals vs non-goals; reject passive voice hiding accountability (“will be handled”)—assign noun.

```bash
# Optional: sanity-check referenced tickets exist (adapt prefix)
grep -Eo '[A-Z]+-[0-9]+' YOUR_PLAN.md | sort -u
```

2. Dataflow diagram in prose: inputs → systems → persistence → outputs; note sync/async boundaries.

3. Migration / rollout ordering: backwards-compatible phases; expand **rollback** for each phase.

```bash
# If migrations live in-repo
find . -path '*migration*' -name '*.sql' 2>/dev/null | head
```

4. Test strategy: unit vs integration vs contract; what breaks first if wrong.

5. Schedule cut-lines: minimum viable slice that proves riskiest assumption—often integration or perf—not CRUD polish.

## Safety

- Plans that widen blast radius (multi-region, kernel flags, deleting data) require explicit downtime/RPO/RTO language—flag gaps instead of guessing.

## References

- Prefer linking existing ADRs over rewriting history—note superseded decisions explicitly.
