package digester_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/digester"
)

func TestStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "digest.json")
	store := digester.NewStore(path)

	results := []diff.Result{
		makeResult("HOST", "matched", "a", "a"),
		makeResult("PORT", "mismatched", "80", "8080"),
	}
	orig := digester.Compute(results)

	if err := store.Save(orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Hash != orig.Hash {
		t.Fatalf("hash mismatch: want %s got %s", orig.Hash, loaded.Hash)
	}
	if loaded.Entries != orig.Entries {
		t.Fatalf("entries mismatch: want %d got %d", orig.Entries, loaded.Entries)
	}
}

func TestStore_Load_NotFound(t *testing.T) {
	store := digester.NewStore("/nonexistent/path/digest.json")
	_, err := store.Load()
	if !errors.Is(err, digester.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	store := digester.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestStore_Save_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "digest.json")
	store := digester.NewStore(path)

	first := digester.Compute([]diff.Result{makeResult("A", "matched", "1", "1")})
	second := digester.Compute([]diff.Result{makeResult("B", "mismatched", "x", "y")})

	if err := store.Save(first); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(second); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Hash != second.Hash {
		t.Fatalf("expected second hash %s, got %s", second.Hash, loaded.Hash)
	}
}
