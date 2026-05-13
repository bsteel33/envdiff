// Package timeline provides chronological tracking of environment diff results.
//
// A Timeline accumulates diff snapshots tagged with a label and timestamp.
// Calling Trend() reduces each snapshot to a TrendPoint containing counts of
// total, missing, and mismatched keys, enabling callers to observe how
// environment drift evolves over time.
//
// Usage:
//
//	tl := &timeline.Timeline{}
//	tl.Add("deploy-123", time.Now(), results)
//	for _, pt := range tl.Trend() {
//		fmt.Printf("%s: missing=%d mismatched=%d\n", pt.Label, pt.Missing, pt.Mismatched)
//	}
package timeline
