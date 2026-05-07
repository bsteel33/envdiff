// Package redactor provides utilities for masking sensitive values
// in diff results before display or export.
package redactor

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// DefaultSensitivePatterns contains common substrings that indicate a key
// holds a sensitive value (case-insensitive match).
var DefaultSensitivePatterns = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"apikey",
	"api_key",
	"private",
	"credentials",
	"auth",
}

const redactedPlaceholder = "***REDACTED***"

// Redactor masks sensitive values in diff results.
type Redactor struct {
	patterns []string
}

// New returns a Redactor using the provided patterns.
// If patterns is nil, DefaultSensitivePatterns is used.
func New(patterns []string) *Redactor {
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}
	return &Redactor{patterns: patterns}
}

// isSensitive returns true if key matches any sensitive pattern.
func (r *Redactor) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// Apply returns a new slice of results with sensitive values replaced
// by the redacted placeholder. Original results are not mutated.
func (r *Redactor) Apply(results []diff.Result) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, res := range results {
		if r.isSensitive(res.Key) {
			if res.LeftValue != "" {
				res.LeftValue = redactedPlaceholder
			}
			if res.RightValue != "" {
				res.RightValue = redactedPlaceholder
			}
		}
		out[i] = res
	}
	return out
}
