// Package annotator attaches human-readable annotations to diff results,
// explaining why each key was flagged and suggesting remediation steps.
package annotator

import "github.com/yourusername/envdiff/internal/diff"

// Annotation holds a short reason and an optional suggestion for a diff result.
type Annotation struct {
	Key        string
	Status     string
	Reason     string
	Suggestion string
}

// Annotate returns an Annotation slice for the given diff results.
func Annotate(results []diff.Result) []Annotation {
	annotations := make([]Annotation, 0, len(results))
	for _, r := range results {
		annotations = append(annotations, annotateOne(r))
	}
	return annotations
}

func annotateOne(r diff.Result) Annotation {
	a := Annotation{
		Key:    r.Key,
		Status: string(r.Status),
	}
	switch r.Status {
	case diff.StatusMissingInLeft:
		a.Reason = "Key exists in the right environment but is absent from the left."
		a.Suggestion = "Add " + r.Key + " to the left .env file."
	case diff.StatusMissingInRight:
		a.Reason = "Key exists in the left environment but is absent from the right."
		a.Suggestion = "Add " + r.Key + " to the right .env file."
	case diff.StatusMismatch:
		a.Reason = "Key is present in both environments but values differ."
		a.Suggestion = "Reconcile the value of " + r.Key + " across environments."
	default:
		a.Reason = "No difference detected."
	}
	return a
}
