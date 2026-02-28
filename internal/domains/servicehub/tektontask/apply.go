package tektontask

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	tektontaskv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektontask/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// Apply creates or updates a Tekton task via the
// TektonTaskCommandController.Apply RPC.
//
// The input is a raw JSON map matching the TektonTask proto shape. It is
// serialized to JSON bytes and deserialized into the typed proto using
// protojson, which handles proto field-name conventions, enums, and
// well-known types correctly.
func Apply(ctx context.Context, serverAddress string, raw map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize Tekton task input: %w", err)
	}

	task := &tektontaskv1.TektonTask{}
	if err := protojson.Unmarshal(jsonBytes, task); err != nil {
		return "", fmt.Errorf("invalid Tekton task structure: %w", err)
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := tektontaskv1.NewTektonTaskCommandControllerClient(conn)
			result, err := client.Apply(ctx, task)
			if err != nil {
				desc := "Tekton task"
				if md := task.GetMetadata(); md != nil {
					name := md.GetName()
					if name == "" {
						name = md.GetSlug()
					}
					if name != "" {
						desc = fmt.Sprintf("Tekton task %q in org %q", name, md.GetOrg())
					}
				}
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
