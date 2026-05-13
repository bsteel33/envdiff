package tagger

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable tagged diff report to w.
func ReportText(w io.Writer, tagged []TaggedResult) {
	if len(tagged) == 0 {
		fmt.Fprintln(w, "No tagged results.")
		return
	}
	for _, t := range tagged {
		labels := joinTags(t.Tags)
		fmt.Fprintf(w, "[%s] %s (%s)\n", labels, t.Result.Key, t.Result.Status)
	}
}

// ReportJSON writes a JSON array of tagged results to w.
func ReportJSON(w io.Writer, tagged []TaggedResult) error {
	type jsonEntry struct {
		Key    string   `json:"key"`
		Status string   `json:"status"`
		Tags   []string `json:"tags"`
	}
	entries := make([]jsonEntry, 0, len(tagged))
	for _, t := range tagged {
		tags := make([]string, len(t.Tags))
		for i, tag := range t.Tags {
			tags[i] = string(tag)
		}
		entries = append(entries, jsonEntry{
			Key:    t.Result.Key,
			Status: string(t.Result.Status),
			Tags:   tags,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func joinTags(tags []Tag) string {
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = string(t)
	}
	return strings.Join(parts, ",")
}
