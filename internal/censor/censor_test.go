package censor_test

import (
	"testing"

	"github.com/user/envdiff/internal/censor"
	"github.com/user/envdiff/internal/diff"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftValue: left, RightValue: right, Status: status}
}

func TestApply_NonSensitiveKeysUnchanged(t *testing.T) {
	input := []diff.Result{
		makeResult("APP_ENV", "production", "staging", diff.Mismatched),
		makeResult("PORT", "8080", "8080", diff.Matched),
	}
	out := censor.Apply(input, nil)
	if out[0].LeftValue != "production" || out[0].RightValue != "staging" {
		t.Errorf("expected non-sensitive values unchanged, got %+v", out[0])
	}
}

func TestApply_SensitiveKeyMasked(t *testing.T) {
	input := []diff.Result{
		makeResult("DB_PASSWORD", "hunter2", "s3cr3t", diff.Mismatched),
	}
	out := censor.Apply(input, nil)
	if out[0].LeftValue != "***" || out[0].RightValue != "***" {
		t.Errorf("expected masked values, got left=%q right=%q", out[0].LeftValue, out[0].RightValue)
	}
}

func TestApply_CustomMask(t *testing.T) {
	input := []diff.Result{
		makeResult("API_KEY", "abc123", "", diff.MissingInRight),
	}
	opts := &censor.Options{Mask: "<REDACTED>"}
	out := censor.Apply(input, opts)
	if out[0].LeftValue != "<REDACTED>" {
		t.Errorf("expected custom mask, got %q", out[0].LeftValue)
	}
	// empty right value should remain empty
	if out[0].RightValue != "" {
		t.Errorf("expected empty right value unchanged, got %q", out[0].RightValue)
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	input := []diff.Result{
		makeResult("STRIPE_TOKEN", "tok_live", "tok_test", diff.Mismatched),
		makeResult("INTERNAL_FLAG", "true", "false", diff.Mismatched),
	}
	opts := &censor.Options{SensitiveSubstrings: []string{"flag"}}
	out := censor.Apply(input, opts)
	// STRIPE_TOKEN should NOT be masked (pattern doesn't include "token")
	if out[0].LeftValue != "tok_live" {
		t.Errorf("expected unchanged, got %q", out[0].LeftValue)
	}
	// INTERNAL_FLAG should be masked
	if out[1].LeftValue != "***" {
		t.Errorf("expected masked, got %q", out[1].LeftValue)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := []diff.Result{
		makeResult("SECRET_KEY", "original", "original", diff.Matched),
	}
	censor.Apply(input, nil)
	if input[0].LeftValue != "original" {
		t.Error("Apply mutated the original slice")
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := censor.Apply(nil, nil)
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d results", len(out))
	}
}
