// Package requiredoc provides the requireDoc analyzer, which reports
// unexported functions and methods that carry no doc comment.
//
// Go's own conventions — and the linters that implement them, revive's
// exported rule and staticcheck's ST1020 family — scope documentation
// requirements to the exported surface, on the reasoning that unexported
// helpers are internal detail. This analyzer deliberately extends the
// requirement inward, for codebases that want every function to state its
// intent regardless of visibility.
//
// Exported symbols are NOT reported here: revive's exported rule already
// owns them, and reporting both would double-report the same defect.
package requiredoc

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports unexported functions and methods declared without a doc
// comment. It is a static composite literal rather than a constructed value
// so that deadcode's reachability analysis can trace run (and the helpers it
// calls) from the address-taken Run field, keeping them out of the
// allowlist.
var Analyzer = &analysis.Analyzer{
	Name: "requireDoc",
	Doc:  "reports unexported functions and methods that have no doc comment",
	Run:  run,
}

// NewAnalyzer returns the requireDoc analyzer. This rule carries no injected
// configuration, so every caller shares one value.
func NewAnalyzer() *analysis.Analyzer {
	return Analyzer
}

// run reports every unexported function or method declaration in the package
// that has no doc comment attached.
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if skip(fn) {
				continue
			}
			if fn.Doc != nil && len(fn.Doc.List) > 0 {
				continue
			}
			pass.Reportf(fn.Pos(), "unexported %s %q has no doc comment", kind(fn), name(fn))
		}
	}
	return nil, nil
}

// skip reports whether fn is outside this rule's scope: exported functions
// belong to revive's exported rule, and init has no name to document.
func skip(fn *ast.FuncDecl) bool {
	if fn.Name == nil {
		return true
	}
	if fn.Name.IsExported() {
		return true
	}
	// init and main are lifecycle entry points: no caller, no meaningful
	// name to document. Their purpose belongs in the package comment, and
	// requiring one here yields "// main is the entry point" -- exactly the
	// tautology this rule exists to avoid producing.
	if fn.Recv == nil && (fn.Name.Name == "init" || fn.Name.Name == "main") {
		return true
	}
	// A method on an exported type is still part of that type's documented
	// surface only when the method itself is exported; unexported methods
	// remain in scope here.
	return false
}

// kind returns "method" for declarations with a receiver and "func"
// otherwise, so the diagnostic names what it found.
func kind(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		return "method"
	}
	return "func"
}

// name returns the identifier used in diagnostics: bare for functions, and
// qualified with the receiver type for methods, so the report is unambiguous
// in packages where several types share a method name.
func name(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return fn.Name.Name
	}
	if recv := receiverType(fn.Recv.List[0].Type); recv != "" {
		return fmt.Sprintf("%s.%s", recv, fn.Name.Name)
	}
	return fn.Name.Name
}

// receiverType renders a receiver's type name, unwrapping a pointer receiver
// and ignoring generic type parameters, which do not affect identity here.
func receiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverType(t.X)
	case *ast.IndexExpr:
		return receiverType(t.X)
	case *ast.IndexListExpr:
		return receiverType(t.X)
	case *ast.Ident:
		return t.Name
	default:
		return ""
	}
}
