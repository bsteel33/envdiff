package linter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ReportText writes all issues to w in a human-readable format.
func ReportText(w io.Writer, issues []Issue) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No lint issues found.")
		return
	}
	sorted := sortedIssues(issues)
	for _, i := range sorted {
		fmt.Fprintln(w, FormatIssue(i))
	}
	fmt.Fprintf(w, "\n%d issue(s) found.\n", len(issues))
}

// ReportJSON writes all issues to w as a JSON array.
func ReportJSON(w io.Writer, issues []Issue) error {
	sorted := sortedIssues(issues)
	type jsonIssue struct {
		Key      string `json:"key"`
		Message  string `json:"message"`
		Severity string `json:"severity"`
	}
	out := make([]jsonIssue, len(sorted))
	for idx, i := range sorted {
		out[idx] = jsonIssue{Key: i.Key, Message: i.Message, Severity: i.Severity}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// ReportSummary writes a compact summary of issues grouped by severity to w.
func ReportSummary(w io.Writer, issues []Issue) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No lint issues found.")
		return
	}
	counts := make(map[string]int)
	for _, i := range issues {
		counts[i.Severity]++
	}
	fmt.Fprintf(w, "Total: %d issue(s) — ", len(issues))
	for _, sev := range []string{"error", "warning", "info"} {
		if n, ok := counts[sev]; ok {
			fmt.Fprintf(w, "%s: %d  ", sev, n)
		}
	}
	fmt.Fprintln(w)
}

func sortedIssues(issues []Issue) []Issue {
	copy_ := make([]Issue, len(issues))
	copy(copy_, issues)
	sort.Slice(copy_, func(i, j int) bool {
		if copy_[i].Key != copy_[j].Key {
			return copy_[i].Key < copy_[j].Key
		}
		return copy_[i].Severity < copy_[j].Severity
	})
	return copy_
}
