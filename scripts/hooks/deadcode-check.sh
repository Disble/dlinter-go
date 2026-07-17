#!/usr/bin/env bash
# Runs deadcode with the same false-positive exceptions as
# .github/workflows/ci.yml's deadcode job. Invoked by lefthook pre-push
# (see lefthook.yml). See scripts/hooks/gofmt-check.sh for why this lives
# in its own script.
set -uo pipefail

if ! command -v deadcode >/dev/null 2>&1; then
  echo "installing deadcode..."
  go install golang.org/x/tools/cmd/deadcode@latest
fi

# Known false-positive exception: golangci-lint's Module Plugin System calls
# plugin.go's registration surface (init, New, plugin.BuildAnalyzers,
# plugin.GetLoadMode) and skeleton.NewAnalyzer from the *external* custom-gcl
# framework, not from cmd/dlinter (this module's only main package). deadcode
# can only trace reachability from main packages inside this module, so it
# reports these as dead even though they are real, externally-invoked API.
# See docs/spikes.md.
out="$(deadcode ./... 2>&1 | grep -v -E '(plugin\.go:.*unreachable func: (init#1|New|plugin\.BuildAnalyzers|plugin\.GetLoadMode)|analyzer\.go:.*unreachable func: NewAnalyzer)')"
if [ -n "$out" ]; then
  echo "$out"
  echo "deadcode found unreachable code (see output above)"
  exit 1
fi
echo "no dead code found (plugin-entrypoint exceptions documented in docs/spikes.md)"
