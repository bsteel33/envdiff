package highlighter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/highlighter"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Left: left, Right: right, Status: status}
}

func TestApply_MatchedKeyUnchanged(t *testing.T) {
	results := []diff.Result{makeResult("PORT", "8080", "8080", diff.StatusMatch)}
	highlights := highlighter.Apply(results, highlighter.Options{Style: highlighter.StylePlain})
	if len(highlights) != 1 {
		t.Fatalf("expected 1 highlight, got %d", len(highlights))
	}
	if highlights[0].Left != "8080" || highlights[0].Right != "8080" {
		t.Errorf("matched values should be unchanged")
	}
}

func TestApply_MissingInLeft_PlainStyle(t *testing.T) {
	results := []diff.Result{makeResult("NEW_KEY", "", "value", diff.StatusMissingInLeft)}
	highlights := highlighter.Apply(results, highlighter.Options{Style: highlighter.StylePlain})
	if !strings.Contains(highlights[0].Left, "missing") {
		t.Errorf("left side should indicate missing, got %q", highlights[0].Left)
	}
	if !strings.Contains(highlights[0].Right, "+value") {
		t.Errorf("right side should show added value, got %q", highlights[0].Right)
	}
}

func TestApply_MissingInRight_PlainStyle(t *testing.T) {
	results := []diff.Result{makeResult("OLD_KEY", "val", "", diff.StatusMissingInRight)}
	highlights := highlighter.Apply(results, highlighter.Options{Style: highlighter.StylePlain})
	if !strings.Contains(highlights[0].Left, "-val") {
		t.Errorf("left side should show removed value, got %q", highlights[0].Left)
	}
	if !strings.Contains(highlights[0].Right, "missing") {
		t.Errorf("right side should indicate missing, got %q", highlights[0].Right)
	}
}

func TestApply_Mismatch_ANSIStyle(t *testing.T) {
	results := []diff.Result{makeResult("DB_HOST", "localhost", "prod-db", diff.StatusMismatch)}
	highlights := highlighter.Apply(results, highlighter.Options{Style: highlighter.StyleANSI})
	// ANSI escape codes should be present
	if !strings.Contains(highlights[0].Left, "\033[") {
		t.Errorf("expected ANSI escape in left value, got %q", highlights[0].Left)
	}
	if !strings.Contains(highlights[0].Right, "\033[") {
		t.Errorf("expected ANSI escape in right value, got %q", highlights[0].Right)
	}
}

func TestApply_Mismatch_MarkdownStyle(t *testing.T) {
	results := []diff.Result{makeResult("SECRET", "abc", "xyz", diff.StatusMismatch)}
	highlights := highlighter.Apply(results, highlighter.Options{Style: highlighter.StyleMarkdown})
	if !strings.HasPrefix(highlights[0].Left, "**") {
		t.Errorf("expected markdown bold in left, got %q", highlights[0].Left)
	}
}

func TestApply_EmptyResults(t *testing.T) {
	highlights := highlighter.Apply(nil, highlighter.Options{Style: highlighter.StylePlain})
	if len(highlights) != 0 {
		t.Errorf("expected empty highlights for nil input")
	}
}

func TestApply_StatusPreserved(t *testing.T) {
	results := []diff.Result{makeResult("X", "a", "b", diff.StatusMismatch)}
	highlights := highlighter.Apply(results, highlighter.Options{Style: highlighter.StylePlain})
	if highlights[0].Status != diff.StatusMismatch {
		t.Errorf("expected status to be preserved")
	}
	if highlights[0].Key != "X" {
		t.Errorf("expected key to be preserved")
	}
}
