package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsSamlRolesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test-acc-meraki-organizations-saml-roles"),
				),
			},

			// Enable SAML on test Organization
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigUpdateOrganizationSaml,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "enabled", "true"),
				),
			},

			// Create and Read Network and Organization Saml Role
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigCreateNetworkAndSamlRole,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_role.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_role.test", "role", "testrole"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_role.test", "org_access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_role.test", "tags.0.tag", "west"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_role.test", "tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_role.test", "networks.0.access", "read-only"),
				),
			},

			// Read OrganizationsSamlRoles
			{
				Config: testAccOrganizationsSamlRolesDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "list.0.role", "testrole"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "list.0.org_access", "read-only"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "list.0.tags.0.tag", "west"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "list.0.tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_roles.test", "list.0.networks.0.access", "read-only"),
				),
			},
		},
	})
}

const testAccOrganizationsSamlRolesDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test-acc-meraki-organizations-saml-roles"
 	api_enabled = true
 } 
 `

const testAccOrganizationsSamlRolesDataSourceConfigUpdateOrganizationSaml = `
resource "meraki_organization" "test" {
}

resource "meraki_organization_saml" "test" {
	depends_on = [
    	resource.meraki_organization.test
  	]
	organization_id = resource.meraki_organization.test.organization_id
	enabled = true
}
`

const testAccOrganizationsSamlRolesDataSourceConfigCreateNetworkAndSamlRole = `

resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
resource "meraki_organizations_saml_role" "test" {	
    depends_on = [
		resource.meraki_organization.test,
		resource.meraki_network.test
	]
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

const testAccOrganizationsSamlRolesDataSourceConfigRead = `
resource "meraki_organization" "test" {}

data "meraki_organizations_saml_roles" "test" {
      organization_id = resource.meraki_organization.test.organization_id
}
`
