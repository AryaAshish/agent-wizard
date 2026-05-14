# Dependency upgrade advisor

Summarize **risk**, **changelogs**, and **verification commands** for bumping dependencies.

## When to use

- You plan to bump `YOUR_PACKAGE` or refresh a lockfile.

## When not to use

- You cannot name the package(s) or ecosystem (npm/go/maven/etc.).

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_ECOSYSTEM` (one of: go, npm, pnpm, yarn, maven, gradle, pip, cargo)
- `YOUR_PACKAGE` (or `YOUR_LOCKFILE_PATH`)

## Outputs

```
UPGRADE_TARGET:
- package: YOUR_PACKAGE
- from: YOUR_FROM_VERSION
- to: YOUR_TO_VERSION

RISK:
- low|medium|high with one-line reason

NOTABLE_CHANGES:
- bullet (max 5) or "- unknown -"

VERIFY:
1. command
2. command

ROLLBACK:
- one sentence (how to revert lockfile / mod)
```

## Steps

1. Detect current version from lock/manifest.

```bash
cd YOUR_REPO_ROOT
grep -RInF 'YOUR_PACKAGE' go.mod package.json package-lock.json pnpm-lock.yaml yarn.lock Cargo.toml pyproject.toml 2>/dev/null | head -n 40
```

2. If semver tags exist locally, inspect recent commits touching the module (optional).

```bash
git log -n 20 --oneline -- go.mod package.json 2>/dev/null
```

3. Map verification to repo scripts.

```bash
ls scripts 2>/dev/null; grep -RIn '"test"' package.json 2>/dev/null | head
```

## Stop and ask

Stop if neither `YOUR_PACKAGE` nor a lockfile path resolves in the repo.

## Reject if

- `NOTABLE_CHANGES` claims specifics without any changelog reference URL or quoted release note supplied by the user.

## Safety

- Do not run `curl | sh`; prefer read-only inspection unless user explicitly requests installs.
