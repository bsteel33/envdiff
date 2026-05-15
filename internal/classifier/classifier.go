// Package classifier categorises diff results into severity buckets
// (critical, warning, info) based on configurable rules.
package classifier

import "github.com/user/envdiff/internal/diff"

// Severity represents how serious a diff result is.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
	SeverityNone     Severity = "none"
)

// ClassifiedResult pairs a diff result with its assigned severity.
type ClassifiedResult struct {
	Result   diff.Result
	Severity Severity
}

// Options controls how results are classified.
type Options struct {
	// CriticalKeys are key names whose absence or mismatch is always critical.
	CriticalKeys []string
	// MissingIsCritical treats any missing key as critical regardless of name.
	MissingIsCritical bool
}

// Classify assigns a Severity to each result according to opts.
// Results with status diff.StatusMatch are assigned SeverityNone.
func Classify(results []diff.Result, opts *Options) []ClassifiedResult {
	if opts == nil {
		opts = &Options{}
	}

	criticalSet := make(map[string]struct{}, len(opts.CriticalKeys))
	for _, k := range opts.CriticalKeys {
		criticalSet[k] = struct{}{}
	}

	out := make([]ClassifiedResult, 0, len(results))
	for _, r := range results {
		out = append(out, ClassifiedResult{
			Result:   r,
			Severity: classify(r, criticalSet, opts),
		})
	}
	return out
}

func classify(r diff.Result, criticalSet map[string]struct{}, opts *Options) Severity {
	if r.Status == diff.StatusMatch {
		return SeverityNone
	}

	if _, ok := criticalSet[r.Key]; ok {
		return SeverityCritical
	}

	isMissing := r.Status == diff.StatusMissingInLeft || r.Status == diff.StatusMissingInRight
	if isMissing && opts.MissingIsCritical {
		return SeverityCritical
	}
	if isMissing {
		return SeverityWarning
	}

	// Mismatch
	return SeverityInfo
}
