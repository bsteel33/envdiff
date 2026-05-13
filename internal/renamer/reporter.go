package renamer

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ReportText writes a human-readable summary of rename suggestions to w.
func ReportText(w io.Writer, suggestions []Suggestion) {
	if len(suggestions) == 0 {
		fmt.Fprintln(w, "No rename suggestions found.")
		return
	}
	sorted := sortedSuggestions(suggestions)
	fmt.Fprintf(w, "Rename suggestions (%d):\n", len(sorted))
	for _, s := range sorted {
		fmt.Fprintf(w, "  %-30s  →  %-30s  (confidence: %.0f%%)\n",
			s.OldKey, s.NewKey, s.Score*100)
	}
}

// ReportJSON writes rename suggestions as a JSON array to w.
func ReportJSON(w io.Writer, suggestions []Suggestion) error {
	type jsonSuggestion struct {
		OldKey     string  `json:"old_key"`
		NewKey     string  `json:"new_key"`
		OldValue   string  `json:"old_value"`
		NewValue   string  `json:"new_value"`
		Confidence float64 `json:"confidence"`
	}

	out := make([]jsonSuggestion, 0, len(suggestions))
	for _, s := range sortedSuggestions(suggestions) {
		out = append(out, jsonSuggestion{
			OldKey:     s.OldKey,
			NewKey:     s.NewKey,
			OldValue:   s.OldValue,
			NewValue:   s.NewValue,
			Confidence: s.Score,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func sortedSuggestions(suggestions []Suggestion) []Suggestion {
	copy_ := make([]Suggestion, len(suggestions))
	copy(copy_, suggestions)
	sort.Slice(copy_, func(i, j int) bool {
		if copy_[i].Score != copy_[j].Score {
			return copy_[i].Score > copy_[j].Score
		}
		return copy_[i].OldKey < copy_[j].OldKey
	})
	return copy_
}
