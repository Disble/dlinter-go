// Package basic exercises the requireDoc analyzer.
package basic

// documentedFunc has a doc comment and must not be reported.
func documentedFunc() {}

func undocumentedFunc() {} // want `unexported func "undocumentedFunc" has no doc comment`

// DocumentedExported has a doc comment and must not be reported.
func DocumentedExported() {}

// UndocumentedExported is reported by revive: exported, not by this rule,
// so requireDoc must stay silent to avoid double-reporting.
func UndocumentedExported() {}

type carrier struct{}

// documentedMethod has a doc comment and must not be reported.
func (c carrier) documentedMethod() {}

func (c carrier) undocumentedMethod() {} // want `unexported method "carrier.undocumentedMethod" has no doc comment`

// init functions are lifecycle hooks with no caller and no name to document;
// documenting them adds nothing, so they are exempt.
func init() {}

// Short helpers are still covered: there is no length exemption. (This
// comment is deliberately separated by a blank line so it does NOT attach
// as a doc comment to the declaration below.)

func tiny() int { return 1 } // want `unexported func "tiny" has no doc comment`
