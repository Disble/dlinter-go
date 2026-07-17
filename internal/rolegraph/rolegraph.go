// Package rolegraph will hold the package-role dependency graph that a
// future dlinter analyzer evaluates to enforce mayDependOn constraints
// (see [Settings] in the root package). This slice only establishes the
// contract's shape; no evaluation logic exists yet.
//
// internal/rolegraph MUST NOT import pkg/analyzers/* or cmd/dlinter: it is
// the "core" role in the self-applied .golangci.yml role graph, and every
// other role may depend on it, never the reverse.
package rolegraph

// Role identifies an architectural role (e.g. "core", "adapter",
// "entrypoint") that a package belongs to in the role graph. The future
// evaluation logic will resolve a package's Role from Settings.Roles and
// check it against the MayDependOn list of the roles its imports belong to.
type Role string
