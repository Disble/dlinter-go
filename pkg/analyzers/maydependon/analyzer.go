// Package maydependon provides the mayDependOn analyzer, which enforces
// the package-role import-direction contract declared in .golangci.yml
// using the internal/rolegraph engine.
package maydependon

import (
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/Disble/dlinter-go/internal/rolegraph"
)

// Analyzer is the package-level *analysis.Analyzer value used to anchor
// this package's reachability for deadcode analysis (mirrors the pattern
// used by the retired skeleton package: cmd/dlinter address-takes this var).
// Production wiring uses NewAnalyzer, which injects a real *rolegraph.Graph;
// this shared var has no Graph and is unused in production wiring.
var Analyzer = &analysis.Analyzer{
	Name:             "mayDependOn",
	Doc:              "reports imports that cross a package-role boundary not permitted by the configured role graph",
	Run:              run,
	RunDespiteErrors: false,
}

// NewAnalyzer returns a new *analysis.Analyzer wired to evaluate imports
// against g. It returns a distinct analyzer value (not the shared Analyzer
// var) so the injected Graph is available to a closure-bound run function,
// per Design D1.
func NewAnalyzer(g *rolegraph.Graph) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "mayDependOn",
		Doc:  "reports imports that cross a package-role boundary not permitted by the configured role graph",
		Run: func(pass *analysis.Pass) (any, error) {
			return runWithGraph(pass, g)
		},
		RunDespiteErrors: false,
	}
}

// run is the Run function bound to the shared Analyzer var. It has no
// injected Graph and reports no diagnostics; it exists only so Analyzer can
// be address-taken for deadcode reachability without requiring a Graph at
// package-init time.
func run(_ *analysis.Pass) (any, error) {
	return nil, nil
}

// runWithGraph evaluates every import in pass against g, reporting a
// diagnostic for each import that crosses a forbidden role boundary.
func runWithGraph(pass *analysis.Pass, g *rolegraph.Graph) (any, error) {
	importerRole, ok := g.Resolve(relativize(pass.Pkg.Path(), pass.Module))
	if !ok {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}

			importedRole, ok := g.Resolve(relativize(importPath, pass.Module))
			if !ok {
				continue
			}

			if importerRole == importedRole {
				continue
			}

			if !g.Allowed(importerRole, importedRole) {
				pass.Reportf(spec.Pos(), "role %q may not depend on role %q (import %q)", importerRole, importedRole, importPath)
			}
		}
	}

	return nil, nil
}

// relativize strips the module prefix from pkgPath, returning the
// module-relative path that Graph.Resolve expects. pass.Module is never
// nil under analysistest (only its Path is empty), so a single
// TrimPrefix call correctly handles both real runtime (Module.Path set)
// and analysistest/testdata (Module.Path == "") without a nil-check
// branch: TrimPrefix with an empty prefix is a no-op.
func relativize(pkgPath string, mod *analysis.Module) string {
	rel := strings.TrimPrefix(pkgPath, mod.Path)
	rel = strings.TrimPrefix(rel, "/")
	return rel
}
