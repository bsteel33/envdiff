// Package drifter measures how much an environment has drifted from a
// reference set of results over time, producing a numeric drift score and
// a human-readable severity label.
package drifter

import (
	"math"

	"github.com/user/envdiff/internal/diff"
)

// Severity represents the level of drift detected.
type Severity string

const (
	SeverityNone     Severity = "none"
	SeverityLow      Severity = "low"
	SeverityModerate Severity = "moderate"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Report holds the outcome of a drift measurement.
type Report struct {
	Total     int      `json:"total"`
	Drifted   int      `json:"drifted"`
	Score     float64  `json:"score"`
	Severity  Severity `json:"severity"`
}

// Measure computes a drift report from a slice of diff results.
// Score is in the range [0, 1] where 0 means no drift and 1 means fully drifted.
func Measure(results []diff.Result) Report {
	if len(results) == 0 {
		return Report{Severity: SeverityNone}
	}

	drifted := 0
	for _, r := range results {
		if r.Status != diff.StatusMatch {
			drifted++
		}
	}

	score := math.Round(float64(drifted)/float64(len(results))*1000) / 1000

	return Report{
		Total:    len(results),
		Drifted:  drifted,
		Score:    score,
		Severity: classify(score),
	}
}

func classify(score float64) Severity {
	switch {
	case score == 0:
		return SeverityNone
	case score < 0.1:
		return SeverityLow
	case score < 0.35:
		return SeverityModerate
	case score < 0.65:
		return SeverityHigh
	default:
		return SeverityCritical
	}
}
