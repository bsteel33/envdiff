package multiplexer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/multiplexer"
)

var base = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "db.prod",
	"SECRET":   "s3cr3t",
}

func TestRun_NoTargets(t *testing.T) {
	got := multiplexer.Run(base, nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestRun_SingleTarget_OrderedByLabel(t *testing.T) {
	targets := map[string]map[string]string{
		"staging": {"APP_ENV": "staging", "DB_HOST": "db.staging"},
	}
	got := multiplexer.Run(base, targets)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Label != "staging" {
		t.Errorf("unexpected label %q", got[0].Label)
	}
	// SECRET is missing in staging
	var foundMissing bool
	for _, r := range got[0].Results {
		if r.Key == "SECRET" && r.Status == diff.MissingInRight {
			foundMissing = true
		}
	}
	if !foundMissing {
		t.Error("expected SECRET to be missing in right (staging)")
	}
}

func TestRun_MultipleTargets_SortedByLabel(t *testing.T) {
	targets := map[string]map[string]string{
		"zeta": {"APP_ENV": "zeta"},
		"alpha": {"APP_ENV": "alpha"},
		"beta": {"APP_ENV": "beta"},
	}
	got := multiplexer.Run(base, targets)
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	want := []string{"alpha", "beta", "zeta"}
	for i, nr := range got {
		if nr.Label != want[i] {
			t.Errorf("position %d: want %q, got %q", i, want[i], nr.Label)
		}
	}
}

func TestFlatten_Empty(t *testing.T) {
	got := multiplexer.Flatten(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestFlatten_HighestSeverityWins(t *testing.T) {
	named := []multiplexer.NamedResult{
		{
			Label: "a",
			Results: []diff.Result{
				{Key: "DB_HOST", Status: diff.MissingInRight},
			},
		},
		{
			Label: "b",
			Results: []diff.Result{
				{Key: "DB_HOST", Status: diff.Mismatched, LeftVal: "a", RightVal: "b"},
			},
		},
	}
	got := multiplexer.Flatten(named)
	if len(got) != 1 {
		t.Fatalf("expected 1 flattened result, got %d", len(got))
	}
	if got[0].Status != diff.Mismatched {
		t.Errorf("expected Mismatched status, got %v", got[0].Status)
	}
}
