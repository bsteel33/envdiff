package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: "ok", LeftVal: "localhost", RightVal: "localhost"},
		{Key: "API_KEY", Status: "missing_right", LeftVal: "secret", RightVal: ""},
		{Key: "PORT", Status: "mismatched", LeftVal: "8080", RightVal: "9090"},
	}
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := sampleResults()
	if err := snapshot.Save(path, "test-snap", results); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if snap.Label != "test-snap" {
		t.Errorf("expected label 'test-snap', got %q", snap.Label)
	}
	if len(snap.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiff_StatusChange(t *testing.T) {
	before := &snapshot.Snapshot{
		CreatedAt: time.Now(),
		Label:     "before",
		Results: []diff.Result{
			{Key: "API_KEY", Status: "missing_right"},
			{Key: "DB_HOST", Status: "ok"},
		},
	}
	after := &snapshot.Snapshot{
		CreatedAt: time.Now(),
		Label:     "after",
		Results: []diff.Result{
			{Key: "API_KEY", Status: "ok"},
			{Key: "DB_HOST", Status: "ok"},
			{Key: "NEW_KEY", Status: "missing_left"},
		},
	}

	changes := snapshot.Diff(before, after)

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}

	changeMap := make(map[string]snapshot.Change)
	for _, c := range changes {
		changeMap[c.Key] = c
	}

	if c, ok := changeMap["API_KEY"]; !ok || c.Before != "missing_right" || c.After != "ok" {
		t.Errorf("unexpected API_KEY change: %+v", changeMap["API_KEY"])
	}
	if c, ok := changeMap["NEW_KEY"]; !ok || c.Before != "absent" || c.After != "missing_left" {
		t.Errorf("unexpected NEW_KEY change: %+v", c)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	results := sampleResults()
	before := &snapshot.Snapshot{Label: "a", Results: results}
	after := &snapshot.Snapshot{Label: "b", Results: results}

	changes := snapshot.Diff(before, after)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}
