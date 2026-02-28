package secretsgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretsgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/secretsgroup/v1"
	"google.golang.org/grpc"
)

// DeleteEntry removes a single secret from a secrets group via the
// SecretsGroupCommandController.DeleteEntry RPC.
//
// The target group can be identified by group_id directly, or by org+slug
// (which triggers an extra lookup to resolve the group ID).
func DeleteEntry(ctx context.Context, serverAddress, groupID, org, groupSlug, entryName string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resolvedID, err := resolveGroupID(ctx, conn, groupID, org, groupSlug)
			if err != nil {
				return "", err
			}

			client := secretsgroupv1.NewSecretsGroupCommandControllerClient(conn)
			result, err := client.DeleteEntry(ctx, &secretsgroupv1.DeleteSecretRequest{
				GroupId:   resolvedID,
				EntryName: entryName,
			})
			if err != nil {
				desc := describeGroup(groupID, org, groupSlug)
				return "", domains.RPCError(err, fmt.Sprintf("delete secret %q from %s", entryName, desc))
			}
			return domains.MarshalJSON(result)
		})
}
