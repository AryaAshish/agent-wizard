# ADR 0001: v1 default project mode

## Status

Accepted

## Decision

Use manifest-only + `sync` as the v1 default workflow.

## Rationale

- Keeps repos smaller by default.
- Preserves deterministic behavior when used with a lockfile.
- Allows organizations to opt into vendoring when required.

## Consequences

- Teams must run `sync` in bootstrap flows unless they vendor skill files.
