# dlinter-go

Architecture governance for Go — as build-time guarantees, not prose.

dlinter-go is the Go sibling of [dlinter-ts-react](https://github.com/Disble/dlinter-ts-react). Instead of reinventing the linting wheel, it builds on the ecosystem's giant shoulders:

- **Host**: [golangci-lint](https://golangci-lint.run) via its [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/) — one custom binary, one config, one CI step.
- **Custom rules**: opinionated architecture and best-practice analyzers built on `golang.org/x/tools/go/analysis` (import-direction contracts, package-role conventions, layer boundaries).
- **Preset**: a curated `.golangci.yml` composing the best of the built-in linters (`unused`, `dupl`, `gocognit`, `depguard`, ...) with tuned, documented severities.
- **Scaffolding**: a `dlinter init` CLI that generates `.golangci.yml`, `.custom-gcl.yml`, git hooks, and CI config.
- **Dead code**: the official [`deadcode`](https://go.dev/blog/deadcode) tool orchestrated as a periodic CI step (whole-program RTA analysis cannot live inside golangci-lint).

## Status

Early development. This repository enforces its own rules on itself — self-governance is the proof that the rules work.

## Local setup

1. Build the self-lint binary (`golangci-lint` v2.1.0 must be on `PATH`):

   ```sh
   golangci-lint custom
   ```

   This reads `.custom-gcl.yml` and produces `./bin/custom-gcl` (`.exe` on
   Windows), a golangci-lint build that bundles `dlinter` alongside every
   standard linter.

2. Install [lefthook](https://github.com/evilmartians/lefthook) and wire the
   local git hooks:

   ```sh
   lefthook install
   ```

   `pre-commit` runs `gofmt -l` on staged files, the self-lint
   (`./bin/custom-gcl run ./...`), and `go test ./...`. `pre-push` runs
   `deadcode` (see below) — mirroring `.github/workflows/ci.yml` so failures
   surface locally before CI.

3. Install [`deadcode`](https://go.dev/blog/deadcode) (only needed for the
   `pre-push` hook and the `deadcode` Makefile target; installed
   automatically by the hook on first push if missing):

   ```sh
   go install golang.org/x/tools/cmd/deadcode@latest
   ```

## License

MIT
