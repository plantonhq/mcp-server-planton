package search

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	searchapiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/apiresource"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// search_api_resources_by_text
// ---------------------------------------------------------------------------

type SearchByTextInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier to search within. Use list_organizations to discover available organizations."`
	Env        string `json:"env,omitempty"          jsonschema:"Environment slug to narrow results. Omit to search across all environments."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text query to match against resource names, descriptions, and indexed fields."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchByTextTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_api_resources_by_text",
		Description: "Full-text search across all API resources within an organization. " +
			"Returns paginated search records with resource IDs, names, kinds, and metadata. " +
			"Optionally filter by environment. Use this for broad discovery when you don't know the resource kind.",
	}
}

func SearchByTextHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchByTextInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchByTextInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		req := &searchapiresource.SearchByTextInput{
			Org:        input.Org,
			Env:        input.Env,
			SearchText: input.SearchText,
			PageInfo:   buildPageInfo(input.PageNum, input.PageSize),
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchapiresource.NewApiResourceSearchQueryControllerClient(conn).SearchByText(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("text search in org %q", input.Org))
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
// search_api_resources_by_kind
// ---------------------------------------------------------------------------

type SearchByKindInput struct {
	Org             string `json:"org"                    jsonschema:"required,Organization identifier to search within."`
	Env             string `json:"env,omitempty"           jsonschema:"Environment slug to narrow results. Omit to search across all environments."`
	ApiResourceKind string `json:"api_resource_kind"      jsonschema:"required,The API resource kind to filter by (snake_case, e.g. 'organization', 'environment', 'kafka_cluster')."`
	SearchText      string `json:"search_text,omitempty"   jsonschema:"Free-text query to further filter results within the specified kind."`
	PageNum         int32  `json:"page_num,omitempty"      jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize        int32  `json:"page_size,omitempty"     jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchByKindTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_api_resources_by_kind",
		Description: "Search API resources of a specific kind within an organization. " +
			"Use this when you know the resource type (e.g. 'kafka_cluster', 'postgres_cluster') " +
			"and want to find instances. Returns paginated search records.",
	}
}

func SearchByKindHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchByKindInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchByKindInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.ApiResourceKind == "" {
			return nil, nil, fmt.Errorf("'api_resource_kind' is required")
		}
		kind, err := domains.ResolveApiResourceKind(input.ApiResourceKind)
		if err != nil {
			return nil, nil, err
		}
		req := &searchapiresource.SearchApiResourcesByKindInput{
			Org:             input.Org,
			Env:             input.Env,
			ApiResourceKind: kind,
			SearchText:      input.SearchText,
			PageInfo:        buildPageInfo(input.PageNum, input.PageSize),
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchapiresource.NewApiResourceSearchQueryControllerClient(conn).SearchByKind(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("kind search %q in org %q", input.ApiResourceKind, input.Org))
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
// get_api_resource_by_org_kind_name
// ---------------------------------------------------------------------------

type GetByOrgKindNameInput struct {
	Org             string `json:"org"                jsonschema:"required,Organization identifier."`
	ApiResourceKind string `json:"api_resource_kind"  jsonschema:"required,The API resource kind (snake_case, e.g. 'kafka_cluster')."`
	Name            string `json:"name"               jsonschema:"required,The resource name (slug) to look up within the organization."`
}

func GetByOrgKindNameTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_api_resource_by_org_kind_name",
		Description: "Look up a single API resource by its organization, kind, and name. " +
			"Returns the search record for the exact match or an error if not found. " +
			"Use this for precise lookups when you know all three identifiers.",
	}
}

func GetByOrgKindNameHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetByOrgKindNameInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetByOrgKindNameInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.ApiResourceKind == "" {
			return nil, nil, fmt.Errorf("'api_resource_kind' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		kind, err := domains.ResolveApiResourceKind(input.ApiResourceKind)
		if err != nil {
			return nil, nil, err
		}
		req := &searchapiresource.GetByOrgByKindByNameRequest{
			Org:             input.Org,
			ApiResourceKind: kind,
			Name:            input.Name,
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchapiresource.NewApiResourceSearchQueryControllerClient(conn).GetByOrgByKindByName(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("%s %q in org %q", input.ApiResourceKind, input.Name, input.Org))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
