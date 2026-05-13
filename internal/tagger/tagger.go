// Package tagger assigns human-readable tags to diff results based on
// configurable rules, making it easier to categorise and filter output.
package tagger

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Tag represents a label attached to a diff result.
type Tag string

const (
	TagSecret  Tag = "secret"
	TagConfig  Tag = "config"
	TagFeature Tag = "feature"
	TagDB      Tag = "database"
	TagUnknown Tag = "unknown"
)

// TaggedResult pairs a diff result with its assigned tags.
type TaggedResult struct {
	Result diff.Result
	Tags   []Tag
}

// Rule maps a key pattern to a tag.
type Rule struct {
	Contains string
	Tag      Tag
}

// DefaultRules returns a sensible set of built-in tagging rules.
func DefaultRules() []Rule {
	return []Rule{
		{Contains: "SECRET", Tag: TagSecret},
		{Contains: "PASSWORD", Tag: TagSecret},
		{Contains: "TOKEN", Tag: TagSecret},
		{Contains: "KEY", Tag: TagSecret},
		{Contains: "DB_", Tag: TagDB},
		{Contains: "DATABASE", Tag: TagDB},
		{Contains: "FEATURE", Tag: TagFeature},
		{Contains: "FLAG", Tag: TagFeature},
	}
}

// Apply tags each result using the provided rules. Results that match no rule
// receive the TagUnknown tag.
func Apply(results []diff.Result, rules []Rule) []TaggedResult {
	out := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		tags := tagsFor(r.Key, rules)
		out = append(out, TaggedResult{Result: r, Tags: tags})
	}
	return out
}

func tagsFor(key string, rules []Rule) []Tag {
	upper := strings.ToUpper(key)
	seen := map[Tag]bool{}
	var tags []Tag
	for _, rule := range rules {
		if strings.Contains(upper, rule.Contains) && !seen[rule.Tag] {
			tags = append(tags, rule.Tag)
			seen[rule.Tag] = true
		}
	}
	if len(tags) == 0 {
		tags = []Tag{TagUnknown}
	}
	return tags
}
