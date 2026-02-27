package cloudresource

import (
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
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

// resolveKinds maps a slice of PascalCase kind strings to their corresponding
// CloudResourceKind enum values. Returns an error on the first unknown kind.
// A nil or empty input returns a nil slice (no filtering).
func resolveKinds(kindStrs []string) ([]cloudresourcekind.CloudResourceKind, error) {
	if len(kindStrs) == 0 {
		return nil, nil
	}
	kinds := make([]cloudresourcekind.CloudResourceKind, len(kindStrs))
	for i, s := range kindStrs {
		k, err := domains.ResolveKind(s)
		if err != nil {
			return nil, err
		}
		kinds[i] = k
	}
	return kinds, nil
}
