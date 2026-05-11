package annotator_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/annotator"
	"github.com/yourusername/envdiff/internal/diff"
)

func makeResult(key string, status diff.Status, left, right string) diff.Result {
	return diff.Result{Key: key, Status: status, LeftValue: left, RightValue: right}
}

func TestAnnotate_Empty(t *testing.T) {
	out := annotator.Annotate(nil)
	if len(out) != 0 {
		t.Fatalf("expected 0 annotations, got %d", len(out))
	}
}

func TestAnnotate_MissingInLeft(t *testing.T) {
	res := []diff.Result{makeResult("DB_HOST", diff.StatusMissingInLeft, "", "localhost")}
	out := annotator.Annotate(res)
	if len(out) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(out))
	}
	a := out[0]
	if a.Key != "DB_HOST" {
		t.Errorf("unexpected key: %s", a.Key)
	}
	if a.Suggestion == "" {
		t.Error("expected non-empty suggestion")
	}
}

func TestAnnotate_MissingInRight(t *testing.T) {
	res := []diff.Result{makeResult("API_KEY", diff.StatusMissingInRight, "secret", "")}
	out := annotator.Annotate(res)
	if out[0].Status != string(diff.StatusMissingInRight) {
		t.Errorf("wrong status: %s", out[0].Status)
	}
	if out[0].Suggestion == "" {
		t.Error("expected suggestion")
	}
}

func TestAnnotate_Mismatch(t *testing.T) {
	res := []diff.Result{makeResult("PORT", diff.StatusMismatch, "3000", "4000")}
	out := annotator.Annotate(res)
	if out[0].Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestAnnotate_Equal(t *testing.T) {
	res := []diff.Result{makeResult("LOG_LEVEL", diff.StatusEqual, "info", "info")}
	out := annotator.Annotate(res)
	if out[0].Suggestion != "" {
		t.Errorf("expected no suggestion for equal key, got: %s", out[0].Suggestion)
	}
}
