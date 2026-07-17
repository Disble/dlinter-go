// Package skeleton provides a throwaway analyzer used to prove the
// golangci-lint Module Plugin System end-to-end (registration, settings
// decoding, and LoadModeTypesInfo) before any real architecture rule exists.
package skeleton

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports every package-level function named skeletonMarker.
// It exists only to prove the plugin plumbing works; it is not a real rule.
var Analyzer = &analysis.Analyzer{
	Name: "skeleton",
	Doc:  "reports package-level functions named skeletonMarker (walking-skeleton marker)",
	Run:  run,
}

// Settings configures the skeleton analyzer. It is intentionally unused:
// the analyzer's behavior does not depend on settings, but NewAnalyzer
// accepts it to establish the constructor-injection pattern real rules
// will reuse.
type Settings struct{}

// NewAnalyzer wraps Analyzer, ignoring settings, to satisfy the
// constructor-injection wiring used by the plugin.
func NewAnalyzer(_ Settings) *analysis.Analyzer {
	return Analyzer
}

func run(pass *analysis.Pass) (any, error) {
	for ident, obj := range pass.TypesInfo.Defs {
		fn, ok := obj.(*types.Func)
		if !ok || fn == nil {
			continue
		}
		if fn.Name() != "skeletonMarker" {
			continue
		}
		if sig, ok := fn.Type().(*types.Signature); !ok || sig.Recv() != nil {
			continue
		}
		pass.Reportf(ident.Pos(), "skeletonMarker: walking-skeleton marker function")
	}
	return nil, nil
}
