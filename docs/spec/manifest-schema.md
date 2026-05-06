# Manifest Schema (frozen v1)

## File

`agentskills.yaml`

## Required Fields

- `schemaVersion`: integer
- `targetDir`: string (default `.agents/skills`)
- `installMode`: string (default `manifest-only`)
- `sources`: list of named source aliases (must exist in `.agent-wizard-config.yaml`)

## Optional Fields

- `skills`: list of skill selectors (`id` or `source/id`)
- `packs`: list of pack ids resolved from the first configured library root
- `profiles`: list of output profiles; omitted profiles synthesize `default` targeting `targetDir`
- `hooks.preSync`, `hooks.postSync`: lists of shell commands (run with `-c`)

## Notes

- Ambiguous bare ids across multiple sources MUST error unless users qualify with `source/id`.
