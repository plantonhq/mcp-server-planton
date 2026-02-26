// Package parse provides shared utilities used by generated cloud resource
// input types. These helpers handle the common cloud_object envelope
// (api_version, kind, metadata, spec) so that generated per-provider code
// only needs to deal with spec validation and normalization.
package parse

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

// ValidateHeader checks that the cloud_object map contains the expected
// api_version and kind values. Returns an error if either is missing or
// does not match.
func ValidateHeader(cloudObject map[string]any, expectedAPIVersion, expectedKind string) error {
	av, _ := cloudObject["api_version"].(string)
	if av == "" {
		return fmt.Errorf("cloud_object missing required field api_version")
	}
	if av != expectedAPIVersion {
		return fmt.Errorf("expected api_version %q, got %q", expectedAPIVersion, av)
	}

	k, _ := cloudObject["kind"].(string)
	if k == "" {
		return fmt.Errorf("cloud_object missing required field kind")
	}
	if k != expectedKind {
		return fmt.Errorf("expected kind %q, got %q", expectedKind, k)
	}

	return nil
}

// ExtractSpecMap extracts the "spec" sub-map from a cloud_object.
// Returns an error if spec is missing or not a map.
func ExtractSpecMap(cloudObject map[string]any) (map[string]any, error) {
	specRaw, ok := cloudObject["spec"]
	if !ok {
		return nil, fmt.Errorf("cloud_object missing required field spec")
	}
	specMap, ok := specRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cloud_object field spec must be an object, got %T", specRaw)
	}
	return specMap, nil
}

// RebuildCloudObject creates a new cloud_object map by copying all top-level
// fields from the original and replacing the "spec" field with the normalized
// spec map. The result is converted to a structpb.Struct ready for use as
// CloudResource.Spec.CloudObject.
func RebuildCloudObject(original map[string]any, normalizedSpec map[string]any) (*structpb.Struct, error) {
	result := make(map[string]any, len(original))
	for k, v := range original {
		result[k] = v
	}
	result["spec"] = normalizedSpec
	return structpb.NewStruct(result)
}
