package domains

import (
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
)

var (
	providerResolver    = NewEnumResolver[cloudresourcekind.CloudResourceProvider](cloudresourcekind.CloudResourceProvider_value, "cloud resource provider", "cloud_resource_provider_unspecified")
	provisionerResolver = NewEnumResolver[shared.IacProvisioner](shared.IacProvisioner_value, "iac provisioner", "iac_provisioner_unspecified")
)

// ResolveProvider maps a lowercase provider string (e.g. "aws", "gcp") to the
// corresponding CloudResourceProvider enum value from the openmcf proto stubs.
func ResolveProvider(s string) (cloudresourcekind.CloudResourceProvider, error) {
	return providerResolver.Resolve(s)
}

// ResolveProvisioner maps a lowercase provisioner string (e.g. "terraform",
// "pulumi", "tofu") to the corresponding IacProvisioner enum value from the
// openmcf proto stubs.
func ResolveProvisioner(s string) (shared.IacProvisioner, error) {
	return provisionerResolver.Resolve(s)
}
