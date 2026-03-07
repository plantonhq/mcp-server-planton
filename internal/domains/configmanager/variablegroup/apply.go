package variablegroup

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Apply creates or updates a variable group via the envelope pattern.
func Apply(ctx context.Context, serverAddress string, groupObject map[string]any) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			var vg variablegroupv1.VariableGroup
			jsonBytes, err := json.Marshal(groupObject)
			if err != nil {
				return "", fmt.Errorf("encoding group object: %w", err)
			}
			if err := protojson.Unmarshal(jsonBytes, &vg); err != nil {
				return "", fmt.Errorf("invalid group object: %w", err)
			}
			client := variablegroupv1.NewVariableGroupCommandControllerClient(conn)
			resp, err := client.Apply(ctx, &vg)
			if err != nil {
				return "", domains.RPCError(err, "variable group")
			}
			return domains.MarshalJSON(resp)
		})
}
