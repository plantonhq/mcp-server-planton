package audit

import (
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
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
			s, domains.JoinEnumValues(apiresourcekind.ApiResourceKind_value, "unspecified"))
	}
	return apiresourcekind.ApiResourceKind(v), nil
}
