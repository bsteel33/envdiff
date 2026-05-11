// Package grouper groups diff results by a shared key prefix (e.g. "DB_", "AWS_").
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Group holds all results that share a common prefix.
type Group struct {
	Prefix  string
	Results []diff.Result
}

// ByPrefix partitions results into groups based on the first segment of each
// key when split by sep (typically "_"). Keys with no separator are placed
// under the empty-string prefix group.
func ByPrefix(results []diff.Result, sep string) []Group {
	index := map[string][]diff.Result{}

	for _, r := range results {
		prefix := prefixOf(r.Key, sep)
		index[prefix] = append(index[prefix], r)
	}

	groups := make([]Group, 0, len(index))
	for prefix, rs := range index {
		groups = append(groups, Group{Prefix: prefix, Results: rs})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})

	return groups
}

// prefixOf returns the portion of key before the first occurrence of sep.
// If sep is not found the whole key is returned as the prefix.
func prefixOf(key, sep string) string {
	if sep == "" {
		return ""
	}
	if idx := strings.Index(key, sep); idx >= 0 {
		return key[:idx]
	}
	return key
}
