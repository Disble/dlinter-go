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

## License

MIT
