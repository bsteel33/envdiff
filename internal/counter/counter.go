// Package counter provides utilities for counting and summarising
// diff result statistics across one or more environment comparisons.
package counter

import "github.com/your-org/envdiff/internal/diff"

// Stats holds aggregate counts for a set of diff results.
type Stats struct {
	Total      int
	Matched    int
	Missing    int // missing in either side
	Mismatched int
	MissingLeft  int
	MissingRight int
}

// Count tallies the diff results and returns a Stats summary.
func Count(results []diff.Result) Stats {
	var s Stats
	for _, r := range results {
		s.Total++
		switch r.Status {
		case diff.StatusMatch:
			s.Matched++
		case diff.StatusMissingInLeft:
			s.Missing++
			s.MissingLeft++
		case diff.StatusMissingInRight:
			s.Missing++
			s.MissingRight++
		case diff.StatusMismatch:
			s.Mismatched++
		}
	}
	return s
}

// DriftRatio returns the fraction of results that are not a clean match.
// Returns 0 if there are no results.
func DriftRatio(s Stats) float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Total-s.Matched) / float64(s.Total)
}

// IsClean reports whether all results are matched (no drift).
func IsClean(s Stats) bool {
	return s.Total > 0 && s.Matched == s.Total
}
