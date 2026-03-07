package secretbackend

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	secretbackendv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/secretbackend/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Apply creates or updates a secret backend via the envelope pattern.
func Apply(ctx context.Context, serverAddress string, backendObject map[string]any) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			var sb secretbackendv1.SecretBackend
			jsonBytes, err := json.Marshal(backendObject)
			if err != nil {
				return "", fmt.Errorf("encoding backend object: %w", err)
			}
			if err := protojson.Unmarshal(jsonBytes, &sb); err != nil {
				return "", fmt.Errorf("invalid backend object: %w", err)
			}
			client := secretbackendv1.NewSecretBackendCommandControllerClient(conn)
			resp, err := client.Apply(ctx, &sb)
			if err != nil {
				return "", domains.RPCError(err, "secret backend")
			}
			RedactSecretBackend(resp)
			return domains.MarshalJSON(resp)
		})
}
