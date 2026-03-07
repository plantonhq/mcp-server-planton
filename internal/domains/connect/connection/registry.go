package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	auth0connectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/auth0providerconnection/v1"
	atlasconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/atlasproviderconnection/v1"
	awsconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/awsproviderconnection/v1"
	azureconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/azureproviderconnection/v1"
	civoconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/civoproviderconnection/v1"
	cloudflareconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/cloudflareproviderconnection/v1"
	cfworkersr2v1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/cloudflareworkerscriptsr2connection/v1"
	confluentconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/confluentproviderconnection/v1"
	containerregistryv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/containerregistryconnection/v1"
	digitaloceanconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/digitaloceanproviderconnection/v1"
	gcpconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/gcpproviderconnection/v1"
	githubconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/githubconnection/v1"
	gitlabconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/gitlabconnection/v1"
	k8sconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/kubernetesproviderconnection/v1"
	mavenconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/mavenconnection/v1"
	npmconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/npmconnection/v1"
	openfgaconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/openfgaproviderconnection/v1"
	pulumiconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/pulumibackendconnection/v1"
	snowflakeconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/snowflakeproviderconnection/v1"
	terraformconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/terraformbackendconnection/v1"
)

type applyFn func(ctx context.Context, conn *grpc.ClientConn, connectionObject map[string]any) (proto.Message, error)
type getFn func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)
type getByOrgBySlugFn func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error)
type deleteFn func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)

// connectionDispatcher binds a connection kind to its per-type gRPC clients.
type connectionDispatcher struct {
	apply          applyFn
	get            getFn
	getByOrgBySlug getByOrgBySlugFn
	del            deleteFn
}

// dispatchers maps PascalCase connection kind to its dispatcher.
var dispatchers = map[string]connectionDispatcher{
	"Auth0ProviderConnection":               auth0Dispatcher(),
	"AtlasProviderConnection":               atlasDispatcher(),
	"AwsProviderConnection":                 awsDispatcher(),
	"AzureProviderConnection":               azureDispatcher(),
	"CivoProviderConnection":                civoDispatcher(),
	"CloudflareProviderConnection":          cloudflareDispatcher(),
	"CloudflareWorkerScriptsR2Connection":   cloudflareWorkerScriptsR2Dispatcher(),
	"ConfluentProviderConnection":           confluentDispatcher(),
	"ContainerRegistryConnection":           containerRegistryDispatcher(),
	"DigitalOceanProviderConnection":        digitalOceanDispatcher(),
	"GcpProviderConnection":                 gcpDispatcher(),
	"GithubConnection":                      githubDispatcher(),
	"GitlabConnection":                      gitlabDispatcher(),
	"KubernetesProviderConnection":          kubernetesDispatcher(),
	"MavenConnection":                       mavenDispatcher(),
	"NpmConnection":                         npmDispatcher(),
	"OpenFgaProviderConnection":             openFgaDispatcher(),
	"PulumiBackendConnection":               pulumiBackendDispatcher(),
	"SnowflakeProviderConnection":           snowflakeDispatcher(),
	"TerraformBackendConnection":            terraformBackendDispatcher(),
}

