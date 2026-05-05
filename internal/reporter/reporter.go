package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable diff report to the given writer.
func Report(w io.Writer, result diff.Result, leftName, rightName string, format Format) {
	switch format {
	case FormatJSON:
		reportJSON(w, result, leftName, rightName)
	default:
		reportText(w, result, leftName, rightName)
	}
}

func reportText(w io.Writer, result diff.Result, leftName, rightName string) {
	if len(result.MissingInRight) == 0 && len(result.MissingInLeft) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	if len(result.MissingInRight) > 0 {
		fmt.Fprintf(w, "Missing in %s:\n", rightName)
		for _, key := range result.MissingInRight {
			fmt.Fprintf(w, "  - %s\n", key)
		}
	}

	if len(result.MissingInLeft) > 0 {
		fmt.Fprintf(w, "Missing in %s:\n", leftName)
		for _, key := range result.MissingInLeft {
			fmt.Fprintf(w, "  + %s\n", key)
		}
	}

	if len(result.Mismatched) > 0 {
		fmt.Fprintln(w, "Mismatched values:")
		for _, m := range result.Mismatched {
			fmt.Fprintf(w, "  ~ %s\n", m.Key)
			fmt.Fprintf(w, "      %s: %s\n", leftName, m.LeftValue)
			fmt.Fprintf(w, "      %s: %s\n", rightName, m.RightValue)
		}
	}
}

func reportJSON(w io.Writer, result diff.Result, leftName, rightName string) {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"missing_in_%s\": [\n", sanitizeJSONKey(rightName)))
	for i, key := range result.MissingInRight {
		comma := ","
		if i == len(result.MissingInRight)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s\"%s\n", key, comma))
	}
	sb.WriteString("  ],\n")
	sb.WriteString(fmt.Sprintf("  \"missing_in_%s\": [\n", sanitizeJSONKey(leftName)))
	for i, key := range result.MissingInLeft {
		comma := ","
		if i == len(result.MissingInLeft)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s\"%s\n", key, comma))
	}
	sb.WriteString("  ],\n")
	sb.WriteString("  \"mismatched\": [\n")
	for i, m := range result.Mismatched {
		comma := ","
		if i == len(result.Mismatched)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    {\"key\": \"%s\", \"%s\": \"%s\", \"%s\": \"%s\"}%s\n",
			m.Key, sanitizeJSONKey(leftName), m.LeftValue, sanitizeJSONKey(rightName), m.RightValue, comma))
	}
	sb.WriteString("  ]\n")
	sb.WriteString("}\n")
	fmt.Fprint(w, sb.String())
}

func sanitizeJSONKey(name string) string {
	return strings.NewReplacer(" ", "_", "/", "_", ".", "_").Replace(name)
}
