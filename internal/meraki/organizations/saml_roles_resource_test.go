package organizations_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlRolesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_saml_roles"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_saml_roles"),
			},

			// Enable SAML on test Organization
			{
				Config: testAccOrganizationsSamlRolesResourceConfigUpdateOrganizationSaml(),
				Check:  OrganizationsSamlTestChecks(true),
			},

			// Create and Read Network
			{
				Config: testAccOrganizationsSamlRolesResourceConfigCreateNetwork(),
				Check:  OrganizationsNetworkTestChecks("test_acc_organizations_saml_roles"),
			},

			// Create and Read Organizations SAML Role
			{
				Config: testAccOrganizationsSamlRolesResourceConfigCreateNetworkAndSamlRole(),
				Check:  OrganizationsSamlRoleResourceTestChecks("testrole", "read-only", "read-only"),
			},

			// Update testing
			{
				Config: testAccOrganizationsSamlRolesResourceConfigUpdateNetworkAndSamlRole(),
				Check:  OrganizationsSamlRoleResourceTestChecks("testrole", "read-only", "read-only"),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_saml_role.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// testAccOrganizationsSamlRolesResourceConfigUpdateOrganizationSaml returns the configuration for enabling SAML on the organization
func testAccOrganizationsSamlRolesResourceConfigUpdateOrganizationSaml() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organization_saml" "test" {
		depends_on = [resource.meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		enabled = true
	}
	`
}

// testAccOrganizationsSamlRolesResourceConfigCreateNetwork returns the configuration for creating the network
func testAccOrganizationsSamlRolesResourceConfigCreateNetwork() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_network" "test" {
		depends_on = [resource.meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		product_types = ["appliance", "switch", "wireless"]
		tags = ["tag1"]
		name = "test_acc_organizations_saml_roles"
		timezone = "America/Los_Angeles"
		notes = "Additional description of the network"
	}
	`
}

// testAccOrganizationsSamlRolesResourceConfigCreateNetworkAndSamlRole returns the configuration for creating the network and SAML role
func testAccOrganizationsSamlRolesResourceConfigCreateNetworkAndSamlRole() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_network" "test" {
		organization_id = resource.meraki_organization.test.organization_id
		product_types = ["appliance", "switch", "wireless"]
	}

	resource "meraki_organizations_saml_role" "test" {	
		depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
		organization_id = resource.meraki_organization.test.organization_id
		role = "testrole"
		tags = []
		org_access = "read-only"
		networks = [{
			id = resource.meraki_network.test.network_id
			access = "read-only"
		}]
	}
	`
}

// testAccOrganizationsSamlRolesResourceConfigUpdateNetworkAndSamlRole returns the configuration for updating the SAML role
func testAccOrganizationsSamlRolesResourceConfigUpdateNetworkAndSamlRole() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_network" "test" {
		organization_id = resource.meraki_organization.test.organization_id
		product_types = ["appliance", "switch", "wireless"]
	}

	resource "meraki_organizations_saml_role" "test" {	
		depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
		organization_id = resource.meraki_organization.test.organization_id
		role = "testrole"
		org_access = "read-only"
		networks = [{
			id = resource.meraki_network.test.network_id
			access = "read-only"
		}]
	}
	`
}

// OrganizationsSamlRoleTestChecks returns the test check functions for verifying the SAML role
func OrganizationsSamlRoleResourceTestChecks(role, orgAccess, networkAccess string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"role":              role,
		"org_access":        orgAccess,
		"networks.0.access": networkAccess,
	}

	return utils.ResourceTestCheck("meraki_organizations_saml_role.test", expectedAttrs)
}

// OrganizationsNetworkTestChecks returns the test check functions for verifying the network resource
func OrganizationsNetworkResourceTestChecks(networkName string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":            networkName,
		"timezone":        "America/Los_Angeles",
		"tags.#":          "1",
		"tags.0":          "tag1",
		"product_types.#": "3",
		"product_types.0": "appliance",
		"product_types.1": "switch",
		"product_types.2": "wireless",
		"notes":           "Additional description of the network",
	}

	return utils.ResourceTestCheck("meraki_network.test", expectedAttrs)
}
