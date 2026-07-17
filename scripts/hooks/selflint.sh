#!/usr/bin/env bash
# Builds custom-gcl lazily (same cache logic as Makefile's `build` target)
# and runs the self-lint. Invoked by lefthook pre-commit (see lefthook.yml).
# See scripts/hooks/gofmt-check.sh for why this lives in its own script.
set -euo pipefail

bin="bin/custom-gcl"
[ "${OS:-}" = "Windows_NT" ] && bin="bin/custom-gcl.exe"

if [ ! -f "$bin" ] || [ .custom-gcl.yml -nt "$bin" ] || [ go.sum -nt "$bin" ]; then
  echo "building custom-gcl..."
  golangci-lint custom
fi

"./$bin" run ./...
