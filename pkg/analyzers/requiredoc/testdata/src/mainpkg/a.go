// Package main exercises the entry-point exemptions.
package main

// main and init are lifecycle entry points: they have no callers and no
// meaningful name to document, so the package comment is where their purpose
// belongs. Requiring a doc comment here yields "// main is the entry point",
// which is the tautology this rule exists to avoid producing.
func main() { helper() }

func init() {}

func helper() {} // want `unexported func "helper" has no doc comment`
