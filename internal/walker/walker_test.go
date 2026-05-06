package walker_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/envdiff/internal/walker"
)

func setupDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	files := []string{
		".env",
		".env.local",
		".env.production",
		"subdir/.env",
		"subdir/nested/.env.test",
		"README.md",
	}
	for _, f := range files {
		full := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("KEY=val\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestWalk_DefaultPatterns(t *testing.T) {
	root := setupDir(t)
	got, err := walker.Walk(root, walker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 5 {
		t.Errorf("expected 5 files, got %d: %v", len(got), got)
	}
}

func TestWalk_CustomPattern(t *testing.T) {
	root := setupDir(t)
	got, err := walker.Walk(root, walker.Options{Patterns: []string{".env.local"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 file, got %d", len(got))
	}
}

func TestWalk_MaxDepth(t *testing.T) {
	root := setupDir(t)
	got, err := walker.Walk(root, walker.Options{MaxDepth: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only root-level files: .env, .env.local, .env.production
	if len(got) != 3 {
		t.Errorf("expected 3 files at depth 1, got %d: %v", len(got), got)
	}
}

func TestWalk_SkipsNonEnvFiles(t *testing.T) {
	root := setupDir(t)
	got, err := walker.Walk(root, walker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(got)
	for _, p := range got {
		if filepath.Base(p) == "README.md" {
			t.Errorf("README.md should not be included")
		}
	}
}

func TestWalk_InvalidRoot(t *testing.T) {
	_, err := walker.Walk("/nonexistent/path/xyz", walker.Options{})
	if err == nil {
		t.Error("expected error for nonexistent root, got nil")
	}
}
