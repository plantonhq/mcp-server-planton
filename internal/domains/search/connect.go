package search

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	searchconnect "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/connect"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// search_connections_by_context
// ---------------------------------------------------------------------------

type SearchConnectionsInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier."`
	Env        string `json:"env,omitempty"          jsonschema:"Environment slug to narrow results."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text query to filter connections by name."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchConnectionsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_connections_by_context",
		Description: "Search connection resources (cloud provider connections) within an organization. " +
			"Connections link Planton to external cloud providers (AWS, GCP, Azure, etc.). " +
			"Optionally filter by environment. Returns paginated search records.",
	}
}

func SearchConnectionsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchConnectionsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchConnectionsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		req := &searchconnect.SearchConnectionApiResourcesByContext{
			Org:        input.Org,
			Env:        input.Env,
			SearchText: input.SearchText,
			PageInfo:   buildPageInfo(input.PageNum, input.PageSize),
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchconnect.NewConnectSearchQueryControllerClient(conn).SearchConnectionApiResourcesByContext(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("connections in org %q", input.Org))
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
// get_connections_by_env
// ---------------------------------------------------------------------------

type GetConnectionsByEnvInput struct {
	Org string `json:"org" jsonschema:"required,Organization identifier."`
	Env string `json:"env" jsonschema:"required,Environment slug to retrieve connections for."`
}

func GetConnectionsByEnvTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_connections_by_env",
		Description: "Get all connections associated with a specific environment, grouped by connection kind. " +
			"Returns connections organized by their type (aws_connection, gcp_connection, etc.) " +
			"with display names and search records for each.",
	}
}

func GetConnectionsByEnvHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetConnectionsByEnvInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetConnectionsByEnvInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Env == "" {
			return nil, nil, fmt.Errorf("'env' is required")
		}
		req := &searchconnect.GetConnectionsByEnvInput{
			Org: input.Org,
			Env: input.Env,
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchconnect.NewConnectSearchQueryControllerClient(conn).GetConnectionsByEnv(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("connections in env %q (org %q)", input.Env, input.Org))
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
// search_runner_registrations_by_org
// ---------------------------------------------------------------------------

type SearchRunnerRegistrationsInput struct {
	Org                        string `json:"org"                             jsonschema:"required,Organization identifier."`
	SearchText                 string `json:"search_text,omitempty"           jsonschema:"Free-text query to filter runners by name."`
	IncludePlatformRunners     bool   `json:"include_platform_runners,omitempty"      jsonschema:"Include platform-managed runners pre-deployed in Planton infrastructure."`
	IncludeOrganizationRunners bool   `json:"include_organization_runners,omitempty"  jsonschema:"Include organization runners deployed by customers. Defaults to true when both flags are false."`
	PageNum                    int32  `json:"page_num,omitempty"              jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize                   int32  `json:"page_size,omitempty"             jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchRunnerRegistrationsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_runner_registrations_by_org",
		Description: "Search runner registrations within an organization. " +
			"Runners execute infrastructure operations (deploy, destroy, etc.). " +
			"By default returns organization runners only. Set include_platform_runners " +
			"to also include Planton-managed runners.",
	}
}

func SearchRunnerRegistrationsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchRunnerRegistrationsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchRunnerRegistrationsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		req := &searchconnect.SearchRunnerRegistrationsByOrgContextInput{
			Org:                        input.Org,
			SearchText:                 input.SearchText,
			PageInfo:                   buildPageInfo(input.PageNum, input.PageSize),
			IncludePlatformRunners:     input.IncludePlatformRunners,
			IncludeOrganizationRunners: input.IncludeOrganizationRunners,
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchconnect.NewConnectSearchQueryControllerClient(conn).SearchRunnerRegistrationsByOrgContext(ctx, req)
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
