package iacmodule

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	infrahubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// SearchInput holds the validated parameters for searching IaC modules.
type SearchInput struct {
	Org          string
	SearchText   string
	Kind         string
	Provisioner  string
	Provider     string
	PageNum      int32
	PageSize     int32
}

// Search queries IaC modules via the InfraHubSearchQueryController.
//
// When Org is provided, the org-context RPC is called with both official and
// organization modules included. When Org is empty, only official (public)
// modules are searched.
//
// Optional Kind, Provisioner, and Provider filters are resolved to their
// respective enum values before being sent to the backend.
func Search(ctx context.Context, serverAddress string, input SearchInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	var kind cloudresourcekind.CloudResourceKind
	if input.Kind != "" {
		k, err := domains.ResolveKind(input.Kind)
		if err != nil {
			return "", err
		}
		kind = k
	}

	var provisioners []shared.IacProvisioner
	if input.Provisioner != "" {
		p, err := domains.ResolveProvisioner(input.Provisioner)
		if err != nil {
			return "", err
		}
		provisioners = []shared.IacProvisioner{p}
	}

	var providers []cloudresourcekind.CloudResourceProvider
	if input.Provider != "" {
		p, err := domains.ResolveProvider(input.Provider)
		if err != nil {
			return "", err
		}
		providers = []cloudresourcekind.CloudResourceProvider{p}
	}

	pageInfo := &rpc.PageInfo{
		Num:  pageNum - 1,
		Size: pageSize,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrahubsearch.NewInfraHubSearchQueryControllerClient(conn)

			if input.Org != "" {
				resp, err := client.SearchIacModulesByOrgContext(ctx, &infrahubsearch.SearchIacModulesByOrgContextInput{
					Org:                         input.Org,
					SearchText:                  input.SearchText,
					PageInfo:                    pageInfo,
					CloudResourceKind:           kind,
					Provisioners:                provisioners,
					Providers:                   providers,
					IsIncludeOfficial:           true,
					IsIncludeOrganizationModules: true,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("iac modules in org %q", input.Org))
				}
				return domains.MarshalJSON(resp)
			}

			resp, err := client.SearchOfficialIacModules(ctx, &infrahubsearch.SearchOfficialIacModulesInput{
				SearchText:        input.SearchText,
				PageInfo:          pageInfo,
				CloudResourceKind: kind,
				Provisioners:      provisioners,
				Providers:         providers,
			})
			if err != nil {
				return "", domains.RPCError(err, "official iac modules")
			}
			return domains.MarshalJSON(resp)
		})
}
