// Package digester computes a deterministic hash digest of a set of diff
// results, enabling change detection between runs without persisting full
// snapshots.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Digest holds the hex-encoded SHA-256 hash of a result set together with
// the number of entries that were hashed.
type Digest struct {
	Hash    string `json:"hash"`
	Entries int    `json:"entries"`
}

// Compute returns a stable Digest for the given results. Results are sorted
// by key before hashing so that insertion order does not affect the output.
func Compute(results []diff.Result) Digest {
	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	h := sha256.New()
	for _, r := range sorted {
		fmt.Fprintf(h, "%s|%s|%s|%s\n", r.Key, r.Status, r.Left, r.Right)
	}

	return Digest{
		Hash:    hex.EncodeToString(h.Sum(nil)),
		Entries: len(sorted),
	}
}

// Equal reports whether two Digests represent the same result set.
func Equal(a, b Digest) bool {
	return a.Hash == b.Hash
}

// Changed returns true when the two digests differ, indicating that at least
// one result has been added, removed, or modified.
func Changed(a, b Digest) bool {
	return !Equal(a, b)
}
