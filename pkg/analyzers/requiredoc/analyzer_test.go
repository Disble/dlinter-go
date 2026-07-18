package requiredoc_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Disble/dlinter-go/pkg/analyzers/requiredoc"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), requiredoc.NewAnalyzer(), "basic", "mainpkg")
}
