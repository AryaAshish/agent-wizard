#!/usr/bin/env sh
set -eu
VERSION="${VERSION:-latest}"
DIR="${INSTALL_DIR:-$HOME/go/bin}"
REPO="${REPO:-AryaAshish/agent-wizard}"
BASE_URL="${BASE_URL:-}"

resolve_version() {
  if [ "$VERSION" != "latest" ]; then
    printf '%s' "${VERSION#v}"
    return
  fi
  if [ -n "$BASE_URL" ]; then
    echo "VERSION must be set when BASE_URL is provided." >&2
    exit 1
  fi
  latest_url="$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest")"
  tag="${latest_url##*/}"
  # Reject bogus segments (e.g. HTML or API URLs) from odd redirect chains.
  case "$tag" in
    ''|*/*|*:*|*%*) tag="" ;;
  esac
  if [ -z "$tag" ]; then
    tag="$(
      curl -fsSL \
        -H 'Accept: application/vnd.github+json' \
        "https://api.github.com/repos/${REPO}/releases/latest" |
        sed -n '/"tag_name"[[:space:]]*:[[:space:]]*"/{
          s/.*"tag_name"[[:space:]]*:[[:space:]]*"//
          s/".*//
          p
          q
        }'
    )"
  fi
  if [ -z "$tag" ]; then
    echo "Could not resolve latest release tag." >&2
    exit 1
  fi
  printf '%s' "${tag#v}"
}

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
esac
case "$os" in
  darwin|linux) ;;
  mingw*|msys*|cygwin*) os="windows" ;;
  *) echo "Unsupported OS: $os" >&2; exit 1 ;;
esac

version="$(resolve_version)"
archive_ext="tar.gz"
bin_name="agent-wizard"
if [ "$os" = "windows" ]; then
  archive_ext="zip"
  bin_name="agent-wizard.exe"
fi
asset="agent-wizard_${version}_${os}_${arch}.${archive_ext}"
base_url="https://github.com/${REPO}/releases/download/v${version}"
if [ -n "$BASE_URL" ]; then
  base_url="$BASE_URL"
fi

echo "Installing agent-wizard v${version} into ${DIR}"
mkdir -p "${DIR}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
curl -fsSL "${base_url}/${asset}" -o "${tmp}/${asset}"
curl -fsSL "${base_url}/checksums.txt" -o "${tmp}/checksums.txt"

expected="$(awk -v f="$asset" '$2==f {print $1; exit}' "${tmp}/checksums.txt")"
if [ -z "$expected" ]; then
  echo "Checksum not found for ${asset}" >&2
  exit 1
fi

if command -v shasum >/dev/null 2>&1; then
  actual="$(shasum -a 256 "${tmp}/${asset}" | awk '{print $1}')"
elif command -v sha256sum >/dev/null 2>&1; then
  actual="$(sha256sum "${tmp}/${asset}" | awk '{print $1}')"
else
  echo "No SHA256 tool found (need shasum or sha256sum)." >&2
  exit 1
fi

if [ "$expected" != "$actual" ]; then
  echo "Checksum verification failed for ${asset}" >&2
  exit 1
fi

if [ "$archive_ext" = "zip" ]; then
  unzip -q "${tmp}/${asset}" -d "${tmp}"
else
  tar -xzf "${tmp}/${asset}" -C "${tmp}"
fi

install -m 0755 "${tmp}/${bin_name}" "${DIR}/agent-wizard"
echo "Installed ${DIR}/agent-wizard"
