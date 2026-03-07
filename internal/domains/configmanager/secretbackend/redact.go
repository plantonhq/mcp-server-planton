package secretbackend

import (
	secretbackendv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/secretbackend/v1"
)

const redacted = "[REDACTED]"

// RedactSecretBackend replaces sensitive credential fields in the spec
// config blocks with a placeholder. This is defense-in-depth — the backend
// already applies mask-on-write for query responses, but command responses
// may contain unmasked data.
func RedactSecretBackend(sb *secretbackendv1.SecretBackend) {
	if sb == nil || sb.Spec == nil {
		return
	}
	spec := sb.Spec

	if c := spec.GetOpenbao(); c != nil && c.Token != "" {
		c.Token = redacted
	}
	if c := spec.GetAwsSecretsManager(); c != nil {
		if c.AccessKeyId != "" {
			c.AccessKeyId = redacted
		}
		if c.SecretAccessKey != "" {
			c.SecretAccessKey = redacted
		}
	}
	if c := spec.GetHashicorpVault(); c != nil && c.Token != "" {
		c.Token = redacted
	}
	if c := spec.GetGcpSecretManager(); c != nil && c.ServiceAccountKeyJson != "" {
		c.ServiceAccountKeyJson = redacted
	}
	if c := spec.GetAzureKeyVault(); c != nil && c.ClientSecret != "" {
		c.ClientSecret = redacted
	}

	if enc := spec.GetEncryption(); enc != nil {
		if c := enc.GetAwsKms(); c != nil {
			if c.AccessKeyId != "" {
				c.AccessKeyId = redacted
			}
			if c.SecretAccessKey != "" {
				c.SecretAccessKey = redacted
			}
		}
		if c := enc.GetGcpKms(); c != nil && c.ServiceAccountKeyJson != "" {
			c.ServiceAccountKeyJson = redacted
		}
		if c := enc.GetAzureKeyVault(); c != nil && c.ClientSecret != "" {
			c.ClientSecret = redacted
		}
	}
}
