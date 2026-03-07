#!/usr/bin/env bash
# Generates Go protobuf + gRPC stubs from planton/apis protos into gen/go/.
# Pulls protos directly from the planton git repo (no local clone needed).
#
# Usage:
#   ./tools/buf-generate-go.sh

set -euo pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "$REPO_ROOT"

echo "Generating Go stubs from planton/apis (remote)..."

rm -rf .buf-generated gen/go/ai

buf generate --template buf.gen.go.yaml

readonly SRC=".buf-generated/github.com/plantonhq/mcp-server-planton/gen/go"

if [ ! -d "$SRC" ]; then
    echo "ERROR: Expected output not found at ${SRC}"
    echo "Contents of .buf-generated/:"
    find .buf-generated -maxdepth 4 -type d 2>/dev/null || true
    rm -rf .buf-generated
    exit 1
fi

mkdir -p gen/go
cp -r "${SRC}/"* gen/go/

rm -rf .buf-generated

echo "Done. Generated stubs are in gen/go/"
echo ""
echo "Stub count: $(find gen/go -name '*.go' | wc -l | tr -d ' ') files"
