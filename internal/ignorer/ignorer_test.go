package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/ignorer"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envdiffignore")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempIgnore: %v", err)
	}
	return path
}

func TestNew_ContainsKeys(t *testing.T) {
	ig := ignorer.New([]string{"SECRET", "TOKEN"})
	if !ig.Contains("SECRET") {
		t.Error("expected SECRET to be ignored")
	}
	if !ig.Contains("TOKEN") {
		t.Error("expected TOKEN to be ignored")
	}
	if ig.Contains("OTHER") {
		t.Error("OTHER should not be ignored")
	}
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempIgnore(t, "# ignore these\nSECRET\nTOKEN\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ig.Contains("SECRET") || !ig.Contains("TOKEN") {
		t.Error("expected SECRET and TOKEN to be ignored")
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempIgnore(t, "\n# comment\n\nKEY_A\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := ig.Keys()
	if len(keys) != 1 || keys[0] != "KEY_A" {
		t.Errorf("expected [KEY_A], got %v", keys)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := ignorer.LoadFile("/nonexistent/.envdiffignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFileOrEmpty_MissingFile(t *testing.T) {
	ig, err := ignorer.LoadFileOrEmpty("/nonexistent/.envdiffignore")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ig.Contains("ANYTHING") {
		t.Error("empty ignorer should not contain any keys")
	}
}

func TestKeys_Sorted(t *testing.T) {
	ig := ignorer.New([]string{"ZEBRA", "ALPHA", "MIDDLE"})
	keys := ig.Keys()
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: want %s, got %s", i, k, keys[i])
		}
	}
}
