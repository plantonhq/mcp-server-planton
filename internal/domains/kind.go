package domains

import (
	"fmt"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
)

// ResolveKind maps a PascalCase kind string (e.g. "AwsEksCluster") to the
// corresponding CloudResourceKind enum value from the openmcf proto stubs.
//
// This is intentionally in the shared domains package so that all domain
// packages (cloudresource, stackjob, preset, etc.) can resolve kinds without
// depending on each other.
func ResolveKind(kindStr string) (cloudresourcekind.CloudResourceKind, error) {
	v, ok := cloudresourcekind.CloudResourceKind_value[kindStr]
	if !ok {
		return 0, fmt.Errorf("unknown cloud resource kind %q — read cloud-resource-kinds://catalog for all valid kinds", kindStr)
	}
	return cloudresourcekind.CloudResourceKind(v), nil
}

var apiResourceKindResolver = NewEnumResolver[apiresourcekind.ApiResourceKind](
	apiresourcekind.ApiResourceKind_value,
	"API resource kind",
	"api_resource_kind_unspecified",
)

// ResolveApiResourceKind maps a snake_case kind string (e.g. "organization",
// "environment") to the corresponding ApiResourceKind enum value.
func ResolveApiResourceKind(s string) (apiresourcekind.ApiResourceKind, error) {
	return apiResourceKindResolver.Resolve(s)
}
