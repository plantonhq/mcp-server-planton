package variablegroup

import (
	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

var scopeResolver = domains.NewEnumResolver[variablegroupv1.VariableGroupSpec_Scope](
	variablegroupv1.VariableGroupSpec_Scope_value, "scope", "scope_unspecified")
