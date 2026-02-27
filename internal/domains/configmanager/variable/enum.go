package variable

import (
	"fmt"
	"sort"
	"strings"

	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
)

// resolveScope maps a user-supplied string ("organization" or "environment")
// to the corresponding VariableSpec_Scope enum value.
func resolveScope(s string) (variablev1.VariableSpec_Scope, error) {
	v, ok := variablev1.VariableSpec_Scope_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown scope %q — valid values: %s",
			s, joinScopeValues())
	}
	return variablev1.VariableSpec_Scope(v), nil
}

// joinScopeValues returns a sorted, comma-separated list of valid scope
// values, excluding the zero-value sentinel.
func joinScopeValues() string {
	vals := make([]string, 0, len(variablev1.VariableSpec_Scope_value)-1)
	for k := range variablev1.VariableSpec_Scope_value {
		if k != "scope_unspecified" {
			vals = append(vals, k)
		}
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}
