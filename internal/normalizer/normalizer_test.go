package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/normalizer"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftValue: left, RightValue: right, Status: status}
}

func TestApply_NoOptions_ReturnsUnchanged(t *testing.T) {
	input := []diff.Result{
		makeResult("DB_HOST", "  localhost  ", "  prod.db  ", diff.Mismatched),
	}
	opts := normalizer.Options{}
	out := normalizer.Apply(input, opts)

	if out[0].LeftValue != "  localhost  " {
		t.Errorf("expected unchanged left value, got %q", out[0].LeftValue)
	}
}

func TestApply_TrimValues(t *testing.T) {
	input := []diff.Result{
		makeResult("API_KEY", "  abc123  ", "\txyz789\n", diff.Mismatched),
	}
	opts := normalizer.Options{TrimValues: true}
	out := normalizer.Apply(input, opts)

	if out[0].LeftValue != "abc123" {
		t.Errorf("expected trimmed left value, got %q", out[0].LeftValue)
	}
	if out[0].RightValue != "xyz789" {
		t.Errorf("expected trimmed right value, got %q", out[0].RightValue)
	}
}

func TestApply_LowercaseKeys(t *testing.T) {
	input := []diff.Result{
		makeResult("DB_HOST", "localhost", "", diff.MissingInRight),
	}
	opts := normalizer.Options{LowercaseKeys: true}
	out := normalizer.Apply(input, opts)

	if out[0].Key != "db_host" {
		t.Errorf("expected lowercase key, got %q", out[0].Key)
	}
}

func TestApply_CollapseWhitespace(t *testing.T) {
	input := []diff.Result{
		makeResult("MSG", "hello   world", "hello world", diff.Mismatched),
	}
	opts := normalizer.Options{CollapseWhitespace: true}
	out := normalizer.Apply(input, opts)

	if out[0].LeftValue != "hello world" {
		t.Errorf("expected collapsed left value, got %q", out[0].LeftValue)
	}
	if out[0].RightValue != "hello world" {
		t.Errorf("expected collapsed right value, got %q", out[0].RightValue)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := []diff.Result{
		makeResult("KEY", "  val  ", "  val  ", diff.Matched),
	}
	opts := normalizer.DefaultOptions()
	normalizer.Apply(input, opts)

	if input[0].LeftValue != "  val  " {
		t.Error("original slice was mutated")
	}
}

func TestDefaultOptions_TrimValues(t *testing.T) {
	opts := normalizer.DefaultOptions()
	if !opts.TrimValues {
		t.Error("expected TrimValues to be true by default")
	}
	if opts.LowercaseKeys {
		t.Error("expected LowercaseKeys to be false by default")
	}
	if opts.CollapseWhitespace {
		t.Error("expected CollapseWhitespace to be false by default")
	}
}
