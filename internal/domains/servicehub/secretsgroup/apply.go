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

// Apply creates or updates a secrets group via the
// SecretsGroupCommandController.Apply RPC.
//
// The input is a raw JSON map matching the SecretsGroup proto shape. It is
// serialized to JSON bytes and deserialized into the typed proto using
// protojson, which handles proto field-name conventions, enums, and
// well-known types correctly.
func Apply(ctx context.Context, serverAddress string, raw map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize secrets group input: %w", err)
	}

	group := &secretsgroupv1.SecretsGroup{}
	if err := protojson.Unmarshal(jsonBytes, group); err != nil {
		return "", fmt.Errorf("invalid secrets group structure: %w", err)
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretsgroupv1.NewSecretsGroupCommandControllerClient(conn)
			result, err := client.Apply(ctx, group)
			if err != nil {
				desc := "secrets group"
				if md := group.GetMetadata(); md != nil {
					name := md.GetName()
					if name == "" {
						name = md.GetSlug()
					}
					if name != "" {
						desc = fmt.Sprintf("secrets group %q in org %q", name, md.GetOrg())
					}
				}
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
