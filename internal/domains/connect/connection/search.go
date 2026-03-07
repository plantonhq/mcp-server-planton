package connection

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
	connectsearch "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/connect"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// connectionKindToAPIResourceKind maps PascalCase connection kinds used in the
// MCP tool interface to the ApiResourceKind enum values used by the search RPC.
var connectionKindToAPIResourceKind = map[string]apiresourcekind.ApiResourceKind{
	"Auth0ProviderConnection":             apiresourcekind.ApiResourceKind_auth0_provider_connection,
	"AtlasProviderConnection":             apiresourcekind.ApiResourceKind_atlas_provider_connection,
	"AwsProviderConnection":               apiresourcekind.ApiResourceKind_aws_provider_connection,
	"AzureProviderConnection":             apiresourcekind.ApiResourceKind_azure_provider_connection,
	"CivoProviderConnection":              apiresourcekind.ApiResourceKind_civo_provider_connection,
	"CloudflareProviderConnection":        apiresourcekind.ApiResourceKind_cloudflare_provider_connection,
	"CloudflareWorkerScriptsR2Connection": apiresourcekind.ApiResourceKind_cloudflare_worker_scripts_r2_connection,
	"ConfluentProviderConnection":         apiresourcekind.ApiResourceKind_confluent_provider_connection,
	"ContainerRegistryConnection":         apiresourcekind.ApiResourceKind_container_registry_connection,
	"DigitalOceanProviderConnection":      apiresourcekind.ApiResourceKind_digital_ocean_provider_connection,
	"GcpProviderConnection":               apiresourcekind.ApiResourceKind_gcp_provider_connection,
	"GithubConnection":                    apiresourcekind.ApiResourceKind_github_connection,
	"GitlabConnection":                    apiresourcekind.ApiResourceKind_gitlab_connection,
	"KubernetesProviderConnection":        apiresourcekind.ApiResourceKind_kubernetes_provider_connection,
	"MavenConnection":                     apiresourcekind.ApiResourceKind_maven_connection,
	"NpmConnection":                       apiresourcekind.ApiResourceKind_npm_connection,
	"OpenFgaProviderConnection":           apiresourcekind.ApiResourceKind_open_fga_provider_connection,
	"PulumiBackendConnection":             apiresourcekind.ApiResourceKind_pulumi_backend_connection,
	"SnowflakeProviderConnection":         apiresourcekind.ApiResourceKind_snowflake_provider_connection,
	"TerraformBackendConnection":          apiresourcekind.ApiResourceKind_terraform_backend_connection,
}

// Search queries the ConnectSearchQueryController for connections visible to
// the caller within an organization, optionally filtered by environment,
// connection kinds, and free-text search.
func Search(ctx context.Context, serverAddress, org, env string, kinds []string, searchText string) (string, error) {
	apiKinds, err := resolveAPIResourceKinds(kinds)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := connectsearch.NewConnectSearchQueryControllerClient(conn)
			resp, err := client.SearchConnectionApiResourcesByContext(ctx, &connectsearch.SearchConnectionApiResourcesByContext{
				Org:        org,
				Env:        env,
				Kinds:      apiKinds,
				SearchText: searchText,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("connections in org %q", org))
			}
			return domains.MarshalJSON(resp)
		})
}

func resolveAPIResourceKinds(kinds []string) ([]apiresourcekind.ApiResourceKind, error) {
	if len(kinds) == 0 {
		return nil, nil
	}
	out := make([]apiresourcekind.ApiResourceKind, 0, len(kinds))
	for _, k := range kinds {
		v, ok := connectionKindToAPIResourceKind[k]
		if !ok {
			return nil, fmt.Errorf("unknown connection kind %q — valid values: %s", k, validConnectionKindNames())
		}
		out = append(out, v)
	}
	return out, nil
}

func validConnectionKindNames() string {
	names := make([]string, 0, len(connectionKindToAPIResourceKind))
	for k := range connectionKindToAPIResourceKind {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
