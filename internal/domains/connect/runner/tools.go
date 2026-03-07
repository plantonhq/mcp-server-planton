package runner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	runnerregistrationv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/runnerregistration/v1"
	connectsearch "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/connect"
)

// ---------------------------------------------------------------------------
// apply_runner_registration
// ---------------------------------------------------------------------------

type ApplyInput struct {
	RegistrationObject map[string]any `json:"registration_object" jsonschema:"required,Full RunnerRegistration object in OpenMCF envelope format."`
}

func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_runner_registration",
		Description: "Create or update a runner registration. " +
			"Runners are compute agents deployed in customer environments that execute infrastructure operations. " +
			"Pass the full RunnerRegistration object as an OpenMCF envelope.",
	}
}

func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyInput) (*mcp.CallToolResult, any, error) {
		if input.RegistrationObject == nil {
			return nil, nil, fmt.Errorf("'registration_object' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				var rr runnerregistrationv1.RunnerRegistration
				jsonBytes, err := json.Marshal(input.RegistrationObject)
				if err != nil {
					return "", fmt.Errorf("encoding registration object: %w", err)
				}
				if err := protojson.Unmarshal(jsonBytes, &rr); err != nil {
					return "", fmt.Errorf("invalid registration object: %w", err)
				}
				client := runnerregistrationv1.NewRunnerRegistrationCommandControllerClient(conn)
				resp, err := client.Apply(ctx, &rr)
				if err != nil {
					return "", domains.RPCError(err, "runner registration")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_runner_registration
// ---------------------------------------------------------------------------

type GetInput struct {
	ID string `json:"id" jsonschema:"required,RunnerRegistration ID."`
}

func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_runner_registration",
		Description: "Retrieve a runner registration by ID. Returns the full registration including status and connectivity info.",
	}
}

func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := runnerregistrationv1.NewRunnerRegistrationQueryControllerClient(conn)
				resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("runner registration %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_runner_registration
// ---------------------------------------------------------------------------

type DeleteInput struct {
	ID string `json:"id" jsonschema:"required,RunnerRegistration ID to delete."`
}

func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_runner_registration",
		Description: "Delete a runner registration by ID. " +
			"WARNING: Credentials using this runner will lose connectivity. " +
			"Ensure no credentials reference this runner before deleting.",
	}
}

func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := runnerregistrationv1.NewRunnerRegistrationCommandControllerClient(conn)
				resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("runner registration %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// generate_runner_credentials
// ---------------------------------------------------------------------------

type GenerateCredentialsInput struct {
	ID string `json:"id" jsonschema:"required,RunnerRegistration ID to generate credentials for."`
}

func GenerateCredentialsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "generate_runner_credentials",
		Description: "Generate initial authentication credentials for a runner registration. " +
			"SECURITY WARNING: The response contains sensitive cryptographic material including " +
			"private keys, certificates, and API keys. Handle the output with extreme care.",
	}
}

func GenerateCredentialsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GenerateCredentialsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GenerateCredentialsInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := runnerregistrationv1.NewRunnerRegistrationCommandControllerClient(conn)
				resp, err := client.GenerateCredentials(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("credentials for runner registration %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// regenerate_runner_credentials
// ---------------------------------------------------------------------------

type RegenerateCredentialsInput struct {
	ID string `json:"id" jsonschema:"required,RunnerRegistration ID to regenerate credentials for."`
}

func RegenerateCredentialsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "regenerate_runner_credentials",
		Description: "Rotate and regenerate authentication credentials for a runner registration. " +
			"This invalidates the previous credentials. The runner must be reconfigured with the new credentials. " +
			"SECURITY WARNING: The response contains sensitive cryptographic material including " +
			"private keys, certificates, and API keys. Handle the output with extreme care.",
	}
}

func RegenerateCredentialsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RegenerateCredentialsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RegenerateCredentialsInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := runnerregistrationv1.NewRunnerRegistrationCommandControllerClient(conn)
				resp, err := client.RegenerateCredentials(ctx, &apiresource.ApiResourceId{Value: input.ID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("credentials for runner registration %q", input.ID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// search_runner_registrations
// ---------------------------------------------------------------------------

type SearchInput struct {
	Org                        string `json:"org" jsonschema:"required,Organization ID to search within."`
	SearchText                 string `json:"search_text,omitempty" jsonschema:"Optional text to filter runners by name."`
	IncludePlatformRunners     bool   `json:"include_platform_runners,omitempty" jsonschema:"Include platform-managed runners in results (default: false)."`
	IncludeOrganizationRunners bool   `json:"include_organization_runners,omitempty" jsonschema:"Include organization runners in results (default: true if both flags are false)."`
}

func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_runner_registrations",
		Description: "Search runner registrations within an organization. " +
			"By default returns only organization runners. Set include_platform_runners=true to also see platform-managed runners.",
	}
}

func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := connectsearch.NewConnectSearchQueryControllerClient(conn)
				resp, err := client.SearchRunnerRegistrationsByOrgContext(ctx, &connectsearch.SearchRunnerRegistrationsByOrgContextInput{
					Org:                        input.Org,
					SearchText:                 input.SearchText,
					IncludePlatformRunners:     input.IncludePlatformRunners,
					IncludeOrganizationRunners: input.IncludeOrganizationRunners,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("runner registrations in org %q", input.Org))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
