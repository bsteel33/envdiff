// Package digester produces stable SHA-256 fingerprints of diff result sets.
//
// A Digest is computed by sorting results by key and hashing each entry's
// key, status, left value, and right value. This makes the digest
// order-independent and sensitive only to meaningful content changes.
//
// Typical usage:
//
//	results := diff.Compare(left, right)
//	d := digester.Compute(results)
//
//	store := digester.NewStore(".envdiff.digest")
//	prev, err := store.Load()
//	if errors.Is(err, digester.ErrNotFound) || digester.Changed(prev, d) {
//		// results have changed since last run
//		_ = store.Save(d)
//	}
package digester
