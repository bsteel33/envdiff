// Package auditor tracks changes between two diff results over time,
// producing a structured audit log of what changed, was added, or resolved.
package auditor

import (
	"time"

	"github.com/user/envdiff/internal/diff"
)

// EventKind describes the nature of an audit event.
type EventKind string

const (
	EventIntroduced EventKind = "introduced" // key difference appeared
	EventResolved   EventKind = "resolved"   // key difference disappeared
	EventChanged    EventKind = "changed"    // difference status changed
	EventPersisted  EventKind = "persisted"  // difference unchanged between runs
)

// Event represents a single audit record for one key.
type Event struct {
	Key       string          `json:"key"`
	Kind      EventKind       `json:"kind"`
	Previous  *diff.Result    `json:"previous,omitempty"`
	Current   *diff.Result    `json:"current,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// Report is the full audit output comparing two snapshots.
type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	Events      []Event   `json:"events"`
}

// Audit compares a previous slice of diff results against a current slice
// and returns a Report describing what changed between them.
func Audit(previous, current []diff.Result) Report {
	now := time.Now().UTC()

	prevIndex := index(previous)
	currIndex := index(current)

	var events []Event

	for key, curr := range currIndex {
		currCopy := curr
		if prev, ok := prevIndex[key]; ok {
			prevCopy := prev
			if prev.Status != curr.Status {
				events = append(events, Event{
					Key:       key,
					Kind:      EventChanged,
					Previous:  &prevCopy,
					Current:   &currCopy,
					Timestamp: now,
				})
			} else {
				events = append(events, Event{
					Key:       key,
					Kind:      EventPersisted,
					Previous:  &prevCopy,
					Current:   &currCopy,
					Timestamp: now,
				})
			}
		} else {
			events = append(events, Event{
				Key:       key,
				Kind:      EventIntroduced,
				Current:   &currCopy,
				Timestamp: now,
			})
		}
	}

	for key, prev := range prevIndex {
		if _, ok := currIndex[key]; !ok {
			prevCopy := prev
			events = append(events, Event{
				Key:       key,
				Kind:      EventResolved,
				Previous:  &prevCopy,
				Timestamp: now,
			})
		}
	}

	return Report{
		GeneratedAt: now,
		Events:      events,
	}
}

func index(results []diff.Result) map[string]diff.Result {
	m := make(map[string]diff.Result, len(results))
	for _, r := range results {
		m[r.Key] = r
	}
	return m
}
