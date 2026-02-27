package domains

import (
	"fmt"

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
