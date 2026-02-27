package domains

import (
	"sort"
	"strings"
)

// JoinEnumValues returns a sorted, comma-separated list of the map's keys,
// excluding the specified zero-value key (e.g. "unspecified" sentinel).
//
// This is used by domain-level enum resolvers to produce user-friendly error
// messages that list all valid values when a lookup fails.
func JoinEnumValues(m map[string]int32, exclude string) string {
	vals := make([]string, 0, len(m)-1)
	for k := range m {
		if k != exclude {
			vals = append(vals, k)
		}
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}
