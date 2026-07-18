// Package rolegraph holds the package-role dependency graph that the
// mayDependOn analyzer evaluates to enforce import-direction constraints
// (see [Settings] in the root package). It resolves a package's role from
// its module-relative path and answers whether one role may depend on
// another.
//
// internal/rolegraph MUST NOT import pkg/analyzers/* or cmd/dlinter: it is
// the "core" role in the self-applied .golangci.yml role graph, and every
// other role may depend on it, never the reverse.
package rolegraph

import "strings"

// Role identifies an architectural role (e.g. "core", "adapter",
// "entrypoint") that a package belongs to in the role graph. A package's
// Role is resolved from Settings.Roles via Graph.Resolve and checked
// against the MayDependOn list of the roles its imports belong to, via
// Graph.Allowed.
type Role string

// RoleDef describes one architectural role: the packages that belong to it
// and the roles it is allowed to depend on. It mirrors the shape of the
// root package's RoleSpec, but lives in this package so that
// internal/rolegraph never has to import the root package (which would
// violate its own core role).
type RoleDef struct {
	Packages    []string
	MayDependOn []string
}

// patternKind classifies a configured package entry so Resolve can apply
// the correct matching rule.
type patternKind int

const (
	patternExact patternKind = iota
	patternRoot
	patternPrefix
)

// pattern is one compiled entry: a role's configured package string,
// classified into a matching rule.
type pattern struct {
	role  Role
	kind  patternKind
	entry string // for patternPrefix, the prefix without trailing slash; for patternExact, the exact relative path
}

// Graph is the compiled package-role dependency graph. It is built once via
// New from the decoded Settings.Roles configuration and then queried via
// Resolve and Allowed.
type Graph struct {
	patterns    []pattern
	mayDependOn map[Role]map[Role]bool
}

// New builds a Graph from role definitions keyed by role name.
//
// Each entry in a RoleDef.Packages list is classified as follows:
//   - "." matches only the module-root package (module-relative path "").
//   - a string ending in "/" matches as a subtree prefix.
//   - any other string matches exactly.
func New(roles map[string]RoleDef) *Graph {
	g := &Graph{
		mayDependOn: make(map[Role]map[Role]bool, len(roles)),
	}

	for name, def := range roles {
		role := Role(name)

		for _, entry := range def.Packages {
			g.patterns = append(g.patterns, classify(role, entry))
		}

		allowed := make(map[Role]bool, len(def.MayDependOn))
		for _, dep := range def.MayDependOn {
			allowed[Role(dep)] = true
		}
		g.mayDependOn[role] = allowed
	}

	return g
}

// classify compiles one configured package entry into a pattern, choosing
// the matching rule from the entry's shape: "." is the module root, a
// trailing slash marks a subtree prefix, and anything else matches exactly.
func classify(role Role, entry string) pattern {
	switch {
	case entry == ".":
		return pattern{role: role, kind: patternRoot}
	case strings.HasSuffix(entry, "/"):
		return pattern{role: role, kind: patternPrefix, entry: strings.TrimSuffix(entry, "/")}
	default:
		return pattern{role: role, kind: patternExact, entry: entry}
	}
}

// Resolve returns the role of the package whose module-relative path is
// rel, and whether a role was found.
//
// Precedence: an exact or root match always wins over any prefix match.
// Among prefix matches, the longest prefix wins. If two prefixes of equal
// length collide (no such collision exists in this repo's current
// .golangci.yml), the tie is broken deterministically by role name.
func (g *Graph) Resolve(rel string) (Role, bool) {
	var (
		bestPrefix    pattern
		bestPrefixLen = -1
		havePrefix    bool
	)

	for _, p := range g.patterns {
		if p.kind != patternPrefix {
			if matches(p, rel) {
				return p.role, true
			}
			continue
		}
		if matches(p, rel) && betterPrefix(p, bestPrefix, bestPrefixLen) {
			bestPrefix = p
			bestPrefixLen = len(p.entry)
			havePrefix = true
		}
	}

	if havePrefix {
		return bestPrefix.role, true
	}
	return "", false
}

// matches reports whether rel is matched by pattern p, per p's kind: a root
// pattern matches only the module-root package ("" ), an exact pattern
// matches only its literal entry, and a prefix pattern matches its entry or
// anything in its subtree.
func matches(p pattern, rel string) bool {
	switch p.kind {
	case patternRoot:
		return rel == ""
	case patternExact:
		return rel == p.entry
	case patternPrefix:
		return rel == p.entry || strings.HasPrefix(rel, p.entry+"/")
	default:
		return false
	}
}

// betterPrefix reports whether candidate should replace best as the current
// longest-prefix match, given best's entry length is bestLen. A longer entry
// always wins; among equal-length entries, the tie is broken deterministically
// by the lower role name.
func betterPrefix(candidate, best pattern, bestLen int) bool {
	if len(candidate.entry) != bestLen {
		return len(candidate.entry) > bestLen
	}
	return candidate.role < best.role
}

// Allowed reports whether role from may depend on role to. A role may
// always depend on itself; otherwise to must appear in from's configured
// MayDependOn list.
func (g *Graph) Allowed(from, to Role) bool {
	if from == to {
		return true
	}
	return g.mayDependOn[from][to]
}
