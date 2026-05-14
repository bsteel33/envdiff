package streamer

import (
	"github.com/user/envdiff/internal/diff"
)

// Collect drains ch and returns all results, or the first error encountered.
// It blocks until the channel is closed by Stream.
func Collect(ch <-chan Event) ([]diff.Result, error) {
	var results []diff.Result
	for e := range ch {
		if e.Err != nil {
			// Drain remaining events to avoid goroutine leak.
			go func() {
				for range ch { //nolint:revive
				}
			}()
			return nil, e.Err
		}
		results = append(results, e.Result)
	}
	return results, nil
}

// CollectFiltered drains ch and returns only results whose status is in the
// provided set. Pass an empty set to collect all statuses.
func CollectFiltered(ch <-chan Event, statuses ...diff.Status) ([]diff.Result, error) {
	allow := make(map[diff.Status]struct{}, len(statuses))
	for _, s := range statuses {
		allow[s] = struct{}{}
	}

	var results []diff.Result
	for e := range ch {
		if e.Err != nil {
			go func() {
				for range ch { //nolint:revive
				}
			}()
			return nil, e.Err
		}
		if len(allow) == 0 {
			results = append(results, e.Result)
			continue
		}
		if _, ok := allow[e.Result.Status]; ok {
			results = append(results, e.Result)
		}
	}
	return results, nil
}
