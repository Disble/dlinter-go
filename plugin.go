// Package dlinter registers the dlinter module plugin with golangci-lint.
// This is a walking skeleton: it proves the self-referential custom-gcl
// build, settings decoding, and LoadModeTypesInfo wiring work end-to-end
// with a single throwaway analyzer, before any real architecture rule
// exists.
package dlinter

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/Disble/dlinter-go/pkg/analyzers/skeleton"
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

// BuildAnalyzers returns the analyzers this plugin contributes. This slice
// only contributes the throwaway skeleton analyzer.
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		skeleton.NewAnalyzer(skeleton.Settings{}),
	}, nil
}

// GetLoadMode returns LoadModeTypesInfo so the skeleton analyzer exercises
// the same (expensive) load path that real architecture rules will need.
func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
