package secret

import (
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
)

var scopeResolver = domains.NewEnumResolver[secretv1.SecretSpec_Scope](
	secretv1.SecretSpec_Scope_value, "scope", "scope_unspecified")
