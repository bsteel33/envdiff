package highlighter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable highlighted diff table to w.
func ReportText(w io.Writer, highlights []Highlight) {
	if len(highlights) == 0 {
		fmt.Fprintln(w, "No differences to highlight.")
		return
	}
	fmt.Fprintf(w, "%-30s %-30s %s\n", "KEY", "LEFT", "RIGHT")
	fmt.Fprintln(w, strings.Repeat("-", 92))
	for _, h := range highlights {
		fmt.Fprintf(w, "%-30s %-30s %s\n", h.Key, h.Left, h.Right)
	}
}

type jsonHighlight struct {
	Key    string `json:"key"`
	Left   string `json:"left"`
	Right  string `json:"right"`
	Status string `json:"status"`
}

// ReportJSON writes highlights as a JSON array to w.
func ReportJSON(w io.Writer, highlights []Highlight) error {
	records := make([]jsonHighlight, len(highlights))
	for i, h := range highlights {
		records[i] = jsonHighlight{
			Key:    h.Key,
			Left:   h.Left,
			Right:  h.Right,
			Status: string(h.Status),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
