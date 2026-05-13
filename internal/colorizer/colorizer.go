// Package colorizer provides ANSI colour helpers for terminal output.
// It maps diff result statuses to consistent colour codes and supports
// a no-colour mode for environments that do not support ANSI escapes.
package colorizer

import "github.com/envdiff/envdiff/internal/diff"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
)

// Options controls colorizer behaviour.
type Options struct {
	// NoColor disables all ANSI escape sequences.
	NoColor bool
}

// Colorizer wraps output strings with ANSI colour codes based on diff status.
type Colorizer struct {
	opts Options
}

// New returns a Colorizer configured with opts.
func New(opts Options) *Colorizer {
	return &Colorizer{opts: opts}
}

// ForStatus returns the string s wrapped in the ANSI colour that corresponds
// to the given diff status. If NoColor is set the string is returned as-is.
func (c *Colorizer) ForStatus(s string, status diff.Status) string {
	if c.opts.NoColor {
		return s
	}
	switch status {
	case diff.StatusMissingInRight:
		return red + s + reset
	case diff.StatusMissingInLeft:
		return green + s + reset
	case diff.StatusMismatch:
		return yellow + s + reset
	case diff.StatusMatch:
		return cyan + s + reset
	default:
		return s
	}
}

// Bold wraps s in the ANSI bold escape. Returns s unchanged when NoColor is set.
func (c *Colorizer) Bold(s string) string {
	if c.opts.NoColor {
		return s
	}
	return bold + s + reset
}

// Label returns a short human-readable label for a diff status, coloured
// according to the status. Useful for table-style output.
func (c *Colorizer) Label(status diff.Status) string {
	var label string
	switch status {
	case diff.StatusMissingInRight:
		label = "MISSING_RIGHT"
	case diff.StatusMissingInLeft:
		label = "MISSING_LEFT"
	case diff.StatusMismatch:
		label = "MISMATCH"
	case diff.StatusMatch:
		label = "MATCH"
	default:
		label = "UNKNOWN"
	}
	return c.ForStatus(label, status)
}
