// Package sorter provides utilities for ordering diff results
// by various criteria such as key name, status, or difference count.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortBy defines the ordering strategy for diff results.
type SortBy int

const (
	// ByKey sorts results alphabetically by key name.
	ByKey SortBy = iota
	// ByStatus groups results by their diff status.
	ByStatus
	// ByStatusThenKey groups by status, then alphabetically within each group.
	ByStatusThenKey
)

// statusOrder assigns a numeric priority to each diff status for ordering.
var statusOrder = map[diff.Status]int{
	diff.Missing:    0,
	diff.Extra:      1,
	diff.Mismatched: 2,
	diff.Equal:      3,
}

// Sort returns a new slice of Results ordered by the given strategy.
// The original slice is not modified.
func Sort(results []diff.Result, by SortBy) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)

	switch by {
	case ByKey:
		sort.Slice(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	case ByStatus:
		sort.SliceStable(out, func(i, j int) bool {
			return statusOrder[out[i].Status] < statusOrder[out[j].Status]
		})
	case ByStatusThenKey:
		sort.Slice(out, func(i, j int) bool {
			si := statusOrder[out[i].Status]
			sj := statusOrder[out[j].Status]
			if si != sj {
				return si < sj
			}
			return out[i].Key < out[j].Key
		})
	}

	return out
}
