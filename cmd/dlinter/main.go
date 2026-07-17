// Command dlinter is the entrypoint skeleton for this repo's linter.
// This slice only establishes the entrypoint role target in the
// self-applied role graph (see .golangci.yml); no scaffolding logic is
// implemented yet.
package main

import (
	"github.com/Disble/dlinter-go/internal/rolegraph"
	"github.com/Disble/dlinter-go/pkg/analyzers/skeleton"
)

// These references exist only to anchor cmd/dlinter's dependency on the
// core and adapter roles, proving the entrypoint role can see both.
var (
	_ = skeleton.Analyzer
	_ rolegraph.Role
)

func main() {}
