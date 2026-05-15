package drifter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/drifter"
)

func makeResult(key string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestMeasure_Empty(t *testing.T) {
	r := drifter.Measure(nil)
	if r.Total != 0 || r.Drifted != 0 || r.Score != 0 || r.Severity != drifter.SeverityNone {
		t.Errorf("unexpected report for empty input: %+v", r)
	}
}

func TestMeasure_AllMatched(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.StatusMatch),
		makeResult("B", diff.StatusMatch),
	}
	r := drifter.Measure(results)
	if r.Score != 0 {
		t.Errorf("expected score 0, got %v", r.Score)
	}
	if r.Severity != drifter.SeverityNone {
		t.Errorf("expected SeverityNone, got %v", r.Severity)
	}
}

func TestMeasure_LowDrift(t *testing.T) {
	results := make([]diff.Result, 20)
	for i := range results {
		results[i] = makeResult("K", diff.StatusMatch)
	}
	results[0] = makeResult("K", diff.StatusMissingInRight)

	r := drifter.Measure(results)
	if r.Severity != drifter.SeverityLow {
		t.Errorf("expected SeverityLow, got %v", r.Severity)
	}
}

func TestMeasure_ModerateDrift(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.StatusMissingInRight),
		makeResult("B", diff.StatusMissingInRight),
		makeResult("C", diff.StatusMatch),
		makeResult("D", diff.StatusMatch),
		makeResult("E", diff.StatusMatch),
		makeResult("F", diff.StatusMatch),
		makeResult("G", diff.StatusMatch),
		makeResult("H", diff.StatusMatch),
	}
	r := drifter.Measure(results)
	if r.Severity != drifter.SeverityModerate {
		t.Errorf("expected SeverityModerate, got %v (score=%v)", r.Severity, r.Score)
	}
}

func TestMeasure_CriticalDrift(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.StatusMissingInRight),
		makeResult("B", diff.StatusMismatched),
		makeResult("C", diff.StatusMissingInLeft),
	}
	r := drifter.Measure(results)
	if r.Score != 1.0 {
		t.Errorf("expected score 1.0, got %v", r.Score)
	}
	if r.Severity != drifter.SeverityCritical {
		t.Errorf("expected SeverityCritical, got %v", r.Severity)
	}
	if r.Drifted != 3 {
		t.Errorf("expected 3 drifted, got %d", r.Drifted)
	}
}
