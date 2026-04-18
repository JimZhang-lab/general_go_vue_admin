#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[1/3] Running backend tests..."
(cd "$ROOT_DIR/server" && go test ./...)

echo "[2/3] Building backend..."
(cd "$ROOT_DIR/server" && go build ./...)

echo "[3/3] Building frontend..."
(cd "$ROOT_DIR/web" && npm run build-only)

echo "Self-check passed."
