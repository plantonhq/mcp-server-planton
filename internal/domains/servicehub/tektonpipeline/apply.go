package tektonpipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	tektonpipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektonpipeline/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// Apply creates or updates a Tekton pipeline via the
// TektonPipelineCommandController.Apply RPC.
//
// The input is a raw JSON map matching the TektonPipeline proto shape. It is
// serialized to JSON bytes and deserialized into the typed proto using
// protojson, which handles proto field-name conventions, enums, and
// well-known types correctly.
func Apply(ctx context.Context, serverAddress string, raw map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Tekton pipeline input: %w", err)
	}

	pipeline := &tektonpipelinev1.TektonPipeline{}
	if err := protojson.Unmarshal(jsonBytes, pipeline); err != nil {
		return "", fmt.Errorf("invalid Tekton pipeline structure: %w", err)
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := tektonpipelinev1.NewTektonPipelineCommandControllerClient(conn)
			result, err := client.Apply(ctx, pipeline)
			if err != nil {
				desc := "Tekton pipeline"
				if md := pipeline.GetMetadata(); md != nil {
					name := md.GetName()
					if name == "" {
						name = md.GetSlug()
					}
					if name != "" {
						desc = fmt.Sprintf("Tekton pipeline %q in org %q", name, md.GetOrg())
					}
				}
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
