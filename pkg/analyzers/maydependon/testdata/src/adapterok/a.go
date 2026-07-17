// Package adapterok belongs to the "adapter" role in the test Graph and
// imports coreok (allowed: adapter mayDependOn core), fmt (stdlib, always
// exempt), and unroled (importer/target permissive: unroled has no
// configured role, so the import is never evaluated against role rules).
package adapterok

import (
	"fmt"

	"coreok"
	"unroled"
)

var (
	_ = coreok.Marker
	_ = unroled.Marker
)

var Marker int

func init() {
	fmt.Sprintln()
}
