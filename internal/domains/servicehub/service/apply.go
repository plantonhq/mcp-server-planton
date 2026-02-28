package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// Apply creates or updates a service via the
// ServiceCommandController.Apply RPC.
//
// The input is a raw JSON map matching the Service proto shape. It is
// serialized to JSON bytes and deserialized into the typed proto using
// protojson, which handles proto field-name conventions, enums, and
// well-known types correctly.
func Apply(ctx context.Context, serverAddress string, raw map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize service input: %w", err)
	}

	svc := &servicev1.Service{}
	if err := protojson.Unmarshal(jsonBytes, svc); err != nil {
		return "", fmt.Errorf("invalid service structure: %w", err)
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := servicev1.NewServiceCommandControllerClient(conn)
			result, err := client.Apply(ctx, svc)
			if err != nil {
				desc := "service"
				if md := svc.GetMetadata(); md != nil {
					name := md.GetName()
					if name == "" {
						name = md.GetSlug()
					}
					if name != "" {
						desc = fmt.Sprintf("service %q in org %q", name, md.GetOrg())
					}
				}
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
