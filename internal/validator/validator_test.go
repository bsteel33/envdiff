package validator_test

import (
	"regexp"
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func mustCompile(t *testing.T, pattern string) *regexp.Regexp {
	t.Helper()
	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("failed to compile pattern %q: %v", pattern, err)
	}
	return re
}

func TestValidate_NoViolations(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	violations := validator.Validate(env, validator.DefaultRules())
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestValidate_EmptyValueViolation(t *testing.T) {
	env := map[string]string{
		"HOST": "",
		"PORT": "8080",
	}
	violations := validator.Validate(env, validator.DefaultRules())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "HOST" {
		t.Errorf("expected violation for HOST, got %q", violations[0].Key)
	}
	if violations[0].Rule != "no-empty-value" {
		t.Errorf("expected rule no-empty-value, got %q", violations[0].Rule)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	rules := []validator.Rule{
		{Name: "digits-only", Pattern: mustCompile(t, `^\d+$`)},
	}
	env := map[string]string{"PORT": "abc"}
	violations := validator.Validate(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "digits-only" {
		t.Errorf("unexpected rule name: %q", violations[0].Rule)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	rules := []validator.Rule{
		{Name: "digits-only", Pattern: mustCompile(t, `^\d+$`)},
	}
	env := map[string]string{"PORT": "8080"}
	violations := validator.Validate(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for matching pattern, got %d", len(violations))
	}
}

func TestValidate_MultipleEmptyValues(t *testing.T) {
	env := map[string]string{
		"HOST": "",
		"PORT": "",
		"NAME": "app",
	}
	violations := validator.Validate(env, validator.DefaultRules())
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestValidateKeys_AllPresent(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	violations := validator.ValidateKeys(env, []string{"A", "B"})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestValidateKeys_MissingKey(t *testing.T) {
	env := map[string]string{"A": "1"}
	violations := validator.ValidateKeys(env, []string{"A", "B"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "B" {
		t.Errorf("expected missing key B, got %q", violations[0].Key)
	}
}

func TestValidateKeys_Empty(t *testing.T) {
	env := map[string]string{}
	violations := validator.ValidateKeys(env, []string{})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}
