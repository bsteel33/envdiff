package reporter

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeResults(statuses ...diff.StatusType) []diff.Result {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{Key: fmt.Sprintf("KEY_%d", i), Status: s}
	}
	return results
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]diff.Result{})
	if s.Total != 0 || s.MissingInRight != 0 || s.MissingInLeft != 0 || s.Mismatched != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestSummarize_Counts(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.MissingInRight},
		{Key: "B", Status: diff.MissingInRight},
		{Key: "C", Status: diff.MissingInLeft},
		{Key: "D", Status: diff.Mismatched},
	}
	s := Summarize(results)
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.MissingInRight != 2 {
		t.Errorf("expected MissingInRight=2, got %d", s.MissingInRight)
	}
	if s.MissingInLeft != 1 {
		t.Errorf("expected MissingInLeft=1, got %d", s.MissingInLeft)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
}

func TestFormatSummary_NoDifferences(t *testing.T) {
	s := Summary{}
	out := FormatSummary(s)
	if out != "No differences found." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatSummary_WithDifferences(t *testing.T) {
	s := Summary{Total: 3, MissingInRight: 1, MissingInLeft: 1, Mismatched: 1}
	out := FormatSummary(s)
	for _, want := range []string{"3 difference(s)", "1 missing in right", "1 missing in left", "1 mismatched"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output %q", want, out)
		}
	}
}

func TestFormatSummary_OnlyMismatched(t *testing.T) {
	s := Summary{Total: 2, Mismatched: 2}
	out := FormatSummary(s)
	if !strings.Contains(out, "2 mismatched") {
		t.Errorf("expected '2 mismatched' in %q", out)
	}
	if strings.Contains(out, "missing") {
		t.Errorf("did not expect 'missing' in %q", out)
	}
}
