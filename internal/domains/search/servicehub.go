package search

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	searchservicehub "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/search/v1/servicehub"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// search_infra_charts_by_org
// ---------------------------------------------------------------------------

type SearchInfraChartsInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier to search within."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text query to filter infra charts by name or description."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

func SearchInfraChartsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_infra_charts_by_org",
		Description: "Search infra charts within an organization. " +
			"Infra charts are reusable infrastructure templates (similar to Helm charts) " +
			"that define deployable cloud resource compositions. " +
			"Returns both official Planton charts and organization-owned charts.",
	}
}

func SearchInfraChartsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchInfraChartsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchInfraChartsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		req := &searchservicehub.SearchInfraChartsByOrgContextInput{
			Org:        input.Org,
			SearchText: input.SearchText,
			PageInfo:   buildPageInfo(input.PageNum, input.PageSize),
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := searchservicehub.NewServiceHubSearchQueryControllerClient(conn).SearchInfraChartsByOrgContext(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("infra charts in org %q", input.Org))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
