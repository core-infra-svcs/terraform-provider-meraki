package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceTrafficShappingUplinkBandWidthResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_appliance_traffic_shapping_uplink_bandWidth"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_traffic_shapping_uplink_bandwidth.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigCreateNetwork,
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

			// Create and Read Network Appliance Traffic Shapping UplinkBandWidth.
			{
				Config: testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigCreateNetworksApplianceTrafficShappingUplinkBandWidth,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.wan1.limit_up", "1000000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.wan1.limit_up", "1000000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.cellular.limit_up", "51200"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.cellular.limit_up", "51200"),
				),
			},

			// Update Network Appliance Traffic Shapping UplinkBandWidth.
			{
				Config: testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigUpdateNetworksApplianceTrafficShappingUplinkBandWidth,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.wan1.limit_up", "900000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.wan1.limit_up", "900000"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.cellular.limit_up", "101200"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shapping_uplink_bandWidth.test", "bandwidth_limits.cellular.limit_up", "101200"),
				),
			},
		},
	})
}

const testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_networks_appliance_traffic_shapping_uplink_bandWidth"
 	api_enabled = true
 }
 `
const testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigCreateNetwork = `
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

const testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigUpdateNetworksApplianceTrafficShappingUplinkBandWidth = `
 resource "meraki_organization" "test" {}
 
 resource "meraki_network" "test" {
	 depends_on = [resource.meraki_organization.test]	
	 product_types = ["appliance", "switch", "wireless"]	
 }
 
 resource "meraki_networks_appliance_traffic_shapping_uplink_bandWidth" "test" {
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

const testAccNetworksApplianceTrafficShappingUplinkBandWidthResourceConfigCreateNetworksApplianceTrafficShappingUplinkBandWidth = `
 resource "meraki_organization" "test" {}
 
 resource "meraki_network" "test" {
	 depends_on = [resource.meraki_organization.test]	
	 product_types = ["appliance", "switch", "wireless"]	
 }
 
 resource "meraki_networks_appliance_traffic_shapping_uplink_bandWidth" "test" {
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
