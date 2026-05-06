# Lockfile Schema (frozen v1)

## File

`agentskills.lock`

## Required Fields

- `schemaVersion`: integer
- `generatedAt`: RFC3339 timestamp
- `entries`: list of locked skill records

## Entry Fields

- `skillId`
- `sourceName`
- `resolvedRef` (commit SHA or immutable version)
- `digest` (optional integrity checksum)

## Notes

- Compatible with [`agentskills.yaml`](manifest-schema.md) selectors (bare ids vs `source/id`).
