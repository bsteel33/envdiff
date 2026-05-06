package linter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestRuleUppercaseKey_Pass(t *testing.T) {
	issues := linter.RuleUppercaseKey("DATABASE_URL", "postgres://")
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestRuleUppercaseKey_Fail(t *testing.T) {
	issues := linter.RuleUppercaseKey("database_url", "postgres://")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", issues[0].Severity)
	}
}

func TestRuleNoSpaceInKey_Fail(t *testing.T) {
	issues := linter.RuleNoSpaceInKey("MY KEY", "val")
	if len(issues) != 1 || issues[0].Severity != "error" {
		t.Fatalf("expected 1 error issue, got %v", issues)
	}
}

func TestRuleNonEmptyValue_Fail(t *testing.T) {
	issues := linter.RuleNonEmptyValue("PORT", "")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestRuleNonEmptyValue_Pass(t *testing.T) {
	issues := linter.RuleNonEmptyValue("PORT", "8080")
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestRuleNoLeadingUnderscore_Fail(t *testing.T) {
	issues := linter.RuleNoLeadingUnderscore("_SECRET", "abc")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"good_key": "",
		"GOOD_KEY":  "value",
		"_BAD":      "x",
	}
	issues := linter.Lint(env, linter.DefaultRules())
	if len(issues) == 0 {
		t.Fatal("expected at least one issue")
	}
}

func TestFormatIssue(t *testing.T) {
	i := linter.Issue{Key: "PORT", Message: "value is empty", Severity: "warn"}
	out := linter.FormatIssue(i)
	if !strings.Contains(out, "[WARN]") || !strings.Contains(out, "PORT") {
		t.Errorf("unexpected format: %s", out)
	}
}
