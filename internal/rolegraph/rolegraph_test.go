package rolegraph_test

import (
	"testing"

	"github.com/Disble/dlinter-go/internal/rolegraph"
)

func TestRole_ZeroValue(t *testing.T) {
	var r rolegraph.Role

	if r != "" {
		t.Fatalf("zero value of Role must be the empty string, got %q", r)
	}
}

func TestGraph_Resolve(t *testing.T) {
	g := rolegraph.New(map[string]rolegraph.RoleDef{
		"core": {
			Packages:    []string{"internal/rolegraph"},
			MayDependOn: nil,
		},
		"adapter": {
			Packages:    []string{"pkg/analyzers/"},
			MayDependOn: []string{"core"},
		},
		"entrypoint": {
			Packages:    []string{".", "cmd/dlinter"},
			MayDependOn: []string{"core", "adapter"},
		},
	})

	tests := []struct {
		name     string
		rel      string
		wantRole rolegraph.Role
		wantOK   bool
	}{
		{
			name:     "exact match",
			rel:      "internal/rolegraph",
			wantRole: "core",
			wantOK:   true,
		},
		{
			name:     "root match",
			rel:      "",
			wantRole: "entrypoint",
			wantOK:   true,
		},
		{
			name:     "subtree prefix match",
			rel:      "pkg/analyzers/maydependon",
			wantRole: "adapter",
			wantOK:   true,
		},
		{
			name:     "unresolved",
			rel:      "internal/other",
			wantRole: "",
			wantOK:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, ok := g.Resolve(tt.rel)
			if ok != tt.wantOK {
				t.Fatalf("Resolve(%q) ok = %v, want %v", tt.rel, ok, tt.wantOK)
			}
			if role != tt.wantRole {
				t.Fatalf("Resolve(%q) role = %q, want %q", tt.rel, role, tt.wantRole)
			}
		})
	}
}

// TestGraph_Resolve_LongestPrefixPrecedence proves that among two
// overlapping subtree-prefix entries, the longest prefix wins. No such
// overlapping-prefix case exists in this repo's real .golangci.yml today,
// so this constructs a synthetic Graph inline to prove the precedence rule
// the engine documents.
func TestGraph_Resolve_LongestPrefixPrecedence(t *testing.T) {
	g := rolegraph.New(map[string]rolegraph.RoleDef{
		"broad": {
			Packages: []string{"pkg/"},
		},
		"narrow": {
			Packages: []string{"pkg/analyzers/"},
		},
	})

	role, ok := g.Resolve("pkg/analyzers/maydependon")
	if !ok {
		t.Fatalf("Resolve() ok = false, want true")
	}
	if role != "narrow" {
		t.Fatalf("Resolve() role = %q, want %q (longest prefix must win)", role, "narrow")
	}
}

func TestGraph_Allowed(t *testing.T) {
	g := rolegraph.New(map[string]rolegraph.RoleDef{
		"core": {
			MayDependOn: nil,
		},
		"adapter": {
			MayDependOn: []string{"core"},
		},
		"entrypoint": {
			MayDependOn: []string{"core", "adapter"},
		},
	})

	tests := []struct {
		name string
		from rolegraph.Role
		to   rolegraph.Role
		want bool
	}{
		{
			name: "same role always allowed",
			from: "core",
			to:   "core",
			want: true,
		},
		{
			name: "role in mayDependOn",
			from: "adapter",
			to:   "core",
			want: true,
		},
		{
			name: "role absent from mayDependOn",
			from: "core",
			to:   "adapter",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.Allowed(tt.from, tt.to)
			if got != tt.want {
				t.Fatalf("Allowed(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
