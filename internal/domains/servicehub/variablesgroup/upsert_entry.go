package variablesgroup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// UpsertEntry adds or updates a single variable in a variables group via the
// VariablesGroupCommandController.UpsertEntry RPC.
//
// The target group can be identified by group_id directly, or by org+slug
// (which triggers an extra lookup to resolve the group ID).
//
// The entry is provided as a raw JSON map and converted to the typed proto
// using protojson, supporting both literal values and value_from references.
func UpsertEntry(ctx context.Context, serverAddress, groupID, org, groupSlug string, entryRaw map[string]any) (string, error) {
	entryBytes, err := json.Marshal(entryRaw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize entry: %w", err)
	}

	entry := &variablesgroupv1.VariablesGroupEntry{}
	if err := protojson.Unmarshal(entryBytes, entry); err != nil {
		return "", fmt.Errorf("invalid entry structure: %w", err)
	}

	if entry.GetName() == "" {
		return "", fmt.Errorf("entry 'name' is required")
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resolvedID, err := resolveGroupID(ctx, conn, groupID, org, groupSlug)
			if err != nil {
				return "", err
			}

			client := variablesgroupv1.NewVariablesGroupCommandControllerClient(conn)
			result, err := client.UpsertEntry(ctx, &variablesgroupv1.UpsertVariableRequest{
				GroupId: resolvedID,
				Entry:   entry,
			})
			if err != nil {
				desc := describeGroup(groupID, org, groupSlug)
				return "", domains.RPCError(err, fmt.Sprintf("upsert variable %q in %s", entry.GetName(), desc))
			}
			return domains.MarshalJSON(result)
		})
}
