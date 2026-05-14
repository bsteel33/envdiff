package counter_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/counter"
	"github.com/your-org/envdiff/internal/diff"
)

func makeResult(key string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestCount_Empty(t *testing.T) {
	s := counter.Count(nil)
	if s.Total != 0 {
		t.Fatalf("expected 0 total, got %d", s.Total)
	}
}

func TestCount_AllMatched(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.StatusMatch),
		makeResult("B", diff.StatusMatch),
	}
	s := counter.Count(results)
	if s.Total != 2 || s.Matched != 2 || s.Missing != 0 || s.Mismatched != 0 {
		t.Fatalf("unexpected stats: %+v", s)
	}
}

func TestCount_Mixed(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.StatusMatch),
		makeResult("B", diff.StatusMissingInRight),
		makeResult("C", diff.StatusMissingInLeft),
		makeResult("D", diff.StatusMismatch),
	}
	s := counter.Count(results)
	if s.Total != 4 {
		t.Fatalf("expected total 4, got %d", s.Total)
	}
	if s.Matched != 1 {
		t.Errorf("expected matched 1, got %d", s.Matched)
	}
	if s.Missing != 2 {
		t.Errorf("expected missing 2, got %d", s.Missing)
	}
	if s.MissingLeft != 1 || s.MissingRight != 1 {
		t.Errorf("expected missingLeft/Right 1/1, got %d/%d", s.MissingLeft, s.MissingRight)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected mismatched 1, got %d", s.Mismatched)
	}
}

func TestDriftRatio_ZeroWhenEmpty(t *testing.T) {
	if r := counter.DriftRatio(counter.Stats{}); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestDriftRatio_HalfDrift(t *testing.T) {
	s := counter.Stats{Total: 4, Matched: 2, Mismatched: 2}
	got := counter.DriftRatio(s)
	if got != 0.5 {
		t.Fatalf("expected 0.5, got %f", got)
	}
}

func TestIsClean_True(t *testing.T) {
	s := counter.Stats{Total: 3, Matched: 3}
	if !counter.IsClean(s) {
		t.Fatal("expected clean")
	}
}

func TestIsClean_False(t *testing.T) {
	s := counter.Stats{Total: 3, Matched: 2, Mismatched: 1}
	if counter.IsClean(s) {
		t.Fatal("expected not clean")
	}
}

func TestIsClean_EmptyIsNotClean(t *testing.T) {
	if counter.IsClean(counter.Stats{}) {
		t.Fatal("empty stats should not be clean")
	}
}
