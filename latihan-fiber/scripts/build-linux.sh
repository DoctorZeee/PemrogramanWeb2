#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
GOOS_TARGET="${GOOS_TARGET:-linux}"
GOARCH_TARGET="${GOARCH_TARGET:-amd64}"
OUTPUT="${OUTPUT:-bin/api}"

cd "$ROOT_DIR"
mkdir -p "$(dirname -- "$OUTPUT")"

printf 'Membangun Fiber untuk %s/%s...\n' "$GOOS_TARGET" "$GOARCH_TARGET"
CGO_ENABLED=0 GOOS="$GOOS_TARGET" GOARCH="$GOARCH_TARGET" \
  go build -trimpath -ldflags='-s -w' -o "$OUTPUT" ./cmd/api

printf 'Binary tersedia: %s/%s\n' "$ROOT_DIR" "$OUTPUT"
