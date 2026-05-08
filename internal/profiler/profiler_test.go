package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/profiler"
)

func makeResults(pairs ...struct {
	key, left, right string
	status           diff.Status
}) []diff.Result {
	out := make([]diff.Result, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, diff.Result{
			Key:        p.key,
			LeftValue:  p.left,
			RightValue: p.right,
			Status:     p.status,
		})
	}
	return out
}

func TestAnalyse_Empty(t *testing.T) {
	p := profiler.Analyse()
	if p.TotalKeys != 0 {
		t.Fatalf("expected 0 total keys, got %d", p.TotalKeys)
	}
	if len(p.TopDiffering) != 0 {
		t.Fatalf("expected empty top differing, got %v", p.TopDiffering)
	}
}

func TestAnalyse_SingleSet(t *testing.T) {
	results := makeResults(
		struct{ key, left, right string; status diff.Status }{"A", "1", "1", diff.StatusEqual},
		struct{ key, left, right string; status diff.Status }{"B", "", "2", diff.StatusMissingInLeft},
		struct{ key, left, right string; status diff.Status }{"C", "3", "", diff.StatusMissingInRight},
		struct{ key, left, right string; status diff.Status }{"D", "4", "5", diff.StatusMismatch},
	)
	p := profiler.Analyse(results)

	if p.TotalKeys != 4 {
		t.Errorf("TotalKeys: want 4, got %d", p.TotalKeys)
	}
	if p.Identical != 1 {
		t.Errorf("Identical: want 1, got %d", p.Identical)
	}
	if p.MissingLeft != 1 {
		t.Errorf("MissingLeft: want 1, got %d", p.MissingLeft)
	}
	if p.MissingRight != 1 {
		t.Errorf("MissingRight: want 1, got %d", p.MissingRight)
	}
	if p.Mismatched != 1 {
		t.Errorf("Mismatched: want 1, got %d", p.Mismatched)
	}
}

func TestAnalyse_MultipleSets_FrequencyAccumulates(t *testing.T) {
	set1 := makeResults(
		struct{ key, left, right string; status diff.Status }{"X", "a", "b", diff.StatusMismatch},
	)
	set2 := makeResults(
		struct{ key, left, right string; status diff.Status }{"X", "a", "b", diff.StatusMismatch},
		struct{ key, left, right string; status diff.Status }{"Y", "", "c", diff.StatusMissingInLeft},
	)
	p := profiler.Analyse(set1, set2)

	if p.TotalKeys != 3 {
		t.Errorf("TotalKeys: want 3, got %d", p.TotalKeys)
	}
	if len(p.TopDiffering) == 0 {
		t.Fatal("expected non-empty TopDiffering")
	}
	if p.TopDiffering[0].Key != "X" {
		t.Errorf("expected top key X, got %s", p.TopDiffering[0].Key)
	}
	if p.TopDiffering[0].Count != 2 {
		t.Errorf("expected count 2 for X, got %d", p.TopDiffering[0].Count)
	}
}

func TestAnalyse_TopDiffering_CappedAtFive(t *testing.T) {
	var results []diff.Result
	for i := 0; i < 10; i++ {
		results = append(results, diff.Result{
			Key:    string(rune('A' + i)),
			Status: diff.StatusMismatch,
		})
	}
	p := profiler.Analyse(results)
	if len(p.TopDiffering) > 5 {
		t.Errorf("TopDiffering should be capped at 5, got %d", len(p.TopDiffering))
	}
}
