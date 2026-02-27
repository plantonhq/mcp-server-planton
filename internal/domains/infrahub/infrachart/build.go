package infrachart

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	infrachartv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrachart/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

// BuildInput holds the validated parameters for building (rendering) an infra
// chart with optional parameter overrides.
type BuildInput struct {
	ChartID string
	Params  map[string]any
}

// Build fetches an infra chart by ID, applies the supplied parameter
// overrides, and calls InfraChartQueryController.Build to render the
// templates into final YAML output and a cloud resource DAG.
//
// Both the Get and Build RPCs share a single gRPC connection.
func Build(ctx context.Context, serverAddress string, input BuildInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrachartv1.NewInfraChartQueryControllerClient(conn)

			chart, err := client.Get(ctx, &apiresource.ApiResourceId{Value: input.ChartID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra chart %q", input.ChartID))
			}

			if err := mergeParams(chart, input.Params); err != nil {
				return "", err
			}

			resp, err := client.Build(ctx, chart)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("build infra chart %q", input.ChartID))
			}
			return domains.MarshalJSON(resp)
		})
}

// mergeParams applies param overrides from the user-supplied map onto the
// chart's existing parameter definitions. Only params whose names appear in
// the map are modified; all others keep their chart defaults.
//
// Returns an error if a supplied param name does not exist in the chart or
// if the value cannot be converted to a structpb.Value.
func mergeParams(chart *infrachartv1.InfraChart, overrides map[string]any) error {
	if len(overrides) == 0 {
		return nil
	}

	if chart.GetSpec() == nil || len(chart.GetSpec().GetParams()) == 0 {
		return fmt.Errorf("infra chart %q has no parameters — cannot apply overrides", chart.GetMetadata().GetName())
	}

	paramsByName := make(map[string]*infrachartv1.InfraChartParam, len(chart.Spec.Params))
	for _, p := range chart.Spec.Params {
		paramsByName[p.Name] = p
	}

	for name, rawValue := range overrides {
		param, ok := paramsByName[name]
		if !ok {
			known := make([]string, 0, len(paramsByName))
			for k := range paramsByName {
				known = append(known, k)
			}
			return fmt.Errorf("unknown parameter %q — available params: %v", name, known)
		}

		val, err := structpb.NewValue(rawValue)
		if err != nil {
			return fmt.Errorf("parameter %q: cannot convert value %v: %w", name, rawValue, err)
		}
		param.Value = val
	}

	return nil
}
