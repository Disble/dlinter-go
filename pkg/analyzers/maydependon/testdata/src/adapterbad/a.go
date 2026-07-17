// Package adapterbad belongs to the "adapter" role in the test Graph and
// imports entry ("entrypoint" role), which adapter is not allowed to
// depend on.
package adapterbad

import (
	"entry" // want `role "adapter" may not depend on role "entrypoint" \(import "entry"\)`
)

var _ = entry.Marker

var Marker int
