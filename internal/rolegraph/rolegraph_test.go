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
