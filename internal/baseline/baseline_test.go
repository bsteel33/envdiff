package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/baseline"
	"github.com/user/envdiff/internal/diff"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCompare_NoBaselinePath(t *testing.T) {
	_, err := baseline.Compare(baseline.Options{TargetPaths: []string{"x.env"}})
	if err == nil {
		t.Fatal("expected error for empty BaselinePath")
	}
}

func TestCompare_NoTargets(t *testing.T) {
	_, err := baseline.Compare(baseline.Options{BaselinePath: "base.env"})
	if err == nil {
		t.Fatal("expected error for empty TargetPaths")
	}
}

func TestCompare_SingleTarget_AllEqual(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	target := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	results, err := baseline.Compare(baseline.Options{
		BaselinePath: base,
		TargetPaths:  []string{target},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	for _, r := range results[0].Results {
		if r.Status != diff.StatusEqual {
			t.Errorf("key %q: expected equal, got %v", r.Key, r.Status)
		}
	}
}

func TestCompare_SingleTarget_WithDifferences(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nONLY_BASE=yes\n")
	target := writeTempEnv(t, "FOO=changed\nONLY_TARGET=yes\n")

	results, err := baseline.Compare(baseline.Options{
		BaselinePath: base,
		TargetPaths:  []string{target},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := baseline.Summary(results)
	if sum[filepath.Base(target)] == 0 && sum[target] == 0 {
		t.Error("expected non-zero diff count for target")
	}
}

func TestSummary_Empty(t *testing.T) {
	sum := baseline.Summary(nil)
	if len(sum) != 0 {
		t.Errorf("expected empty summary, got %v", sum)
	}
}

func TestCompare_MultipleTargets(t *testing.T) {
	base := writeTempEnv(t, "KEY=value\n")
	t1 := writeTempEnv(t, "KEY=value\n")
	t2 := writeTempEnv(t, "KEY=other\n")

	results, err := baseline.Compare(baseline.Options{
		BaselinePath: base,
		TargetPaths:  []string{t1, t2},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestCompare_BaselineFileNotFound(t *testing.T) {
	_, err := baseline.Compare(baseline.Options{
		BaselinePath: "/nonexistent/path/base.env",
		TargetPaths:  []string{"target.env"},
	})
	if err == nil {
		t.Fatal("expected error when baseline file does not exist")
	}
}

func TestCompare_TargetFileNotFound(t *testing.T) {
	base := writeTempEnv(t, "KEY=value\n")
	_, err := baseline.Compare(baseline.Options{
		BaselinePath: base,
		TargetPaths:  []string{"/nonexistent/path/target.env"},
	})
	if err == nil {
		t.Fatal("expected error when target file does not exist")
	}
}
