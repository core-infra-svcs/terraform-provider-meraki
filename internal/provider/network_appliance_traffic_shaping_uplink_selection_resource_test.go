package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceTrafficShapingUplinkSelectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworkApplianceTrafficShapingUplinkSelectionResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_network_appliance_traffic_shaping_uplink_selection"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_traffic_shaping_uplink_selection.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworkApplianceTrafficShapingUplinkSelectionResourceConfigCreateNetwork,
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

			// Update and Read Network Settings.
			{
				Config: testAccNetworkApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworkApplianceTrafficShapingUplinkSelection,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_selection.test", "active_active_auto_vpn_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_selection.test", "default_uplink", "wan1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_selection.test", "load_balancing_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_selection.test", "failover_and_failback.immediate.enabled", "false"),
				),
			},
		},
	})
}

const testAccNetworkApplianceTrafficShapingUplinkSelectionResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_network_appliance_traffic_shaping_uplink_selection"
 	api_enabled = true
 } 
 `

const testAccNetworkApplianceTrafficShapingUplinkSelectionResourceConfigCreateNetwork = `
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

const testAccNetworkApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworkApplianceTrafficShapingUplinkSelection = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]	
	product_types = ["appliance", "switch", "wireless"]	
}

resource "meraki_networks_appliance_traffic_shaping_uplink_selection" "test" {
	depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
    network_id = "N_784752235069351032"
	load_balancing_enabled = false
	vpn_traffic_uplink_preferences = []
	failover_and_failback = {
    immediate =  {}
    }
	wan_traffic_uplink_preferences = []	
}
`
