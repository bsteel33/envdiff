package filter

import "strings"

// Options holds configuration for filtering diff results.
type Options struct {
	// Prefix restricts results to keys with the given prefix (case-insensitive).
	Prefix string
	// ExcludeKeys is a set of exact key names to omit from results.
	ExcludeKeys []string
	// OnlyMissing, when true, drops mismatched entries and keeps only missing ones.
	OnlyMissing bool
}

// Result mirrors the diff.Result structure to avoid circular imports.
type Result struct {
	Key          string
	MissingInLeft  bool
	MissingInRight bool
	LeftValue    string
	RightValue   string
}

// Apply filters a slice of Results according to the provided Options.
func Apply(results []Result, opts Options) []Result {
	excludeSet := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = struct{}{}
	}

	filtered := make([]Result, 0, len(results))
	for _, r := range results {
		if _, excluded := excludeSet[r.Key]; excluded {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(strings.ToUpper(r.Key), strings.ToUpper(opts.Prefix)) {
			continue
		}
		if opts.OnlyMissing && !r.MissingInLeft && !r.MissingInRight {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}
