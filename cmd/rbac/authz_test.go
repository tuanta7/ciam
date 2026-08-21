package main

import (
	"context"
	"testing"
)

func TestAuthorize(t *testing.T) {
	cases := []struct {
		role, scope string
		want        bool
	}{
		{"admin", "clients:delete", true},
		{"editor", "clients:write", true},
		{"viewer", "clients:read", true},
		{"viewer", "clients:write", false},
		{"unknown", "clients:read", false},
	}

	for _, c := range cases {
		got, err := Authorize(context.Background(), c.role, c.scope)
		if err != nil {
			t.Fatalf("Authorize(%q, %q) error: %v", c.role, c.scope, err)
		}
		if got != c.want {
			t.Errorf("Authorize(%q, %q) = %v, want %v", c.role, c.scope, got, c.want)
		}
	}
}
