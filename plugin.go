// Package dlinter registers the dlinter module plugin with golangci-lint.
// It decodes the package-role graph from Settings, builds a
// *rolegraph.Graph from it, and contributes the mayDependOn analyzer that
// enforces the resulting import-direction contract.
package dlinter

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/Disble/dlinter-go/internal/rolegraph"
	"github.com/Disble/dlinter-go/pkg/analyzers/maydependon"
	"github.com/Disble/dlinter-go/pkg/analyzers/requiredoc"
)

func init() {
	register.Plugin("dlinter", New)
}

// plugin implements register.LinterPlugin.
type plugin struct {
	settings Settings
}

// New decodes the plugin settings and constructs the dlinter plugin.
func New(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](conf)
	if err != nil {
		return nil, err
	}
	return &plugin{settings: settings}, nil
}

// BuildAnalyzers translates the decoded Settings.Roles into a
// *rolegraph.Graph and returns the mayDependOn analyzer injected with that
// Graph, plus the requireDoc analyzer when Settings.RequireDoc is set.
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	roles := make(map[string]rolegraph.RoleDef, len(p.settings.Roles))
	for name, spec := range p.settings.Roles {
		roles[name] = rolegraph.RoleDef{
			Packages:    spec.Packages,
			MayDependOn: spec.MayDependOn,
		}
	}

	g := rolegraph.New(roles)

	analyzers := []*analysis.Analyzer{
		maydependon.NewAnalyzer(g),
	}
	// requireDoc defaults to ON when the setting is absent: dlinter is
	// opinionated, and its defaults are where the opinion lives. Kept inline
	// rather than extracted, because a helper reachable only from this
	// framework-invoked method would need its own deadcode exception.
	if p.settings.RequireDoc == nil || *p.settings.RequireDoc {
		analyzers = append(analyzers, requiredoc.NewAnalyzer())
	}
	return analyzers, nil
}

// GetLoadMode returns LoadModeTypesInfo, which mayDependOn requires to
// resolve full type/package information for each analyzed package.
func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
