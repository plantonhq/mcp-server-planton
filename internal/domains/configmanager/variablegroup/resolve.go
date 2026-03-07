package variablegroup

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"

	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ResolveEntryInput defines the parameters for the resolve_variable_group_entry tool.
type ResolveEntryInput struct {
	Org       string `json:"org"        jsonschema:"required,Organization identifier."`
	Scope     string `json:"scope"      jsonschema:"required,Group scope: 'organization' or 'environment'."`
	Slug      string `json:"slug"       jsonschema:"required,Variable group slug."`
	EntryName string `json:"entry_name" jsonschema:"required,Name of the entry whose value to resolve."`
}

// ResolveEntryTool returns the MCP tool definition for resolve_variable_group_entry.
func ResolveEntryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_variable_group_entry",
		Description: "Quick value lookup for a single entry in a variable group. " +
			"Returns just the plain string value, not the full resource. " +
			"Identifies the group by org+scope+slug and the entry by name.",
	}
}

// ResolveEntryHandler returns the typed tool handler for resolve_variable_group_entry.
func ResolveEntryHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveEntryInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveEntryInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Scope == "" {
			return nil, nil, fmt.Errorf("'scope' is required")
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		if input.EntryName == "" {
			return nil, nil, fmt.Errorf("'entry_name' is required")
		}
		scope, err := scopeResolver.Resolve(input.Scope)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := variablegroupv1.NewVariableGroupQueryControllerClient(conn)
				resp, err := client.Resolve(ctx, &variablegroupv1.ResolveVariableGroupEntryRequest{
					Org:       input.Org,
					Scope:     scope,
					Slug:      input.Slug,
					EntryName: input.EntryName,
				})
				if err != nil {
					return "", domains.RPCError(err,
						fmt.Sprintf("entry %q in variable group %q (scope=%s) in org %q",
							input.EntryName, input.Slug, input.Scope, input.Org))
				}
				return resp.GetValue(), nil
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
