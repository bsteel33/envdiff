// Package summarizer produces a human-readable summary report across
// multiple diff result sets, aggregating counts and drift metrics.
package summarizer

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// SetSummary holds aggregated statistics for a single named result set.
type SetSummary struct {
	Label      string
	Total      int
	Matched    int
	Missing    int
	Mismatched int
	DriftPct   float64
}

// Report holds the full summarized output across all sets.
type Report struct {
	Sets       []SetSummary
	GrandTotal int
	GrandDrift float64
}

// Summarize computes a Report from a labelled map of result slices.
func Summarize(sets map[string][]diff.Result) Report {
	summaries := make([]SetSummary, 0, len(sets))

	var totalKeys, totalDrift int

	for label, results := range sets {
		s := SetSummary{Label: label, Total: len(results)}
		for _, r := range results {
			switch r.Status {
			case diff.StatusMatch:
				s.Matched++
			case diff.StatusMissingInLeft, diff.StatusMissingInRight:
				s.Missing++
			case diff.StatusMismatch:
				s.Mismatched++
			}
		}
		if s.Total > 0 {
			s.DriftPct = float64(s.Missing+s.Mismatched) / float64(s.Total) * 100
		}
		totalKeys += s.Total
		totalDrift += s.Missing + s.Mismatched
		summaries = append(summaries, s)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Label < summaries[j].Label
	})

	var grandDrift float64
	if totalKeys > 0 {
		grandDrift = float64(totalDrift) / float64(totalKeys) * 100
	}

	return Report{
		Sets:       summaries,
		GrandTotal: totalKeys,
		GrandDrift: grandDrift,
	}
}

// Render writes a plain-text summary of the Report to w.
func Render(w io.Writer, r Report) {
	if len(r.Sets) == 0 {
		fmt.Fprintln(w, "No result sets to summarize.")
		return
	}

	fmt.Fprintln(w, strings.Repeat("-", 60))
	fmt.Fprintf(w, "%-30s %6s %8s %10s %8s\n", "Label", "Total", "Matched", "Missing", "Drift%")
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, s := range r.Sets {
		fmt.Fprintf(w, "%-30s %6d %8d %10d %7.1f%%\n",
			s.Label, s.Total, s.Matched, s.Missing+s.Mismatched, s.DriftPct)
	}

	fmt.Fprintln(w, strings.Repeat("-", 60))
	fmt.Fprintf(w, "%-30s %6d %29.1f%%\n", "GRAND TOTAL", r.GrandTotal, r.GrandDrift)
}
