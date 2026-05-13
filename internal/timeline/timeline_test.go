package timeline

import (
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
)

func makeEntry(label string, ts time.Time, results []diff.Result) Entry {
	return Entry{Label: label, Timestamp: ts, Results: results}
}

func TestAdd_And_EntriesSortedChronologically(t *testing.T) {
	tl := &Timeline{}
	t2 := time.Now()
	t1 := t2.Add(-time.Hour)
	tl.Add("second", t2, nil)
	tl.Add("first", t1, nil)

	entries := tl.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Label != "first" {
		t.Errorf("expected first entry label 'first', got %q", entries[0].Label)
	}
}

func TestTrend_Empty(t *testing.T) {
	tl := &Timeline{}
	if pts := tl.Trend(); len(pts) != 0 {
		t.Errorf("expected empty trend, got %d points", len(pts))
	}
}

func TestTrend_CountsCorrectly(t *testing.T) {
	tl := &Timeline{}
	results := []diff.Result{
		{Key: "A", Status: diff.MissingInRight},
		{Key: "B", Status: diff.Mismatched},
		{Key: "C", Status: diff.Equal},
		{Key: "D", Status: diff.MissingInLeft},
	}
	tl.Add("v1", time.Now(), results)
	pts := tl.Trend()
	if len(pts) != 1 {
		t.Fatalf("expected 1 trend point, got %d", len(pts))
	}
	p := pts[0]
	if p.Total != 4 {
		t.Errorf("expected total 4, got %d", p.Total)
	}
	if p.Missing != 2 {
		t.Errorf("expected missing 2, got %d", p.Missing)
	}
	if p.Mismatched != 1 {
		t.Errorf("expected mismatched 1, got %d", p.Mismatched)
	}
}

func TestTrend_MultipleEntries(t *testing.T) {
	tl := &Timeline{}
	now := time.Now()
	tl.Add("v1", now.Add(-time.Hour), []diff.Result{
		{Key: "X", Status: diff.MissingInRight},
	})
	tl.Add("v2", now, []diff.Result{})
	pts := tl.Trend()
	if pts[0].Missing != 1 {
		t.Errorf("v1 missing should be 1, got %d", pts[0].Missing)
	}
	if pts[1].Missing != 0 {
		t.Errorf("v2 missing should be 0, got %d", pts[1].Missing)
	}
}
