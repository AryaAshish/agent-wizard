# Resolver Precedence (frozen v1)

## Default Behavior

- Resolution is deterministic.
- Ambiguous skill IDs across sources return an error.
- Explicitly namespaced skill IDs may be used for override behavior.

## Source Priority

1. Project-local overrides
2. Organization source
3. Community source

## Conflict Handling

- Default: fail fast with explicit error.
- Future option: opt-in override mode via namespaced identifiers.
