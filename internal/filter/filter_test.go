package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/filter"
)

func makeResults() []filter.Result {
	return []filter.Result{
		{Key: "APP_HOST", MissingInRight: true},
		{Key: "APP_PORT", LeftValue: "8080", RightValue: "9090"},
		{Key: "DB_URL", MissingInLeft: true},
		{Key: "SECRET_KEY", LeftValue: "abc", RightValue: "xyz"},
		{Key: "DEBUG", MissingInRight: true},
	}
}

func TestApply_NoOptions(t *testing.T) {
	results := makeResults()
	out := filter.Apply(results, filter.Options{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{Prefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Key != "APP_HOST" && r.Key != "APP_PORT" {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{ExcludeKeys: []string{"SECRET_KEY", "DEBUG"}})
	for _, r := range out {
		if r.Key == "SECRET_KEY" || r.Key == "DEBUG" {
			t.Errorf("key %q should have been excluded", r.Key)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{OnlyMissing: true})
	for _, r := range out {
		if !r.MissingInLeft && !r.MissingInRight {
			t.Errorf("key %q is not missing but was included", r.Key)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 missing results, got %d", len(out))
	}
}

func TestApply_CombinedOptions(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{
		Prefix:      "APP_",
		OnlyMissing: true,
	})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %q", out[0].Key)
	}
}
