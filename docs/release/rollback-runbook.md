# Release rollback runbook

If a release breaks users:

1. Pin users to previous tag/binary in internal docs.
2. If lockfile schema is compatible, keep `agentskills.lock` unchanged and roll back CLI only.
3. If schema bumped, document `migrate` path and restore manifests from `.bak` backups when safe.
