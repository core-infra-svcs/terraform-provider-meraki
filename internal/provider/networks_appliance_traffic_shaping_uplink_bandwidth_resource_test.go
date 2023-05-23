package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceTrafficShapingUplinkBandWidthResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_appliance_traffic_shaping_uplink_bandwidth"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateNetwork,
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

			// Create and Read Network Appliance Traffic Shaping UplinkBandWidth.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateNetworksApplianceTrafficShapingUplinkBandWidth,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.wan1.limit_up", "1000000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.wan1.limit_up", "1000000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.cellular.limit_up", "51200"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.cellular.limit_up", "51200"),
				),
			},

			// Update Network Appliance Traffic Shaping UplinkBandWidth.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateNetworksApplianceTrafficShapingUplinkBandWidth,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.wan1.limit_up", "900000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.wan1.limit_up", "900000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.cellular.limit_up", "101200"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limits.cellular.limit_up", "101200"),
				),
			},
		},
	})
}

const testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_networks_appliance_traffic_shaping_uplink_bandwidth"
 	api_enabled = true
 }
 `
const testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateNetwork = `
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

const testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateNetworksApplianceTrafficShapingUplinkBandWidth = `
 resource "meraki_organization" "test" {}
 
 resource "meraki_network" "test" {
	 depends_on = [resource.meraki_organization.test]	
	 product_types = ["appliance", "switch", "wireless"]	
 }
 
 resource "meraki_networks_appliance_traffic_shaping_uplink_bandwidth" "test" {
	 depends_on = [resource.meraki_organization.test,
	 resource.meraki_network.test]
	 network_id = resource.meraki_network.test.network_id
	 bandwidth_limits = {
		cellular = {
			limit_up = 101200
            limit_down = 101200
		}
		wan2 = {}
		wan1 = {
			limit_up = 900000
            limit_down = 900000
		}	
	} 
 }
 `

const testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateNetworksApplianceTrafficShapingUplinkBandWidth = `
 resource "meraki_organization" "test" {}
 
 resource "meraki_network" "test" {
	 depends_on = [resource.meraki_organization.test]	
	 product_types = ["appliance", "switch", "wireless"]	
 }
 
 resource "meraki_networks_appliance_traffic_shaping_uplink_bandwidth" "test" {
	 depends_on = [resource.meraki_organization.test,
	 resource.meraki_network.test]
	 network_id = resource.meraki_network.test.network_id
	 bandwidth_limits = {
		cellular = {
			limit_up = 51200
            limit_down = 51200
		}
		wan2 = {}
		wan1 = {
			limit_up = 1000000
            limit_down = 1000000
		}	
	} 
 }
 `
