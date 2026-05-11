package grouper_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func sampleGroups() []grouper.Group {
	return []grouper.Group{
		{
			Prefix: "DB",
			Results: []diff.Result{
				{Key: "DB_HOST", Status: "ok"},
				{Key: "DB_PORT", Status: "mismatched"},
			},
		},
		{
			Prefix: "AWS",
			Results: []diff.Result{
				{Key: "AWS_KEY", Status: "missing_in_right"},
			},
		},
	}
}

func TestReportText_NoGroups(t *testing.T) {
	var buf bytes.Buffer
	grouper.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No groups") {
		t.Errorf("expected 'No groups' message, got: %s", buf.String())
	}
}

func TestReportText_WithGroups(t *testing.T) {
	var buf bytes.Buffer
	grouper.ReportText(&buf, sampleGroups())
	out := buf.String()

	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] header in output")
	}
	if !strings.Contains(out, "[AWS]") {
		t.Errorf("expected [AWS] header in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if !strings.Contains(out, "missing in right") {
		t.Errorf("expected 'missing in right' label in output")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := grouper.ReportJSON(&buf, sampleGroups()); err != nil {
		t.Fatalf("ReportJSON returned error: %v", err)
	}

	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 groups in JSON, got %d", len(out))
	}
	for _, g := range out {
		if _, ok := g["prefix"]; !ok {
			t.Error("missing 'prefix' field in JSON group")
		}
		if _, ok := g["count"]; !ok {
			t.Error("missing 'count' field in JSON group")
		}
		if _, ok := g["results"]; !ok {
			t.Error("missing 'results' field in JSON group")
		}
	}
}
