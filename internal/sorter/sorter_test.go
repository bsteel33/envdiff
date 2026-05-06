package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", Status: diff.Equal, Left: "1", Right: "1"},
		{Key: "ALPHA", Status: diff.Mismatched, Left: "old", Right: "new"},
		{Key: "MANGO", Status: diff.Missing, Left: "v", Right: ""},
		{Key: "BETA", Status: diff.Extra, Left: "", Right: "x"},
		{Key: "APPLE", Status: diff.Missing, Left: "y", Right: ""},
	}
}

func TestSort_ByKey(t *testing.T) {
	results := makeResults()
	sorted := sorter.Sort(results, sorter.ByKey)

	expected := []string{"ALPHA", "APPLE", "BETA", "MANGO", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_ByStatus(t *testing.T) {
	results := makeResults()
	sorted := sorter.Sort(results, sorter.ByStatus)

	// Missing should come before Extra, before Mismatched, before Equal
	statusSeq := []diff.Status{
		diff.Missing, diff.Missing, diff.Extra, diff.Mismatched, diff.Equal,
	}
	for i, r := range sorted {
		if r.Status != statusSeq[i] {
			t.Errorf("index %d: got status %v, want %v", i, r.Status, statusSeq[i])
		}
	}
}

func TestSort_ByStatusThenKey(t *testing.T) {
	results := makeResults()
	sorted := sorter.Sort(results, sorter.ByStatusThenKey)

	// Missing: APPLE, MANGO | Extra: BETA | Mismatched: ALPHA | Equal: ZEBRA
	expected := []string{"APPLE", "MANGO", "BETA", "ALPHA", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got key %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	results := makeResults()
	originalFirst := results[0].Key

	sorter.Sort(results, sorter.ByKey)

	if results[0].Key != originalFirst {
		t.Errorf("original slice was mutated: first key changed from %q to %q", originalFirst, results[0].Key)
	}
}

func TestSort_Empty(t *testing.T) {
	result := sorter.Sort([]diff.Result{}, sorter.ByKey)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(result))
	}
}
