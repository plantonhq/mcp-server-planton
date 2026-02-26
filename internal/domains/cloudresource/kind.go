package cloudresource

import (
	"fmt"

	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
)

// extractKindFromCloudObject reads the "kind" field from a cloud_object map.
// Returns an error if the field is missing or not a string.
func extractKindFromCloudObject(co map[string]any) (string, error) {
	v, ok := co["kind"]
	if !ok {
		return "", fmt.Errorf("cloud_object missing required field \"kind\"")
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("cloud_object field \"kind\" must be a non-empty string")
	}
	return s, nil
}

// resolveKind maps a PascalCase kind string (e.g. "AwsEksCluster") to the
// corresponding CloudResourceKind enum value from the openmcf proto stubs.
func resolveKind(kindStr string) (cloudresourcekind.CloudResourceKind, error) {
	v, ok := cloudresourcekind.CloudResourceKind_value[kindStr]
	if !ok {
		return 0, fmt.Errorf("unknown cloud resource kind %q — use the cloud-resource-schema resource template to discover valid kinds", kindStr)
	}
	return cloudresourcekind.CloudResourceKind(v), nil
}
