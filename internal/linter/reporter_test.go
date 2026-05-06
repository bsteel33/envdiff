package linter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestReportText_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	linter.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No lint issues") {
		t.Errorf("expected no-issues message, got: %s", buf.String())
	}
}

func TestReportText_WithIssues(t *testing.T) {
	issues := []linter.Issue{
		{Key: "PORT", Message: "value is empty", Severity: "warn"},
		{Key: "MY KEY", Message: "key must not contain spaces", Severity: "error"},
	}
	var buf bytes.Buffer
	linter.ReportText(&buf, issues)
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
	if !strings.Contains(out, "2 issue(s)") {
		t.Errorf("expected issue count in output, got: %s", out)
	}
}

func TestReportJSON_Structure(t *testing.T) {
	issues := []linter.Issue{
		{Key: "SECRET", Message: "value is empty", Severity: "warn"},
	}
	var buf bytes.Buffer
	if err := linter.ReportJSON(&buf, issues); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0]["key"] != "SECRET" {
		t.Errorf("expected key=SECRET, got %s", out[0]["key"])
	}
	if out[0]["severity"] != "warn" {
		t.Errorf("expected severity=warn, got %s", out[0]["severity"])
	}
}

func TestReportJSON_Sorted(t *testing.T) {
	issues := []linter.Issue{
		{Key: "Z_KEY", Message: "m", Severity: "warn"},
		{Key: "A_KEY", Message: "m", Severity: "warn"},
	}
	var buf bytes.Buffer
	_ = linter.ReportJSON(&buf, issues)
	var out []map[string]string
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out[0]["key"] != "A_KEY" {
		t.Errorf("expected A_KEY first, got %s", out[0]["key"])
	}
}
