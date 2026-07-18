package dlinter

// Settings is the top-level configuration decoded from .golangci.yml for the
// dlinter module linter. It defines the package-role graph that the future
// mayDependOn rule will enforce; this slice only proves the decode path.
type Settings struct {
	Roles map[string]RoleSpec `json:"roles"`

	// RequireDoc enables the requireDoc analyzer, which reports unexported
	// functions and methods with no doc comment. It is opt-in because it
	// goes beyond Go's own convention of documenting only the exported
	// surface; see pkg/analyzers/requiredoc for the reasoning.
	RequireDoc bool `json:"requireDoc"`
}

// RoleSpec describes one architectural role: the packages that belong to it
// and the roles it is allowed to depend on.
type RoleSpec struct {
	Packages    []string `json:"packages"`
	MayDependOn []string `json:"mayDependOn"`
}
