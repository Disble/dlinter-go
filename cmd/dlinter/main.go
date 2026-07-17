// Command dlinter is the entrypoint for this repo's linter. It has no
// scaffolding logic of its own: the real work happens in the
// golangci-lint Module Plugin System, which loads plugin.go via
// custom-gcl. This package exists to occupy the "entrypoint" role target
// in the self-applied role graph (see .golangci.yml).
package main

import (
	"github.com/Disble/dlinter-go/internal/rolegraph"
	"github.com/Disble/dlinter-go/pkg/analyzers/maydependon"
)

// These references exist only to anchor cmd/dlinter's dependency on the
// core and adapter roles, proving the entrypoint role can see both, and to
// keep maydependon.Analyzer reachable from a real main package for
// deadcode analysis (golangci-lint's external plugin framework invokes
// maydependon.NewAnalyzer, which deadcode cannot trace from outside this
// module).
var (
	_ = maydependon.Analyzer
	_ rolegraph.Role
)

func main() {}
