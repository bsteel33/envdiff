package deduplicator_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduplicator"
	"github.com/user/envdiff/internal/diff"
)

func makeResult(key string, status diff.Status, left, right string) diff.Result {
	return diff.Result{Key: key, Status: status, LeftValue: left, RightValue: right}
}

func TestApply_Empty(t *testing.T) {
	out := deduplicator.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(out))
	}
}

func TestApply_NoDuplicates(t *testing.T) {
	input := []diff.Result{
		makeResult("AAA", diff.StatusMatch, "x", "x"),
		makeResult("BBB", diff.StatusMismatch, "x", "y"),
	}
	out := deduplicator.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_DuplicateKeepHigherSeverity(t *testing.T) {
	// Same key appears twice: once as Match, once as Mismatch.
	// Mismatch should win.
	input := []diff.Result{
		makeResult("KEY", diff.StatusMatch, "v", "v"),
		makeResult("KEY", diff.StatusMismatch, "v", "w"),
	}
	out := deduplicator.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Status != diff.StatusMismatch {
		t.Errorf("expected Mismatch status, got %v", out[0].Status)
	}
}

func TestApply_DuplicateTieKeepsFirst(t *testing.T) {
	// Two MissingInRight entries for the same key — first wins.
	input := []diff.Result{
		makeResult("KEY", diff.StatusMissingInRight, "first", ""),
		makeResult("KEY", diff.StatusMissingInRight, "second", ""),
	}
	out := deduplicator.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].LeftValue != "first" {
		t.Errorf("expected first entry to be kept, got %q", out[0].LeftValue)
	}
}

func TestApply_OutputSortedByKey(t *testing.T) {
	input := []diff.Result{
		makeResult("ZZZ", diff.StatusMatch, "a", "a"),
		makeResult("AAA", diff.StatusMatch, "b", "b"),
		makeResult("MMM", diff.StatusMismatch, "c", "d"),
	}
	out := deduplicator.Apply(input)
	keys := []string{out[0].Key, out[1].Key, out[2].Key}
	want := []string{"AAA", "MMM", "ZZZ"}
	for i, k := range want {
		if keys[i] != k {
			t.Errorf("position %d: want %q, got %q", i, k, keys[i])
		}
	}
}
