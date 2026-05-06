package ignorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignorer"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch},
		{Key: "SECRET", Status: diff.StatusMissingInRight},
		{Key: "TOKEN", Status: diff.StatusMismatched},
		{Key: "PORT", Status: diff.StatusMissingInLeft},
	}
}

func TestFilterResults_NilIgnorer(t *testing.T) {
	results := makeResults()
	got := ignorer.FilterResults(results, nil)
	if len(got) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(got))
	}
}

func TestFilterResults_RemovesIgnoredKeys(t *testing.T) {
	ig := ignorer.New([]string{"SECRET", "TOKEN"})
	got := ignorer.FilterResults(makeResults(), ig)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Key == "SECRET" || r.Key == "TOKEN" {
			t.Errorf("key %s should have been filtered out", r.Key)
		}
	}
}

func TestFilterResults_DoesNotMutateOriginal(t *testing.T) {
	original := makeResults()
	ig := ignorer.New([]string{"DB_HOST"})
	_ = ignorer.FilterResults(original, ig)
	if len(original) != 4 {
		t.Error("original slice was mutated")
	}
}

func TestFilterResults_EmptyIgnorer(t *testing.T) {
	ig := ignorer.New(nil)
	got := ignorer.FilterResults(makeResults(), ig)
	if len(got) != 4 {
		t.Errorf("expected 4 results with empty ignorer, got %d", len(got))
	}
}
