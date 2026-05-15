// Package comparer provides multi-file environment comparison across a named
// set of environments, producing a unified view of key presence and value
// consistency.
package comparer

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// EnvMap is a labelled collection of key/value pairs representing one
// environment (e.g. "staging", "production").
type EnvMap struct {
	Label string
	Vars  map[string]string
}

// KeyReport summarises how a single key behaves across all environments.
type KeyReport struct {
	Key     string
	// Values holds each environment label mapped to the value found there.
	// An absent entry means the key is missing in that environment.
	Values  map[string]string
	// Consistent is true when every environment that contains the key shares
	// the same value.
	Consistent bool
	// PresentIn lists the environment labels where the key exists.
	PresentIn  []string
	// MissingIn lists the environment labels where the key is absent.
	MissingIn  []string
}

// Report is the full multi-environment comparison result.
type Report struct {
	Keys   []KeyReport
	Labels []string
}

// Compare accepts two or more EnvMaps and returns a Report describing how
// every key behaves across all supplied environments.
func Compare(envs []EnvMap) Report {
	if len(envs) == 0 {
		return Report{}
	}

	labels := make([]string, len(envs))
	for i, e := range envs {
		labels[i] = e.Label
	}

	// Collect the union of all keys.
	keySet := map[string]struct{}{}
	for _, e := range envs {
		for k := range e.Vars {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	reports := make([]KeyReport, 0, len(keys))
	for _, key := range keys {
		kr := buildKeyReport(key, envs)
		reports = append(reports, kr)
	}

	return Report{Keys: reports, Labels: labels}
}

func buildKeyReport(key string, envs []EnvMap) KeyReport {
	kr := KeyReport{
		Key:    key,
		Values: make(map[string]string),
	}

	var seenValues []string
	for _, e := range envs {
		v, ok := e.Vars[key]
		if ok {
			kr.Values[e.Label] = v
			kr.PresentIn = append(kr.PresentIn, e.Label)
			seenValues = append(seenValues, v)
		} else {
			kr.MissingIn = append(kr.MissingIn, e.Label)
		}
	}

	kr.Consistent = allEqual(seenValues)
	return kr
}

func allEqual(vals []string) bool {
	if len(vals) == 0 {
		return true
	}
	for _, v := range vals[1:] {
		if v != vals[0] {
			return false
		}
	}
	return true
}

// ToDiffResults converts a Report into a flat slice of diff.Result so that
// existing filter/reporter pipelines can consume multi-env comparisons.
func ToDiffResults(r Report, base, target string) []diff.Result {
	var out []diff.Result
	for _, kr := range r.Keys {
		baseVal, baseOK := kr.Values[base]
		targetVal, targetOK := kr.Values[target]

		switch {
		case baseOK && !targetOK:
			out = append(out, diff.Result{Key: kr.Key, Left: baseVal, Status: diff.MissingInRight})
		case !baseOK && targetOK:
			out = append(out, diff.Result{Key: kr.Key, Right: targetVal, Status: diff.MissingInLeft})
		case baseOK && targetOK && baseVal != targetVal:
			out = append(out, diff.Result{Key: kr.Key, Left: baseVal, Right: targetVal, Status: diff.Mismatched})
		default:
			out = append(out, diff.Result{Key: kr.Key, Left: baseVal, Right: targetVal, Status: diff.Match})
		}
	}
	return out
}
