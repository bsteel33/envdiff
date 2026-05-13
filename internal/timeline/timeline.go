// Package timeline tracks how diff results change across multiple snapshots
// over time, providing a chronological view of environment drift.
package timeline

import (
	"sort"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Entry represents a single point-in-time observation of diff results.
type Entry struct {
	Timestamp time.Time
	Label     string
	Results   []diff.Result
}

// Timeline holds an ordered sequence of entries.
type Timeline struct {
	entries []Entry
}

// Add appends a new entry to the timeline.
func (t *Timeline) Add(label string, ts time.Time, results []diff.Result) {
	t.entries = append(t.entries, Entry{
		Timestamp: ts,
		Label:     label,
		Results:   results,
	})
}

// Entries returns all entries sorted chronologically.
func (t *Timeline) Entries() []Entry {
	sorted := make([]Entry, len(t.entries))
	copy(sorted, t.entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})
	return sorted
}

// Trend returns a slice of TrendPoint summarising each entry.
type TrendPoint struct {
	Timestamp  time.Time
	Label      string
	Total      int
	Missing    int
	Mismatched int
}

// Trend computes a TrendPoint for every entry in chronological order.
func (t *Timeline) Trend() []TrendPoint {
	entries := t.Entries()
	points := make([]TrendPoint, 0, len(entries))
	for _, e := range entries {
		var missing, mismatched int
		for _, r := range e.Results {
			switch r.Status {
			case diff.MissingInLeft, diff.MissingInRight:
				missing++
			case diff.Mismatched:
				mismatched++
			}
		}
		points = append(points, TrendPoint{
			Timestamp:  e.Timestamp,
			Label:      e.Label,
			Total:      len(e.Results),
			Missing:    missing,
			Mismatched: mismatched,
		})
	}
	return points
}
