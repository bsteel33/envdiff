// Package ignorer provides functionality to load and apply .envdiffignore
// files, allowing users to suppress specific keys from diff results.
package ignorer

import (
	"bufio"
	"os"
	"strings"
)

// Ignorer holds the set of keys that should be ignored during diffing.
type Ignorer struct {
	keys map[string]struct{}
}

// New creates an Ignorer with the given set of keys to ignore.
func New(keys []string) *Ignorer {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.TrimSpace(k)] = struct{}{}
	}
	return &Ignorer{keys: m}
}

// LoadFile reads an ignore file (one key per line, # for comments)
// and returns an Ignorer populated with those keys.
func LoadFile(path string) (*Ignorer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var keys []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keys = append(keys, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return New(keys), nil
}

// LoadFileOrEmpty attempts to load an ignore file. If the file does not
// exist it silently returns an empty Ignorer, making it optional.
func LoadFileOrEmpty(path string) (*Ignorer, error) {
	ig, err := LoadFile(path)
	if os.IsNotExist(err) {
		return New(nil), nil
	}
	return ig, err
}

// Contains reports whether the given key should be ignored.
func (ig *Ignorer) Contains(key string) bool {
	_, ok := ig.keys[key]
	return ok
}

// Keys returns a sorted slice of all ignored keys.
func (ig *Ignorer) Keys() []string {
	out := make([]string, 0, len(ig.keys))
	for k := range ig.keys {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
