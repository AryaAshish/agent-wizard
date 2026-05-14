#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DIST_DIR="$(mktemp -d)"
PROJECT_DIR="$(mktemp -d)"
HOME_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$DIST_DIR" "$PROJECT_DIR" "$HOME_DIR"
}
trap cleanup EXIT

VERSION="${VERSION:-v0.0.0-smoke}"
(
  cd "$ROOT_DIR"
  VERSION="$VERSION" bash scripts/release/build_assets.sh "$DIST_DIR"
)

export HOME="$HOME_DIR"
export INSTALL_DIR="$HOME_DIR/bin"
export BASE_URL="file://${DIST_DIR}"
export VERSION="${VERSION#v}"

bash "$ROOT_DIR/install.sh"
"$HOME_DIR/bin/agent-wizard" --version

export PATH="$HOME_DIR/bin:$PATH"
export AGENT_WIZARD_VERSION="${VERSION#v}"
export AGENT_WIZARD_CACHE_DIR="$HOME_DIR/.cache/agent-wizard/npm"
mkdir -p "$AGENT_WIZARD_CACHE_DIR/v${VERSION#v}"
case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
  darwin) os_name="darwin" ;;
  linux) os_name="linux" ;;
  *) os_name="windows" ;;
esac
case "$(uname -m)" in
  x86_64|amd64) arch_name="amd64" ;;
  arm64|aarch64) arch_name="arm64" ;;
  *) arch_name="amd64" ;;
esac
asset="agent-wizard_${VERSION#v}_${os_name}_${arch_name}.tar.gz"
tar -xzf "$DIST_DIR/$asset" -C "$AGENT_WIZARD_CACHE_DIR/v${VERSION#v}"
(
  cd "$ROOT_DIR/npm/agent-wizard"
  node bin/agent-wizard.js --version
)

(
  cd "$PROJECT_DIR"
  git init -q 2>/dev/null || true
  "$HOME_DIR/bin/agent-wizard" init </dev/null
  "$HOME_DIR/bin/agent-wizard" list --source-name community --filter pr-review | grep -q pr-review
  "$HOME_DIR/bin/agent-wizard" add pr-review --source community
  "$HOME_DIR/bin/agent-wizard" sync
  "$HOME_DIR/bin/agent-wizard" sync
  test -f .agents/skills/pr-review/SKILL.md
)

echo "distribution smoke passed"
