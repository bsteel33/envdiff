// Package linter implements style and convention checks for .env file keys
// and values. It provides a set of built-in Rule functions that can be
// composed and extended, as well as text and JSON reporters for surfacing
// lint issues to end users.
//
// Usage:
//
//	rules := linter.DefaultRules()
//	issues := linter.Lint(envMap, rules)
//	linter.ReportText(os.Stdout, issues)
//
// Rules are simple functions with the signature func(key, value string) []Issue,
// making it straightforward to add project-specific checks.
package linter
