# Conventional commits and release notes

Apply Conventional Commits (`feat:`, `fix:`, etc.) so changelog automation and semver bumps stay predictable—especially across squashed merges.

## When to use

- Libraries or CLIs consumed externally where release notes matter.
- Teams adopting semantic-release or manual Keep a Changelog updates fed by git history.

## When not to use

- Internal-only forks where commit noise exceeds signal—still prefix breaking changes explicitly.

## Inputs

- Allowed commit types/scopes for your org (`feat(api): ...`).
- Squash-vs-merge policy on default branch.

## Outputs

- Commit message examples + CHANGELOG section stub matching unreleased bucket.

## Steps

1. Write imperative subject ≤72 chars; body explains motivation and rollout notes.

```bash
git log --oneline -n 15
```

2. Breaking changes **must** announce token `BREAKING CHANGE:` in footer or use `feat!:` / `fix!:` style per team convention.

```bash
# Preview unreleased commits since last tag
git describe --tags --abbrev=0 2>/dev/null || echo "No tag yet."
git log $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~50)..HEAD --oneline
```

3. Map types to changelog buckets: feat→Added, fix→Fixed, perf→Changed—automate later; manual pass acceptable early.

```bash
rg "^feat:|^fix:|^perf:|^docs:|^chore:" <<< "$(git log -n 50 --pretty=%s)" || true
```

## Safety

- Don’t rewrite published history on shared default branches—communicate migrations instead.

## References

- Pair with `launch-ready` when tagging—ensure CHANGELOG unreleased section drains into tagged release.
