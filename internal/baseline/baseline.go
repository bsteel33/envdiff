// Package baseline provides functionality to establish and compare
// a reference .env file against one or more target environments.
package baseline

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/loader"
)

// Result holds the comparison outcome between the baseline and a target.
type Result struct {
	// Label is a human-readable name for the target (e.g. filename or env name).
	Label string
	// Results contains the per-key diff results.
	Results []diff.Result
}

// Options configures how the baseline comparison is performed.
type Options struct {
	// BaselinePath is the path to the reference .env file.
	BaselinePath string
	// TargetPaths is a list of .env files to compare against the baseline.
	TargetPaths []string
}

// Compare loads the baseline file and compares each target against it,
// returning one Result per target.
func Compare(opts Options) ([]Result, error) {
	if opts.BaselinePath == "" {
		return nil, fmt.Errorf("baseline: BaselinePath must not be empty")
	}
	if len(opts.TargetPaths) == 0 {
		return nil, fmt.Errorf("baseline: at least one TargetPath is required")
	}

	baseEnv, err := loader.Load(opts.BaselinePath)
	if err != nil {
		return nil, fmt.Errorf("baseline: loading baseline %q: %w", opts.BaselinePath, err)
	}

	var results []Result
	for _, target := range opts.TargetPaths {
		targetEnv, err := loader.Load(target)
		if err != nil {
			return nil, fmt.Errorf("baseline: loading target %q: %w", target, err)
		}
		cmp := diff.Compare(baseEnv, targetEnv)
		results = append(results, Result{
			Label:   target,
			Results: cmp,
		})
	}
	return results, nil
}

// Summary returns a map from target label to the count of differing keys.
func Summary(results []Result) map[string]int {
	out := make(map[string]int, len(results))
	for _, r := range results {
		count := 0
		for _, res := range r.Results {
			if res.Status != diff.StatusEqual {
				count++
			}
		}
		out[r.Label] = count
	}
	return out
}