// supportedKinds returns a comma-separated list of all registered connection
// kinds, for use in error messages.
func supportedKinds() string {
	names := make([]string, 0, len(dispatchers))
	for k := range dispatchers {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// unmarshalConnectionObject JSON-encodes a connection_object map and decodes
// it into a typed proto message via protojson. This bridges the generic MCP
// input to the per-type gRPC message, correctly handling enums, nested
// messages, and proto field name conventions.
func unmarshalConnectionObject(connectionObject map[string]any, msg proto.Message) error {
	jsonBytes, err := json.Marshal(connectionObject)
	if err != nil {
		return fmt.Errorf("encoding connection object: %w", err)
	}
	if err := protojson.Unmarshal(jsonBytes, msg); err != nil {
		return fmt.Errorf("invalid connection object: %w", err)
	}
	return nil
}

// extractKind reads the "kind" field from a connection_object map.
func extractKind(connectionObject map[string]any) (string, error) {
	v, ok := connectionObject["kind"]
	if !ok {
		return "", fmt.Errorf("connection_object is missing required 'kind' field")
	}
	kind, ok := v.(string)
	if !ok || kind == "" {
		return "", fmt.Errorf("connection_object 'kind' must be a non-empty string")
	}
	return kind, nil
}

// ---------------------------------------------------------------------------
// Auth0
// ---------------------------------------------------------------------------

func auth0Dispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg auth0connectionv1.Auth0ProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return auth0connectionv1.NewAuth0ProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return auth0connectionv1.NewAuth0ProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return auth0connectionv1.NewAuth0ProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return auth0connectionv1.NewAuth0ProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Atlas (MongoDB Atlas)
// ---------------------------------------------------------------------------

func atlasDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg atlasconnectionv1.AtlasProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return atlasconnectionv1.NewAtlasProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return atlasconnectionv1.NewAtlasProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return atlasconnectionv1.NewAtlasProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return atlasconnectionv1.NewAtlasProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// AWS
// ---------------------------------------------------------------------------

func awsDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg awsconnectionv1.AwsProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return awsconnectionv1.NewAwsProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return awsconnectionv1.NewAwsProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return awsconnectionv1.NewAwsProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return awsconnectionv1.NewAwsProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Azure
// ---------------------------------------------------------------------------

func azureDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg azureconnectionv1.AzureProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return azureconnectionv1.NewAzureProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return azureconnectionv1.NewAzureProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return azureconnectionv1.NewAzureProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return azureconnectionv1.NewAzureProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Civo
// ---------------------------------------------------------------------------

func civoDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg civoconnectionv1.CivoProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return civoconnectionv1.NewCivoProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return civoconnectionv1.NewCivoProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return civoconnectionv1.NewCivoProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return civoconnectionv1.NewCivoProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Cloudflare
// ---------------------------------------------------------------------------

func cloudflareDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg cloudflareconnectionv1.CloudflareProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return cloudflareconnectionv1.NewCloudflareProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return cloudflareconnectionv1.NewCloudflareProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return cloudflareconnectionv1.NewCloudflareProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return cloudflareconnectionv1.NewCloudflareProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Cloudflare Worker Scripts R2
// ---------------------------------------------------------------------------

func cloudflareWorkerScriptsR2Dispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg cfworkersr2v1.CloudflareWorkerScriptsR2Connection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return cfworkersr2v1.NewCloudflareWorkerScriptsR2ConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return cfworkersr2v1.NewCloudflareWorkerScriptsR2ConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return cfworkersr2v1.NewCloudflareWorkerScriptsR2ConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return cfworkersr2v1.NewCloudflareWorkerScriptsR2ConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Confluent
// ---------------------------------------------------------------------------

func confluentDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg confluentconnectionv1.ConfluentProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return confluentconnectionv1.NewConfluentProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return confluentconnectionv1.NewConfluentProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return confluentconnectionv1.NewConfluentProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return confluentconnectionv1.NewConfluentProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Container Registry (was Docker)
// ---------------------------------------------------------------------------

func containerRegistryDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg containerregistryv1.ContainerRegistryConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return containerregistryv1.NewContainerRegistryConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return containerregistryv1.NewContainerRegistryConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return containerregistryv1.NewContainerRegistryConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return containerregistryv1.NewContainerRegistryConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// DigitalOcean
// ---------------------------------------------------------------------------

func digitalOceanDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg digitaloceanconnectionv1.DigitalOceanProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return digitaloceanconnectionv1.NewDigitalOceanProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return digitaloceanconnectionv1.NewDigitalOceanProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return digitaloceanconnectionv1.NewDigitalOceanProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return digitaloceanconnectionv1.NewDigitalOceanProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// GCP
// ---------------------------------------------------------------------------

func gcpDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg gcpconnectionv1.GcpProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return gcpconnectionv1.NewGcpProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gcpconnectionv1.NewGcpProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return gcpconnectionv1.NewGcpProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gcpconnectionv1.NewGcpProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// GitHub
// ---------------------------------------------------------------------------

func githubDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg githubconnectionv1.GithubConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return githubconnectionv1.NewGithubConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return githubconnectionv1.NewGithubConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return githubconnectionv1.NewGithubConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return githubconnectionv1.NewGithubConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// GitLab
// ---------------------------------------------------------------------------

func gitlabDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg gitlabconnectionv1.GitlabConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return gitlabconnectionv1.NewGitlabConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gitlabconnectionv1.NewGitlabConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return gitlabconnectionv1.NewGitlabConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gitlabconnectionv1.NewGitlabConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Kubernetes
// ---------------------------------------------------------------------------

func kubernetesDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg k8sconnectionv1.KubernetesProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return k8sconnectionv1.NewKubernetesProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return k8sconnectionv1.NewKubernetesProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return k8sconnectionv1.NewKubernetesProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return k8sconnectionv1.NewKubernetesProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Maven
// ---------------------------------------------------------------------------

func mavenDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg mavenconnectionv1.MavenConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return mavenconnectionv1.NewMavenConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return mavenconnectionv1.NewMavenConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return mavenconnectionv1.NewMavenConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return mavenconnectionv1.NewMavenConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// NPM
// ---------------------------------------------------------------------------

func npmDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg npmconnectionv1.NpmConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return npmconnectionv1.NewNpmConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return npmconnectionv1.NewNpmConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return npmconnectionv1.NewNpmConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return npmconnectionv1.NewNpmConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// OpenFGA
// ---------------------------------------------------------------------------

func openFgaDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg openfgaconnectionv1.OpenFgaProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return openfgaconnectionv1.NewOpenFgaProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return openfgaconnectionv1.NewOpenFgaProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return openfgaconnectionv1.NewOpenFgaProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return openfgaconnectionv1.NewOpenFgaProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Pulumi Backend
// ---------------------------------------------------------------------------

func pulumiBackendDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg pulumiconnectionv1.PulumiBackendConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return pulumiconnectionv1.NewPulumiBackendConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return pulumiconnectionv1.NewPulumiBackendConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return pulumiconnectionv1.NewPulumiBackendConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return pulumiconnectionv1.NewPulumiBackendConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Snowflake
// ---------------------------------------------------------------------------

func snowflakeDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg snowflakeconnectionv1.SnowflakeProviderConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return snowflakeconnectionv1.NewSnowflakeProviderConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return snowflakeconnectionv1.NewSnowflakeProviderConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return snowflakeconnectionv1.NewSnowflakeProviderConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return snowflakeconnectionv1.NewSnowflakeProviderConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}

// ---------------------------------------------------------------------------
// Terraform Backend
// ---------------------------------------------------------------------------

func terraformBackendDispatcher() connectionDispatcher {
	return connectionDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var msg terraformconnectionv1.TerraformBackendConnection
			if err := unmarshalConnectionObject(co, &msg); err != nil {
				return nil, err
			}
			return terraformconnectionv1.NewTerraformBackendConnectionCommandControllerClient(conn).Apply(ctx, &msg)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return terraformconnectionv1.NewTerraformBackendConnectionQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return terraformconnectionv1.NewTerraformBackendConnectionQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return terraformconnectionv1.NewTerraformBackendConnectionCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
	}
}
