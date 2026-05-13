package timeline

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// ReportText writes a human-readable timeline trend table to w.
func ReportText(t *Timeline, w io.Writer) {
	points := t.Trend()
	if len(points) == 0 {
		fmt.Fprintln(w, "No timeline entries.")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tLABEL\tTOTAL\tMISSING\tMISMATCHED")
	for _, p := range points {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%d\t%d\n",
			p.Timestamp.Format("2006-01-02T15:04:05Z"),
			p.Label,
			p.Total,
			p.Missing,
			p.Mismatched,
		)
	}
	tw.Flush()
}

// ReportJSON writes the timeline trend as a JSON array to w.
func ReportJSON(t *Timeline, w io.Writer) error {
	type jsonPoint struct {
		Timestamp  string `json:"timestamp"`
		Label      string `json:"label"`
		Total      int    `json:"total"`
		Missing    int    `json:"missing"`
		Mismatched int    `json:"mismatched"`
	}
	points := t.Trend()
	out := make([]jsonPoint, 0, len(points))
	for _, p := range points {
		out = append(out, jsonPoint{
			Timestamp:  p.Timestamp.Format("2006-01-02T15:04:05Z"),
			Label:      p.Label,
			Total:      p.Total,
			Missing:    p.Missing,
			Mismatched: p.Mismatched,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
