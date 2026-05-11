package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func makeResult(key, status string) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestByPrefix_Empty(t *testing.T) {
	groups := grouper.ByPrefix(nil, "_")
	if len(groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(groups))
	}
}

func TestByPrefix_SingleGroup(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", "missing_in_right"),
		makeResult("DB_PORT", "ok"),
		makeResult("DB_NAME", "mismatched"),
	}
	groups := grouper.ByPrefix(results, "_")
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Prefix != "DB" {
		t.Errorf("expected prefix DB, got %q", groups[0].Prefix)
	}
	if len(groups[0].Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(groups[0].Results))
	}
}

func TestByPrefix_MultipleGroups(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", "ok"),
		makeResult("AWS_KEY", "mismatched"),
		makeResult("AWS_SECRET", "missing_in_right"),
		makeResult("PORT", "ok"),
	}
	groups := grouper.ByPrefix(results, "_")
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	// groups are sorted alphabetically by prefix
	if groups[0].Prefix != "AWS" {
		t.Errorf("expected first group AWS, got %q", groups[0].Prefix)
	}
	if groups[1].Prefix != "DB" {
		t.Errorf("expected second group DB, got %q", groups[1].Prefix)
	}
	if groups[2].Prefix != "PORT" {
		t.Errorf("expected third group PORT, got %q", groups[2].Prefix)
	}
}

func TestByPrefix_NoSeparatorInKey(t *testing.T) {
	results := []diff.Result{
		makeResult("PORT", "ok"),
		makeResult("HOST", "mismatched"),
	}
	groups := grouper.ByPrefix(results, "_")
	// Each key has no underscore, so each becomes its own group equal to the full key.
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestByPrefix_EmptySeparator(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", "ok"),
		makeResult("AWS_KEY", "ok"),
	}
	groups := grouper.ByPrefix(results, "")
	if len(groups) != 1 {
		t.Fatalf("expected 1 group (empty prefix), got %d", len(groups))
	}
	if groups[0].Prefix != "" {
		t.Errorf("expected empty prefix, got %q", groups[0].Prefix)
	}
}
