package summarizer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summarizer"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	out := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Result{Key: "KEY", Status: s}
	}
	return out
}

func TestSummarize_Empty(t *testing.T) {
	r := summarizer.Summarize(nil)
	if len(r.Sets) != 0 {
		t.Fatalf("expected 0 sets, got %d", len(r.Sets))
	}
	if r.GrandTotal != 0 {
		t.Errorf("expected GrandTotal 0, got %d", r.GrandTotal)
	}
	if r.GrandDrift != 0 {
		t.Errorf("expected GrandDrift 0, got %f", r.GrandDrift)
	}
}

func TestSummarize_SingleSet_AllMatched(t *testing.T) {
	sets := map[string][]diff.Result{
		"prod": makeResults(diff.StatusMatch, diff.StatusMatch),
	}
	r := summarizer.Summarize(sets)
	if len(r.Sets) != 1 {
		t.Fatalf("expected 1 set")
	}
	s := r.Sets[0]
	if s.Matched != 2 || s.Missing != 0 || s.Mismatched != 0 {
		t.Errorf("unexpected counts: %+v", s)
	}
	if s.DriftPct != 0 {
		t.Errorf("expected 0%% drift, got %.2f", s.DriftPct)
	}
}

func TestSummarize_SingleSet_WithDrift(t *testing.T) {
	sets := map[string][]diff.Result{
		"staging": makeResults(
			diff.StatusMatch,
			diff.StatusMissingInRight,
			diff.StatusMismatch,
			diff.StatusMismatch,
		),
	}
	r := summarizer.Summarize(sets)
	s := r.Sets[0]
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Missing)
	}
	if s.Mismatched != 2 {
		t.Errorf("expected Mismatched=2, got %d", s.Mismatched)
	}
	want := 75.0
	if s.DriftPct != want {
		t.Errorf("expected DriftPct=%.1f, got %.1f", want, s.DriftPct)
	}
}

func TestSummarize_MultipleSets_SortedByLabel(t *testing.T) {
	sets := map[string][]diff.Result{
		"zebra": makeResults(diff.StatusMatch),
		"alpha": makeResults(diff.StatusMismatch),
	}
	r := summarizer.Summarize(sets)
	if r.Sets[0].Label != "alpha" || r.Sets[1].Label != "zebra" {
		t.Errorf("sets not sorted: %v, %v", r.Sets[0].Label, r.Sets[1].Label)
	}
}

func TestRender_NoSets(t *testing.T) {
	var buf bytes.Buffer
	summarizer.Render(&buf, summarizer.Report{})
	if !strings.Contains(buf.String(), "No result sets") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestRender_WithSets(t *testing.T) {
	sets := map[string][]diff.Result{
		"prod": makeResults(diff.StatusMatch, diff.StatusMissingInLeft),
	}
	r := summarizer.Summarize(sets)
	var buf bytes.Buffer
	summarizer.Render(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Errorf("expected label 'prod' in output")
	}
	if !strings.Contains(out, "GRAND TOTAL") {
		t.Errorf("expected GRAND TOTAL row in output")
	}
}
