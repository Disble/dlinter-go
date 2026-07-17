package maydependon_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Disble/dlinter-go/internal/rolegraph"
	"github.com/Disble/dlinter-go/pkg/analyzers/maydependon"
)

func TestAnalyzer(t *testing.T) {
	g := rolegraph.New(map[string]rolegraph.RoleDef{
		"core": {
			Packages:    []string{"coreok", "corebad"},
			MayDependOn: nil,
		},
		"adapter": {
			Packages:    []string{"adapterok", "adapterbad"},
			MayDependOn: []string{"core"},
		},
		"entrypoint": {
			Packages:    []string{"entry"},
			MayDependOn: []string{"core"},
		},
	})

	analyzer := maydependon.NewAnalyzer(g)

	analysistest.Run(t, analysistest.TestData(), analyzer,
		"coreok", "corebad", "adapterok", "adapterbad", "entry", "unroled")
}
