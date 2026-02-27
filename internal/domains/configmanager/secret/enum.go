package secret

import (
	"fmt"
	"sort"
	"strings"

	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
)

// resolveScope maps a user-supplied string ("organization" or "environment")
// to the corresponding SecretSpec_Scope enum value.
func resolveScope(s string) (secretv1.SecretSpec_Scope, error) {
	v, ok := secretv1.SecretSpec_Scope_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown scope %q — valid values: %s",
			s, joinScopeValues())
	}
	return secretv1.SecretSpec_Scope(v), nil
}

// joinScopeValues returns a sorted, comma-separated list of valid scope
// values, excluding the zero-value sentinel.
func joinScopeValues() string {
	vals := make([]string, 0, len(secretv1.SecretSpec_Scope_value)-1)
	for k := range secretv1.SecretSpec_Scope_value {
		if k != "scope_unspecified" {
			vals = append(vals, k)
		}
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}
