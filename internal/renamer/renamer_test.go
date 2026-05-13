package renamer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftValue: left, RightValue: right, Status: status}
}

func TestDetect_NoResults(t *testing.T) {
	got := Detect(nil, 0.8)
	if len(got) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(got))
	}
}

func TestDetect_IdenticalValues(t *testing.T) {
	results := []diff.Result{
		makeResult("OLD_DB_HOST", "localhost", "", diff.MissingInRight),
		makeResult("NEW_DB_HOST", "", "localhost", diff.MissingInLeft),
	}
	got := Detect(results, 0.8)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if got[0].OldKey != "OLD_DB_HOST" || got[0].NewKey != "NEW_DB_HOST" {
		t.Errorf("unexpected suggestion: %+v", got[0])
	}
	if got[0].Score != 1.0 {
		t.Errorf("expected score 1.0, got %f", got[0].Score)
	}
}

func TestDetect_BelowThreshold(t *testing.T) {
	results := []diff.Result{
		makeResult("FOO", "abc", "", diff.MissingInRight),
		makeResult("BAR", "", "xyz", diff.MissingInLeft),
	}
	got := Detect(results, 0.8)
	if len(got) != 0 {
		t.Errorf("expected 0 suggestions below threshold, got %d", len(got))
	}
}

func TestDetect_IgnoresMismatched(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", "localhost", "remotehost", diff.Mismatched),
	}
	got := Detect(results, 0.8)
	if len(got) != 0 {
		t.Errorf("expected 0 suggestions for mismatched, got %d", len(got))
	}
}

func TestDetect_DefaultThreshold(t *testing.T) {
	results := []diff.Result{
		makeResult("OLD_KEY", "value123", "", diff.MissingInRight),
		makeResult("NEW_KEY", "", "value123", diff.MissingInLeft),
	}
	// threshold 0 should fall back to 0.8
	got := Detect(results, 0)
	if len(got) != 1 {
		t.Errorf("expected 1 suggestion with default threshold, got %d", len(got))
	}
}

func TestReportText_NoSuggestions(t *testing.T) {
	var buf bytes.Buffer
	ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No rename suggestions") {
		t.Errorf("expected no-suggestions message, got: %s", buf.String())
	}
}

func TestReportText_WithSuggestions(t *testing.T) {
	suggestions := []Suggestion{
		{OldKey: "OLD_HOST", NewKey: "NEW_HOST", OldValue: "localhost", NewValue: "localhost", Score: 1.0},
	}
	var buf bytes.Buffer
	ReportText(&buf, suggestions)
	out := buf.String()
	if !strings.Contains(out, "OLD_HOST") || !strings.Contains(out, "NEW_HOST") {
		t.Errorf("expected keys in output, got: %s", out)
	}
	if !strings.Contains(out, "100%") {
		t.Errorf("expected confidence in output, got: %s", out)
	}
}

func TestReportJSON_Structure(t *testing.T) {
	suggestions := []Suggestion{
		{OldKey: "A", NewKey: "B", OldValue: "v", NewValue: "v", Score: 0.9},
	}
	var buf bytes.Buffer
	if err := ReportJSON(&buf, suggestions); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"old_key", "new_key", "confidence"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got: %s", want, out)
		}
	}
}

func TestValueSimilarity_Equal(t *testing.T) {
	if s := valueSimilarity("hello", "hello"); s != 1.0 {
		t.Errorf("expected 1.0, got %f", s)
	}
}

func TestValueSimilarity_Empty(t *testing.T) {
	if s := valueSimilarity("", ""); s != 1.0 {
		t.Errorf("expected 1.0 for two empty strings, got %f", s)
	}
	if s := valueSimilarity("abc", ""); s != 0.0 {
		t.Errorf("expected 0.0 for empty vs non-empty, got %f", s)
	}
}
