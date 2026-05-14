# Go module CI baseline

Keep `go test` reliable in CI: caching, race policy, `-count=1`, module proxy hygiene.

## When to use

- `go.mod` present; adding or tightening CI.

## When not to use

- Non-Go repo (use ecosystem-specific skill).

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_GO_VERSION` (policy string, e.g. `1.22`)
- Integration tests: `YOUR_DOCKER_COMPOSE_FILE` or `none`

## Outputs

```
RECOMMENDATIONS:
- bullet (ordered)

WORKFLOW_SNIPPET_LINES:
- uses: actions/setup-go@YOUR_PIN
- with: go-version: YOUR_GO_VERSION, cache: true
- run: go test ./... -count=1

CACHE_KEY_NOTE:
- one line

RACE_POLICY:
- on|off|subset with reason

RISKS:
- bullet or "- none -"
```

## Steps

1. Local baseline.

```bash
cd YOUR_REPO_ROOT
go version
go mod verify
go test ./... -count=1 -short 2>&1 | tail -n 50
```

2. Race where affordable.

```bash
go test ./... -count=1 -race -timeout=10m 2>&1 | tail -n 50
```

3. Cache key material.

```bash
shasum -a 256 go.sum 2>/dev/null || sha256sum go.sum
```

## Stop and ask

Stop if `go.mod` is missing at `YOUR_REPO_ROOT`.

## Reject if

- Recommend `-race` on entire suite without timeout/shard when `go list ./...` implies very large packages—note `subset` instead.

## Safety

- Never echo `GOPROXY` tokens; CI secrets only.

## References

- Pair with `github-actions-matrix` for YAML placement.
