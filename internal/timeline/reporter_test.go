package timeline

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
)

func sampleTimeline() *Timeline {
	tl := &Timeline{}
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tl.Add("prod-v1", base, []diff.Result{
		{Key: "DB_URL", Status: diff.MissingInRight},
		{Key: "API_KEY", Status: diff.Mismatched},
	})
	tl.Add("prod-v2", base.Add(24*time.Hour), []diff.Result{
		{Key: "API_KEY", Status: diff.Equal},
	})
	return tl
}

func TestReportText_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	ReportText(&Timeline{}, &buf)
	if !strings.Contains(buf.String(), "No timeline entries") {
		t.Errorf("expected no-entries message, got: %s", buf.String())
	}
}

func TestReportText_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	ReportText(sampleTimeline(), &buf)
	out := buf.String()
	if !strings.Contains(out, "prod-v1") {
		t.Errorf("expected label prod-v1 in output")
	}
	if !strings.Contains(out, "prod-v2") {
		t.Errorf("expected label prod-v2 in output")
	}
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING header in output")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := ReportJSON(sampleTimeline(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, field := range []string{"timestamp", "label", "total", "missing", "mismatched"} {
		if _, ok := result[0][field]; !ok {
			t.Errorf("missing field %q in JSON output", field)
		}
	}
	if result[0]["missing"].(float64) != 1 {
		t.Errorf("expected missing=1 for first entry")
	}
}
