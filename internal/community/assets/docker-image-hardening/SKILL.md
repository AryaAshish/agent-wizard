# Dockerfile hardening pass

Produce a safer default container image for Go or generic services: multi-stage builds, non-root user, pinned bases, minimal runtime packages.

## When to use

- Shipping user-facing APIs or workers where container breakout impact is meaningful.
- Replacing `FROM ubuntu:latest` patterns without reproducibility.

## When not to use

- Internal ephemeral jobs where startup speed dominates and threat model excludes multi-tenant escape—still avoid `latest` tags.

## Inputs

- Base image policy (distroless, alpine, wolfi, etc.).
- Port and user UID requirements from platform (k8s securityContext).

## Outputs

- Dockerfile diff with rationale for each layer removed or merged.

## Steps

1. Multi-stage build: compile/build in builder; copy only binaries + static assets to runtime.

```dockerfile
# Example pattern — adapt paths and binary names
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
USER nonroot:nonroot
ENTRYPOINT ["/server"]
```

2. Pin digest optionally when reproducibility matters beyond semver tags.

```bash
docker pull gcr.io/distroless/static-debian12:nonroot
docker inspect --format '{{index .RepoDigests 0}}' gcr.io/distroless/static-debian12:nonroot
```

3. Drop package managers from runtime unless mandatory—fewer CVE surface updates.

## Safety

- Secrets via build-args leak into image history—use runtime injection (env, mounted files, KMS).
- Running as root in production containers—block unless documented exception.

## References

- Scan images in CI (`grype`, `trivy`) after Dockerfile changes—automate separately.
