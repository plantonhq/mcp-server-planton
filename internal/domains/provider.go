package domains

import (
	"fmt"

	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
)

// ResolveProvider maps a lowercase provider string (e.g. "aws", "gcp") to the
// corresponding CloudResourceProvider enum value from the openmcf proto stubs.
func ResolveProvider(s string) (cloudresourcekind.CloudResourceProvider, error) {
	v, ok := cloudresourcekind.CloudResourceProvider_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown cloud resource provider %q — valid values: %s",
			s, JoinEnumValues(cloudresourcekind.CloudResourceProvider_value, "cloud_resource_provider_unspecified"))
	}
	return cloudresourcekind.CloudResourceProvider(v), nil
}

// ResolveProvisioner maps a lowercase provisioner string (e.g. "terraform",
// "pulumi", "tofu") to the corresponding IacProvisioner enum value from the
// openmcf proto stubs.
func ResolveProvisioner(s string) (shared.IacProvisioner, error) {
	v, ok := shared.IacProvisioner_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown iac provisioner %q — valid values: %s",
			s, JoinEnumValues(shared.IacProvisioner_value, "iac_provisioner_unspecified"))
	}
	return shared.IacProvisioner(v), nil
}
