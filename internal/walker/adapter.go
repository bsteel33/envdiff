package walker

import (
	"fmt"
	"strings"
)

// DiscoverPairs walks root and groups discovered files by their environment
// suffix, returning pairs suitable for comparison. Files named exactly ".env"
// are treated as the base environment.
//
// For example: .env and .env.production become ("base", ".env", ".env.production").
type Pair struct {
	Label string
	Base  string
	Other string
}

// DiscoverPairs finds all .env files under root and pairs each variant with
// the base ".env" file if one exists in the same directory.
func DiscoverPairs(root string, opts Options) ([]Pair, error) {
	paths, err := Walk(root, opts)
	if err != nil {
		return nil, err
	}

	// Group by directory.
	dirMap := make(map[string][]string)
	for _, p := range paths {
		dir := dirOf(p)
		dirMap[dir] = append(dirMap[dir], p)
	}

	var pairs []Pair
	for dir, files := range dirMap {
		base := findBase(dir, files)
		if base == "" {
			continue
		}
		for _, f := range files {
			if f == base {
				continue
			}
			label := labelFor(f)
			pairs = append(pairs, Pair{Label: label, Base: base, Other: f})
		}
	}
	return pairs, nil
}

func dirOf(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return "."
	}
	return strings.Join(parts[:len(parts)-1], "/")
}

func findBase(dir string, files []string) string {
	target := fmt.Sprintf("%s/.env", dir)
	for _, f := range files {
		if f == target || f == ".env" {
			return f
		}
	}
	return ""
}

func labelFor(path string) string {
	base := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		base = path[idx+1:]
	}
	parts := strings.SplitN(base, ".", 3)
	if len(parts) == 3 {
		return parts[2]
	}
	return base
}
