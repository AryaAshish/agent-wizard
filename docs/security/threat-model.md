# Threat Model (Initial Outline)

## Scope

`agent-wizard` skill discovery, resolution, and sync workflows.

## Primary Risks

1. Malicious skill content from untrusted sources.
2. Source spoofing or tampering in remote fetch paths.
3. Archive extraction path traversal (`../` escape).
4. Unsafe sync writes that overwrite unintended files.

## Planned Mitigations

- Source trust policy and allowlists.
- Integrity metadata in lockfile entries.
- Safe extraction and path sanitization.
- Atomic writes and explicit prune mode for destructive operations.

## Verification Strategy

- Security regression tests for traversal and spoofed metadata.
- CI security scanners (CodeQL, vuln scanners, secret scanning).
