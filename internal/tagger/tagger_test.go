package tagger_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/tagger"
)

func makeResult(key string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestApply_UnknownTag(t *testing.T) {
	results := []diff.Result{makeResult("SOME_RANDOM_VAR", diff.StatusMatch)}
	tagged := tagger.Apply(results, tagger.DefaultRules())
	if len(tagged) != 1 {
		t.Fatalf("expected 1 result, got %d", len(tagged))
	}
	if tagged[0].Tags[0] != tagger.TagUnknown {
		t.Errorf("expected TagUnknown, got %v", tagged[0].Tags[0])
	}
}

func TestApply_SecretTag(t *testing.T) {
	results := []diff.Result{makeResult("DB_PASSWORD", diff.StatusMismatch)}
	tagged := tagger.Apply(results, tagger.DefaultRules())
	if len(tagged) != 1 {
		t.Fatalf("expected 1 result, got %d", len(tagged))
	}
	hasSecret := false
	hasDB := false
	for _, tag := range tagged[0].Tags {
		if tag == tagger.TagSecret {
			hasSecret = true
		}
		if tag == tagger.TagDB {
			hasDB = true
		}
	}
	if !hasSecret {
		t.Error("expected TagSecret")
	}
	if !hasDB {
		t.Error("expected TagDB")
	}
}

func TestApply_FeatureTag(t *testing.T) {
	results := []diff.Result{makeResult("FEATURE_FLAG_X", diff.StatusMissingInRight)}
	tagged := tagger.Apply(results, tagger.DefaultRules())
	if len(tagged) != 1 {
		t.Fatalf("expected 1 result, got %d", len(tagged))
	}
	if tagged[0].Tags[0] != tagger.TagFeature {
		t.Errorf("expected TagFeature, got %v", tagged[0].Tags[0])
	}
}

func TestApply_EmptyResults(t *testing.T) {
	tagged := tagger.Apply(nil, tagger.DefaultRules())
	if len(tagged) != 0 {
		t.Errorf("expected empty slice, got %d", len(tagged))
	}
}

func TestApply_CustomRules(t *testing.T) {
	rules := []tagger.Rule{
		{Contains: "CUSTOM", Tag: "custom"},
	}
	results := []diff.Result{makeResult("MY_CUSTOM_VAR", diff.StatusMatch)}
	tagged := tagger.Apply(results, rules)
	if tagged[0].Tags[0] != "custom" {
		t.Errorf("expected custom tag, got %v", tagged[0].Tags[0])
	}
}
