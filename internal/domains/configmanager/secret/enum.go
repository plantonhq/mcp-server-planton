package secret

import (
	secretv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/secret/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

var scopeResolver = domains.NewEnumResolver[secretv1.SecretSpec_Scope](
	secretv1.SecretSpec_Scope_value, "scope", "scope_unspecified")
