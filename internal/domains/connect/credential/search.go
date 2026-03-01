package credential

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	connectsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/connect"
	"google.golang.org/grpc"
)

// credentialKindToAPIResourceKind maps PascalCase credential kinds used in
// the MCP tool interface to the ApiResourceKind enum values used by the
// search RPC.
var credentialKindToAPIResourceKind = map[string]apiresourcekind.ApiResourceKind{
	"Auth0Credential":             apiresourcekind.ApiResourceKind_auth0_credential,
	"AwsCredential":               apiresourcekind.ApiResourceKind_aws_credential,
	"AzureCredential":             apiresourcekind.ApiResourceKind_azure_credential,
	"CivoCredential":              apiresourcekind.ApiResourceKind_civo_credential,
	"CloudflareCredential":        apiresourcekind.ApiResourceKind_cloudflare_credential,
	"ConfluentCredential":         apiresourcekind.ApiResourceKind_confluent_credential,
	"DigitalOceanCredential":      apiresourcekind.ApiResourceKind_digital_ocean_credential,
	"DockerCredential":            apiresourcekind.ApiResourceKind_docker_credential,
	"GcpCredential":               apiresourcekind.ApiResourceKind_gcp_credential,
	"GithubCredential":            apiresourcekind.ApiResourceKind_github_credential,
	"GitlabCredential":            apiresourcekind.ApiResourceKind_gitlab_credential,
	"KubernetesClusterCredential": apiresourcekind.ApiResourceKind_kubernetes_cluster_credential,
	"MavenCredential":             apiresourcekind.ApiResourceKind_maven_credential,
	"MongodbAtlasCredential":      apiresourcekind.ApiResourceKind_mongodb_atlas_credential,
	"NpmCredential":               apiresourcekind.ApiResourceKind_npm_credential,
	"OpenFgaCredential":           apiresourcekind.ApiResourceKind_open_fga_credential,
	"PulumiBackendCredential":     apiresourcekind.ApiResourceKind_pulumi_backend_credential,
	"SnowflakeCredential":         apiresourcekind.ApiResourceKind_snowflake_credential,
	"TerraformBackendCredential":  apiresourcekind.ApiResourceKind_terraform_backend_credential,
}

// Search queries the ConnectSearchQueryController for credentials visible
// to the caller within an organization, optionally filtered by environment,
// credential kinds, and free-text search.
func Search(ctx context.Context, serverAddress, org, env string, kinds []string, searchText string) (string, error) {
	apiKinds, err := resolveAPIResourceKinds(kinds)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := connectsearch.NewConnectSearchQueryControllerClient(conn)
			resp, err := client.SearchCredentialApiResourcesByContext(ctx, &connectsearch.SearchCredentialApiResourcesByContext{
				Org:        org,
				Env:        env,
				Kinds:      apiKinds,
				SearchText: searchText,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("credentials in org %q", org))
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
		v, ok := credentialKindToAPIResourceKind[k]
		if !ok {
			return nil, fmt.Errorf("unknown credential kind %q — valid values: %s", k, validCredentialKindNames())
		}
		out = append(out, v)
	}
	return out, nil
}

func validCredentialKindNames() string {
	names := make([]string, 0, len(credentialKindToAPIResourceKind))
	for k := range credentialKindToAPIResourceKind {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
