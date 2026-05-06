// Package snapshot provides functionality to save and load env diff results
// as snapshots, enabling comparison of diff states over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot holds a saved diff result with metadata.
type Snapshot struct {
	CreatedAt time.Time     `json:"created_at"`
	Label     string        `json:"label"`
	Results   []diff.Result `json:"results"`
}

// Save writes a snapshot of the given results to the specified file path.
func Save(path, label string, results []diff.Result) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the specified file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return &snap, nil
}

// Diff compares two snapshots and returns keys whose status changed.
func Diff(before, after *Snapshot) []Change {
	beforeMap := indexResults(before.Results)
	afterMap := indexResults(after.Results)

	var changes []Change

	for key, bRes := range beforeMap {
		aRes, ok := afterMap[key]
		if !ok {
			changes = append(changes, Change{Key: key, Before: bRes.Status, After: "removed"})
			continue
		}
		if bRes.Status != aRes.Status {
			changes = append(changes, Change{Key: key, Before: bRes.Status, After: aRes.Status})
		}
	}

	for key, aRes := range afterMap {
		if _, ok := beforeMap[key]; !ok {
			changes = append(changes, Change{Key: key, Before: "absent", After: aRes.Status})
		}
	}

	return changes
}

// Change represents a status change for a single key between two snapshots.
type Change struct {
	Key    string `json:"key"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func indexResults(results []diff.Result) map[string]diff.Result {
	m := make(map[string]diff.Result, len(results))
	for _, r := range results {
		m[r.Key] = r
	}
	return m
}
