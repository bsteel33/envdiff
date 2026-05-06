package validator_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestReportViolations_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	n := validator.ReportViolations(&buf, nil)
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
	if !strings.Contains(buf.String(), "No validation violations") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestReportViolations_WithViolations(t *testing.T) {
	violations := []validator.Violation{
		{Key: "SECRET", Rule: "no-empty-value", Message: `key "SECRET" has an empty value`},
		{Key: "PORT", Rule: "digits-only", Message: `key "PORT" value "abc" does not match rule "digits-only"`},
	}
	var buf bytes.Buffer
	n := validator.ReportViolations(&buf, violations)
	if n != 2 {
		t.Errorf("expected 2 violations reported, got %d", n)
	}
	out := buf.String()
	if !strings.Contains(out, "no-empty-value") {
		t.Errorf("expected rule name in output, got: %q", out)
	}
	if !strings.Contains(out, "digits-only") {
		t.Errorf("expected rule name in output, got: %q", out)
	}
	if !strings.Contains(out, "2 validation violation") {
		t.Errorf("expected count in output, got: %q", out)
	}
}
