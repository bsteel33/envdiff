// Package renamer provides utilities for detecting and suggesting key renames
// across .env file diff results. A rename is inferred when a key disappears
// from one side and a new key appears on the other with a similar value.
package renamer

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Suggestion represents a probable key rename between two environments.
type Suggestion struct {
	OldKey   string
	NewKey   string
	OldValue string
	NewValue string
	Score    float64 // 0.0–1.0 similarity confidence
}

// Detect analyses diff results and returns rename suggestions.
// It pairs keys missing in the right side with keys missing in the left side
// whose values are sufficiently similar.
func Detect(results []diff.Result, threshold float64) []Suggestion {
	if threshold <= 0 || threshold > 1 {
		threshold = 0.8
	}

	var missingRight []diff.Result // present in left, absent in right
	var missingLeft []diff.Result  // present in right, absent in left

	for _, r := range results {
		switch r.Status {
		case diff.MissingInRight:
			missingRight = append(missingRight, r)
		case diff.MissingInLeft:
			missingLeft = append(missingLeft, r)
		}
	}

	var suggestions []Suggestion
	used := make(map[int]bool)

	for _, right := range missingRight {
		best := -1
		bestScore := 0.0
		for i, left := range missingLeft {
			if used[i] {
				continue
			}
			s := valueSimilarity(right.LeftValue, left.RightValue)
			if s >= threshold && s > bestScore {
				bestScore = s
				best = i
			}
		}
		if best >= 0 {
			used[best] = true
			suggestions = append(suggestions, Suggestion{
				OldKey:   right.Key,
				NewKey:   missingLeft[best].Key,
				OldValue: right.LeftValue,
				NewValue: missingLeft[best].RightValue,
				Score:    bestScore,
			})
		}
	}

	return suggestions
}

// valueSimilarity returns a simple similarity score between two strings
// based on the proportion of shared characters (Dice coefficient on bigrams).
func valueSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	biA := bigrams(a)
	biB := bigrams(b)
	if len(biA) == 0 && len(biB) == 0 {
		return 1.0
	}
	if len(biA) == 0 || len(biB) == 0 {
		return 0.0
	}
	intersection := 0
	for k, v := range biA {
		if vb, ok := biB[k]; ok {
			if v < vb {
				intersection += v
			} else {
				intersection += vb
			}
		}
	}
	return float64(2*intersection) / float64(len(biA)+len(biB))
}

func bigrams(s string) map[string]int {
	m := make(map[string]int)
	for i := 0; i+1 < len(s); i++ {
		m[s[i:i+2]]++
	}
	return m
}
