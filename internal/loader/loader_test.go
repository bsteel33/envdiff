package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func TestLoad_SingleFile(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	ef, err := loader.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", ef.Vars["FOO"])
	}
	if ef.Vars["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", ef.Vars["BAZ"])
	}
}

func TestLoad_MergeFiles_LaterWins(t *testing.T) {
	p1 := writeTempEnv(t, "FOO=first\nSHARED=base\n")
	p2 := writeTempEnv(t, "BAR=second\nSHARED=override\n")
	ef, err := loader.Load(p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Vars["FOO"] != "first" {
		t.Errorf("expected FOO=first, got %q", ef.Vars["FOO"])
	}
	if ef.Vars["BAR"] != "second" {
		t.Errorf("expected BAR=second, got %q", ef.Vars["BAR"])
	}
	if ef.Vars["SHARED"] != "override" {
		t.Errorf("expected SHARED=override, got %q", ef.Vars["SHARED"])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := loader.Load("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoPaths(t *testing.T) {
	_, err := loader.Load()
	if err == nil {
		t.Fatal("expected error when no paths provided, got nil")
	}
}

func TestLoad_PathLabel_Single(t *testing.T) {
	p := writeTempEnv(t, "X=1\n")
	ef, err := loader.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Path != p {
		t.Errorf("expected path %q, got %q", p, ef.Path)
	}
}
