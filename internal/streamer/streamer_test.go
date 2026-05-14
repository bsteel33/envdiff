package streamer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/streamer"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func collectEvents(ch <-chan streamer.Event) ([]diff.Result, error) {
	var results []diff.Result
	for e := range ch {
		if e.Err != nil {
			return nil, e.Err
		}
		results = append(results, e.Result)
	}
	return results, nil
}

func TestStream_NoDifferences(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	right := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	ch := make(chan streamer.Event, 10)
	streamer.Stream(left, right, nil, ch)

	results, err := collectEvents(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Status != diff.StatusMatch {
			t.Errorf("expected match for %q, got %s", r.Key, r.Status)
		}
	}
}

func TestStream_MissingInRight(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nONLY_LEFT=yes\n")
	right := writeTempEnv(t, "FOO=bar\n")

	ch := make(chan streamer.Event, 10)
	streamer.Stream(left, right, nil, ch)

	results, err := collectEvents(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, r := range results {
		if r.Key == "ONLY_LEFT" && r.Status == diff.StatusMissingInRight {
			found = true
		}
	}
	if !found {
		t.Error("expected ONLY_LEFT to be missing in right")
	}
}

func TestStream_MismatchedValue(t *testing.T) {
	left := writeTempEnv(t, "API_KEY=abc\n")
	right := writeTempEnv(t, "API_KEY=xyz\n")

	ch := make(chan streamer.Event, 10)
	streamer.Stream(left, right, nil, ch)

	results, err := collectEvents(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != diff.StatusMismatch {
		t.Errorf("expected mismatch, got %+v", results)
	}
}

func TestStream_FileNotFound(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nonexistent.env")
	existing := writeTempEnv(t, "FOO=bar\n")

	ch := make(chan streamer.Event, 10)
	streamer.Stream(missing, existing, nil, ch)

	_, err := collectEvents(ch)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestStream_SkipsCommentsAndBlanks(t *testing.T) {
	content := "# comment\n\nFOO=bar\n"
	left := writeTempEnv(t, content)
	right := writeTempEnv(t, content)

	ch := make(chan streamer.Event, 10)
	streamer.Stream(left, right, nil, ch)

	results, err := collectEvents(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
