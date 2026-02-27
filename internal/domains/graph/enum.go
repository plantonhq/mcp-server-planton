package graph

import (
	"fmt"
	"sort"
	"strings"

	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
)

// resolveNodeTypes maps user-supplied strings (e.g. "cloud_resource",
// "credential") to the corresponding GraphNode_Type enum values.
// Returns all resolved values, or an error on the first unknown string.
func resolveNodeTypes(names []string) ([]graphv1.GraphNode_Type, error) {
	out := make([]graphv1.GraphNode_Type, 0, len(names))
	for _, name := range names {
		v, ok := graphv1.GraphNode_Type_value[name]
		if !ok {
			return nil, fmt.Errorf("unknown node type %q — valid values: %s",
				name, joinEnumValues(graphv1.GraphNode_Type_value, "type_unspecified"))
		}
		out = append(out, graphv1.GraphNode_Type(v))
	}
	return out, nil
}

// resolveRelationshipTypes maps user-supplied strings (e.g. "depends_on",
// "uses_credential") to the corresponding GraphRelationship_Type enum values.
// Returns all resolved values, or an error on the first unknown string.
func resolveRelationshipTypes(names []string) ([]graphv1.GraphRelationship_Type, error) {
	out := make([]graphv1.GraphRelationship_Type, 0, len(names))
	for _, name := range names {
		v, ok := graphv1.GraphRelationship_Type_value[name]
		if !ok {
			return nil, fmt.Errorf("unknown relationship type %q — valid values: %s",
				name, joinEnumValues(graphv1.GraphRelationship_Type_value, "type_unspecified"))
		}
		out = append(out, graphv1.GraphRelationship_Type(v))
	}
	return out, nil
}

// resolveChangeType maps a user-supplied string ("delete" or "update") to the
// corresponding GetImpactAnalysisInput_ChangeType enum value.
func resolveChangeType(s string) (graphv1.GetImpactAnalysisInput_ChangeType, error) {
	v, ok := graphv1.GetImpactAnalysisInput_ChangeType_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown change type %q — valid values: %s",
			s, joinEnumValues(graphv1.GetImpactAnalysisInput_ChangeType_value, "change_type_unspecified"))
	}
	return graphv1.GetImpactAnalysisInput_ChangeType(v), nil
}

// joinEnumValues returns a sorted, comma-separated list of the map's keys,
// excluding the specified zero-value key (e.g. "type_unspecified" sentinel).
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
