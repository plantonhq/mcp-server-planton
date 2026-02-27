package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing stack jobs.
// Optional string fields use their zero value ("") to indicate "not set".
// Optional int32 fields use 0 to indicate "use default".
type ListInput struct {
	Org               string
	Env               string
	CloudResourceKind string
	CloudResourceID   string
	OperationType     string
	Status            string
	Result            string
	PageNum           int32
	PageSize          int32
}

// List queries stack jobs matching the given filters via the
// StackJobQueryController.ListByFilters RPC.
//
// Enum string fields are resolved to their proto values before the call.
// Pagination defaults are applied when page_num or page_size are zero.
func List(ctx context.Context, serverAddress string, input ListInput) (string, error) {
	req := &stackjobv1.ListStackJobsByFiltersQueryInput{
		Org:             input.Org,
		Env:             input.Env,
		CloudResourceId: input.CloudResourceID,
		PageInfo:        pageInfo(input.PageNum, input.PageSize),
	}

	if input.CloudResourceKind != "" {
		k, err := resolveKind(input.CloudResourceKind)
		if err != nil {
			return "", err
		}
		req.CloudResourceKind = k
	}

	if input.OperationType != "" {
		op, err := resolveOperationType(input.OperationType)
		if err != nil {
			return "", err
		}
		req.StackJobOperation = op
	}

	if input.Status != "" {
		s, err := resolveExecutionStatus(input.Status)
		if err != nil {
			return "", err
		}
		req.Status = s
	}

	if input.Result != "" {
		r, err := resolveExecutionResult(input.Result)
		if err != nil {
			return "", err
		}
		req.Result = r
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.ListByFilters(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("stack jobs in org %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}

// pageInfo returns a PageInfo with sensible defaults applied.
func pageInfo(num, size int32) *rpc.PageInfo {
	if num <= 0 {
		num = defaultPageNum
	}
	if size <= 0 {
		size = defaultPageSize
	}
	return &rpc.PageInfo{Num: num, Size: size}
}
