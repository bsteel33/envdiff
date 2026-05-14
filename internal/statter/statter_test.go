package statter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/statter"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	out := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Result{Key: fmt.Sprintf("KEY_%d", i), Status: s}
	}
	return out
}

func TestCompute_Empty(t *testing.T) {
	s := statter.Compute("empty", nil)
	if s.Total != 0 || s.MatchRate != 0 || s.DriftScore != 0 {
		t.Fatalf("expected zero stats, got %+v", s)
	}
}

func TestCompute_AllMatched(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMatch},
	}
	s := statter.Compute("all-match", results)
	if s.Matched != 2 || s.MatchRate != 1.0 || s.DriftScore != 0.0 {
		t.Fatalf("unexpected stats: %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMismatch},
		{Key: "C", Status: diff.StatusMissingInLeft},
		{Key: "D", Status: diff.StatusMissingInRight},
	}
	s := statter.Compute("mixed", results)
	if s.Total != 4 {
		t.Fatalf("expected Total=4, got %d", s.Total)
	}
	if s.Matched != 1 {
		t.Fatalf("expected Matched=1, got %d", s.Matched)
	}
	if s.Mismatched != 1 {
		t.Fatalf("expected Mismatched=1, got %d", s.Mismatched)
	}
	if s.MissingLeft != 1 || s.MissingRight != 1 {
		t.Fatalf("expected MissingLeft=1 MissingRight=1, got %+v", s)
	}
	wantMatch := 0.25
	if s.MatchRate != wantMatch {
		t.Fatalf("expected MatchRate=%.2f, got %.2f", wantMatch, s.MatchRate)
	}
	wantDrift := 0.75
	if s.DriftScore != wantDrift {
		t.Fatalf("expected DriftScore=%.2f, got %.2f", wantDrift, s.DriftScore)
	}
}

func TestComputeAll_SortedByDriftDesc(t *testing.T) {
	sets := map[string][]diff.Result{
		"low": {
			{Key: "A", Status: diff.StatusMatch},
			{Key: "B", Status: diff.StatusMatch},
		},
		"high": {
			{Key: "A", Status: diff.StatusMismatch},
			{Key: "B", Status: diff.StatusMissingInRight},
		},
	}
	all := statter.ComputeAll(sets)
	if len(all) != 2 {
		t.Fatalf("expected 2 stats, got %d", len(all))
	}
	if all[0].Label != "high" {
		t.Fatalf("expected first label=high (highest drift), got %s", all[0].Label)
	}
}

func TestCompute_LabelPreserved(t *testing.T) {
	s := statter.Compute("my-env", nil)
	if s.Label != "my-env" {
		t.Fatalf("expected label my-env, got %s", s.Label)
	}
}
