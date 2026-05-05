package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestReportText_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}
	Report(&buf, result, "dev", "prod", FormatText)
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no differences message, got: %s", buf.String())
	}
}

func TestReportText_MissingInRight(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInRight: []string{"SECRET_KEY", "DB_PASS"},
	}
	Report(&buf, result, "dev", "prod", FormatText)
	out := buf.String()
	if !strings.Contains(out, "Missing in prod") {
		t.Errorf("expected 'Missing in prod', got: %s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", out)
	}
}

func TestReportText_MissingInLeft(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInLeft: []string{"NEW_FEATURE_FLAG"},
	}
	Report(&buf, result, "dev", "prod", FormatText)
	out := buf.String()
	if !strings.Contains(out, "Missing in dev") {
		t.Errorf("expected 'Missing in dev', got: %s", out)
	}
	if !strings.Contains(out, "NEW_FEATURE_FLAG") {
		t.Errorf("expected NEW_FEATURE_FLAG in output, got: %s", out)
	}
}

func TestReportText_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Mismatched: []diff.Mismatch{
			{Key: "LOG_LEVEL", LeftValue: "debug", RightValue: "error"},
		},
	}
	Report(&buf, result, "dev", "prod", FormatText)
	out := buf.String()
	if !strings.Contains(out, "Mismatched values") {
		t.Errorf("expected 'Mismatched values', got: %s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected LOG_LEVEL in output, got: %s", out)
	}
	if !strings.Contains(out, "debug") || !strings.Contains(out, "error") {
		t.Errorf("expected both values in output, got: %s", out)
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInRight: []string{"API_KEY"},
		MissingInLeft:  []string{"NEW_VAR"},
		Mismatched: []diff.Mismatch{
			{Key: "HOST", LeftValue: "localhost", RightValue: "prod.example.com"},
		},
	}
	Report(&buf, result, "dev", "prod", FormatJSON)
	out := buf.String()
	if !strings.Contains(out, "missing_in_prod") {
		t.Errorf("expected 'missing_in_prod' key, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in JSON output, got: %s", out)
	}
	if !strings.Contains(out, "mismatched") {
		t.Errorf("expected 'mismatched' key in JSON output, got: %s", out)
	}
}
