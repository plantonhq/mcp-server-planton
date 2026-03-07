package search

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	searchrm "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/resourcemanager"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// get_context_hierarchy
// ---------------------------------------------------------------------------

type GetContextHierarchyInput struct {
	SearchText string `json:"search_text,omitempty" jsonschema:"Free-text query to filter organizations and environments by name."`
}

func GetContextHierarchyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_context_hierarchy",
		Description: "Get the hierarchy of organizations and environments that the current user has access to. " +
			"Returns a tree of organizations, each containing their environments with IDs, names, and slugs. " +
			"Use this to discover available organizations and environments before other operations.",
	}
}

func GetContextHierarchyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetContextHierarchyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetContextHierarchyInput) (*mcp.CallToolResult, any, error) {
		req := &searchrm.GetContextHierarchyInput{
			SearchText: input.SearchText,
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchrm.NewResourceManagerSearchQueryControllerClient(conn).GetContextHierarchy(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "context hierarchy")
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
// search_quick_actions
// ---------------------------------------------------------------------------

type SearchQuickActionsInput struct {
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text query to filter quick actions by name or description."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchQuickActionsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_quick_actions",
		Description: "Search available quick actions. Quick actions are pre-defined operations " +
			"that can be executed against resources (e.g. deploy, restart, scale). " +
			"Returns paginated search records.",
	}
}

func SearchQuickActionsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchQuickActionsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchQuickActionsInput) (*mcp.CallToolResult, any, error) {
		req := &searchrm.SearchQuickActionsRequest{
			SearchText: input.SearchText,
			PageInfo:   buildPageInfo(input.PageNum, input.PageSize),
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchrm.NewResourceManagerSearchQueryControllerClient(conn).SearchQuickActions(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "quick actions")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
