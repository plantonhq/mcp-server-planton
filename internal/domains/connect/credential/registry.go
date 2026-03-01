package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	auth0credentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/auth0credential/v1"
	awscredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/awscredential/v1"
	azurecredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/azurecredential/v1"
	civocredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/civocredential/v1"
	cloudflarecredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/cloudflarecredential/v1"
	confluentcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/confluentcredential/v1"
	digitaloceancredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/digitaloceancredential/v1"
	dockercredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/dockercredential/v1"
	gcpcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/gcpcredential/v1"
	githubcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/githubcredential/v1"
	gitlabcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/gitlabcredential/v1"
	kubernetesclustercredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/kubernetesclustercredential/v1"
	mavencredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/mavencredential/v1"
	mongodbatlascredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/mongodbatlascredential/v1"
	npmcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/npmcredential/v1"
	openfgacredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/openfgacredential/v1"
	pulumibackendcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/pulumibackendcredential/v1"
	snowflakecredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/snowflakecredential/v1"
	terraformbackendcredentialv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/connect/terraformbackendcredential/v1"
)

type applyFn func(ctx context.Context, conn *grpc.ClientConn, credentialObject map[string]any) (proto.Message, error)
type getFn func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)
type getByOrgBySlugFn func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error)
type deleteFn func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)

// credentialDispatcher binds a credential kind to its per-type gRPC clients
// and defines which spec fields contain secrets.
type credentialDispatcher struct {
	apply           applyFn
	get             getFn
	getByOrgBySlug  getByOrgBySlugFn
	del             deleteFn
	sensitiveFields []string
}

// dispatchers maps PascalCase credential kind to its dispatcher.
var dispatchers = map[string]credentialDispatcher{
	"Auth0Credential":             auth0Dispatcher(),
	"AwsCredential":               awsDispatcher(),
	"AzureCredential":             azureDispatcher(),
	"CivoCredential":              civoDispatcher(),
	"CloudflareCredential":        cloudflareDispatcher(),
	"ConfluentCredential":         confluentDispatcher(),
	"DigitalOceanCredential":      digitalOceanDispatcher(),
	"DockerCredential":            dockerDispatcher(),
	"GcpCredential":               gcpDispatcher(),
	"GithubCredential":            githubDispatcher(),
	"GitlabCredential":            gitlabDispatcher(),
	"KubernetesClusterCredential": kubernetesClusterDispatcher(),
	"MavenCredential":             mavenDispatcher(),
	"MongodbAtlasCredential":      mongodbAtlasDispatcher(),
	"NpmCredential":               npmDispatcher(),
	"OpenFgaCredential":           openFgaDispatcher(),
	"PulumiBackendCredential":     pulumiBackendDispatcher(),
	"SnowflakeCredential":         snowflakeDispatcher(),
	"TerraformBackendCredential":  terraformBackendDispatcher(),
}

// supportedKinds returns a comma-separated list of all registered credential
// kinds, for use in error messages.
func supportedKinds() string {
	names := make([]string, 0, len(dispatchers))
	for k := range dispatchers {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// unmarshalCredentialObject JSON-encodes a credential_object map and decodes
// it into a typed proto message via protojson. This bridges the generic MCP
// input to the per-type gRPC message, correctly handling enums, nested
// messages, and proto field name conventions.
func unmarshalCredentialObject(credentialObject map[string]any, msg proto.Message) error {
	jsonBytes, err := json.Marshal(credentialObject)
	if err != nil {
		return fmt.Errorf("encoding credential object: %w", err)
	}
	if err := protojson.Unmarshal(jsonBytes, msg); err != nil {
		return fmt.Errorf("invalid credential object: %w", err)
	}
	return nil
}

// extractKind reads the "kind" field from a credential_object map.
func extractKind(credentialObject map[string]any) (string, error) {
	v, ok := credentialObject["kind"]
	if !ok {
		return "", fmt.Errorf("credential_object is missing required 'kind' field")
	}
	kind, ok := v.(string)
	if !ok || kind == "" {
		return "", fmt.Errorf("credential_object 'kind' must be a non-empty string")
	}
	return kind, nil
}

// ---------------------------------------------------------------------------
// Auth0
// ---------------------------------------------------------------------------

func auth0Dispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred auth0credentialv1.Auth0Credential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return auth0credentialv1.NewAuth0CredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return auth0credentialv1.NewAuth0CredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return auth0credentialv1.NewAuth0CredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return auth0credentialv1.NewAuth0CredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.client_secret",
		},
	}
}

