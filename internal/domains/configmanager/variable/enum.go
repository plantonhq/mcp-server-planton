package variable

import (
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
)

var scopeResolver = domains.NewEnumResolver[variablev1.VariableSpec_Scope](
	variablev1.VariableSpec_Scope_value, "scope", "scope_unspecified")
