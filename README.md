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

## Use in your project

dlinter ships as a [golangci-lint module plugin](https://golangci-lint.run/docs/plugins/module-plugins/): you build a custom golangci-lint binary that bundles it, then enable it like any other linter.

1. Create `.custom-gcl.yml` in your repo:

   ```yaml
   version: v2.1.0
   plugins:
     - module: github.com/Disble/dlinter-go
       version: v0.1.0
   ```

2. Build the custom binary (requires Go and `golangci-lint` v2.1.0 on `PATH`):

   ```sh
   golangci-lint custom
   ```

   This produces `./custom-gcl` (`.exe` on Windows) — golangci-lint plus the `dlinter` linter.

3. Declare your architecture in `.golangci.yml`. Assign each package a role and state which roles it may depend on:

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

   An import that crosses a boundary not listed in `mayDependOn` fails the lint:

   ```
   role "core" may not depend on role "adapter" (import "example.com/app/internal/adapters/db")
   ```

Path matching: `"."` matches only the module root, a trailing `/` matches the whole subtree (longest prefix wins), anything else matches exactly. Packages with no role are ignored, as are stdlib and external imports — the rules constrain only your own module's graph.

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
   (`./bin/custom-gcl run ./...`), `go test ./...`, and `deadcode` —
   mirroring `.github/workflows/ci.yml` so failures surface locally
   before CI.

3. Install [`deadcode`](https://go.dev/blog/deadcode) (used by the
   `pre-commit` hook and the `deadcode` Makefile target; installed
   automatically by the hook on first commit if missing):

   ```sh
   go install golang.org/x/tools/cmd/deadcode@latest
   ```

## License

MIT
