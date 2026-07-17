# dlinter-go

[![CI](https://github.com/Disble/dlinter-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Disble/dlinter-go/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Disble/dlinter-go)](https://github.com/Disble/dlinter-go/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/Disble/dlinter-go.svg)](https://pkg.go.dev/github.com/Disble/dlinter-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Architecture governance for Go — as build-time guarantees, not prose.

Declare which packages belong to which architectural role (`core`, `adapter`, `entrypoint`, or your own) and which roles each may depend on. Any import that crosses a forbidden boundary fails the lint:

```
internal/rolegraph/rolegraph.go:15:2: mayDependOn: role "core" may not depend on role "adapter" (import "example.com/app/internal/adapters/db")
```

dlinter ships as a [golangci-lint module plugin](https://golangci-lint.run/docs/plugins/module-plugins/), so it runs inside the toolchain you already use — one binary, one config, one CI step. This repository enforces its own rules on itself; self-governance is the proof that the rules work.

## Quick start

1. Create `.custom-gcl.yml` in your repo:

   ```yaml
   version: v2.1.0
   plugins:
     - module: github.com/Disble/dlinter-go
       version: v0.1.0
   ```

2. Build the custom binary (requires Go and [golangci-lint](https://golangci-lint.run) v2.1.0):

   ```sh
   golangci-lint custom
   ```

3. Declare your architecture in `.golangci.yml`:

   ```yaml
   version: "2"

   linters:
     enable:
       - dlinter
     settings:
       custom:
         dlinter:
           type: module
           description: package-role / import-direction contracts
           settings:
             roles:
               core:
                 packages:
                   - internal/domain        # exact match
                 mayDependOn: []
               adapter:
                 packages:
                   - internal/adapters/     # trailing slash = subtree
                 mayDependOn:
                   - core
               entrypoint:
                 packages:
                   - "."                    # module root
                   - cmd/app
                 mayDependOn:
                   - core
                   - adapter
   ```

4. Run it:

   ```sh
   ./custom-gcl run ./...
   ```

## Configuration

Each role lists the packages that belong to it and the roles it may depend on.

| Pattern | Matches |
|---------|---------|
| `internal/domain` | exactly that package |
| `internal/adapters/` (trailing `/`) | the whole subtree; longest prefix wins |
| `"."` | only the module root package |

Rules that keep the check predictable:

- A role may always depend on itself.
- Packages with no role are ignored, as are stdlib and external imports — the rules constrain only your own module's import graph.
- Exact and root matches always win over prefix matches.

## Why golangci-lint instead of a standalone tool?

golangci-lint is already the orchestrator of the Go linting ecosystem (dead code via `unused`, duplication via `dupl`, complexity via `gocognit`, 100+ linters). dlinter adds the missing piece — enforceable architecture contracts — without asking your team to adopt another binary, config format, or CI step.

## Contributing

Development happens through the standard fork-and-PR flow. Local setup:

1. Build the self-lint binary (`golangci-lint` v2.1.0 on `PATH`):

   ```sh
   golangci-lint custom
   ```

   This produces `./bin/custom-gcl`, a golangci-lint build that bundles `dlinter` alongside every standard linter.

2. Install [lefthook](https://github.com/evilmartians/lefthook) and wire the git hooks:

   ```sh
   lefthook install
   ```

   `pre-commit` runs `gofmt`, the self-lint (`./bin/custom-gcl run ./...`), `go test ./...`, and [`deadcode`](https://go.dev/blog/deadcode) — mirroring CI so failures surface locally first.

3. Run the tests:

   ```sh
   go test ./...
   ```

   Analyzer behavior is specified with [`analysistest`](https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest) fixtures under `pkg/analyzers/*/testdata/`; new rules follow the same pattern.

Conventional commits are required (`feat:`, `fix:`, `docs:`, ...). See [docs/spikes.md](docs/spikes.md) for recorded design decisions and known gotchas (self-referential plugin build, deadcode allowlist rationale).

## License

[MIT](LICENSE)
