package audit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
)

// resolveApiResourceKind maps a user-supplied string (e.g. "cloud_resource",
// "infra_project") to the corresponding ApiResourceKind enum value.
//
// All valid enum values are accepted — the backend will reject kinds that
// don't support versioning. The error message lists valid values to guide
// the agent.
func resolveApiResourceKind(s string) (apiresourcekind.ApiResourceKind, error) {
	v, ok := apiresourcekind.ApiResourceKind_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown resource kind %q — valid values: %s",
			s, joinEnumValues(apiresourcekind.ApiResourceKind_value, "unspecified"))
	}
	return apiresourcekind.ApiResourceKind(v), nil
}

// joinEnumValues returns a sorted, comma-separated list of the map's keys,
// excluding the specified zero-value key (e.g. "unspecified" sentinel).
func joinEnumValues(m map[string]int32, exclude string) string {
	vals := make([]string, 0, len(m)-1)
	for k := range m {
		if k != exclude {
			vals = append(vals, k)
		}
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}
