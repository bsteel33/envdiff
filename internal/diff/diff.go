package diff

// Result holds the comparison result between two env file maps.
type Result struct {
	MissingInRight []string          // keys present in left but not in right
	MissingInLeft  []string          // keys present in right but not in left
	Mismatched     map[string][2]string // keys present in both but with different values
}

// Compare takes two maps of env key-value pairs and returns a Result
// describing the differences between them.
func Compare(left, right map[string]string) Result {
	result := Result{
		Mismatched: make(map[string][2]string),
	}

	for key, leftVal := range left {
		rightVal, ok := right[key]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, key)
			continue
		}
		if leftVal != rightVal {
			result.Mismatched[key] = [2]string{leftVal, rightVal}
		}
	}

	for key := range right {
		if _, ok := left[key]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, key)
		}
	}

	sortStrings(result.MissingInRight)
	sortStrings(result.MissingInLeft)

	return result
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInRight) > 0 ||
		len(r.MissingInLeft) > 0 ||
		len(r.Mismatched) > 0
}

// sortStrings sorts a slice of strings in place using a simple insertion sort.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
