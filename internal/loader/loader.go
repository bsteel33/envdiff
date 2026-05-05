package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its path and parsed key-value pairs.
type EnvFile struct {
	Path string
	Vars map[string]string
}

// Load reads and parses one or more .env files, merging their contents.
// Later files take precedence over earlier ones for duplicate keys.
func Load(paths ...string) (*EnvFile, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("loader: no paths provided")
	}

	merged := make(map[string]string)

	for _, p := range paths {
		clean := filepath.Clean(p)
		if _, err := os.Stat(clean); os.IsNotExist(err) {
			return nil, fmt.Errorf("loader: file not found: %s", clean)
		}

		vars, err := parser.ParseFile(clean)
		if err != nil {
			return nil, fmt.Errorf("loader: failed to parse %s: %w", clean, err)
		}

		for k, v := range vars {
			merged[k] = v
		}
	}

	label := paths[0]
	if len(paths) > 1 {
		label = fmt.Sprintf("%s (+%d more)", paths[0], len(paths)-1)
	}

	return &EnvFile{
		Path: label,
		Vars: merged,
	}, nil
}

// MustLoad is like Load but panics on error. Useful in test helpers.
func MustLoad(paths ...string) *EnvFile {
	ef, err := Load(paths...)
	if err != nil {
		panic(err)
	}
	return ef
}
