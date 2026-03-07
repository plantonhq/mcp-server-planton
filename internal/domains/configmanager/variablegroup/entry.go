package variablegroup

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"

	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// upsert_variable_group_entry
// ---------------------------------------------------------------------------

// UpsertEntryInput defines the parameters for the upsert_variable_group_entry tool.
type UpsertEntryInput struct {
	GroupID     string `json:"group_id"               jsonschema:"required,Variable group ID."`
	Name        string `json:"name"                   jsonschema:"required,Entry name (e.g. DATABASE_HOST). Must be unique within the group."`
	Value       string `json:"value"                  jsonschema:"required,Entry value."`
	Description string `json:"description,omitempty"   jsonschema:"Optional description of what this entry is used for."`
}

// UpsertEntryTool returns the MCP tool definition for upsert_variable_group_entry.
func UpsertEntryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "upsert_variable_group_entry",
		Description: "Add or update a single entry in a variable group. " +
			"If an entry with the same name already exists, it is updated; otherwise a new entry is added. " +
			"Returns the full updated variable group.",
	}
}

// UpsertEntryHandler returns the typed tool handler for upsert_variable_group_entry.
func UpsertEntryHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpsertEntryInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpsertEntryInput) (*mcp.CallToolResult, any, error) {
		if input.GroupID == "" {
			return nil, nil, fmt.Errorf("'group_id' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := variablegroupv1.NewVariableGroupCommandControllerClient(conn)
				resp, err := client.UpsertEntry(ctx, &variablegroupv1.UpsertVariableGroupEntryRequest{
					GroupId: input.GroupID,
					Entry: &variablegroupv1.VariableGroupEntry{
						Name:        input.Name,
						Value:       input.Value,
						Description: input.Description,
					},
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("entry %q in variable group %q", input.Name, input.GroupID))
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
// delete_variable_group_entry
// ---------------------------------------------------------------------------

// DeleteEntryInput defines the parameters for the delete_variable_group_entry tool.
type DeleteEntryInput struct {
	GroupID   string `json:"group_id"    jsonschema:"required,Variable group ID."`
	EntryName string `json:"entry_name"  jsonschema:"required,Name of the entry to remove."`
}

// DeleteEntryTool returns the MCP tool definition for delete_variable_group_entry.
func DeleteEntryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_variable_group_entry",
		Description: "Remove a single entry from a variable group by name. " +
			"Returns the full updated variable group.",
	}
}

// DeleteEntryHandler returns the typed tool handler for delete_variable_group_entry.
func DeleteEntryHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteEntryInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteEntryInput) (*mcp.CallToolResult, any, error) {
		if input.GroupID == "" {
			return nil, nil, fmt.Errorf("'group_id' is required")
		}
		if input.EntryName == "" {
			return nil, nil, fmt.Errorf("'entry_name' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := variablegroupv1.NewVariableGroupCommandControllerClient(conn)
				resp, err := client.DeleteEntry(ctx, &variablegroupv1.DeleteVariableGroupEntryRequest{
					GroupId:   input.GroupID,
					EntryName: input.EntryName,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("entry %q in variable group %q", input.EntryName, input.GroupID))
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
// refresh_variable_group_entry
// ---------------------------------------------------------------------------

// RefreshEntryInput defines the parameters for the refresh_variable_group_entry tool.
type RefreshEntryInput struct {
	GroupID   string `json:"group_id"    jsonschema:"required,Variable group ID."`
	EntryName string `json:"entry_name"  jsonschema:"required,Name of the entry to refresh from its source."`
}

// RefreshEntryTool returns the MCP tool definition for refresh_variable_group_entry.
func RefreshEntryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "refresh_variable_group_entry",
		Description: "Refresh a single entry's value from its external source reference. " +
			"Fails if the entry has no source configured. " +
			"Returns the full updated variable group.",
	}
}

// RefreshEntryHandler returns the typed tool handler for refresh_variable_group_entry.
func RefreshEntryHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RefreshEntryInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RefreshEntryInput) (*mcp.CallToolResult, any, error) {
		if input.GroupID == "" {
			return nil, nil, fmt.Errorf("'group_id' is required")
		}
		if input.EntryName == "" {
			return nil, nil, fmt.Errorf("'entry_name' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := variablegroupv1.NewVariableGroupCommandControllerClient(conn)
				resp, err := client.RefreshEntry(ctx, &variablegroupv1.RefreshVariableGroupEntryRequest{
					GroupId:   input.GroupID,
					EntryName: input.EntryName,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("entry %q in variable group %q", input.EntryName, input.GroupID))
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
// refresh_all_variable_group_entries
// ---------------------------------------------------------------------------

// RefreshAllInput defines the parameters for the refresh_all_variable_group_entries tool.
type RefreshAllInput struct {
	GroupID string `json:"group_id" jsonschema:"required,Variable group ID."`
}

// RefreshAllTool returns the MCP tool definition for refresh_all_variable_group_entries.
func RefreshAllTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "refresh_all_variable_group_entries",
		Description: "Refresh all entries in a variable group that have external source references. " +
			"Entries without a source are skipped. " +
			"Returns the full updated variable group.",
	}
}

// RefreshAllHandler returns the typed tool handler for refresh_all_variable_group_entries.
func RefreshAllHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RefreshAllInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RefreshAllInput) (*mcp.CallToolResult, any, error) {
		if input.GroupID == "" {
			return nil, nil, fmt.Errorf("'group_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := variablegroupv1.NewVariableGroupCommandControllerClient(conn)
				resp, err := client.RefreshAll(ctx, &variablegroupv1.VariableGroupId{Value: input.GroupID})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("variable group %q", input.GroupID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
