// Package patcher generates patch suggestions from diff results.
//
// Given a slice of diff.Result values, patcher.Generate produces a list
// of Patch structs describing the minimal changes needed to reconcile
// the left environment with the right environment:
//
//   - Keys missing in the left file are marked for addition.
//   - Keys missing in the right file are marked for removal.
//   - Keys present in both but with differing values are marked for update.
//
// Example usage:
//
//	results := diff.Compare(left, right)
//	patches := patcher.Generate(results)
//	for _, p := range patches {
//		fmt.Println(p)
//	}
//
// Use RenderEnv to produce a minimal .env snippet containing only the
// keys that need to be added or updated.
package patcher
