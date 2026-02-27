package variable

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
	"google.golang.org/grpc"
)

// ApplyInput holds the explicit parameters for creating or updating a variable.
type ApplyInput struct {
	Name        string
	Org         string
	Scope       variablev1.VariableSpec_Scope
	Env         string
	Description string
	Value       string
}

// Apply creates or updates a variable via the
// VariableCommandController.Apply RPC.
//
// The Variable proto is constructed from explicit parameters rather than JSON
// passthrough because the schema is simple and stable. This provides better
// agent UX and validation than requiring the full proto structure.
func Apply(ctx context.Context, serverAddress string, input ApplyInput) (string, error) {
	variable := &variablev1.Variable{
		Metadata: &apiresource.ApiResourceMetadata{
			Name: input.Name,
			Org:  input.Org,
			Env:  input.Env,
		},
		Spec: &variablev1.VariableSpec{
			Scope:       input.Scope,
			Description: input.Description,
			Value:       input.Value,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := variablev1.NewVariableCommandControllerClient(conn)
			result, err := client.Apply(ctx, variable)
			if err != nil {
				desc := fmt.Sprintf("variable %q in org %q", input.Name, input.Org)
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
