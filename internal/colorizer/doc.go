// Package colorizer maps envdiff diff statuses to ANSI terminal colour codes.
//
// Usage:
//
//	c := colorizer.New(colorizer.Options{NoColor: false})
//	fmt.Println(c.ForStatus("DB_PASSWORD", diff.StatusMissingInRight))
//	fmt.Println(c.Label(diff.StatusMismatch))
//
// When NoColor is true all methods return the input string unchanged, making
// it safe to use in non-interactive environments (CI, log files, etc.).
package colorizer
