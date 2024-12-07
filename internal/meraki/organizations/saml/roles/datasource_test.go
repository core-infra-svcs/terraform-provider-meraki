package roles_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlRolesDataSource(t *testing.T) {
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
				Config: testAccOrganizationsSamlRolesDataSourceConfigUpdateOrganizationSaml(),
				Check:  OrganizationsSamlTestChecks(true),
			},

			// Create and Read Network
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigCreateNetwork(),
				Check:  OrganizationsNetworkTestChecks("test_acc_organizations_saml_roles"),
			},

			// Create and Read Organization SAML Role
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigCreateSamlRole(),
				Check:  OrganizationsSamlRoleTestChecks("testrole", "read-only", "west", "read-only"),
			},

			// Read Organizations SAML Roles
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigRead(),
				Check:  OrganizationsSamlRolesDataSourceTestChecks(),
			},
		},
	})
}

// testAccOrganizationsSamlRolesDataSourceConfigUpdateOrganizationSaml returns the configuration for enabling SAML on the organization
func testAccOrganizationsSamlRolesDataSourceConfigUpdateOrganizationSaml() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organization_saml" "test" {
		depends_on = [resource.meraki_organization.test]
		id = resource.meraki_organization.test.organization_id
		enabled = true
	}
	`
}

// testAccOrganizationsSamlRolesDataSourceConfigCreateNetwork returns the configuration for creating the network
func testAccOrganizationsSamlRolesDataSourceConfigCreateNetwork() string {
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

// testAccOrganizationsSamlRolesDataSourceConfigCreateSamlRole returns the configuration for creating the SAML role
func testAccOrganizationsSamlRolesDataSourceConfigCreateSamlRole() string {
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
		tags = [{
			tag = "west"
			access = "read-only"
		}]
		networks = [{
			id = resource.meraki_network.test.network_id
			access = "read-only"
		}]
	}
	`
}

// testAccOrganizationsSamlRolesDataSourceConfigRead returns the configuration for reading the SAML roles data source
func testAccOrganizationsSamlRolesDataSourceConfigRead() string {
	return `
	resource "meraki_organization" "test" {}

	data "meraki_organizations_saml_roles" "test" {
		id = resource.meraki_organization.test.organization_id
	}
	`
}

// OrganizationsNetworkTestChecks returns the test check functions for verifying the network resource
func OrganizationsNetworkTestChecks(networkName string) resource.TestCheckFunc {
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

// OrganizationsSamlRoleTestChecks returns the test check functions for verifying the SAML role
func OrganizationsSamlRoleTestChecks(role, orgAccess, tag, tagAccess string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"role":              role,
		"org_access":        orgAccess,
		"tags.0.tag":        tag,
		"tags.0.access":     tagAccess,
		"networks.0.access": "read-only",
	}

	return utils.ResourceTestCheck("meraki_organizations_saml_role.test", expectedAttrs)
}

// OrganizationsSamlRolesDataSourceTestChecks returns the test check functions for verifying the SAML roles data source
func OrganizationsSamlRolesDataSourceTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.#":                   "1",
		"list.0.role":              "testrole",
		"list.0.org_access":        "read-only",
		"list.0.tags.0.tag":        "west",
		"list.0.tags.0.access":     "read-only",
		"list.0.networks.0.access": "read-only",
	}

	return utils.ResourceTestCheck("data.meraki_organizations_saml_roles.test", expectedAttrs)
}
