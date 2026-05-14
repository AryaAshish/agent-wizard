# GitHub Actions matrix workflow

Stand up or refactor GitHub Actions: matrix, setup caching, concurrency, test invocation.

## When to use

- Bootstrapping `.github/workflows/` for a repo.
- Replacing ad-hoc CI without caching or cancel-in-progress.

## When not to use

- Self-hosted runner topology undefined (labels, secrets) and user cannot supply them.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_WORKFLOW_PATH` (e.g. `.github/workflows/ci.yml`)
- Languages: `YOUR_PRIMARY_LANG` (e.g. go)
- `YOUR_GO_VERSION` if Go

## Outputs

```
FILE_ACTION:
create|update YOUR_WORKFLOW_PATH

YAML_BLOCKS:
CONCURRENCY:
- key line
MATRIX:
- os rows | tool versions
STEPS:
- ordered step summaries (no secrets)

SECRETS:
- required secret names only | none

SUPPLY_CHAIN:
- action pin note (major@sha or version policy)
```

## Steps

1. Inspect existing workflows.

```bash
cd YOUR_REPO_ROOT
find .github/workflows -name '*.yml' -o -name '*.yaml' 2>/dev/null | head
test -f YOUR_WORKFLOW_PATH && head -n 80 YOUR_WORKFLOW_PATH || true
```

2. Align test command with repo (Go example).

```bash
test -f go.mod && head -n 5 go.mod || echo "non-Go: set YOUR_TEST_CMD manually"
```

3. Concurrency + matrix template (adapt names).

```bash
grep -nE 'concurrency:|strategy:|matrix:' YOUR_WORKFLOW_PATH 2>/dev/null || true
```

## Stop and ask

Stop if `.github/` must be created and user did not confirm default branch name for `concurrency.group`—ask for `YOUR_DEFAULT_BRANCH`.

## Reject if

- YAML embeds long-lived PATs or private SSH keys.
- Matrix explodes (>6 cells) without `fail-fast` or cost justification in `YAML_BLOCKS`.

## Safety

- Prefer `GITHUB_TOKEN` with least permissions; OIDC where applicable.

## References

- Pair with `go-ci-module` for Go flags (`-race`, `-short`, timeouts).
