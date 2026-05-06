// Package validator provides rules-based validation for env key-value pairs.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule applied to env values.
type Rule struct {
	Name    string
	Pattern *regexp.Regexp
	Required bool
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

// DefaultRules returns a set of sensible default validation rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "no-empty-value",
			Pattern: regexp.MustCompile(`^.+$`),
			Required: true,
		},
	}
}

// Validate checks the given env map against the provided rules and returns
// any violations found.
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for key, value := range env {
		for _, rule := range rules {
			if rule.Required && strings.TrimSpace(value) == "" {
				violations = append(violations, Violation{
					Key:     key,
					Value:   value,
					Rule:    rule.Name,
					Message: fmt.Sprintf("key %q has an empty value", key),
				})
				continue
			}
			if rule.Pattern != nil && !rule.Pattern.MatchString(value) {
				violations = append(violations, Violation{
					Key:     key,
					Value:   value,
					Rule:    rule.Name,
					Message: fmt.Sprintf("key %q value %q does not match rule %q", key, value, rule.Name),
				})
			}
		}
	}

	return violations
}

// ValidateKeys checks that all required keys are present in the env map.
func ValidateKeys(env map[string]string, required []string) []Violation {
	var violations []Violation
	for _, key := range required {
		if _, ok := env[key]; !ok {
			violations = append(violations, Violation{
				Key:     key,
				Rule:    "required-key",
				Message: fmt.Sprintf("required key %q is missing", key),
			})
		}
	}
	return violations
}
