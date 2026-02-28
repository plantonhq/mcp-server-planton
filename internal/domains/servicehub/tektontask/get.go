package tektontask

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	tektontaskv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektontask/v1"
	"google.golang.org/grpc"
)

// Get retrieves a Tekton task by ID or by org+slug via the
// TektonTaskQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(ApiResourceId) directly.
//   - Slug path: calls GetByOrgAndName(GetByOrgAndNameInput) with the slug
//     passed as the Name field (the server converts name to slug internally).
func Get(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			task, err := resolveTask(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(task)
		})
}

// resolveTask fetches the full TektonTask proto by ID or by org+slug.
func resolveTask(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (*tektontaskv1.TektonTask, error) {
	client := tektontaskv1.NewTektonTaskQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("Tekton task %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgAndName(ctx, &tektontaskv1.GetByOrgAndNameInput{
		Org:  org,
		Name: slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("Tekton task %q in org %q", slug, org))
	}
	return resp, nil
}

// describeTask returns a human-readable description of the Tekton task for
// use in error messages.
func describeTask(id, org, slug string) string {
	if id != "" {
		return fmt.Sprintf("Tekton task %q", id)
	}
	return fmt.Sprintf("Tekton task %q in org %q", slug, org)
}
