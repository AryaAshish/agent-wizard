# Dockerfile hardening pass

Safer default container: multi-stage, non-root, pinned bases, minimal runtime packages.

## When to use

- Shipping APIs/workers where breakout impact matters.
- Replacing `FROM *:latest` without digest policy.

## When not to use

- No `Dockerfile` path and no image build context named.

## Inputs

- `YOUR_REPO_ROOT`
- `YOUR_DOCKERFILE_PATH`
- Base policy: distroless | alpine | wolfi | other

## Outputs

```
CHANGES:
- layer or instruction | rationale (one line each)

USER_AND_PORTS:
- USER line target | exposed ports note

DIGEST_PIN:
- image@sha256:... or "semver tag only with reason"

REMOVED_RUNTIME_PKGS:
- bullet or "- none -"

BLOCKERS:
- bullet or "- none -"
```

## Steps

1. Read current Dockerfile.

```bash
cd YOUR_REPO_ROOT
test -f YOUR_DOCKERFILE_PATH && sed -n '1,200p' YOUR_DOCKERFILE_PATH || echo "missing YOUR_DOCKERFILE_PATH"
```

2. Scan for risk patterns.

```bash
grep -nE 'FROM |USER |^RUN |curl |wget |apt-get|apk add|secrets?' YOUR_DOCKERFILE_PATH || true
```

3. Optional digest pin (when reproducibility required).

```bash
docker pull YOUR_BASE_IMAGE_TAG 2>/dev/null && docker inspect --format '{{index .RepoDigests 0}}' YOUR_BASE_IMAGE_TAG || echo "Skipping digest: no docker"
```

## Stop and ask

Stop if `YOUR_DOCKERFILE_PATH` does not exist.

## Reject if

- Final image runs as root in production path without `BLOCKERS` documenting accepted exception.
- Build-time secrets via `ARG` for credentials without `BLOCKERS` + mitigation.

## Safety

- No `docker run` with host mounts unless user requests; planning focus.

## References

- Scan with `grype`/`trivy` in CI separately from this review.
