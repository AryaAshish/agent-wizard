# Release checklist

- [ ] `CHANGELOG.md` updated
- [ ] `go test ./...` green on CI matrix
- [ ] `scripts/verify_docs.sh` green
- [ ] Compatibility notes updated in `docs/compat.md`
- [ ] Threat model reviewed for any new remote-fetch surface
- [ ] Rollback drill notes updated in `docs/release/rollback-runbook.md`
