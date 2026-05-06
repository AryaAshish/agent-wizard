#!/usr/bin/env sh
set -eu
VERSION="${VERSION:-latest}"
DIR="${INSTALL_DIR:-$HOME/go/bin}"

echo "Installing github.com/aryaashish/agent-wizard@${VERSION} into ${DIR}"
mkdir -p "${DIR}"

GOBIN="${DIR}" go install "github.com/aryaashish/agent-wizard@${VERSION}"
