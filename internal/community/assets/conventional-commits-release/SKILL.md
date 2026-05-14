# Conventional commits and release notes

Apply Conventional Commits so changelog automation and semver stay predictable. For **tag-range user-facing** notes, use bundled `release-notes-from-commits` after messages are clean.

## When to use

- External consumers care about semver/changelog.
- Squash or merge policy on default branch is defined.

## When not to use

- You cannot access `git log` on the repo.

## Inputs

- `YOUR_REPO_ROOT`
- Allowed types/scopes (e.g. `feat(api):`).
- `YOUR_SQUASH_POLICY`: squash | merge | mixed

## Outputs

```
MESSAGE_RULES:
- subject rule (≤72 chars imperative)
- breaking rule (footer token or !: style)

EXAMPLE_SUBJECT:
feat(YOUR_SCOPE): YOUR_SUBJECT

EXAMPLE_BREAKING_FOOTER:
BREAKING CHANGE: YOUR_BREAK_DESC

CHANGELOG_STUB_HEADINGS:
- ## Unreleased
- ### Added | Fixed | Changed (only if commits justify)

COMMIT_SCAN:
- notable violations (max 10) or "- none -"
```

## Steps

1. Recent style sample.

```bash
cd YOUR_REPO_ROOT
git log --oneline -n 20
```

2. Since last tag.

```bash
git describe --tags --abbrev=0 2>/dev/null || echo "YOUR_FALLBACK_REF"
git log $(git describe --tags --abbrev=0 2>/dev/null || echo YOUR_FALLBACK_REF)..HEAD --oneline | head -n 80
```

3. Type prefix scan.

```bash
git log -n 80 --pretty=%s | grep -E '^(feat|fix|perf|docs|chore)(\(|):' | head -n 40 || true
```

## Stop and ask

Stop if not a git checkout (`YOUR_REPO_ROOT` has no `.git`).

## Reject if

- `CHANGELOG_STUB` invents shipped features not present in `COMMIT_SCAN` subjects.

## Safety

- Do not recommend `git push --force` on shared default branch.
