# ADR 0002: resolver conflict policy

## Status

Accepted

## Decision

Resolver defaults to error on ambiguous skill IDs across sources.

## Rationale

- Prevents hidden source-order bugs.
- Makes dependency provenance explicit.

## Consequences

- Users must namespace conflicting skills or remove ambiguity.
