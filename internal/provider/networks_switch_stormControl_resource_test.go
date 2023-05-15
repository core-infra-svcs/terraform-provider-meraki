package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchStormControlResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccNetworksSwitchStormControlResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_networks_switch_storm_control"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksSwitchStormControlResourceConfigCreateNetwork,
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

			// TODO: Update and Read Network Switch Settings Rule.
			{
				Config: testAccNetworksSwitchStormControlResourceConfigUpdateRule,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},

			// Update and Read Network Switch Settings.
			{
				Config: testAccNetworksSwitchStormControlResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_storm_control.test", "multicast_threshold", "30"),
					resource.TestCheckResourceAttr("meraki_networks_switch_storm_control.test", "broadcast_threshold", "30"),
					resource.TestCheckResourceAttr("meraki_networks_switch_storm_control.test", "unknown_unicast_threshold", "30"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccNetworksSwitchStormControlResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_networks_switch_storm_control"
 	api_enabled = true
 } 
 `

const testAccNetworksSwitchStormControlResourceConfigCreateNetwork = `
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
`

// TODO: add rule
const testAccNetworksSwitchStormControlResourceConfigUpdateRule = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
}

`

const testAccNetworksSwitchStormControlResourceConfigUpdate = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]	
	product_types = ["appliance", "switch", "wireless"]	
}

resource "meraki_networks_switch_storm_control" "test" {
	depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id

    broadcast_threshold = 30
	multicast_threshold = 30
	unknown_unicast_threshold = 30
}
`
