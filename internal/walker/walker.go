// Package walker provides utilities for discovering .env files
// within a directory tree, optionally matching glob patterns.
package walker

import (
	"os"
	"path/filepath"
	"strings"
)

// Options controls how the directory walk is performed.
type Options struct {
	// Patterns is a list of glob patterns to match filenames (e.g. ".env*").
	// If empty, defaults to [".env", ".env.*"].
	Patterns []string
	// MaxDepth limits recursion depth. 0 means unlimited.
	MaxDepth int
}

var defaultPatterns = []string{".env", ".env.*"}

// Walk traverses root and returns all .env file paths that match opts.
func Walk(root string, opts Options) ([]string, error) {
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = defaultPatterns
	}

	var results []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if opts.MaxDepth > 0 {
				rel, relErr := filepath.Rel(root, path)
				if relErr == nil {
					depth := strings.Count(rel, string(os.PathSeparator))
					if rel != "." && depth >= opts.MaxDepth {
						return filepath.SkipDir
					}
				}
			}
			return nil
		}

		name := d.Name()
		for _, pattern := range patterns {
			matched, matchErr := filepath.Match(pattern, name)
			if matchErr != nil {
				return matchErr
			}
			if matched {
				results = append(results, path)
				break
			}
		}
		return nil
	})

	return results, err
}
