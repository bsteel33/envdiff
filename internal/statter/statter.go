// Package statter computes per-environment statistics from diff results,
// providing a quick overview of key counts, match rates, and drift indicators.
package statter

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Stats holds computed statistics for a single environment comparison.
type Stats struct {
	Label        string
	Total        int
	Matched      int
	MissingLeft  int
	MissingRight int
	Mismatched   int
	MatchRate    float64 // 0.0 – 1.0
	DriftScore   float64 // higher means more drift
}

// String returns a human-readable one-liner for the stats.
func (s Stats) String() string {
	return fmt.Sprintf("%s: total=%d matched=%d missingLeft=%d missingRight=%d mismatched=%d matchRate=%.2f driftScore=%.2f",
		s.Label, s.Total, s.Matched, s.MissingLeft, s.MissingRight, s.Mismatched, s.MatchRate, s.DriftScore)
}

// Compute derives Stats from a slice of diff.Result values.
// label is an arbitrary identifier for the comparison (e.g. a file path or env name).
func Compute(label string, results []diff.Result) Stats {
	s := Stats{Label: label, Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case diff.StatusMatch:
			s.Matched++
		case diff.StatusMissingInLeft:
			s.MissingLeft++
		case diff.StatusMissingInRight:
			s.MissingRight++
		case diff.StatusMismatch:
			s.Mismatched++
		}
	}
	if s.Total > 0 {
		s.MatchRate = float64(s.Matched) / float64(s.Total)
		s.DriftScore = float64(s.MissingLeft+s.MissingRight+s.Mismatched) / float64(s.Total)
	}
	return s
}

// ComputeAll derives Stats for multiple named result sets and returns them
// sorted by DriftScore descending (most drifted first).
func ComputeAll(sets map[string][]diff.Result) []Stats {
	out := make([]Stats, 0, len(sets))
	for label, results := range sets {
		out = append(out, Compute(label, results))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DriftScore != out[j].DriftScore {
			return out[i].DriftScore > out[j].DriftScore
		}
		return out[i].Label < out[j].Label
	})
	return out
}
