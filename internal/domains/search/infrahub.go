package search

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	searchinfrahub "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/infrahub"
	searchcloudresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/infrahub/cloudresource"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// search_iac_modules_by_org
// ---------------------------------------------------------------------------

type SearchIacModulesInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier to search within."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text query to filter IaC modules by name or description."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchIacModulesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_iac_modules_by_org",
		Description: "Search IaC (Infrastructure as Code) modules within an organization. " +
			"IaC modules are reusable infrastructure building blocks (Pulumi, Terraform, etc.) " +
			"that can be used by infra projects. Returns paginated search records.",
	}
}

func SearchIacModulesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchIacModulesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchIacModulesInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		req := &searchinfrahub.SearchIacModulesByOrgContextInput{
			Org:        input.Org,
			SearchText: input.SearchText,
			PageInfo:   buildPageInfo(input.PageNum, input.PageSize),
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchinfrahub.NewInfraHubSearchQueryControllerClient(conn).SearchIacModulesByOrgContext(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("IaC modules in org %q", input.Org))
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
// lookup_cloud_resource
// ---------------------------------------------------------------------------

type LookupCloudResourceInput struct {
	Org               string `json:"org"                  jsonschema:"required,Organization identifier."`
	Env               string `json:"env"                  jsonschema:"required,Environment slug."`
	CloudResourceKind string `json:"cloud_resource_kind"  jsonschema:"required,Cloud resource kind in PascalCase (e.g. 'PostgresCluster', 'RedisCluster', 'AwsEksCluster'). Use cloud-resource-kinds://catalog to list valid kinds."`
	Name              string `json:"name"                 jsonschema:"required,Name (slug) of the cloud resource."`
}

func LookupCloudResourceTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "lookup_cloud_resource",
		Description: "Look up a specific cloud resource by org, environment, kind, and name. " +
			"Returns the search record for the exact match. " +
			"Use this for precise lookups when you know all four identifiers.",
	}
}

func LookupCloudResourceHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *LookupCloudResourceInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *LookupCloudResourceInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Env == "" {
			return nil, nil, fmt.Errorf("'env' is required")
		}
		if input.CloudResourceKind == "" {
			return nil, nil, fmt.Errorf("'cloud_resource_kind' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		kind, err := domains.ResolveKind(input.CloudResourceKind)
		if err != nil {
			return nil, nil, err
		}
		req := &searchcloudresource.LookupCloudResourceInput{
			Org:               input.Org,
			Env:               input.Env,
			CloudResourceKind: kind,
			Name:              input.Name,
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchcloudresource.NewCloudResourceSearchQueryControllerClient(conn).LookupCloudResource(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("cloud resource %s %q in env %q (org %q)", input.CloudResourceKind, input.Name, input.Env, input.Org))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
