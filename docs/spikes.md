# Spikes

This document records the empirical evidence for the two spikes that gate
the rest of the foundation-architecture slice: the self-referential
`custom-gcl` build (Spike A) and the Go 1.26 version pin (Spike B).

## Spike A: Self-Referential Dogfood Build

**Requirement**: `.custom-gcl.yml` must declare `plugins[0].path: .` so
`golangci-lint custom` builds `custom-gcl` from the working tree without a
module proxy version.

**Setup**:
- `.custom-gcl.yml`:
  ```yaml
  version: v2.1.0
  name: custom-gcl
  destination: ./bin
  plugins:
    - module: github.com/Disble/dlinter-go
      path: .
  ```
- Root `plugin.go` (package `dlinter`) registers via
  `register.Plugin("dlinter", New)`.

**Result**: PASS on the first attempt.

**Command**: `golangci-lint custom`
**Output**: exit code 0, no output, produced `./bin/custom-gcl.exe`. No
circular-module errors.

**Open question resolved — explicit `import:` line**: NOT needed. With
`path: .` set, the plugin entry's `import` defaults to the module root
package (as documented by the plugin-module-linter example). What WAS
required, that the design draft did not call out, is an explicit `module:`
field — `golangci-lint custom` fails fast with `field 'module' is required`
if omitted. We set `module: github.com/Disble/dlinter-go` (this repo's own
module path) alongside `path: .`.

## Spike B: Go 1.26 Version Pin

**Requirement**: the golangci-lint `version:` pinned in `.custom-gcl.yml`
must be empirically confirmed to build and run against this Go 1.26 module
before being locked.

**Candidate**: `v2.1.0` (first candidate tried, from the design's `v2.12.2+`
suggestion range — v2.1.0 was the version actually available/resolved by
the builder toolchain at spike time).

**Result**: PASS on the first candidate. No version bump was required.

**Commands**:
```
golangci-lint custom
./bin/custom-gcl.exe --version
```
**Output**:
```
golangci-lint has version v2.1.0-custom-gcl built with go1.26.1 from ? on <timestamp>
```

**Locked version**: `v2.1.0` (recorded in `.custom-gcl.yml`).

## Version-Sync Constraint

`go.mod` declares `github.com/golangci/plugin-module-register` and
`golang.org/x/tools` as direct dependencies of this repo's plugin code.
`.custom-gcl.yml`'s `version:` pin selects the golangci-lint release whose
own `go.mod` transitively determines a *compatible* (not necessarily
identical) version of `golang.org/x/tools/go/analysis`, since the
`analysis.Analyzer`/`analysis.Pass` types must match between our plugin
code and the custom-gcl builder's compiled `golangci-lint` core.

**Constraint**: whenever the `.custom-gcl.yml` `version:` pin is bumped,
re-verify (via `golangci-lint custom` + `go test ./...`) that
`golang.org/x/tools` and `github.com/golangci/plugin-module-register` in
`go.mod` still produce a plugin that builds and behaves correctly against
the new pinned golangci-lint version. **No automated check enforces this
in this slice** — it is a manual step for whoever bumps the pin.

## `.golangci.yml` key-path confirmation (Task 5)

See the "Self-Applied Config" section below, added after Task 5 lands.

## Self-lint scope gap (Task 5)

See the "Self-Lint Scope Gap" section below, added after Task 5 lands.
