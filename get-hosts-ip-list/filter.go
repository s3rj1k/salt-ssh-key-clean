package main

import (
	"sort"
)

// FilterStringSlice is used to deduplicate slice of strings and remove empty elements.
// https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable
func FilterStringSlice(s []string) []string {
	if len(s) == 0 {
		return nil
	}

	sort.Slice(
		s, func(i, j int) bool {
			return s[i] < s[j]
		},
	)

	j := 0

	for i := 1; i < len(s); i++ {
		if s[i] == s[j] {
			continue
		}

		j++

		s[j] = s[i]
	}

	if len(s) > 0 && s[0] == "" { // remove empty element
		return s[1 : j+1]
	}

	return s[:j+1]
}
