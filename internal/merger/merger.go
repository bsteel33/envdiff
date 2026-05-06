// Package merger provides functionality to merge diff results from multiple
// environment comparisons into a unified report.
package merger

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// PairResult holds the diff result for a specific pair of environment files.
type PairResult struct {
	Left  string
	Right string
	Diffs []diff.Result
}

// MergedResult aggregates all unique keys and their presence across all pairs.
type MergedResult struct {
	// Keys maps each key to the set of pairs where it appeared as a difference.
	Keys map[string][]PairResult
	// Pairs holds all individual pair results.
	Pairs []PairResult
}

// Merge combines multiple PairResults into a single MergedResult,
// grouping differences by key name for easier analysis.
func Merge(pairs []PairResult) MergedResult {
	keyMap := make(map[string][]PairResult)

	for _, pair := range pairs {
		seen := make(map[string]bool)
		for _, r := range pair.Diffs {
			if !seen[r.Key] {
				seen[r.Key] = true
				keyMap[r.Key] = append(keyMap[r.Key], pair)
			}
		}
	}

	return MergedResult{
		Keys:  keyMap,
		Pairs: pairs,
	}
}

// KeysWithMostDifferences returns keys sorted by how many pairs they appear in,
// descending. Useful for identifying the most problematic keys.
func (m MergedResult) KeysWithMostDifferences() []string {
	type entry struct {
		key   string
		count int
	}

	entries := make([]entry, 0, len(m.Keys))
	for k, pairs := range m.Keys {
		entries = append(entries, entry{key: k, count: len(pairs)})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].key < entries[j].key
	})

	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.key
	}
	return keys
}
