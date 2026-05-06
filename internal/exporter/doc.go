// Package exporter provides functionality to write envdiff comparison results
// to a file or stdout in various formats (env, json).
//
// Usage:
//
//	results := diff.Compare(left, right)
//	err := exporter.Export(results, "output.env", exporter.Options{
//		Format: exporter.FormatEnv,
//	})
//
// Supported formats:
//   - FormatEnv  (default): writes KEY=VALUE lines compatible with .env files.
//   - FormatJSON: writes a JSON array with key, status, left_value, right_value fields.
//
// The OnlyMissing option restricts output to keys that are absent from the
// right-hand environment, useful for generating a template of required variables.
package exporter
