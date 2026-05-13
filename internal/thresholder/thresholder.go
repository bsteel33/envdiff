// Package thresholder evaluates diff results against configurable thresholds
// and reports whether the results exceed acceptable limits.
package thresholder

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Options holds the threshold limits for each result category.
// A zero value means no limit is enforced for that category.
type Options struct {
	MaxMissing    int // combined missing-in-left + missing-in-right
	MaxMismatched int // keys present in both but with different values
	MaxTotal      int // total number of non-matching results
}

// Violation describes a single threshold breach.
type Violation struct {
	Field   string
	Limit   int
	Actual  int
	Message string
}

// Result is the outcome of a threshold evaluation.
type Result struct {
	Passed     bool
	Violations []Violation
}

// Evaluate checks results against the supplied Options and returns a Result.
// If opts is nil, all checks are skipped and Passed is true.
func Evaluate(results []diff.Result, opts *Options) Result {
	if opts == nil {
		return Result{Passed: true}
	}

	var missing, mismatched int
	for _, r := range results {
		switch r.Status {
		case diff.MissingInLeft, diff.MissingInRight:
			missing++
		case diff.Mismatched:
			mismatched++
		}
	}
	total := missing + mismatched

	var violations []Violation

	if opts.MaxMissing > 0 && missing > opts.MaxMissing {
		violations = append(violations, Violation{
			Field:   "missing",
			Limit:   opts.MaxMissing,
			Actual:  missing,
			Message: fmt.Sprintf("missing keys %d exceeds limit %d", missing, opts.MaxMissing),
		})
	}

	if opts.MaxMismatched > 0 && mismatched > opts.MaxMismatched {
		violations = append(violations, Violation{
			Field:   "mismatched",
			Limit:   opts.MaxMismatched,
			Actual:  mismatched,
			Message: fmt.Sprintf("mismatched keys %d exceeds limit %d", mismatched, opts.MaxMismatched),
		})
	}

	if opts.MaxTotal > 0 && total > opts.MaxTotal {
		violations = append(violations, Violation{
			Field:   "total",
			Limit:   opts.MaxTotal,
			Actual:  total,
			Message: fmt.Sprintf("total differences %d exceeds limit %d", total, opts.MaxTotal),
		})
	}

	return Result{
		Passed:     len(violations) == 0,
		Violations: violations,
	}
}
