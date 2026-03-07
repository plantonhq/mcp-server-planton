package cloudops

import (
	"fmt"

	cloudopspb "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// BuildContext constructs a CloudOpsRequestContext from tool input parameters.
//
// Two mutually exclusive access modes are supported:
//
//   - Cloud resource mode (cloudResourceKind + cloudResourceSlug + env):
//     the control plane resolves the provider connection from the cloud resource spec.
//
//   - Provider connection mode (connection, optionally with env):
//     the caller specifies the credential slug directly.
//
// At least one access mode must be provided. If both are given, an error is returned.
func BuildContext(org, env, cloudResourceKind, cloudResourceSlug, connection string) (*cloudopspb.CloudOpsRequestContext, error) {
	if org == "" {
		return nil, fmt.Errorf("'org' is required")
	}

	hasCloudResource := cloudResourceKind != "" || cloudResourceSlug != ""
	hasConnection := connection != ""

	if hasCloudResource && hasConnection {
		return nil, fmt.Errorf("provide either cloud resource identifiers (cloud_resource_kind + cloud_resource_slug) or a connection slug — not both")
	}

	ctx := &cloudopspb.CloudOpsRequestContext{Org: org}

	switch {
	case hasCloudResource:
		if cloudResourceKind == "" {
			return nil, fmt.Errorf("'cloud_resource_kind' is required when using cloud resource access mode")
		}
		if cloudResourceSlug == "" {
			return nil, fmt.Errorf("'cloud_resource_slug' is required when using cloud resource access mode")
		}
		if env == "" {
			return nil, fmt.Errorf("'env' is required when using cloud resource access mode")
		}
		kind, err := domains.ResolveKind(cloudResourceKind)
		if err != nil {
			return nil, err
		}
		ctx.AccessMode = &cloudopspb.CloudOpsRequestContext_CloudResource{
			CloudResource: &cloudopspb.CloudResourceAccess{
				Env:  env,
				Kind: kind,
				Slug: cloudResourceSlug,
			},
		}

	case hasConnection:
		ctx.AccessMode = &cloudopspb.CloudOpsRequestContext_ProviderConnection{
			ProviderConnection: &cloudopspb.ProviderConnectionAccess{
				Env:        env,
				Connection: connection,
			},
		}

	default:
		return nil, fmt.Errorf("provide either cloud resource identifiers (cloud_resource_kind + cloud_resource_slug + env) or a connection slug")
	}

	return ctx, nil
}
