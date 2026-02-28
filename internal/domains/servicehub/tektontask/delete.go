package tektontask

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"

	tektontaskv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektontask/v1"
)

// Delete removes a Tekton task template via the
// TektonTaskCommandController.Delete RPC.
//
// The delete RPC requires the full TektonTask entity as input, so this
// function first resolves the task by ID or org+slug via the query controller,
// then passes the resolved entity to Delete. Both calls share a single gRPC
// connection.
func Delete(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			task, err := resolveTask(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeTask(id, org, slug)
			client := tektontaskv1.NewTektonTaskCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, task)
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
