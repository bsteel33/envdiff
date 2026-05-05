package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/filter"
)

func TestFromDiffResults_RoundTrip(t *testing.T) {
	original := []diff.Result{
		{Key: "FOO", MissingInRight: true, LeftValue: "bar"},
		{Key: "BAZ", MissingInLeft: true, RightValue: "qux"},
		{Key: "MISMATCH", LeftValue: "a", RightValue: "b"},
	}

	filtered := filter.FromDiffResults(original)
	if len(filtered) != len(original) {
		t.Fatalf("expected %d items after conversion, got %d", len(original), len(filtered))
	}

	for i, r := range filtered {
		if r.Key != original[i].Key {
			t.Errorf("[%d] Key mismatch: want %q got %q", i, original[i].Key, r.Key)
		}
		if r.MissingInLeft != original[i].MissingInLeft {
			t.Errorf("[%d] MissingInLeft mismatch", i)
		}
		if r.MissingInRight != original[i].MissingInRight {
			t.Errorf("[%d] MissingInRight mismatch", i)
		}
		if r.LeftValue != original[i].LeftValue {
			t.Errorf("[%d] LeftValue mismatch", i)
		}
		if r.RightValue != original[i].RightValue {
			t.Errorf("[%d] RightValue mismatch", i)
		}
	}

	// Round-trip back to diff.Result
	back := filter.ToDiffResults(filtered)
	if len(back) != len(original) {
		t.Fatalf("expected %d items after back-conversion, got %d", len(original), len(back))
	}
	for i, r := range back {
		if r.Key != original[i].Key {
			t.Errorf("round-trip [%d] Key mismatch: want %q got %q", i, original[i].Key, r.Key)
		}
	}
}

func TestFromDiffResults_Empty(t *testing.T) {
	out := filter.FromDiffResults([]diff.Result{})
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(out))
	}
}
