package annotator

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable annotation report to w.
func ReportText(w io.Writer, annotations []Annotation) {
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No annotations.")
		return
	}
	for _, a := range annotations {
		if a.Status == "equal" || a.Status == "" {
			continue
		}
		fmt.Fprintf(w, "[%s] %s\n", strings.ToUpper(a.Status), a.Key)
		fmt.Fprintf(w, "  Reason:     %s\n", a.Reason)
		if a.Suggestion != "" {
			fmt.Fprintf(w, "  Suggestion: %s\n", a.Suggestion)
		}
	}
}

// ReportJSON writes annotations as a JSON array to w.
func ReportJSON(w io.Writer, annotations []Annotation) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(annotations)
}
