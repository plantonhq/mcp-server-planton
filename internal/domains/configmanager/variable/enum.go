package variable

import (
	variablev1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variable/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

var scopeResolver = domains.NewEnumResolver[variablev1.VariableSpec_Scope](
	variablev1.VariableSpec_Scope_value, "scope", "scope_unspecified")
