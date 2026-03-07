package team

import (
	"context"
	"fmt"

	teamv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/team/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// UpdateFields holds the optional fields that can be modified on an
// existing team. Nil/empty values mean "leave unchanged"; members being
// non-nil (even empty) means "replace the member list".
type UpdateFields struct {
	Name        string
	Description string
	Members     *[]MemberInput
}

// Update performs a read-modify-write on an existing team. Both RPCs share
// a single gRPC connection within one WithConnection callback.
func Update(ctx context.Context, serverAddress, teamID string, fields UpdateFields) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			queryClient := teamv1.NewTeamQueryControllerClient(conn)
			t, err := queryClient.Get(ctx, &teamv1.TeamId{Value: teamID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("team %q", teamID))
			}

			if fields.Name != "" {
				t.Metadata.Name = fields.Name
			}
			if t.Spec == nil {
				t.Spec = &teamv1.TeamSpec{}
			}
			if fields.Description != "" {
				t.Spec.Description = fields.Description
			}
			if fields.Members != nil {
				protoMembers, err := toProtoMembers(*fields.Members)
				if err != nil {
					return "", err
				}
				t.Spec.Members = protoMembers
			}

			cmdClient := teamv1.NewTeamCommandControllerClient(conn)
			resp, err := cmdClient.Update(ctx, t)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("team %q", teamID))
			}
			return domains.MarshalJSON(resp)
		})
}
