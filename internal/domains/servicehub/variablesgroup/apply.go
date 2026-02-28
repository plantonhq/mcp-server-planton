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

// Apply creates or updates a variables group via the
// VariablesGroupCommandController.Apply RPC.
//
// The input is a raw JSON map matching the VariablesGroup proto shape. It is
// serialized to JSON bytes and deserialized into the typed proto using
// protojson, which handles proto field-name conventions, enums, and
// well-known types correctly.
func Apply(ctx context.Context, serverAddress string, raw map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize variables group input: %w", err)
	}

	group := &variablesgroupv1.VariablesGroup{}
	if err := protojson.Unmarshal(jsonBytes, group); err != nil {
		return "", fmt.Errorf("invalid variables group structure: %w", err)
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := variablesgroupv1.NewVariablesGroupCommandControllerClient(conn)
			result, err := client.Apply(ctx, group)
			if err != nil {
				desc := "variables group"
				if md := group.GetMetadata(); md != nil {
					name := md.GetName()
					if name == "" {
						name = md.GetSlug()
					}
					if name != "" {
						desc = fmt.Sprintf("variables group %q in org %q", name, md.GetOrg())
					}
				}
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
