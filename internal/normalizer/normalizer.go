// Package normalizer provides utilities to normalize .env key-value pairs
// before comparison, such as trimming whitespace and case-folding keys.
package normalizer

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls which normalizations are applied.
type Options struct {
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// LowercaseKeys folds all keys to lowercase before comparison.
	LowercaseKeys bool
	// CollapseWhitespace replaces runs of whitespace in values with a single space.
	CollapseWhitespace bool
}

// DefaultOptions returns a sensible default: trim values only.
func DefaultOptions() Options {
	return Options{
		TrimValues: true,
		LowercaseKeys: false,
		CollapseWhitespace: false,
	}
}

// Apply normalizes a slice of diff.Result according to the given Options.
// It returns a new slice; the originals are not mutated.
func Apply(results []diff.Result, opts Options) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, r := range results {
		out[i] = normalizeResult(r, opts)
	}
	return out
}

func normalizeResult(r diff.Result, opts Options) diff.Result {
	key := r.Key
	if opts.LowercaseKeys {
		key = strings.ToLower(key)
	}

	leftVal := r.LeftValue
	rightVal := r.RightValue

	if opts.TrimValues {
		leftVal = strings.TrimSpace(leftVal)
		rightVal = strings.TrimSpace(rightVal)
	}

	if opts.CollapseWhitespace {
		leftVal = collapseWS(leftVal)
		rightVal = collapseWS(rightVal)
	}

	return diff.Result{
		Key:        key,
		LeftValue:  leftVal,
		RightValue: rightVal,
		Status:     r.Status,
	}
}

// collapseWS replaces consecutive whitespace characters with a single space.
func collapseWS(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}
