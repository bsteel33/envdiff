// Package config handles loading and validating envdiff CLI configuration
// from flags, environment variables, or a config file.
package config

import (
	"errors"
	"strings"
)

// Format represents an output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatEnv  Format = "env"
)

// Config holds the resolved configuration for an envdiff run.
type Config struct {
	// Files is the ordered list of .env file paths to compare.
	Files []string

	// OutputFormat controls how results are rendered.
	OutputFormat Format

	// Prefix filters keys to only those starting with this value.
	Prefix string

	// ExcludeKeys is a set of key names to omit from results.
	ExcludeKeys []string

	// OnlyMissing restricts output to keys absent in one or more files.
	OnlyMissing bool

	// ExportPath, if non-empty, writes output to this file instead of stdout.
	ExportPath string

	// Strict causes the process to exit non-zero when any differences exist.
	Strict bool
}

// Validate checks that the Config is internally consistent and usable.
func Validate(c *Config) error {
	if len(c.Files) < 2 {
		return errors.New("at least two .env files are required for comparison")
	}
	switch c.OutputFormat {
	case FormatText, FormatJSON, FormatEnv:
		// valid
	case "":
		c.OutputFormat = FormatText
	default:
		return errors.New("unsupported output format: " + string(c.OutputFormat))
	}
	return nil
}

// ParseFormat converts a raw string to a Format, returning an error for
// unrecognised values.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatText:
		return FormatText, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatEnv:
		return FormatEnv, nil
	}
	return "", errors.New("unknown format: " + s)
}
