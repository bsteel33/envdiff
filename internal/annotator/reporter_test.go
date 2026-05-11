package annotator_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envdiff/internal/annotator"
	"github.com/yourusername/envdiff/internal/diff"
)

func sampleAnnotations() []annotator.Annotation {
	return annotator.Annotate([]diff.Result{
		makeResult("DB_PASS", diff.StatusMissingInRight, "hunter2", ""),
		makeResult("APP_ENV", diff.StatusMismatch, "dev", "prod"),
	})
}

func TestReportText_NoAnnotations(t *testing.T) {
	var buf bytes.Buffer
	annotator.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No annotations") {
		t.Errorf("expected no-annotations message, got: %s", buf.String())
	}
}

func TestReportText_WithAnnotations(t *testing.T) {
	var buf bytes.Buffer
	annotator.ReportText(&buf, sampleAnnotations())
	out := buf.String()
	if !strings.Contains(out, "DB_PASS") {
		t.Error("expected DB_PASS in output")
	}
	if !strings.Contains(out, "Suggestion") {
		t.Error("expected Suggestion label in output")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	err := annotator.ReportJSON(&buf, sampleAnnotations())
	if err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
	if out[0]["Key"] == "" {
		t.Error("expected Key field in JSON output")
	}
}
