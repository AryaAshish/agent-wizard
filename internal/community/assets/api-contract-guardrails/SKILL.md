# API contract guardrails

Check **breaking changes** for HTTP/RPC/CLI surfaces: paths, payloads, status codes, and versioning notes.

## When to use

- A PR changes handlers, DTOs, OpenAPI/Proto, route tables, or client SDKs.

## When not to use

- Internal-only refactor with no exported route/schema change.

## Inputs

- `YOUR_BASE_REF` / `YOUR_HEAD_REF` (or patch listing changed API files).
- Optional: `YOUR_PUBLIC_BASE_URL` for doc examples only (no live calls required).
- `YOUR_API_PATHS` glob or directory list for greps (e.g. `internal/api ./proto`).

## Outputs

```
BREAKING_CHANGES:
- [yes|no] one-line justification

CLIENT_IMPACT:
- bullet or "- none -"

COMPATIBILITY_STRATEGY:
- version bump | deprecation window | feature flag | "- n/a -"

CHECKLIST:
- [ ] path/method change
- [ ] request field removed/renamed
- [ ] response field removed/renamed
- [ ] error shape change
- [ ] auth requirement change

DOC_UPDATES:
- file list or "- none -"
```

## Steps

1. List API-touched files.

```bash
git diff YOUR_BASE_REF...YOUR_HEAD_REF --name-only | grep -Ei 'openapi|swagger|proto|graphql|route|handler|api' || true
```

2. Scan for signature shifts.

```bash
git diff YOUR_BASE_REF...YOUR_HEAD_REF -- YOUR_API_PATHS | head -n 400
grep -RInE '@(Get|Post|Put|Patch|Delete)|rpc |service ' -- YOUR_API_PATHS 2>/dev/null | head -n 60
```

3. Mark breaking only with a cited hunk or schema diff line.

## Stop and ask

Stop if no API-related paths appear in the diff and user did not name `YOUR_API_PATHS`.

## Reject if

- `BREAKING_CHANGES: yes` without a cited removed/renamed field or route.

## Safety

- Read-only `git`/`grep`; do not call production endpoints.
