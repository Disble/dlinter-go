// Package corebad belongs to the "core" role in the test Graph and imports
// adapterok ("adapter" role), which core is not allowed to depend on.
package corebad

import (
	"adapterok" // want `role "core" may not depend on role "adapter" \(import "adapterok"\)`
)

var _ = adapterok.Marker
