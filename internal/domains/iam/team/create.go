package team

import (
	"context"
	"fmt"

	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
	teamv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/team/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// memberTypeResolver maps user-supplied member type strings to the enum.
var memberTypeResolver = domains.NewEnumResolver[apiresourcekind.ApiResourceKind](
	apiresourcekind.ApiResourceKind_value,
	"member type",
	"api_resource_kind_unspecified",
)

// Create provisions a new team via TeamCommandController.Create.
func Create(ctx context.Context, serverAddress, org, name, description string, members []MemberInput) (string, error) {
	protoMembers, err := toProtoMembers(members)
	if err != nil {
		return "", err
	}

	t := &teamv1.Team{
		ApiVersion: "iam.planton.ai/v1",
		Kind:       "Team",
		Metadata:   &apiresource.ApiResourceMetadata{Org: org, Name: name},
		Spec: &teamv1.TeamSpec{
			Description: description,
			Members:     protoMembers,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := teamv1.NewTeamCommandControllerClient(conn)
			resp, err := client.Create(ctx, t)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("team %q in org %q", name, org))
			}
			return domains.MarshalJSON(resp)
		})
}

func toProtoMembers(members []MemberInput) ([]*teamv1.TeamMember, error) {
	if len(members) == 0 {
		return nil, nil
	}
	out := make([]*teamv1.TeamMember, 0, len(members))
	for _, m := range members {
		mt, err := memberTypeResolver.Resolve(m.MemberType)
		if err != nil {
			return nil, err
		}
		out = append(out, &teamv1.TeamMember{
			MemberType: mt,
			MemberId:   m.MemberID,
		})
	}
	return out, nil
}
