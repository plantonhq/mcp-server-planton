package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// pipelineFileEntry is a JSON-friendly representation of a
// ServiceRepoPipelineFile with content decoded from bytes to a plain string.
type pipelineFileEntry struct {
	Path        string `json:"path"`
	SHA         string `json:"sha,omitempty"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// ListFiles discovers Tekton pipeline files in a service's Git repository
// via the PipelineQueryController.ListServiceRepoPipelineFiles RPC.
//
// The response content bytes are decoded to UTF-8 strings so that agents
// receive human-readable YAML/pipeline content instead of base64.
func ListFiles(ctx context.Context, serverAddress, serviceID, branch string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineQueryControllerClient(conn)
			resp, err := client.ListServiceRepoPipelineFiles(ctx, &pipelinev1.ListServiceRepoPipelineFilesInput{
				ServiceId: serviceID,
				Branch:    branch,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("pipeline files for service %q", serviceID))
			}

			entries := make([]pipelineFileEntry, 0, len(resp.GetEntries()))
			for _, f := range resp.GetEntries() {
				entries = append(entries, pipelineFileEntry{
					Path:        f.GetPath(),
					SHA:         f.GetSha(),
					Content:     string(f.GetContent()),
					Encoding:    f.GetEncoding(),
					DisplayName: f.GetDisplayName(),
				})
			}

			b, err := json.MarshalIndent(entries, "", "  ")
			if err != nil {
				return "", fmt.Errorf("marshal pipeline files: %w", err)
			}
			return string(b), nil
		})
}
