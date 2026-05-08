// Package profiler analyses a set of diff results and produces
// a statistical profile: key counts, status breakdown, and the
// most-frequently-differing keys across multiple comparisons.
package profiler

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Profile holds aggregate statistics derived from one or more
// sets of diff results.
type Profile struct {
	TotalKeys    int            `json:"total_keys"`
	MissingLeft  int            `json:"missing_in_left"`
	MissingRight int            `json:"missing_in_right"`
	Mismatched   int            `json:"mismatched"`
	Identical    int            `json:"identical"`
	TopDiffering []KeyFrequency `json:"top_differing"`
}

// KeyFrequency pairs a key name with how many times it appeared
// in a non-identical state across all analysed result sets.
type KeyFrequency struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// Analyse builds a Profile from one or more slices of diff.Result.
// All result sets are merged before counting.
func Analyse(resultSets ...[]diff.Result) Profile {
	freq := make(map[string]int)
	var p Profile

	for _, results := range resultSets {
		for _, r := range results {
			p.TotalKeys++
			switch r.Status {
			case diff.StatusMissingInLeft:
				p.MissingLeft++
				freq[r.Key]++
			case diff.StatusMissingInRight:
				p.MissingRight++
				freq[r.Key]++
			case diff.StatusMismatch:
				p.Mismatched++
				freq[r.Key]++
			case diff.StatusEqual:
				p.Identical++
			}
		}
	}

	p.TopDiffering = topN(freq, 5)
	return p
}

// topN returns up to n KeyFrequency entries sorted by descending count,
// with alphabetical key order as a tiebreaker.
func topN(freq map[string]int, n int) []KeyFrequency {
	list := make([]KeyFrequency, 0, len(freq))
	for k, c := range freq {
		list = append(list, KeyFrequency{Key: k, Count: c})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Key < list[j].Key
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}
