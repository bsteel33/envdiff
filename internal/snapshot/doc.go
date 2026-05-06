// Package snapshot provides save/load functionality for envdiff results,
// allowing users to persist a point-in-time view of env differences and
// compare snapshots across runs to track how the diff state evolves.
//
// Typical usage:
//
//	// Save current diff results as a snapshot
//	err := snapshot.Save("./baseline.snap.json", "baseline", results)
//
//	// Load a previously saved snapshot
//	snap, err := snapshot.Load("./baseline.snap.json")
//
//	// Compare two snapshots to see what changed
//	changes := snapshot.Diff(snapBefore, snapAfter)
package snapshot
