package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchStormcontrolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read NetworksSwitchStormcontrol
			{
				Config: testAccNetworksSwitchStormcontrolDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_switch_storm_control.test", "id", "example-id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccNetworksSwitchStormcontrolDataSourceConfigRead = `
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

data "meraki_networks_switch_storm_control" "test" {
    network_id = resource.meraki_network.test.network_id
}
`
