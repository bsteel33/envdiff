package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Compare(left, right)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"FOO": "bar", "ONLY_LEFT": "value"}
	right := map[string]string{"FOO": "bar"}

	result := Compare(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("expected ONLY_LEFT missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "bar", "ONLY_RIGHT": "value"}

	result := Compare(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("expected ONLY_RIGHT missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"FOO": "old_value"}
	right := map[string]string{"FOO": "new_value"}

	result := Compare(left, right)

	vals, ok := result.Mismatched["FOO"]
	if !ok {
		t.Fatal("expected FOO to be in mismatched")
	}
	if vals[0] != "old_value" || vals[1] != "new_value" {
		t.Errorf("unexpected mismatch values: %v", vals)
	}
}

func TestCompare_MixedDifferences(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "same"}
	right := map[string]string{"A": "99", "D": "4", "C": "same"}

	result := Compare(left, right)

	if !result.HasDifferences() {
		t.Fatal("expected differences")
	}
	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "B" {
		t.Errorf("expected B missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "D" {
		t.Errorf("expected D missing in left, got %v", result.MissingInLeft)
	}
	if _, ok := result.Mismatched["A"]; !ok {
		t.Error("expected A to be mismatched")
	}
}
