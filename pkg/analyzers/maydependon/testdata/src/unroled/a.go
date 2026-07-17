// Package unroled has no configured role in the test Graph and imports
// adapterbad. Because the importer's role is unresolved, the import is
// never evaluated against role rules (permissive).
package unroled

import (
	"adapterbad"
)

var _ = adapterbad.Marker

var Marker int
