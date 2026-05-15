// Package censor replaces sensitive key values with a configurable mask
// before any output is produced, ensuring secrets never leak into reports.
package censor

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// DefaultSensitiveSubstrings is the list of substrings that, when found in a
// key name (case-insensitive), cause the value to be censored.
var DefaultSensitiveSubstrings = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"private", "credential", "auth", "cert", "key",
}

// Options controls how censoring is applied.
type Options struct {
	// Mask is the string substituted for sensitive values. Defaults to "***".
	Mask string
	// SensitiveSubstrings overrides DefaultSensitiveSubstrings when non-nil.
	SensitiveSubstrings []string
}

// Apply returns a new slice of Results with sensitive values replaced by the
// configured mask. The original slice is never modified.
func Apply(results []diff.Result, opts *Options) []diff.Result {
	mask, patterns := resolve(opts)

	out := make([]diff.Result, len(results))
	for i, r := range results {
		if isSensitive(r.Key, patterns) {
			r.LeftValue = maskValue(r.LeftValue, mask)
			r.RightValue = maskValue(r.RightValue, mask)
		}
		out[i] = r
	}
	return out
}

func resolve(opts *Options) (string, []string) {
	if opts == nil {
		return "***", DefaultSensitiveSubstrings
	}
	mask := opts.Mask
	if mask == "" {
		mask = "***"
	}
	patterns := opts.SensitiveSubstrings
	if patterns == nil {
		patterns = DefaultSensitiveSubstrings
	}
	return mask, patterns
}

func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func maskValue(v, mask string) string {
	if v == "" {
		return v
	}
	return mask
}
