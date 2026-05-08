package auditor_test

import (
	"testing"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/diff"
)

func makeResult(key, status, left, right string) diff.Result {
	return diff.Result{Key: key, Status: status, LeftValue: left, RightValue: right}
}

func TestAudit_EmptyBoth(t *testing.T) {
	report := auditor.Audit(nil, nil)
	if len(report.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(report.Events))
	}
}

func TestAudit_Introduced(t *testing.T) {
	curr := []diff.Result{makeResult("NEW_KEY", "missing_in_left", "", "value")}
	report := auditor.Audit(nil, curr)

	if len(report.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(report.Events))
	}
	ev := report.Events[0]
	if ev.Kind != auditor.EventIntroduced {
		t.Errorf("expected introduced, got %s", ev.Kind)
	}
	if ev.Key != "NEW_KEY" {
		t.Errorf("unexpected key %s", ev.Key)
	}
	if ev.Previous != nil {
		t.Error("previous should be nil for introduced event")
	}
}

func TestAudit_Resolved(t *testing.T) {
	prev := []diff.Result{makeResult("OLD_KEY", "missing_in_right", "val", "")}
	report := auditor.Audit(prev, nil)

	if len(report.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(report.Events))
	}
	ev := report.Events[0]
	if ev.Kind != auditor.EventResolved {
		t.Errorf("expected resolved, got %s", ev.Kind)
	}
	if ev.Current != nil {
		t.Error("current should be nil for resolved event")
	}
}

func TestAudit_Changed(t *testing.T) {
	prev := []diff.Result{makeResult("KEY", "missing_in_right", "v", "")}
	curr := []diff.Result{makeResult("KEY", "mismatched", "v", "w")}
	report := auditor.Audit(prev, curr)

	if len(report.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(report.Events))
	}
	ev := report.Events[0]
	if ev.Kind != auditor.EventChanged {
		t.Errorf("expected changed, got %s", ev.Kind)
	}
}

func TestAudit_Persisted(t *testing.T) {
	result := makeResult("KEY", "mismatched", "a", "b")
	report := auditor.Audit([]diff.Result{result}, []diff.Result{result})

	if len(report.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(report.Events))
	}
	if report.Events[0].Kind != auditor.EventPersisted {
		t.Errorf("expected persisted, got %s", report.Events[0].Kind)
	}
}

func TestAudit_GeneratedAtSet(t *testing.T) {
	report := auditor.Audit(nil, nil)
	if report.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should not be zero")
	}
}
