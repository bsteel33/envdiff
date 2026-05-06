// Package merger provides utilities for combining and analysing diff results
// across multiple environment file comparisons.
//
// When comparing more than two environment files it is useful to understand
// which keys are problematic across the board rather than just in a single
// pair. The merger package accepts a slice of [PairResult] values — each
// representing one left/right comparison — and produces a [MergedResult] that
// groups differences by key name.
//
// Typical usage:
//
//	pairs := []merger.PairResult{
//		{Left: ".env.dev", Right: ".env.prod", Diffs: devProdDiffs},
//		{Left: ".env.dev", Right: ".env.staging", Diffs: devStagingDiffs},
//	}
//	merged := merger.Merge(pairs)
//	for _, key := range merged.KeysWithMostDifferences() {
//		fmt.Printf("%s appears in %d pair(s)\n", key, len(merged.Keys[key]))
//	}
package merger
