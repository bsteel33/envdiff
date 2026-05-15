// Package censor provides value-masking for diff results whose keys match
// known sensitive patterns (passwords, tokens, API keys, etc.).
//
// Usage:
//
//	results := diff.Compare(left, right)
//	safe := censor.Apply(results, nil) // uses DefaultSensitiveSubstrings
//
// Custom patterns and mask strings can be supplied via Options:
//
//	safe := censor.Apply(results, &censor.Options{
//		Mask:                "<hidden>",
//		SensitiveSubstrings: []string{"secret", "token"},
//	})
//
// Apply never modifies the input slice; it always returns a fresh copy.
package censor
