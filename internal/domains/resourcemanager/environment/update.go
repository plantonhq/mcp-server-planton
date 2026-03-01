package environment

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	environmentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1"
	"google.golang.org/grpc"
)

// UpdateFields holds the optional fields that can be modified on an existing
// environment. Zero-value (empty string) means "leave unchanged".
type UpdateFields struct {
	Name        string
	Description string
}

// Update performs a read-modify-write on an existing environment.
//
// The function first fetches the current environment by ID via the query
// controller, applies any non-empty fields from UpdateFields, then writes the
// result back via the command controller's Update RPC. Both calls share a
// single gRPC connection within one WithConnection callback.
func Update(ctx context.Context, serverAddress, envID string, fields UpdateFields) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			queryClient := environmentv1.NewEnvironmentQueryControllerClient(conn)
			env, err := queryClient.Get(ctx, &environmentv1.EnvironmentId{Value: envID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment %q", envID))
			}

			if fields.Name != "" {
				env.Metadata.Name = fields.Name
			}
			if env.Spec == nil {
				env.Spec = &environmentv1.EnvironmentSpec{}
			}
			if fields.Description != "" {
				env.Spec.Description = fields.Description
			}

			cmdClient := environmentv1.NewEnvironmentCommandControllerClient(conn)
			resp, err := cmdClient.Update(ctx, env)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment %q", envID))
			}
			return domains.MarshalJSON(resp)
		})
}
