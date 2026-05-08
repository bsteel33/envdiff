package auditor_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/diff"
)

func fixedReport(events []auditor.Event) auditor.Report {
	return auditor.Report{
		GeneratedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Events:      events,
	}
}

func resultPtr(key, status, left, right string) *diff.Result {
	r := makeResult(key, status, left, right)
	return &r
}

func TestReportText_NoEvents(t *testing.T) {
	var buf bytes.Buffer
	auditor.ReportText(&buf, fixedReport(nil))
	if !strings.Contains(buf.String(), "No audit events") {
		t.Errorf("expected no-events message, got: %s", buf.String())
	}
}

func TestReportText_WithEvents(t *testing.T) {
	events := []auditor.Event{
		{Key: "DB_HOST", Kind: auditor.EventIntroduced, Current: resultPtr("DB_HOST", "missing_in_left", "", "localhost")},
		{Key: "OLD_KEY", Kind: auditor.EventResolved, Previous: resultPtr("OLD_KEY", "missing_in_right", "v", "")},
	}
	var buf bytes.Buffer
	auditor.ReportText(&buf, fixedReport(events))
	out := buf.String()

	if !strings.Contains(out, "[+]") {
		t.Error("expected introduced marker [+]")
	}
	if !strings.Contains(out, "[-]") {
		t.Error("expected resolved marker [-]")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Error("expected date in header")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	events := []auditor.Event{
		{Key: "API_KEY", Kind: auditor.EventChanged,
			Previous: resultPtr("API_KEY", "missing_in_right", "x", ""),
			Current:  resultPtr("API_KEY", "mismatched", "x", "y")},
	}
	var buf bytes.Buffer
	if err := auditor.ReportJSON(&buf, fixedReport(events)); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["generated_at"]; !ok {
		t.Error("expected generated_at field")
	}
	evs, ok := out["events"].([]interface{})
	if !ok || len(evs) != 1 {
		t.Fatalf("expected 1 event in JSON, got %v", out["events"])
	}
}
