# Go module CI baseline

Keep `go test` reliable in CI for services using Go modules: caching, race detector policy, `-count=1`, and module proxy hygiene.

## When to use

- Adding or tightening CI for a Go repo (`go.mod` present).
- Flaky tests or mysterious cache hits hide real failures.

## When not to use

- Libraries consumed only via Go workspace with bespoke vendor rules—adapt proxy/vendor flags first.

## Inputs

- Go version requirement (must match `go` directive tolerance).
- Whether integration tests need Docker services (compose vs skip tags).

## Outputs

- Copy-paste workflow snippet aligned with team policy (race on/off, short vs full suite).

## Steps

1. Verify modules tidy locally—fail CI if `go.sum` drift creeps in.

```bash
go version
go mod verify
go test ./... -count=1 -short 2>&1 | tail -n 40
```

2. Optional race build where runtime allows (may exclude `-short` packages).

```bash
go test ./... -count=1 -race -timeout=10m 2>&1 | tail -n 40
```

3. CI cache keys: hash `go.sum` only—avoid caching `$GOPATH/pkg/mod` across Go upgrades blindly.

```bash
sha256sum go.sum | awk '{print $1}'
```

## Safety

- `-race` multiplies CPU—don’t enable on enormous suites without shards or timeouts.
- Never echo module proxy tokens—use CI secret stores only.

## References

- Pair with `github-actions-matrix` for YAML placement.
