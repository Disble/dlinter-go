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

			if ForbiddenImport(g, pass.Module.Path, importerRole, importPath) {
				importedRole, _ := g.Resolve(trimModule(importPath, pass.Module.Path))
				pass.Reportf(spec.Pos(), "role %q may not depend on role %q (import %q)", importerRole, importedRole, importPath)
			}
		}
	}

	return nil, nil
}

// ForbiddenImport reports whether importerRole may not depend on the role
// resolved for importPath under module modPath, per g. An unresolved import
// role (no configured role matches it) is never forbidden — the graph
// constrains only this module's own packages. Exported so it can be unit
// tested directly, without the analysistest overhead runWithGraph otherwise
// requires.
func ForbiddenImport(g *rolegraph.Graph, modPath string, importerRole rolegraph.Role, importPath string) bool {
	importedRole, ok := g.Resolve(trimModule(importPath, modPath))
	if !ok {
		return false
	}
	if importerRole == importedRole {
		return false
	}
	return !g.Allowed(importerRole, importedRole)
}

// relativize strips the module prefix from pkgPath, returning the
// module-relative path that Graph.Resolve expects.
func relativize(pkgPath string, mod *analysis.Module) string {
	return trimModule(pkgPath, mod.Path)
}

// trimModule strips modPath, then a leading "/", from pkgPath. pass.Module
// is never nil under analysistest (only its Path is empty), so this
// correctly handles both real runtime (Module.Path set) and
// analysistest/testdata (Module.Path == "") without a nil-check branch:
// TrimPrefix with an empty prefix is a no-op.
func trimModule(pkgPath, modPath string) string {
	rel := strings.TrimPrefix(pkgPath, modPath)
	return strings.TrimPrefix(rel, "/")
}
