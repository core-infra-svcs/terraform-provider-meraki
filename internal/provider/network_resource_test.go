package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsNetworkResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing (network)
			{
				Config: testAccOrganizationsNetworkResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					//resource.TestCheckResourceAttr("meraki_network.test", "product_types", "[\"appliance\", \"switch\", \"wireless\"]"),
					//resource.TestCheckResourceAttr("meraki_network.test", "enrollment_string", "my-enrollment-string"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationsNetworkResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),

					// this checks the number of tags
					//resource.TestCheckResourceAttr("data.meraki_network.product_tags", "list.#", "2"),

					//resource.TestCheckResourceAttr("meraki_network.test", "tags", "[\"tag1\", \"tag2\"]"),
					//resource.TestCheckResourceAttr("meraki_network.test", "product_types", "[\"appliance\", \"switch\", \"wireless\"]"),
					//resource.TestCheckResourceAttr("meraki_network.test", "enrollment_string", "my-enrollment-string"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

/*
const testAccOrganizationResourceConfigOrganization = `
resource "meraki_organization" "testOrg1" {
	name = "testOrg1"
	api_enabled = true
}

output "testOrg1" {
  value = resource.meraki_organization.testOrg1.organization_id
}
`

resource "meraki_organization" "testOrg" {
	name = "testOrg1"
	api_enabled = true
}

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.testOrg"]
	product_types = ["appliance"]
	organization_id = resource.meraki_organization.testOrg.organization_id
	name = "Main Office"
	timezone = "America/Los_Angeles"
	enrollment_string = "my-enrollment-string"
	notes = "Additional description of the network"
}

*/

const testAccOrganizationsNetworkResourceConfig = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccOrganizationsNetworkResourceConfigUpdate = `
resource "meraki_network" "test" {
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`
