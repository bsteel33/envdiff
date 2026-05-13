// Package deduplicator removes duplicate diff results, keeping the entry
// with the highest severity when the same key appears more than once.
package deduplicator

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// severityRank assigns a numeric rank to each diff status so that
// more-severe statuses win when deduplicating.
var severityRank = map[diff.Status]int{
	diff.StatusMatch:          0,
	diff.StatusMissingInRight: 2,
	diff.StatusMissingInLeft:  2,
	diff.StatusMismatch:       3,
}

// Apply returns a new slice of Results with duplicate keys removed.
// When the same key appears multiple times the entry with the highest
// severity rank is retained; ties are broken by keeping the first
// occurrence in the original order.
func Apply(results []diff.Result) []diff.Result {
	if len(results) == 0 {
		return []diff.Result{}
	}

	// Track the best (highest-rank) result seen for each key.
	type indexed struct {
		result diff.Result
		pos    int
	}
	best := make(map[string]indexed, len(results))

	for i, r := range results {
		prev, seen := best[r.Key]
		if !seen {
			best[r.Key] = indexed{r, i}
			continue
		}
		if severityRank[r.Status] > severityRank[prev.result.Status] {
			best[r.Key] = indexed{r, prev.pos} // keep original position
		}
	}

	// Rebuild slice in original insertion order.
	out := make([]diff.Result, 0, len(best))
	seen := make(map[string]bool, len(best))
	for _, r := range results {
		if seen[r.Key] {
			continue
		}
		seen[r.Key] = true
		out = append(out, best[r.Key].result)
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})

	return out
}
