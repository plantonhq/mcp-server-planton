package secretsgroup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretsgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/secretsgroup/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// UpsertEntry adds or updates a single secret in a secrets group via the
// SecretsGroupCommandController.UpsertEntry RPC.
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

	entry := &secretsgroupv1.SecretsGroupEntry{}
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

			client := secretsgroupv1.NewSecretsGroupCommandControllerClient(conn)
			result, err := client.UpsertEntry(ctx, &secretsgroupv1.UpsertSecretRequest{
				GroupId: resolvedID,
				Entry:   entry,
			})
			if err != nil {
				desc := describeGroup(groupID, org, groupSlug)
				return "", domains.RPCError(err, fmt.Sprintf("upsert secret %q in %s", entry.GetName(), desc))
			}
			return domains.MarshalJSON(result)
		})
}