// ---------------------------------------------------------------------------
// AWS
// ---------------------------------------------------------------------------

func awsDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred awscredentialv1.AwsCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return awscredentialv1.NewAwsCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return awscredentialv1.NewAwsCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return awscredentialv1.NewAwsCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return awscredentialv1.NewAwsCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.secret_access_key",
			"spec.session_token",
		},
	}
}

// ---------------------------------------------------------------------------
// GCP
// ---------------------------------------------------------------------------

func gcpDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred gcpcredentialv1.GcpCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return gcpcredentialv1.NewGcpCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gcpcredentialv1.NewGcpCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return gcpcredentialv1.NewGcpCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gcpcredentialv1.NewGcpCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.service_account_key_base64",
			"spec.refresh_token",
		},
	}
}

// ---------------------------------------------------------------------------
// GitHub
// ---------------------------------------------------------------------------

func githubDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred githubcredentialv1.GithubCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return githubcredentialv1.NewGithubCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return githubcredentialv1.NewGithubCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return githubcredentialv1.NewGithubCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return githubcredentialv1.NewGithubCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: nil,
	}
}

// ---------------------------------------------------------------------------
// Azure
// ---------------------------------------------------------------------------

func azureDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred azurecredentialv1.AzureCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return azurecredentialv1.NewAzureCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return azurecredentialv1.NewAzureCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return azurecredentialv1.NewAzureCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return azurecredentialv1.NewAzureCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.client_secret",
			"spec.refresh_token",
		},
	}
}

// ---------------------------------------------------------------------------
// Kubernetes Cluster
// ---------------------------------------------------------------------------

func kubernetesClusterDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred kubernetesclustercredentialv1.KubernetesClusterCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return kubernetesclustercredentialv1.NewKubernetesClusterCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return kubernetesclustercredentialv1.NewKubernetesClusterCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return kubernetesclustercredentialv1.NewKubernetesClusterCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return kubernetesclustercredentialv1.NewKubernetesClusterCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.gcp_gke.service_account_key_base64",
			"spec.gcp_gke.cluster_ca_data",
			"spec.digital_ocean_doks.kube_config",
		},
	}
}

// ---------------------------------------------------------------------------
// Civo
// ---------------------------------------------------------------------------

func civoDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred civocredentialv1.CivoCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return civocredentialv1.NewCivoCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return civocredentialv1.NewCivoCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return civocredentialv1.NewCivoCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return civocredentialv1.NewCivoCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.api_token",
			"spec.object_store_secret_key",
		},
	}
}

// ---------------------------------------------------------------------------
// Cloudflare
// ---------------------------------------------------------------------------

func cloudflareDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred cloudflarecredentialv1.CloudflareCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return cloudflarecredentialv1.NewCloudflareCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return cloudflarecredentialv1.NewCloudflareCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return cloudflarecredentialv1.NewCloudflareCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return cloudflarecredentialv1.NewCloudflareCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.api_token",
			"spec.api_key",
			"spec.r2.secret_access_key",
		},
	}
}

// ---------------------------------------------------------------------------
// Confluent
// ---------------------------------------------------------------------------

func confluentDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred confluentcredentialv1.ConfluentCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return confluentcredentialv1.NewConfluentCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return confluentcredentialv1.NewConfluentCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return confluentcredentialv1.NewConfluentCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return confluentcredentialv1.NewConfluentCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.api_secret",
		},
	}
}

// ---------------------------------------------------------------------------
// DigitalOcean
// ---------------------------------------------------------------------------

func digitalOceanDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred digitaloceancredentialv1.DigitalOceanCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return digitaloceancredentialv1.NewDigitalOceanCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return digitaloceancredentialv1.NewDigitalOceanCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return digitaloceancredentialv1.NewDigitalOceanCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return digitaloceancredentialv1.NewDigitalOceanCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.api_token",
			"spec.spaces_secret_key",
		},
	}
}

// ---------------------------------------------------------------------------
// Docker
// ---------------------------------------------------------------------------

func dockerDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred dockercredentialv1.DockerCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return dockercredentialv1.NewDockerCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return dockercredentialv1.NewDockerCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return dockercredentialv1.NewDockerCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return dockercredentialv1.NewDockerCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.gcp_artifact_registry.service_account_key_base64",
			"spec.aws_elastic_container_registry.secret_access_key",
			"spec.github_container_registry.personal_access_token",
		},
	}
}

// ---------------------------------------------------------------------------
// GitLab
// ---------------------------------------------------------------------------

func gitlabDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred gitlabcredentialv1.GitlabCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return gitlabcredentialv1.NewGitlabCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gitlabcredentialv1.NewGitlabCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return gitlabcredentialv1.NewGitlabCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return gitlabcredentialv1.NewGitlabCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.access_token",
			"spec.refresh_token",
		},
	}
}

// ---------------------------------------------------------------------------
// Maven
// ---------------------------------------------------------------------------

func mavenDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred mavencredentialv1.MavenCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return mavencredentialv1.NewMavenCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return mavencredentialv1.NewMavenCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return mavencredentialv1.NewMavenCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return mavencredentialv1.NewMavenCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.gcp_artifact_registry.service_account_key_base64",
			"spec.github_packages.personal_access_token",
		},
	}
}

// ---------------------------------------------------------------------------
// MongoDB Atlas
// ---------------------------------------------------------------------------

func mongodbAtlasDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred mongodbatlascredentialv1.MongodbAtlasCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return mongodbatlascredentialv1.NewMongodbAtlasCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return mongodbatlascredentialv1.NewMongodbAtlasCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return mongodbatlascredentialv1.NewMongodbAtlasCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return mongodbatlascredentialv1.NewMongodbAtlasCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.private_key",
		},
	}
}

// ---------------------------------------------------------------------------
// NPM
// ---------------------------------------------------------------------------

func npmDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred npmcredentialv1.NpmCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return npmcredentialv1.NewNpmCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return npmcredentialv1.NewNpmCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return npmcredentialv1.NewNpmCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return npmcredentialv1.NewNpmCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.gcp_artifact_registry.service_account_key_base64",
			"spec.github_packages.personal_access_token",
		},
	}
}

// ---------------------------------------------------------------------------
// OpenFGA
// ---------------------------------------------------------------------------

func openFgaDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred openfgacredentialv1.OpenFgaCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return openfgacredentialv1.NewOpenFgaCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return openfgacredentialv1.NewOpenFgaCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return openfgacredentialv1.NewOpenFgaCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return openfgacredentialv1.NewOpenFgaCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.api_token",
			"spec.client_secret",
		},
	}
}

// ---------------------------------------------------------------------------
// Pulumi Backend
// ---------------------------------------------------------------------------

func pulumiBackendDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred pulumibackendcredentialv1.PulumiBackendCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return pulumibackendcredentialv1.NewPulumiBackendCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return pulumibackendcredentialv1.NewPulumiBackendCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return pulumibackendcredentialv1.NewPulumiBackendCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return pulumibackendcredentialv1.NewPulumiBackendCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.http.access_token",
			"spec.s3.aws_secret_access_key",
			"spec.gcs.service_account_key_base64",
			"spec.azurerm.storage_account_key",
			"spec.secrets_passphrase",
		},
	}
}

// ---------------------------------------------------------------------------
// Snowflake
// ---------------------------------------------------------------------------

func snowflakeDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred snowflakecredentialv1.SnowflakeCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return snowflakecredentialv1.NewSnowflakeCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return snowflakecredentialv1.NewSnowflakeCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return snowflakecredentialv1.NewSnowflakeCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return snowflakecredentialv1.NewSnowflakeCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.password",
		},
	}
}

// ---------------------------------------------------------------------------
// Terraform Backend
// ---------------------------------------------------------------------------

func terraformBackendDispatcher() credentialDispatcher {
	return credentialDispatcher{
		apply: func(ctx context.Context, conn *grpc.ClientConn, co map[string]any) (proto.Message, error) {
			var cred terraformbackendcredentialv1.TerraformBackendCredential
			if err := unmarshalCredentialObject(co, &cred); err != nil {
				return nil, err
			}
			return terraformbackendcredentialv1.NewTerraformBackendCredentialCommandControllerClient(conn).Apply(ctx, &cred)
		},
		get: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return terraformbackendcredentialv1.NewTerraformBackendCredentialQueryControllerClient(conn).
				Get(ctx, &apiresource.ApiResourceId{Value: id})
		},
		getByOrgBySlug: func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error) {
			return terraformbackendcredentialv1.NewTerraformBackendCredentialQueryControllerClient(conn).
				GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{Org: org, Slug: slug})
		},
		del: func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error) {
			return terraformbackendcredentialv1.NewTerraformBackendCredentialCommandControllerClient(conn).
				Delete(ctx, &apiresource.ApiResourceId{Value: id})
		},
		sensitiveFields: []string{
			"spec.s3.aws_secret_access_key",
			"spec.gcs.service_account_key_base64",
			"spec.azurerm.client_secret",
		},
	}
}
