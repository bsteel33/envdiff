package comparer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable multi-environment comparison table to w.
func ReportText(w io.Writer, r Report) {
	if len(r.Keys) == 0 {
		fmt.Fprintln(w, "No keys found.")
		return
	}

	// Header
	fmt.Fprintf(w, "%-30s  %-12s  %s\n", "KEY", "CONSISTENT", strings.Join(r.Labels, "  "))
	fmt.Fprintln(w, strings.Repeat("-", 72))

	for _, kr := range r.Keys {
		consistent := "yes"
		if !kr.Consistent {
			consistent = "NO"
		}
		vals := make([]string, len(r.Labels))
		for i, lbl := range r.Labels {
			v, ok := kr.Values[lbl]
			if !ok {
				v = "<missing>"
			}
			vals[i] = v
		}
		fmt.Fprintf(w, "%-30s  %-12s  %s\n", kr.Key, consistent, strings.Join(vals, "  "))
	}
}

// jsonKeyReport is the JSON-serialisable form of a KeyReport.
type jsonKeyReport struct {
	Key        string            `json:"key"`
	Consistent bool              `json:"consistent"`
	PresentIn  []string          `json:"present_in"`
	MissingIn  []string          `json:"missing_in"`
	Values     map[string]string `json:"values"`
}

// ReportJSON writes a JSON array of key reports to w.
func ReportJSON(w io.Writer, r Report) error {
	out := make([]jsonKeyReport, len(r.Keys))
	for i, kr := range r.Keys {
		present := kr.PresentIn
		if present == nil {
			present = []string{}
		}
		missing := kr.MissingIn
		if missing == nil {
			missing = []string{}
		}
		out[i] = jsonKeyReport{
			Key:        kr.Key,
			Consistent: kr.Consistent,
			PresentIn:  present,
			MissingIn:  missing,
			Values:     kr.Values,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
