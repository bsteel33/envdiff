// Package templater generates .env.example files from existing env maps,
// replacing values with placeholder descriptions.
package templater

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how the template is rendered.
type Options struct {
	// Placeholder is used as the value for every key. Defaults to "<value>".
	Placeholder string
	// CommentPrefix is prepended to each key as a header comment when non-empty.
	CommentPrefix string
	// SensitivePatterns are substrings that, when found in a key (case-insensitive),
	// cause the placeholder to indicate sensitivity.
	SensitivePatterns []string
}

var defaultSensitivePatterns = []string{"SECRET", "PASSWORD", "PASS", "TOKEN", "KEY", "PRIVATE"}

// Generate writes a .env.example template to w from the provided env map.
// Keys are sorted alphabetically for deterministic output.
func Generate(w io.Writer, env map[string]string, opts Options) error {
	if opts.Placeholder == "" {
		opts.Placeholder = "<value>"
	}
	if len(opts.SensitivePatterns) == 0 {
		opts.SensitivePatterns = defaultSensitivePatterns
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		placeholder := opts.Placeholder
		if isSensitive(k, opts.SensitivePatterns) {
			placeholder = "<secret>"
		}
		if opts.CommentPrefix != "" {
			if _, err := fmt.Fprintf(w, "# %s%s\n", opts.CommentPrefix, k); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, placeholder); err != nil {
			return err
		}
	}
	return nil
}

// isSensitive returns true if key contains any of the given patterns (case-insensitive).
func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
