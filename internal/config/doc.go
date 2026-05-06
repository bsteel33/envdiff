// Package config defines the Config struct that captures all runtime options
// for an envdiff invocation, along with helpers to validate and parse those
// options.
//
// Typical usage:
//
//	cfg := &config.Config{
//		Files:        []string{".env.staging", ".env.production"},
//		OutputFormat: config.FormatJSON,
//		Strict:       true,
//	}
//	if err := config.Validate(cfg); err != nil {
//		log.Fatal(err)
//	}
//
// The Config struct is intentionally decoupled from flag parsing so that it
// can be populated from any source (CLI flags, a YAML config file, tests, etc.)
// without modification to this package.
package config
