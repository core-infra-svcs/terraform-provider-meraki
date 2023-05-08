package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchStormcontrolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Update and Read Network Switch Settings.
			{
				Config: testAccNetworksSwitchStormcontrolResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_storm_control.test", "multicast_threshold", "30"),
					resource.TestCheckResourceAttr("meraki_networks_switch_storm_control.test", "broadcast_threshold", "30"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccNetworksSwitchStormcontrolResourceConfigUpdate = `
resource "meraki_organization" "test" {
  name = "test_meraki_organizations"
  api_enabled = true
}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "Main Office"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}

resource "meraki_networks_switch_storm_control" "test" {
	network_id  = resource.meraki_network.test.network_id
    broadcast_threshold = 30
	multicast_threshold = 30
	unknown_unicast_threshold = 30
}
`
