package cloudresource

import (
	"fmt"

	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
)

// extractMetadata maps the cloud_object["metadata"] sub-map to an
// ApiResourceMetadata proto. Required fields: name, org, env. Optional
// fields are set when present and silently omitted otherwise.
func extractMetadata(cloudObject map[string]any) (*apiresource.ApiResourceMetadata, error) {
	raw, ok := cloudObject["metadata"]
	if !ok {
		return nil, fmt.Errorf("cloud_object missing required field \"metadata\"")
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cloud_object field \"metadata\" must be an object, got %T", raw)
	}

	md := &apiresource.ApiResourceMetadata{}

	name, _ := m["name"].(string)
	if name == "" {
		return nil, fmt.Errorf("metadata.name is required")
	}
	md.Name = name

	org, _ := m["org"].(string)
	if org == "" {
		return nil, fmt.Errorf("metadata.org is required")
	}
	md.Org = org

	env, _ := m["env"].(string)
	if env == "" {
		return nil, fmt.Errorf("metadata.env is required")
	}
	md.Env = env

	if slug, ok := m["slug"].(string); ok && slug != "" {
		md.Slug = slug
	}

	if id, ok := m["id"].(string); ok && id != "" {
		md.Id = id
	}

	if labels, ok := m["labels"].(map[string]any); ok {
		md.Labels = toStringMap(labels)
	}

	if annotations, ok := m["annotations"].(map[string]any); ok {
		md.Annotations = toStringMap(annotations)
	}

	if tags, ok := m["tags"].([]any); ok {
		md.Tags = toStringSlice(tags)
	}

	if ver, ok := m["version"].(map[string]any); ok {
		md.Version = &apiresource.ApiResourceMetadataVersion{}
		if msg, ok := ver["message"].(string); ok {
			md.Version.Message = msg
		}
	}

	return md, nil
}

func toStringMap(m map[string]any) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}

func toStringSlice(a []any) []string {
	out := make([]string, 0, len(a))
	for _, v := range a {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
