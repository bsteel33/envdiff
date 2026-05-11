// Package baseline implements baseline comparison for envdiff.
//
// A "baseline" is a reference .env file (e.g. .env.example or .env.production)
// that is treated as the source of truth. One or more target environment files
// are then compared against it to surface keys that are missing, extra, or
// carry different values.
//
// Basic usage:
//
//	results, err := baseline.Compare(baseline.Options{
//		BaselinePath: ".env.example",
//		TargetPaths:  []string{".env.staging", ".env.production"},
//	})
//
// The returned []Result slice contains one entry per target, each holding the
// full []diff.Result produced by internal/diff.Compare. Use Summary to get a
// quick count of differing keys per target.
package baseline
