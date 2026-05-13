package highlighter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/highlighter"
)

func sampleHighlights() []highlighter.Highlight {
	results := []diff.Result{
		makeResult("HOST", "localhost", "prod", diff.StatusMismatch),
		makeResult("PORT", "3000", "", diff.StatusMissingInRight),
	}
	return highlighter.Apply(results, highlighter.Options{Style: highlighter.StylePlain})
}

func TestReportText_NoHighlights(t *testing.T) {
	var buf bytes.Buffer
	highlighter.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got %q", buf.String())
	}
}

func TestReportText_WithHighlights(t *testing.T) {
	var buf bytes.Buffer
	highlighter.ReportText(&buf, sampleHighlights())
	out := buf.String()
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected KEY header or HOST key in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT key in output")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	err := highlighter.ReportJSON(&buf, sampleHighlights())
	if err != nil {
		t.Fatalf("ReportJSON returned error: %v", err)
	}
	var records []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
	for _, r := range records {
		for _, field := range []string{"key", "left", "right", "status"} {
			if _, ok := r[field]; !ok {
				t.Errorf("missing field %q in JSON record", field)
			}
		}
	}
}

func TestReportJSON_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	err := highlighter.ReportJSON(&buf, []highlighter.Highlight{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array in output")
	}
}
