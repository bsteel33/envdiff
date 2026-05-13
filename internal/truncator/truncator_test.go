package truncator_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/truncator"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftVal: left, RightVal: right, Status: status}
}

func TestApply_ShortValuesUnchanged(t *testing.T) {
	input := []diff.Result{
		makeResult("KEY", "short", "also short", diff.StatusMismatch),
	}
	out := truncator.Apply(input, nil)
	if out[0].LeftVal != "short" || out[0].RightVal != "also short" {
		t.Errorf("expected values unchanged, got %q / %q", out[0].LeftVal, out[0].RightVal)
	}
}

func TestApply_LongValueTruncated(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz0123456789"
	input := []diff.Result{
		makeResult("K", long, "v", diff.StatusMismatch),
	}
	out := truncator.Apply(input, &truncator.Options{MaxLen: 10, Ellipsis: "…"})
	if got := out[0].LeftVal; got != "abcdefghij…" {
		t.Errorf("unexpected truncation: %q", got)
	}
	if out[0].RightVal != "v" {
		t.Errorf("short value should be unchanged, got %q", out[0].RightVal)
	}
}

func TestApply_DefaultOptions(t *testing.T) {
	val := make([]rune, 100)
	for i := range val {
		val[i] = 'x'
	}
	input := []diff.Result{
		makeResult("K", string(val), "", diff.StatusMissingInRight),
	}
	out := truncator.Apply(input, nil)
	expected := string(val[:truncator.DefaultMaxLen]) + truncator.DefaultEllipsis
	if out[0].LeftVal != expected {
		t.Errorf("default truncation failed: len=%d", len([]rune(out[0].LeftVal)))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	long := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	input := []diff.Result{
		makeResult("K", long, long, diff.StatusMismatch),
	}
	origLeft := input[0].LeftVal
	truncator.Apply(input, &truncator.Options{MaxLen: 5})
	if input[0].LeftVal != origLeft {
		t.Error("Apply mutated the original slice")
	}
}

func TestApply_ExactLengthNotTruncated(t *testing.T) {
	val := "exactly10x" // 10 runes
	input := []diff.Result{
		makeResult("K", val, "", diff.StatusMissingInRight),
	}
	out := truncator.Apply(input, &truncator.Options{MaxLen: 10})
	if out[0].LeftVal != val {
		t.Errorf("value at exact limit should not be truncated, got %q", out[0].LeftVal)
	}
}
