# Changelog

All notable changes to this project are documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-07-17

### Added

- `mayDependOn` analyzer: enforces package-role import-direction contracts
  declared in `.golangci.yml` (roles with `packages` and `mayDependOn` lists;
  exact, subtree-prefix, and module-root path matching).
- golangci-lint Module Plugin registration (`dlinter`), consumable via
  `.custom-gcl.yml` and `golangci-lint custom` (golangci-lint v2.1.0).
- Self-applied role graph: this repository lints itself with its own binary
  (core `internal/rolegraph`, adapter `pkg/analyzers/`, entrypoint root and
  `cmd/dlinter`).
- Three-stage CI (unit tests + analysistest, self-lint, deadcode) and local
  lefthook gates (pre-commit: gofmt, self-lint, tests; pre-push: deadcode).

[0.1.0]: https://github.com/Disble/dlinter-go/releases/tag/v0.1.0
