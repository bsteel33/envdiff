package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/redactor"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftValue: left, RightValue: right, Status: status}
}

func TestApply_NoSensitiveKeys(t *testing.T) {
	r := redactor.New(nil)
	input := []diff.Result{
		makeResult("APP_NAME", "myapp", "myapp", diff.StatusMatch),
		makeResult("PORT", "8080", "9090", diff.StatusMismatch),
	}
	out := r.Apply(input)
	if out[0].LeftValue != "myapp" {
		t.Errorf("expected myapp, got %s", out[0].LeftValue)
	}
	if out[1].LeftValue != "8080" || out[1].RightValue != "9090" {
		t.Errorf("non-sensitive values should not be redacted")
	}
}

func TestApply_RedactsSensitiveValues(t *testing.T) {
	r := redactor.New(nil)
	input := []diff.Result{
		makeResult("DB_PASSWORD", "hunter2", "s3cr3t", diff.StatusMismatch),
		makeResult("API_TOKEN", "abc123", "", diff.StatusMissingRight),
	}
	out := r.Apply(input)
	if out[0].LeftValue != "***REDACTED***" {
		t.Errorf("expected redacted left value, got %s", out[0].LeftValue)
	}
	if out[0].RightValue != "***REDACTED***" {
		t.Errorf("expected redacted right value, got %s", out[0].RightValue)
	}
	if out[1].LeftValue != "***REDACTED***" {
		t.Errorf("expected redacted left value for token, got %s", out[1].LeftValue)
	}
	if out[1].RightValue != "" {
		t.Errorf("empty right value should stay empty, got %s", out[1].RightValue)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r := redactor.New(nil)
	input := []diff.Result{
		makeResult("SECRET_KEY", "original", "original", diff.StatusMatch),
	}
	_ = r.Apply(input)
	if input[0].LeftValue != "original" {
		t.Error("Apply must not mutate the original slice")
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	r := redactor.New([]string{"internal"})
	input := []diff.Result{
		makeResult("INTERNAL_URL", "http://internal", "http://other", diff.StatusMismatch),
		makeResult("PUBLIC_URL", "http://pub", "http://pub2", diff.StatusMismatch),
	}
	out := r.Apply(input)
	if out[0].LeftValue != "***REDACTED***" {
		t.Errorf("custom pattern should redact INTERNAL_URL")
	}
	if out[1].LeftValue == "***REDACTED***" {
		t.Errorf("PUBLIC_URL should not be redacted with custom pattern")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	r := redactor.New(nil)
	input := []diff.Result{
		makeResult("MYSECRET", "val", "val", diff.StatusMatch),
		makeResult("MY_Auth_Key", "val", "val", diff.StatusMatch),
	}
	out := r.Apply(input)
	for _, res := range out {
		if res.LeftValue != "***REDACTED***" {
			t.Errorf("key %s should be redacted (case-insensitive)", res.Key)
		}
	}
}
