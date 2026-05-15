package digester_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/digester"
)

func makeResult(key, status, left, right string) diff.Result {
	return diff.Result{Key: key, Status: status, Left: left, Right: right}
}

func TestCompute_EmptyResults(t *testing.T) {
	d := digester.Compute(nil)
	if d.Hash == "" {
		t.Fatal("expected non-empty hash for empty input")
	}
	if d.Entries != 0 {
		t.Fatalf("expected 0 entries, got %d", d.Entries)
	}
}

func TestCompute_Deterministic(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", "matched", "localhost", "localhost"),
		makeResult("API_KEY", "mismatched", "abc", "xyz"),
	}
	d1 := digester.Compute(results)
	d2 := digester.Compute(results)
	if d1.Hash != d2.Hash {
		t.Fatalf("expected same hash on repeated calls, got %s vs %s", d1.Hash, d2.Hash)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := []diff.Result{
		makeResult("A", "matched", "1", "1"),
		makeResult("B", "missing_in_right", "2", ""),
	}
	b := []diff.Result{
		makeResult("B", "missing_in_right", "2", ""),
		makeResult("A", "matched", "1", "1"),
	}
	if digester.Compute(a).Hash != digester.Compute(b).Hash {
		t.Fatal("hash should be order-independent")
	}
}

func TestCompute_DifferentResultsProduceDifferentHash(t *testing.T) {
	a := []diff.Result{makeResult("X", "matched", "v1", "v1")}
	b := []diff.Result{makeResult("X", "mismatched", "v1", "v2")}
	if digester.Compute(a).Hash == digester.Compute(b).Hash {
		t.Fatal("expected different hashes for different results")
	}
}

func TestEqual_And_Changed(t *testing.T) {
	results := []diff.Result{makeResult("K", "matched", "v", "v")}
	d1 := digester.Compute(results)
	d2 := digester.Compute(results)
	if !digester.Equal(d1, d2) {
		t.Fatal("expected Equal to be true")
	}
	if digester.Changed(d1, d2) {
		t.Fatal("expected Changed to be false")
	}
}

func TestChanged_WhenDifferent(t *testing.T) {
	d1 := digester.Compute([]diff.Result{makeResult("A", "matched", "1", "1")})
	d2 := digester.Compute([]diff.Result{makeResult("A", "mismatched", "1", "2")})
	if !digester.Changed(d1, d2) {
		t.Fatal("expected Changed to be true")
	}
}
