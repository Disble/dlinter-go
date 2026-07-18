package rolegraph

import "testing"

// Internal test file (package rolegraph, not rolegraph_test): matches and
// betterPrefix are unexported helpers extracted from Graph.Resolve per
// Design D7. They are pure and trivially unit-testable, which is the point
// of the extraction — no reason to widen the public API just to test them.

func TestMatches(t *testing.T) {
	tests := []struct {
		name string
		p    pattern
		rel  string
		want bool
	}{
		{
			name: "root pattern matches only empty rel",
			p:    pattern{role: "entrypoint", kind: patternRoot},
			rel:  "",
			want: true,
		},
		{
			name: "root pattern does not match non-empty rel",
			p:    pattern{role: "entrypoint", kind: patternRoot},
			rel:  "cmd/dlinter",
			want: false,
		},
		{
			name: "exact pattern matches identical rel",
			p:    pattern{role: "core", kind: patternExact, entry: "internal/rolegraph"},
			rel:  "internal/rolegraph",
			want: true,
		},
		{
			name: "exact pattern does not match a subtree of itself",
			p:    pattern{role: "core", kind: patternExact, entry: "internal/rolegraph"},
			rel:  "internal/rolegraph/sub",
			want: false,
		},
		{
			name: "prefix pattern matches the entry itself",
			p:    pattern{role: "adapter", kind: patternPrefix, entry: "pkg/analyzers"},
			rel:  "pkg/analyzers",
			want: true,
		},
		{
			name: "prefix pattern matches a subtree member",
			p:    pattern{role: "adapter", kind: patternPrefix, entry: "pkg/analyzers"},
			rel:  "pkg/analyzers/maydependon",
			want: true,
		},
		{
			name: "prefix pattern does not match a sibling with a shared string prefix",
			p:    pattern{role: "adapter", kind: patternPrefix, entry: "pkg/analyzers"},
			rel:  "pkg/analyzersextra",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matches(tt.p, tt.rel)
			if got != tt.want {
				t.Fatalf("matches(%+v, %q) = %v, want %v", tt.p, tt.rel, got, tt.want)
			}
		})
	}
}

func TestBetterPrefix(t *testing.T) {
	tests := []struct {
		name      string
		candidate pattern
		best      pattern
		bestLen   int
		want      bool
	}{
		{
			name:      "longer entry wins",
			candidate: pattern{role: "narrow", kind: patternPrefix, entry: "pkg/analyzers"},
			best:      pattern{role: "broad", kind: patternPrefix, entry: "pkg"},
			bestLen:   len("pkg"),
			want:      true,
		},
		{
			name:      "shorter entry loses",
			candidate: pattern{role: "broad", kind: patternPrefix, entry: "pkg"},
			best:      pattern{role: "narrow", kind: patternPrefix, entry: "pkg/analyzers"},
			bestLen:   len("pkg/analyzers"),
			want:      false,
		},
		{
			name:      "equal length breaks tie by lower role name",
			candidate: pattern{role: "a-role", kind: patternPrefix, entry: "pkg"},
			best:      pattern{role: "z-role", kind: patternPrefix, entry: "pkg"},
			bestLen:   len("pkg"),
			want:      true,
		},
		{
			name:      "equal length, higher role name does not win",
			candidate: pattern{role: "z-role", kind: patternPrefix, entry: "pkg"},
			best:      pattern{role: "a-role", kind: patternPrefix, entry: "pkg"},
			bestLen:   len("pkg"),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := betterPrefix(tt.candidate, tt.best, tt.bestLen)
			if got != tt.want {
				t.Fatalf("betterPrefix(%+v, %+v, %d) = %v, want %v", tt.candidate, tt.best, tt.bestLen, got, tt.want)
			}
		})
	}
}
