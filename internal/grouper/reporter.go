package grouper

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable grouped summary to w.
func ReportText(w io.Writer, groups []Group) {
	if len(groups) == 0 {
		fmt.Fprintln(w, "No groups to report.")
		return
	}
	for _, g := range groups {
		label := g.Prefix
		if label == "" {
			label = "(no prefix)"
		}
		fmt.Fprintf(w, "[%s] — %d key(s)\n", label, len(g.Results))
		for _, r := range g.Results {
			fmt.Fprintf(w, "  %-40s %s\n", r.Key, statusLabel(r.Status))
		}
	}
}

// jsonGroup is the JSON representation of a single group.
type jsonGroup struct {
	Prefix  string       `json:"prefix"`
	Count   int          `json:"count"`
	Results []jsonResult `json:"results"`
}

type jsonResult struct {
	Key    string `json:"key"`
	Status string `json:"status"`
}

// ReportJSON writes groups as a JSON array to w.
func ReportJSON(w io.Writer, groups []Group) error {
	out := make([]jsonGroup, 0, len(groups))
	for _, g := range groups {
		jrs := make([]jsonResult, len(g.Results))
		for i, r := range g.Results {
			jrs[i] = jsonResult{Key: r.Key, Status: r.Status}
		}
		out = append(out, jsonGroup{
			Prefix:  g.Prefix,
			Count:   len(g.Results),
			Results: jrs,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func statusLabel(status string) string {
	switch strings.ToLower(status) {
	case "ok":
		return "✓"
	case "missing_in_right":
		return "✗ missing in right"
	case "missing_in_left":
		return "✗ missing in left"
	case "mismatched":
		return "≠ mismatched"
	default:
		return status
	}
}
