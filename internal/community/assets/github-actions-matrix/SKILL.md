# GitHub Actions matrix workflow

Stand up or refactor a GitHub Actions workflow for Go (or polyglot) repos: OS/arch matrix, Go/setup caching, concurrency cancellation, and artifact hygiene.

## When to use

- Bootstrapping CI on GitHub with reproducible installs across Linux/macOS/Windows (subset optional).
- Replacing hand-rolled checkout scripts that skip caching.

## When not to use

- Self-hosted runners with proprietary secrets topology—adapt runner labels first.

## Inputs

- Minimum Go version matrix rows you truly need (fewer rows = faster feedback).
- Required secrets (`GITHUB_TOKEN` vs none).

## Outputs

- Workflow YAML skeleton checked into `.github/workflows/ci.yml`.

## Steps

1. Scaffold workflow with concurrency cancel-in-progress per branch.

```yaml
# .github/workflows/ci.yml (snippet — merge into full jobs block)
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true
```

2. Matrix strategy—avoid explosive combinations; prefer `ubuntu-latest` for default PR CI.

```yaml
strategy:
  fail-fast: false
  matrix:
    os: [ubuntu-latest]
    go-version: ["1.22.x"]
```

3. Cache modules via `actions/setup-go` with `cache: true` keyed by `go.sum`.

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: ${{ matrix.go-version }}
    cache: true
```

4. Run tests non-interactively—mirror local `-count=1`.

```yaml
- run: go test ./... -count=1
```

## Safety

- Never embed repo SSH keys or long-lived PATs in YAML—use OIDC or scoped secrets with rotation.
- Third-party actions: pin major versions consciously; review supply-chain implications.

## References

- Combine with `go-ci-module` for Go-specific flags (`-race`, `-short`).
