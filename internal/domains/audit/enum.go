package audit

import (
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
)

var apiResourceKindResolver = domains.NewEnumResolver[apiresourcekind.ApiResourceKind](
	apiresourcekind.ApiResourceKind_value, "resource kind", "unspecified")
