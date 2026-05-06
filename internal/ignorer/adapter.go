package ignorer

import "github.com/user/envdiff/internal/diff"

// FilterResults removes any diff.Result entries whose Key is present in the
// Ignorer. The original slice is not mutated; a new slice is returned.
func FilterResults(results []diff.Result, ig *Ignorer) []diff.Result {
	if ig == nil {
		return results
	}
	out := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if !ig.Contains(r.Key) {
			out = append(out, r)
		}
	}
	return out
}
