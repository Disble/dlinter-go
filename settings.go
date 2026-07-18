package dlinter

// Settings is the top-level configuration decoded from .golangci.yml for the
// dlinter module linter. It defines the package-role graph that the future
// mayDependOn rule will enforce; this slice only proves the decode path.
type Settings struct {
	Roles map[string]RoleSpec `json:"roles"`

	// RequireDoc controls the requireDoc analyzer, which reports unexported
	// functions and methods with no doc comment.
	//
	// It is ENABLED by default. dlinter is an opinionated linter: its
	// defaults state the opinion, and a rule that ships off is a rule most
	// projects never discover. Set it to false to opt out.
	//
	// The pointer distinguishes "unset" (use the default) from an explicit
	// false, which a plain bool cannot express.
	RequireDoc *bool `json:"requireDoc"`
}

// RoleSpec describes one architectural role: the packages that belong to it
// and the roles it is allowed to depend on.
type RoleSpec struct {
	Packages    []string `json:"packages"`
	MayDependOn []string `json:"mayDependOn"`
}
