package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchStpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksSwitchStpResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_switch_stp"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksSwitchStpResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksSwitchStpResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_stp.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_switch_stp.test", "stp_bridge_priority.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_stp.test", "rstp_enabled", "true"),
				),
			},

			//{
			//	ResourceName:      "meraki_networks_switch_stp.test",
			//	ImportState:       true,
			//	ImportStateVerify: false,
			//	ImportStateId:     "1234567890, 0987654321",
			//},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccNetworksSwitchStpResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_switch_stp"
 	api_enabled = true
 }
 `

const testAccNetworksSwitchStpResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "Main Office"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}
`

const testAccNetworksSwitchStpResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_switch_stp" "test" {
 depends_on = [
	resource.meraki_network.test
 ]
 network_id = resource.meraki_network.test.network_id
 rstp_enabled = true
 stp_bridge_priority = [{
	stp_priority = 4096
 }]
}
`
