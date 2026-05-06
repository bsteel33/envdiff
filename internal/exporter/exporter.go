package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for exported env files.
type Format string

const (
	FormatEnv  Format = "env"
	FormatJSON Format = "json"
)

// Options controls export behavior.
type Options struct {
	Format    Format
	OnlyMissing bool
}

// Export writes the diff results to the given file path in the specified format.
// If path is empty, output is written to stdout.
func Export(results []diff.Result, path string, opts Options) error {
	var output string
	var err error

	switch opts.Format {
	case FormatJSON:
		output, err = toJSON(results, opts)
	default:
		output, err = toEnv(results, opts)
	}

	if err != nil {
		return fmt.Errorf("exporter: format error: %w", err)
	}

	if path == "" {
		fmt.Print(output)
		return nil
	}

	return os.WriteFile(path, []byte(output), 0644)
}

func toEnv(results []diff.Result, opts Options) (string, error) {
	var sb strings.Builder
	sorted := sortedResults(results)
	for _, r := range sorted {
		if opts.OnlyMissing && r.Status != diff.Missing {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", r.Key, r.LeftValue))
	}
	return sb.String(), nil
}

func toJSON(results []diff.Result, opts Options) (string, error) {
	type entry struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		Left   string `json:"left_value,omitempty"`
		Right  string `json:"right_value,omitempty"`
	}

	entries := make([]entry, 0, len(results))
	for _, r := range sortedResults(results) {
		if opts.OnlyMissing && r.Status != diff.Missing {
			continue
		}
		entries = append(entries, entry{
			Key:    r.Key,
			Status: string(r.Status),
			Left:   r.LeftValue,
			Right:  r.RightValue,
		})
	}

	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func sortedResults(results []diff.Result) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})
	return out
}
