// Package multiplexer fans out a single base env map against multiple
// target env maps, returning a slice of named diff results.
package multiplexer

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// NamedResult pairs a target label with its diff results.
type NamedResult struct {
	Label   string
	Results []diff.Result
}

// Run compares base against each entry in targets and returns one
// NamedResult per target, ordered by label.
//
// base is the reference environment (e.g. .env.production).
// targets maps a human-readable label to the parsed env map for that
// environment.
func Run(base map[string]string, targets map[string]map[string]string) []NamedResult {
	if len(targets) == 0 {
		return nil
	}

	labels := make([]string, 0, len(targets))
	for label := range targets {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	out := make([]NamedResult, 0, len(labels))
	for _, label := range labels {
		results := diff.Compare(base, targets[label])
		out = append(out, NamedResult{
			Label:   label,
			Results: results,
		})
	}
	return out
}

// Flatten merges all NamedResults into a single deduplicated slice of
// diff.Result values. Where the same key appears in multiple targets
// the entry with the highest-severity status is kept (mismatched >
// missing > ok).
func Flatten(named []NamedResult) []diff.Result {
	type entry struct {
		result diff.Result
		weight int
	}

	weight := func(s diff.Status) int {
		switch s {
		case diff.Mismatched:
			return 3
		case diff.MissingInLeft, diff.MissingInRight:
			return 2
		default:
			return 1
		}
	}

	index := make(map[string]entry)
	for _, nr := range named {
		for _, r := range nr.Results {
			w := weight(r.Status)
			if existing, ok := index[r.Key]; !ok || w > existing.weight {
				index[r.Key] = entry{result: r, weight: w}
			}
		}
	}

	out := make([]diff.Result, 0, len(index))
	for _, e := range index {
		out = append(out, e.result)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}
