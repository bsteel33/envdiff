// Package linter provides style and convention checks for .env file keys.
package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a single linting problem found in an env file.
type Issue struct {
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// Rule is a function that inspects a key-value pair and returns issues.
type Rule func(key, value string) []Issue

// DefaultRules returns the standard set of linting rules.
func DefaultRules() []Rule {
	return []Rule{
		RuleUppercaseKey,
		RuleNoSpaceInKey,
		RuleNonEmptyValue,
		RuleNoLeadingUnderscore,
	}
}

var upperRe = regexp.MustCompile(`^[A-Z0-9_]+$`)

// RuleUppercaseKey warns when a key contains lowercase letters.
func RuleUppercaseKey(key, _ string) []Issue {
	if !upperRe.MatchString(key) {
		return []Issue{{Key: key, Message: "key should be UPPER_SNAKE_CASE", Severity: "warn"}}
	}
	return nil
}

// RuleNoSpaceInKey errors when a key contains whitespace.
func RuleNoSpaceInKey(key, _ string) []Issue {
	if strings.ContainsAny(key, " \t") {
		return []Issue{{Key: key, Message: "key must not contain spaces or tabs", Severity: "error"}}
	}
	return nil
}

// RuleNonEmptyValue warns when a value is empty.
func RuleNonEmptyValue(key, value string) []Issue {
	if strings.TrimSpace(value) == "" {
		return []Issue{{Key: key, Message: "value is empty", Severity: "warn"}}
	}
	return nil
}

// RuleNoLeadingUnderscore warns when a key starts with an underscore.
func RuleNoLeadingUnderscore(key, _ string) []Issue {
	if strings.HasPrefix(key, "_") {
		return []Issue{{Key: key, Message: "key should not start with an underscore", Severity: "warn"}}
	}
	return nil
}

// Lint runs all provided rules against the given env map and returns all issues.
func Lint(env map[string]string, rules []Rule) []Issue {
	var issues []Issue
	for k, v := range env {
		for _, rule := range rules {
			issues = append(issues, rule(k, v)...)
		}
	}
	return issues
}

// FormatIssue returns a human-readable string for a single Issue.
func FormatIssue(i Issue) string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Severity), i.Key, i.Message)
}
