package comparer_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparer"
	"github.com/user/envdiff/internal/diff"
)

func makeEnv(label string, vars map[string]string) comparer.EnvMap {
	return comparer.EnvMap{Label: label, Vars: vars}
}

func TestCompare_Empty(t *testing.T) {
	r := comparer.Compare(nil)
	if len(r.Keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestCompare_SingleEnv(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("prod", map[string]string{"FOO": "bar", "BAZ": "qux"}),
	}
	r := comparer.Compare(envs)
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
	if r.Keys[0].Key != "BAZ" {
		t.Errorf("expected sorted first key BAZ, got %s", r.Keys[0].Key)
	}
}

func TestCompare_ConsistentAcrossEnvs(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("staging", map[string]string{"PORT": "8080"}),
		makeEnv("prod", map[string]string{"PORT": "8080"}),
	}
	r := comparer.Compare(envs)
	if len(r.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(r.Keys))
	}
	if !r.Keys[0].Consistent {
		t.Error("expected key to be consistent")
	}
}

func TestCompare_InconsistentValues(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("staging", map[string]string{"PORT": "8080"}),
		makeEnv("prod", map[string]string{"PORT": "443"}),
	}
	r := comparer.Compare(envs)
	if r.Keys[0].Consistent {
		t.Error("expected key to be inconsistent")
	}
}

func TestCompare_MissingInSomeEnvs(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("staging", map[string]string{"SECRET": "abc"}),
		makeEnv("prod", map[string]string{}),
	}
	r := comparer.Compare(envs)
	kr := r.Keys[0]
	if len(kr.MissingIn) != 1 || kr.MissingIn[0] != "prod" {
		t.Errorf("expected prod to be missing, got %v", kr.MissingIn)
	}
	if len(kr.PresentIn) != 1 || kr.PresentIn[0] != "staging" {
		t.Errorf("expected staging to be present, got %v", kr.PresentIn)
	}
}

func TestToDiffResults_Match(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("base", map[string]string{"KEY": "val"}),
		makeEnv("target", map[string]string{"KEY": "val"}),
	}
	r := comparer.Compare(envs)
	results := comparer.ToDiffResults(r, "base", "target")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != diff.Match {
		t.Errorf("expected Match, got %v", results[0].Status)
	}
}

func TestToDiffResults_MissingInTarget(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("base", map[string]string{"KEY": "val"}),
		makeEnv("target", map[string]string{}),
	}
	r := comparer.Compare(envs)
	results := comparer.ToDiffResults(r, "base", "target")
	if results[0].Status != diff.MissingInRight {
		t.Errorf("expected MissingInRight, got %v", results[0].Status)
	}
}

func TestToDiffResults_Mismatch(t *testing.T) {
	envs := []comparer.EnvMap{
		makeEnv("base", map[string]string{"KEY": "a"}),
		makeEnv("target", map[string]string{"KEY": "b"}),
	}
	r := comparer.Compare(envs)
	results := comparer.ToDiffResults(r, "base", "target")
	if results[0].Status != diff.Mismatched {
		t.Errorf("expected Mismatched, got %v", results[0].Status)
	}
}
