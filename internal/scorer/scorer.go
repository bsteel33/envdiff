// Package scorer assigns a numeric drift score to a set of diff results,
// giving callers a quick way to quantify how far two environments have
// diverged from each other.
package scorer

import "github.com/user/envdiff/internal/diff"

// Weights controls how many points each kind of difference contributes.
type Weights struct {
	MissingInLeft  float64
	MissingInRight float64
	Mismatched     float64
}

// DefaultWeights returns the standard scoring weights.
func DefaultWeights() Weights {
	return Weights{
		MissingInLeft:  1.0,
		MissingInRight: 1.0,
		Mismatched:     2.0,
	}
}

// Score holds the computed drift score and a breakdown by category.
type Score struct {
	Total          float64
	MissingInLeft  int
	MissingInRight int
	Mismatched     int
}

// Grade returns a human-readable grade for the total score.
func (s Score) Grade() string {
	switch {
	case s.Total == 0:
		return "A"
	case s.Total <= 3:
		return "B"
	case s.Total <= 8:
		return "C"
	case s.Total <= 15:
		return "D"
	default:
		return "F"
	}
}

// Compute calculates a drift score for the provided results using the given
// weights. Pass DefaultWeights() for standard behaviour.
func Compute(results []diff.Result, w Weights) Score {
	var s Score
	for _, r := range results {
		switch r.Status {
		case diff.MissingInLeft:
			s.MissingInLeft++
		case diff.MissingInRight:
			s.MissingInRight++
		case diff.Mismatched:
			s.Mismatched++
		}
	}
	s.Total = float64(s.MissingInLeft)*w.MissingInLeft +
		float64(s.MissingInRight)*w.MissingInRight +
		float64(s.Mismatched)*w.Mismatched
	return s
}
