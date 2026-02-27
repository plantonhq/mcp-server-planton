package infraproject

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	infraprojectv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infraproject/v1"
	"google.golang.org/grpc"
)

// Undeploy tears down all cloud resources associated with an infra project
// via the InfraProjectCommandController.Undeploy RPC while keeping the
// project record on the platform.
//
// Two identification paths are supported:
//   - ID path: calls Undeploy directly with the given ID.
//   - Slug path: first resolves org+slug to a project ID via the query
//     controller, then calls Undeploy. Both calls share a single gRPC
//     connection.
//
// The returned InfraProject includes an updated status.pipeline_id pointing
// to the triggered undeploy pipeline.
func Undeploy(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveProjectID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeProject(id, org, slug)
			client := infraprojectv1.NewInfraProjectCommandControllerClient(conn)
			result, err := client.Undeploy(ctx, &infraprojectv1.InfraProjectId{Value: resourceID})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
