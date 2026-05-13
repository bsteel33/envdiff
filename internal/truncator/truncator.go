// Package truncator shortens long env values in output to a configurable
// maximum length, appending an ellipsis so readers know the value was cut.
package truncator

import "github.com/your-org/envdiff/internal/diff"

const (
	DefaultMaxLen  = 64
	DefaultEllipsis = "..."
)

// Options controls truncation behaviour.
type Options struct {
	// MaxLen is the maximum number of runes to keep (default: DefaultMaxLen).
	MaxLen int
	// Ellipsis is appended when a value is truncated (default: "...").
	Ellipsis string
}

func (o *Options) maxLen() int {
	if o == nil || o.MaxLen <= 0 {
		return DefaultMaxLen
	}
	return o.MaxLen
}

func (o *Options) ellipsis() string {
	if o == nil || o.Ellipsis == "" {
		return DefaultEllipsis
	}
	return o.Ellipsis
}

// Apply returns a new slice of diff.Result with long values truncated.
// The original slice is never mutated.
func Apply(results []diff.Result, opts *Options) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, r := range results {
		out[i] = diff.Result{
			Key:      r.Key,
			Status:   r.Status,
			LeftVal:  truncate(r.LeftVal, opts),
			RightVal: truncate(r.RightVal, opts),
		}
	}
	return out
}

func truncate(s string, opts *Options) string {
	runes := []rune(s)
	max := opts.maxLen()
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + opts.ellipsis()
}
