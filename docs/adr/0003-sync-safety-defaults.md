# ADR 0003: sync safety defaults

## Status

Accepted

## Decision

`sync` must support dry-run mode and avoid destructive pruning by default.

## Rationale

- Reduces accidental data loss.
- Enables CI policy checks with non-mutating behavior.

## Consequences

- Users who want pruning must explicitly opt in.
