package reporter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Summary holds aggregated counts for a diff result set.
type Summary struct {
	MissingInRight int
	MissingInLeft  int
	Mismatched     int
	Total          int
}

// Summarize computes a Summary from a slice of diff.Result.
func Summarize(results []diff.Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case diff.MissingInRight:
			s.MissingInRight++
		case diff.MissingInLeft:
			s.MissingInLeft++
		case diff.Mismatched:
			s.Mismatched++
		}
	}
	return s
}

// FormatSummary returns a human-readable summary line.
func FormatSummary(s Summary) string {
	if s.Total == 0 {
		return "No differences found."
	}
	parts := []string{}
	if s.MissingInRight > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in right", s.MissingInRight))
	}
	if s.MissingInLeft > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in left", s.MissingInLeft))
	}
	if s.Mismatched > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", s.Mismatched))
	}
	return fmt.Sprintf("Found %d difference(s): %s.", s.Total, strings.Join(parts, ", "))
}
