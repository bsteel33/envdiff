package patcher_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/patcher"
)

func makeResult(key, left, right string, status diff.Status) diff.Result {
	return diff.Result{Key: key, LeftValue: left, RightValue: right, Status: status}
}

func TestGenerate_EmptyResults(t *testing.T) {
	patches := patcher.Generate(nil)
	if len(patches) != 0 {
		t.Fatalf("expected 0 patches, got %d", len(patches))
	}
}

func TestGenerate_MissingInLeft(t *testing.T) {
	results := []diff.Result{makeResult("NEW_KEY", "", "value1", diff.MissingInLeft)}
	patches := patcher.Generate(results)
	if len(patches) != 1 {
		t.Fatalf("expected 1 patch, got %d", len(patches))
	}
	p := patches[0]
	if p.Action != "add" || p.Key != "NEW_KEY" || p.NewValue != "value1" {
		t.Errorf("unexpected patch: %+v", p)
	}
}

func TestGenerate_MissingInRight(t *testing.T) {
	results := []diff.Result{makeResult("OLD_KEY", "val", "", diff.MissingInRight)}
	patches := patcher.Generate(results)
	if len(patches) != 1 {
		t.Fatalf("expected 1 patch, got %d", len(patches))
	}
	p := patches[0]
	if p.Action != "remove" || p.Key != "OLD_KEY" || p.OldValue != "val" {
		t.Errorf("unexpected patch: %+v", p)
	}
}

func TestGenerate_Mismatched(t *testing.T) {
	results := []diff.Result{makeResult("HOST", "localhost", "prod.example.com", diff.Mismatched)}
	patches := patcher.Generate(results)
	if len(patches) != 1 {
		t.Fatalf("expected 1 patch, got %d", len(patches))
	}
	p := patches[0]
	if p.Action != "update" || p.OldValue != "localhost" || p.NewValue != "prod.example.com" {
		t.Errorf("unexpected patch: %+v", p)
	}
}

func TestGenerate_SkipsMatch(t *testing.T) {
	results := []diff.Result{makeResult("PORT", "8080", "8080", diff.Match)}
	patches := patcher.Generate(results)
	if len(patches) != 0 {
		t.Errorf("expected no patches for matching key, got %d", len(patches))
	}
}

func TestPatch_String(t *testing.T) {
	cases := []struct {
		patch    patcher.Patch
		contains string
	}{
		{patcher.Patch{Key: "A", Action: "add", NewValue: "1"}, "+ A=1"},
		{patcher.Patch{Key: "B", Action: "remove", OldValue: "2"}, "- B=2"},
		{patcher.Patch{Key: "C", Action: "update", OldValue: "x", NewValue: "y"}, "~ C:"},
	}
	for _, tc := range cases {
		s := tc.patch.String()
		if !strings.Contains(s, tc.contains) {
			t.Errorf("String() = %q, want to contain %q", s, tc.contains)
		}
	}
}

func TestRenderEnv_OnlyAddAndUpdate(t *testing.T) {
	patches := []patcher.Patch{
		{Key: "ADD_ME", Action: "add", NewValue: "new"},
		{Key: "UPDATE_ME", Action: "update", OldValue: "old", NewValue: "updated"},
		{Key: "REMOVE_ME", Action: "remove", OldValue: "gone"},
	}
	out := patcher.RenderEnv(patches)
	if !strings.Contains(out, "ADD_ME=new") {
		t.Errorf("expected ADD_ME in output, got: %s", out)
	}
	if !strings.Contains(out, "UPDATE_ME=updated") {
		t.Errorf("expected UPDATE_ME in output, got: %s", out)
	}
	if strings.Contains(out, "REMOVE_ME") {
		t.Errorf("REMOVE_ME should not appear in rendered env output")
	}
}
