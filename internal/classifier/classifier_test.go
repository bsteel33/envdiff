package classifier_test

import (
	"testing"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/diff"
)

func makeResult(key string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftValue: "a", RightValue: "b", Status: status}
}

func TestClassify_MatchIsNone(t *testing.T) {
	results := []diff.Result{makeResult("FOO", diff.StatusMatch)}
	got := classifier.Classify(results, nil)
	if got[0].Severity != classifier.SeverityNone {
		t.Errorf("expected none, got %s", got[0].Severity)
	}
}

func TestClassify_MismatchIsInfo(t *testing.T) {
	results := []diff.Result{makeResult("FOO", diff.StatusMismatch)}
	got := classifier.Classify(results, nil)
	if got[0].Severity != classifier.SeverityInfo {
		t.Errorf("expected info, got %s", got[0].Severity)
	}
}

func TestClassify_MissingIsWarningByDefault(t *testing.T) {
	results := []diff.Result{
		makeResult("A", diff.StatusMissingInLeft),
		makeResult("B", diff.StatusMissingInRight),
	}
	got := classifier.Classify(results, nil)
	for _, cr := range got {
		if cr.Severity != classifier.SeverityWarning {
			t.Errorf("key %s: expected warning, got %s", cr.Result.Key, cr.Severity)
		}
	}
}

func TestClassify_MissingIsCritical_WhenFlagSet(t *testing.T) {
	results := []diff.Result{makeResult("SECRET", diff.StatusMissingInRight)}
	opts := &classifier.Options{MissingIsCritical: true}
	got := classifier.Classify(results, opts)
	if got[0].Severity != classifier.SeverityCritical {
		t.Errorf("expected critical, got %s", got[0].Severity)
	}
}

func TestClassify_CriticalKeyOverridesDefault(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_PASSWORD", diff.StatusMismatch),
		makeResult("APP_NAME", diff.StatusMismatch),
	}
	opts := &classifier.Options{CriticalKeys: []string{"DB_PASSWORD"}}
	got := classifier.Classify(results, opts)

	if got[0].Severity != classifier.SeverityCritical {
		t.Errorf("DB_PASSWORD: expected critical, got %s", got[0].Severity)
	}
	if got[1].Severity != classifier.SeverityInfo {
		t.Errorf("APP_NAME: expected info, got %s", got[1].Severity)
	}
}

func TestClassify_EmptyResults(t *testing.T) {
	got := classifier.Classify(nil, nil)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}

func TestClassify_NilOptsDoesNotPanic(t *testing.T) {
	results := []diff.Result{makeResult("X", diff.StatusMissingInLeft)}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	classifier.Classify(results, nil)
}
