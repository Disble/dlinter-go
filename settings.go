package dlinter

// Settings is the top-level configuration decoded from .golangci.yml for the
// dlinter module linter. It defines the package-role graph that the future
// mayDependOn rule will enforce; this slice only proves the decode path.
type Settings struct {
	Roles map[string]RoleSpec `json:"roles"`
}

// RoleSpec describes one architectural role: the packages that belong to it
// and the roles it is allowed to depend on.
type RoleSpec struct {
	Packages    []string `json:"packages"`
	MayDependOn []string `json:"mayDependOn"`
}
