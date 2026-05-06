// Package patcher provides functionality to generate patch suggestions
// for missing or mismatched keys across .env files.
package patcher

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Patch represents a suggested change for a single key.
type Patch struct {
	Key      string
	Action   string // "add", "update", "remove"
	OldValue string
	NewValue string
}

// String returns a human-readable representation of the patch.
func (p Patch) String() string {
	switch p.Action {
	case "add":
		return fmt.Sprintf("+ %s=%s", p.Key, p.NewValue)
	case "remove":
		return fmt.Sprintf("- %s=%s", p.Key, p.OldValue)
	case "update":
		return fmt.Sprintf("~ %s: %q -> %q", p.Key, p.OldValue, p.NewValue)
	default:
		return fmt.Sprintf("? %s", p.Key)
	}
}

// Generate produces a list of Patch suggestions based on diff results.
// The patches describe what changes would make the left env match the right env.
func Generate(results []diff.Result) []Patch {
	patches := make([]Patch, 0, len(results))
	for _, r := range results {
		switch r.Status {
		case diff.MissingInLeft:
			patches = append(patches, Patch{
				Key:      r.Key,
				Action:   "add",
				NewValue: r.RightValue,
			})
		case diff.MissingInRight:
			patches = append(patches, Patch{
				Key:      r.Key,
				Action:   "remove",
				OldValue: r.LeftValue,
			})
		case diff.Mismatched:
			patches = append(patches, Patch{
				Key:      r.Key,
				Action:   "update",
				OldValue: r.LeftValue,
				NewValue: r.RightValue,
			})
		}
	}
	return patches
}

// RenderEnv returns a minimal .env-style string with all patch suggestions applied.
func RenderEnv(patches []Patch) string {
	var sb strings.Builder
	for _, p := range patches {
		if p.Action == "add" || p.Action == "update" {
			fmt.Fprintf(&sb, "%s=%s\n", p.Key, p.NewValue)
		}
	}
	return sb.String()
}
