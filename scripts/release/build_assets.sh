#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DIST_DIR="${1:-$ROOT_DIR/dist}"
VERSION="${VERSION:-}"
COMMIT="${COMMIT:-$(git -C "$ROOT_DIR" rev-parse --short HEAD)}"
DATE="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

if [[ -z "$VERSION" ]]; then
  VERSION="$(git -C "$ROOT_DIR" describe --tags --abbrev=0 2>/dev/null || true)"
fi
if [[ -z "$VERSION" ]]; then
  VERSION="dev"
fi
VERSION="${VERSION#v}"

mkdir -p "$DIST_DIR"
rm -f "$DIST_DIR"/agent-wizard_* "$DIST_DIR"/checksums.txt

build_target() {
  local goos="$1"
  local goarch="$2"
  local ext=""
  local archive_ext="tar.gz"
  local bin_name="agent-wizard"
  if [[ "$goos" == "windows" ]]; then
    ext=".exe"
    bin_name="agent-wizard.exe"
    archive_ext="zip"
  fi
  local base="agent-wizard_${VERSION}_${goos}_${goarch}"
  local out_dir
  out_dir="$(mktemp -d)"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build \
      -ldflags="-s -w -X github.com/aryaashish/agent-wizard/internal/buildinfo.Version=v${VERSION} -X github.com/aryaashish/agent-wizard/internal/buildinfo.Commit=${COMMIT} -X github.com/aryaashish/agent-wizard/internal/buildinfo.Date=${DATE}" \
      -o "${out_dir}/${bin_name}" \
      "$ROOT_DIR"
  if [[ "$archive_ext" == "zip" ]]; then
    (cd "$out_dir" && zip -q "$DIST_DIR/${base}.${archive_ext}" "${bin_name}")
  else
    tar -C "$out_dir" -czf "$DIST_DIR/${base}.${archive_ext}" "${bin_name}"
  fi
  rm -rf "$out_dir"
}

build_target darwin amd64
build_target darwin arm64
build_target linux amd64
build_target linux arm64
build_target windows amd64

(
  cd "$DIST_DIR"
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 agent-wizard_* > checksums.txt
  else
    sha256sum agent-wizard_* > checksums.txt
  fi
)

echo "Built release assets in $DIST_DIR"
