package comparer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparer"
)

func sampleReport() comparer.Report {
	return comparer.Compare([]comparer.EnvMap{
		makeEnv("dev", map[string]string{"PORT": "3000", "SECRET": "abc"}),
		makeEnv("prod", map[string]string{"PORT": "443"}),
	})
}

func TestReportText_NoKeys(t *testing.T) {
	var buf bytes.Buffer
	comparer.ReportText(&buf, comparer.Report{})
	if !strings.Contains(buf.String(), "No keys found") {
		t.Errorf("expected 'No keys found', got: %s", buf.String())
	}
}

func TestReportText_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	comparer.ReportText(&buf, sampleReport())
	out := buf.String()
	if !strings.Contains(out, "KEY") {
		t.Error("expected KEY column header")
	}
	if !strings.Contains(out, "CONSISTENT") {
		t.Error("expected CONSISTENT column header")
	}
}

func TestReportText_ShowsMissingPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	comparer.ReportText(&buf, sampleReport())
	if !strings.Contains(buf.String(), "<missing>") {
		t.Error("expected <missing> placeholder for absent key")
	}
}

func TestReportText_MarksInconsistentKeys(t *testing.T) {
	r := comparer.Compare([]comparer.EnvMap{
		makeEnv("a", map[string]string{"X": "1"}),
		makeEnv("b", map[string]string{"X": "2"}),
	})
	var buf bytes.Buffer
	comparer.ReportText(&buf, r)
	if !strings.Contains(buf.String(), "NO") {
		t.Error("expected inconsistent key to be marked NO")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := comparer.ReportJSON(&buf, sampleReport()); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected at least one entry in JSON output")
	}
	for _, entry := range out {
		if _, ok := entry["key"]; !ok {
			t.Error("expected 'key' field in JSON entry")
		}
		if _, ok := entry["consistent"]; !ok {
			t.Error("expected 'consistent' field in JSON entry")
		}
		if _, ok := entry["values"]; !ok {
			t.Error("expected 'values' field in JSON entry")
		}
	}
}

func TestReportJSON_MissingInArrayNonNil(t *testing.T) {
	r := comparer.Compare([]comparer.EnvMap{
		makeEnv("only", map[string]string{"K": "v"}),
	})
	var buf bytes.Buffer
	_ = comparer.ReportJSON(&buf, r)
	var out []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)
	missingIn, ok := out[0]["missing_in"].([]interface{})
	if !ok {
		t.Fatal("missing_in should be an array")
	}
	if missingIn == nil {
		t.Error("missing_in should not be null")
	}
}
