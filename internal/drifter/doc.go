// Package drifter provides utilities for quantifying how far a set of
// environment variables has diverged from a reference baseline.
//
// Usage:
//
//	results := diff.Compare(base, target)
//	report  := drifter.Measure(results)
//	fmt.Printf("drift score: %.3f (%s)\n", report.Score, report.Severity)
//
// The Score field is a value in [0, 1] representing the fraction of keys
// that differ between the two environments.  The Severity label maps that
// fraction onto a human-readable tier:
//
//	none      – 0 % drift
//	low       – up to 10 %
//	moderate  – up to 35 %
//	high      – up to 65 %
//	critical  – above 65 %
package drifter
