package tektonpipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	tektonpipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektonpipeline/v1"
	"google.golang.org/grpc"
)

// Get retrieves a Tekton pipeline by ID or by org+slug via the
// TektonPipelineQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(ApiResourceId) directly.
//   - Slug path: calls GetByOrgAndName(GetByOrgAndNameInput) with the slug
//     passed as the Name field (the server converts name to slug internally).
func Get(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			pipeline, err := resolvePipeline(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(pipeline)
		})
}

// resolvePipeline fetches the full TektonPipeline proto by ID or by org+slug.
func resolvePipeline(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (*tektonpipelinev1.TektonPipeline, error) {
	client := tektonpipelinev1.NewTektonPipelineQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("Tekton pipeline %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgAndName(ctx, &tektonpipelinev1.GetByOrgAndNameInput{
		Org:  org,
		Name: slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("Tekton pipeline %q in org %q", slug, org))
	}
	return resp, nil
}

// describePipeline returns a human-readable description of the Tekton
// pipeline for use in error messages.
func describePipeline(id, org, slug string) string {
	if id != "" {
		return fmt.Sprintf("Tekton pipeline %q", id)
	}
	return fmt.Sprintf("Tekton pipeline %q in org %q", slug, org)
}
