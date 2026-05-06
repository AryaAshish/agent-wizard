#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${TAP_REPO:-}" || -z "${TAP_TOKEN:-}" || -z "${VERSION:-}" ]]; then
  echo "TAP_REPO, TAP_TOKEN, and VERSION are required."
  exit 1
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

VERSION="${VERSION#v}"
FORMULA_NAME="agent-wizard.rb"
URL_BASE="https://github.com/AryaAshish/agent-wizard/releases/download/v${VERSION}"
ARCHIVE="agent-wizard_${VERSION}_darwin_arm64.tar.gz"
CHECKSUMS_URL="${URL_BASE}/checksums.txt"

curl -fsSL "$CHECKSUMS_URL" -o "$WORK_DIR/checksums.txt"
SHA="$(awk '/agent-wizard_'"${VERSION}"'_darwin_arm64\.tar\.gz$/ {print $1; exit}' "$WORK_DIR/checksums.txt")"
if [[ -z "$SHA" ]]; then
  echo "Could not find checksum for ${ARCHIVE}"
  exit 1
fi

git clone "https://x-access-token:${TAP_TOKEN}@github.com/${TAP_REPO}.git" "$WORK_DIR/tap"
mkdir -p "$WORK_DIR/tap/Formula"
cat > "$WORK_DIR/tap/Formula/${FORMULA_NAME}" <<EOF
class AgentWizard < Formula
  desc "Local-first CLI for reusable agent skills"
  homepage "https://github.com/AryaAshish/agent-wizard"
  url "${URL_BASE}/${ARCHIVE}"
  sha256 "${SHA}"
  version "${VERSION}"

  def install
    bin.install "agent-wizard"
  end

  test do
    output = shell_output("#{bin}/agent-wizard --version")
    assert_match "agent-wizard", output
  end
end
EOF

cd "$WORK_DIR/tap"
git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git add "Formula/${FORMULA_NAME}"
git commit -m "Update agent-wizard to v${VERSION}" || {
  echo "No formula changes."
  exit 0
}
git push origin HEAD
echo "Homebrew tap updated to v${VERSION}"
