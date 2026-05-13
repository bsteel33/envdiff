package thresholder_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/thresholder"
)

func makeResult(key string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestEvaluate_NilOptions_AlwaysPasses(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.MissingInLeft),
		makeResult("B", diff.Mismatched),
	}
	out := thresholder.Evaluate(results, nil)
	if !out.Passed {
		t.Fatal("expected Passed=true when opts is nil")
	}
	if len(out.Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(out.Violations))
	}
}

func TestEvaluate_NoViolations(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.MissingInRight),
		makeResult("B", diff.Matched),
	}
	opts := &thresholder.Options{MaxMissing: 5, MaxMismatched: 5, MaxTotal: 10}
	out := thresholder.Evaluate(results, opts)
	if !out.Passed {
		t.Fatalf("expected Passed=true, got violations: %v", out.Violations)
	}
}

func TestEvaluate_MissingExceedsLimit(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.MissingInLeft),
		makeResult("B", diff.MissingInRight),
		makeResult("C", diff.MissingInLeft),
	}
	opts := &thresholder.Options{MaxMissing: 2}
	out := thresholder.Evaluate(results, opts)
	if out.Passed {
		t.Fatal("expected Passed=false")
	}
	if len(out.Violations) != 1 || out.Violations[0].Field != "missing" {
		t.Fatalf("expected one missing violation, got %v", out.Violations)
	}
	if out.Violations[0].Actual != 3 {
		t.Fatalf("expected Actual=3, got %d", out.Violations[0].Actual)
	}
}

func TestEvaluate_MismatchedExceedsLimit(t *testing.T) {
	results := []diff.Result{
		makeResult("X", diff.Mismatched),
		makeResult("Y", diff.Mismatched),
	}
	opts := &thresholder.Options{MaxMismatched: 1}
	out := thresholder.Evaluate(results, opts)
	if out.Passed {
		t.Fatal("expected Passed=false")
	}
	if len(out.Violations) != 1 || out.Violations[0].Field != "mismatched" {
		t.Fatalf("expected one mismatched violation, got %v", out.Violations)
	}
}

func TestEvaluate_TotalExceedsLimit(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.MissingInLeft),
		makeResult("B", diff.Mismatched),
		makeResult("C", diff.MissingInRight),
	}
	opts := &thresholder.Options{MaxTotal: 2}
	out := thresholder.Evaluate(results, opts)
	if out.Passed {
		t.Fatal("expected Passed=false")
	}
	if len(out.Violations) != 1 || out.Violations[0].Field != "total" {
		t.Fatalf("expected one total violation, got %v", out.Violations)
	}
	if out.Violations[0].Actual != 3 {
		t.Fatalf("expected Actual=3, got %d", out.Violations[0].Actual)
	}
}

func TestEvaluate_MultipleViolations(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.MissingInLeft),
		makeResult("B", diff.MissingInRight),
		makeResult("C", diff.Mismatched),
		makeResult("D", diff.Mismatched),
	}
	opts := &thresholder.Options{MaxMissing: 1, MaxMismatched: 1, MaxTotal: 3}
	out := thresholder.Evaluate(results, opts)
	if out.Passed {
		t.Fatal("expected Passed=false")
	}
	if len(out.Violations) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(out.Violations))
	}
}
