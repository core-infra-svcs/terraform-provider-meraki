package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdminResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationsAdminResourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_admin"),
				),
			},

			// Create and Read testing (network)
			{
				Config: testAccOrganizationsAdminResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read testing (admin)
			{
				Config: testAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "name", "testAdmin"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "email", "meraki_organizations_admin_test1@example.com"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "org_access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "authentication_method", "Email"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "has_api_key", "false"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "west"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded organizationId + adminId.
			/*
				{
						ResourceName:      "meraki_organizations_admin.test",
						ImportState:       true,
						ImportStateVerify: false,
						ImportStateId:     "657525545596096508, 657525545596237587",
					},
			*/

			// Update testing
			{
				Config: testUpdatedAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "east"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
				),
			},

			// Delete testing automatically occurs in TestCase
			// This test can result in orphaned resources as organizations cannot be deleted with admins still present.
		},
	})
}

const testAccOrganizationsAdminResourceConfigCreateOrg = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_admin"
 	api_enabled = true
 }
 `

const testAccOrganizationsAdminResourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccOrganizationsAdminResourceConfig = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_organizations_admin" "test" {
	depends_on = ["meraki_organization.test", "meraki_network.test"]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "testAdmin"
	email       = "meraki_organizations_admin_test1@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			  {
			   tag = "west"
			   access = "read-only"
			  }]
    networks    = [{
                  id = resource.meraki_network.test.network_id
                  access = "read-only"
                }]
}
`

const testUpdatedAccOrganizationsAdminResourceConfig = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_organizations_admin" "test" {
	depends_on = ["meraki_organization.test", "meraki_network.test"]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "testAdmin"
	email       = "meraki_organizations_admin_test1@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			{
				tag = "east"
				access = "read-only"
			}]
    networks    = [{
                  id = resource.meraki_network.test.network_id
                  access = "read-only"
                }]
}
`
