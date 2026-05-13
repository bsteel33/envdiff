package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Left: left, Right: right, Status: status}
}

func TestCompute_Empty(t *testing.T) {
	s := scorer.Compute(nil, scorer.DefaultWeights())
	if s.Total != 0 {
		t.Fatalf("expected 0, got %v", s.Total)
	}
	if s.Grade() != "A" {
		t.Fatalf("expected grade A, got %s", s.Grade())
	}
}

func TestCompute_OnlyMismatched(t *testing.T) {
	results := []diff.Result{
		makeResult("KEY", "a", "b", diff.Mismatched),
		makeResult("KEY2", "x", "y", diff.Mismatched),
	}
	s := scorer.Compute(results, scorer.DefaultWeights())
	if s.Mismatched != 2 {
		t.Fatalf("expected 2 mismatched, got %d", s.Mismatched)
	}
	if s.Total != 4.0 {
		t.Fatalf("expected total 4.0, got %v", s.Total)
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := []diff.Result{
		makeResult("A", "", "v", diff.MissingInLeft),
		makeResult("B", "v", "", diff.MissingInRight),
		makeResult("C", "x", "y", diff.Mismatched),
	}
	s := scorer.Compute(results, scorer.DefaultWeights())
	if s.MissingInLeft != 1 || s.MissingInRight != 1 || s.Mismatched != 1 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	// 1*1 + 1*1 + 1*2 = 4
	if s.Total != 4.0 {
		t.Fatalf("expected 4.0, got %v", s.Total)
	}
}

func TestGrade_Boundaries(t *testing.T) {
	cases := []struct {
		total float64
		want  string
	}{
		{0, "A"},
		{3, "B"},
		{8, "C"},
		{15, "D"},
		{16, "F"},
	}
	for _, tc := range cases {
		s := scorer.Score{Total: tc.total}
		if g := s.Grade(); g != tc.want {
			t.Errorf("total=%.0f: want grade %s, got %s", tc.total, tc.want, g)
		}
	}
}

func TestCompute_CustomWeights(t *testing.T) {
	w := scorer.Weights{MissingInLeft: 0, MissingInRight: 0, Mismatched: 5}
	results := []diff.Result{
		makeResult("X", "", "v", diff.MissingInLeft),
		makeResult("Y", "v", "w", diff.Mismatched),
	}
	s := scorer.Compute(results, w)
	if s.Total != 5.0 {
		t.Fatalf("expected 5.0, got %v", s.Total)
	}
}

func TestCompute_OnlyMissingInLeft(t *testing.T) {
	results := []diff.Result{
		makeResult("A", "", "v1", diff.MissingInLeft),
		makeResult("B", "", "v2", diff.MissingInLeft),
		makeResult("C", "", "v3", diff.MissingInLeft),
	}
	s := scorer.Compute(results, scorer.DefaultWeights())
	if s.MissingInLeft != 3 {
		t.Fatalf("expected 3 missing-in-left, got %d", s.MissingInLeft)
	}
	if s.MissingInRight != 0 {
		t.Fatalf("expected 0 missing-in-right, got %d", s.MissingInRight)
	}
	if s.Mismatched != 0 {
		t.Fatalf("expected 0 mismatched, got %d", s.Mismatched)
	}
	// default weight for MissingInLeft is 1, so total = 3*1 = 3
	if s.Total != 3.0 {
		t.Fatalf("expected 3.0, got %v", s.Total)
	}
}
