// Package profiler provides statistical analysis of envdiff results.
//
// Given one or more slices of diff.Result (e.g. produced by comparing
// multiple environment pairs), Analyse returns a Profile that summarises:
//
//   - Total key count and per-status breakdown (identical, missing in
//     left/right, mismatched).
//   - The top-5 most frequently differing keys across all provided
//     result sets, useful for spotting systemic configuration drift.
//
// Example usage:
//
//	results1 := diff.Compare(envA, envB)
//	results2 := diff.Compare(envA, envC)
//	profile  := profiler.Analyse(results1, results2)
//	fmt.Printf("mismatched keys: %d\n", profile.Mismatched)
package profiler
