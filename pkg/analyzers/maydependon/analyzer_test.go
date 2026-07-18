package maydependon_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Disble/dlinter-go/internal/rolegraph"
	"github.com/Disble/dlinter-go/pkg/analyzers/maydependon"
)

func TestForbiddenImport(t *testing.T) {
	g := rolegraph.New(map[string]rolegraph.RoleDef{
		"core": {
			Packages:    []string{"core"},
			MayDependOn: nil,
		},
		"adapter": {
			Packages:    []string{"adapter"},
			MayDependOn: []string{"core"},
		},
	})

	tests := []struct {
		name         string
		importerRole rolegraph.Role
		importPath   string
		want         bool
	}{
		{
			name:         "allowed dependency is not forbidden",
			importerRole: "adapter",
			importPath:   "mod/core",
			want:         false,
		},
		{
			name:         "disallowed dependency is forbidden",
			importerRole: "core",
			importPath:   "mod/adapter",
			want:         true,
		},
		{
			name:         "unresolved import role is not forbidden",
			importerRole: "core",
			importPath:   "mod/unroled",
			want:         false,
		},
		{
			name:         "same role as importer is never forbidden",
			importerRole: "core",
			importPath:   "mod/core",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maydependon.ForbiddenImport(g, "mod", tt.importerRole, tt.importPath)
			if got != tt.want {
				t.Fatalf("ForbiddenImport(%q, %q) = %v, want %v", tt.importerRole, tt.importPath, got, tt.want)
			}
		})
	}
}

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
