package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/merger"
)

func makePair(left, right string, keys ...string) merger.PairResult {
	diffs := make([]diff.Result, len(keys))
	for i, k := range keys {
		diffs[i] = diff.Result{Key: k, Status: diff.MissingInRight}
	}
	return merger.PairResult{Left: left, Right: right, Diffs: diffs}
}

func TestMerge_Empty(t *testing.T) {
	result := merger.Merge(nil)
	if len(result.Keys) != 0 {
		t.Errorf("expected empty keys map, got %d entries", len(result.Keys))
	}
	if len(result.Pairs) != 0 {
		t.Errorf("expected empty pairs, got %d", len(result.Pairs))
	}
}

func TestMerge_SinglePair(t *testing.T) {
	pair := makePair(".env.dev", ".env.prod", "DB_HOST", "API_KEY")
	result := merger.Merge([]merger.PairResult{pair})

	if len(result.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result.Keys))
	}
	if _, ok := result.Keys["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in keys map")
	}
}

func TestMerge_MultiplePairs_SharedKey(t *testing.T) {
	p1 := makePair(".env.dev", ".env.prod", "SHARED_KEY", "ONLY_IN_P1")
	p2 := makePair(".env.dev", ".env.staging", "SHARED_KEY", "ONLY_IN_P2")
	result := merger.Merge([]merger.PairResult{p1, p2})

	pairs := result.Keys["SHARED_KEY"]
	if len(pairs) != 2 {
		t.Errorf("expected SHARED_KEY in 2 pairs, got %d", len(pairs))
	}
	if len(result.Keys["ONLY_IN_P1"]) != 1 {
		t.Error("expected ONLY_IN_P1 in exactly 1 pair")
	}
}

func TestKeysWithMostDifferences_Order(t *testing.T) {
	p1 := makePair("a", "b", "HOT_KEY", "COLD_KEY")
	p2 := makePair("a", "c", "HOT_KEY")
	p3 := makePair("a", "d", "HOT_KEY")
	result := merger.Merge([]merger.PairResult{p1, p2, p3})

	keys := result.KeysWithMostDifferences()
	if len(keys) == 0 {
		t.Fatal("expected non-empty keys list")
	}
	if keys[0] != "HOT_KEY" {
		t.Errorf("expected HOT_KEY first, got %s", keys[0])
	}
}

func TestKeysWithMostDifferences_TieBreakAlphabetical(t *testing.T) {
	p1 := makePair("a", "b", "ZEBRA", "APPLE")
	result := merger.Merge([]merger.PairResult{p1})

	keys := result.KeysWithMostDifferences()
	if len(keys) < 2 {
		t.Fatal("expected at least 2 keys")
	}
	if keys[0] != "APPLE" {
		t.Errorf("expected APPLE before ZEBRA on tie, got %s", keys[0])
	}
}
