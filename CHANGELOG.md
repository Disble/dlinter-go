# Changelog

All notable changes to this project are documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.1] - 2026-07-18

### Fixed

- **Preset rules no longer contradict each other.** `revive: file-length-limit`
  now runs with `skipComments: true`. Previously it counted comment lines while
  `requireDoc` mandated a doc comment on every unexported function, so obeying
  one rule pushed a file toward violating the other — two files in a real
  133-file codebase crossed the 400-line limit purely by adding the docs the
  preset demands. File size measures code mass, which is what the lever exists
  for; `funlen` already ignores comments for the same reason.

## [0.3.0] - 2026-07-18

### Added

- **`requireDoc` rule**: reports unexported functions and methods declared
  without a doc comment. **Enabled by default** — dlinter is opinionated, and
  a rule that ships off is a rule most projects never discover. `init` and
  `main` are exempt; exported symbols remain with `revive: exported` so the
  two never double-report. Opt out with `requireDoc: false`. See the README
  for the caveats before adopting.
- `godoclint` in the recommended preset: `revive: exported` enforces that doc
  comments exist, `godoclint` enforces that the ones written are well-formed.

## [0.2.0] - 2026-07-18

### Added

- **Discipline harness**: `recommended.golangci.yml` now bundles a coordinated
  set of levers whose job is to make sprawl expensive to hide — file size
  (`revive: file-length-limit`, 400), function length (`funlen`, 60/40),
  cognitive complexity (`gocognit`, 15), nesting depth (`nestif`, 4), and
  suppression hygiene (`nolintlint`: every `//nolint` must name its linter and
  explain itself).
- `docs/threat-model.md`: what the harness does and does not guarantee, plus a
  governance baseline for keeping the config itself under review.
- README section documenting each lever's cheapest dishonest escape, so the
  weak bars are never mistaken for the strong ones.

### Changed

- **Pinned golangci-lint bumped to v2.12.2** (from v2.1.0). Versions up to
  v2.11 build their internal clone with `-c advice.detachedHead=false` as a
  single argument, which git 2.54 and newer reject
  (`invalid key:  advice.detachedHead`), breaking `golangci-lint custom` on
  current runners and any machine with modern git. **Consumers pinning v2.1.0
  should update `.custom-gcl.yml`.**
- Preset excludes `funlen` and `dupl` on `_test.go`, calibrated against a real
  133-file codebase where 88% of `funlen` findings were tests. `gocognit` and
  `nestif` still apply to tests deliberately.
- Internal decompositions of `Graph.Resolve` and `runWithGraph`, forced by this
  repo's own tightened thresholds — refactored, not exempted.

### Notes

- A custom package-cohesion ("god-package") analyzer was evaluated and
  deliberately **not** shipped: no candidate metric survived its false-positive
  story, and calibration across 33 real packages found a smooth continuum
  rather than a bimodal distribution, leaving no honest threshold to pick.

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

[0.3.1]: https://github.com/Disble/dlinter-go/releases/tag/v0.3.1
[0.3.0]: https://github.com/Disble/dlinter-go/releases/tag/v0.3.0
[0.2.0]: https://github.com/Disble/dlinter-go/releases/tag/v0.2.0
[0.1.0]: https://github.com/Disble/dlinter-go/releases/tag/v0.1.0
