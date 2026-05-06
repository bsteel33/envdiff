package validator

import (
	"fmt"
	"io"
	"strings"
)

// ReportViolations writes a human-readable summary of violations to w.
// Returns the number of violations reported.
func ReportViolations(w io.Writer, violations []Violation) int {
	if len(violations) == 0 {
		fmt.Fprintln(w, "✔  No validation violations found.")
		return 0
	}

	fmt.Fprintf(w, "✖  %d validation violation(s) found:\n", len(violations))
	fmt.Fprintln(w, strings.Repeat("-", 50))

	for _, v := range violations {
		fmt.Fprintf(w, "  [%s] %s\n", v.Rule, v.Message)
	}

	fmt.Fprintln(w, strings.Repeat("-", 50))
	return len(violations)
}
