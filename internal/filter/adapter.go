// Package filter provides utilities for narrowing down diff results
// based on key prefixes, exclusion lists, and result type.
package filter

import "github.com/yourorg/envdiff/internal/diff"

// FromDiffResults converts a slice of diff.Result into filter.Result
// so that the filter package remains decoupled from the diff package
// while still being usable from the top-level command.
func FromDiffResults(in []diff.Result) []Result {
	out := make([]Result, len(in))
	for i, r := range in {
		out[i] = Result{
			Key:            r.Key,
			MissingInLeft:  r.MissingInLeft,
			MissingInRight: r.MissingInRight,
			LeftValue:      r.LeftValue,
			RightValue:     r.RightValue,
		}
	}
	return out
}

// ToDiffResults converts filter.Result back to diff.Result for use
// with the reporter package.
func ToDiffResults(in []Result) []diff.Result {
	out := make([]diff.Result, len(in))
	for i, r := range in {
		out[i] = diff.Result{
			Key:            r.Key,
			MissingInLeft:  r.MissingInLeft,
			MissingInRight: r.MissingInRight,
			LeftValue:      r.LeftValue,
			RightValue:     r.RightValue,
		}
	}
	return out
}
