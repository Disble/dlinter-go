// Package entry belongs to the "entrypoint" role in the test Graph and
// imports coreok (allowed: entrypoint mayDependOn core).
package entry

import (
	"coreok"
)

var _ = coreok.Marker

var Marker int
