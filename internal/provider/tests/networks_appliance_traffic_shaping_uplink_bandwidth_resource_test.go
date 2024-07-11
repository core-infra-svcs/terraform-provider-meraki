package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceTrafficShapingUplinkBandWidthResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_traffic_shaping_uplink_bandwidth"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Update Network Appliance Traffic Shaping UplinkBandWidth.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateNetworksApplianceTrafficShapingUplinkBandWidth(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limit_wan1_limit_up", "100"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limit_wan1_limit_down", "600"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limit_wan2_limit_up", "100"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limit_wan2_limit_down", "600"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limit_cellular_limit_up", "101200"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", "bandwidth_limit_cellular_limit_up", "101200"),
				),
			},
		},
	})
}

func testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_networks_appliance_traffic_shaping_uplink_bandwidth"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

func testAccNetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateNetworksApplianceTrafficShapingUplinkBandWidth(orgId string, serial string) string {
	result := fmt.Sprintf(`
 
 resource "meraki_network" "test" {
	 organization_id = "%s"
	 product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_networks_appliance_traffic_shaping_uplink_bandwidth"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
 }

 resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

 resource "meraki_networks_appliance_traffic_shaping_uplink_bandwidth" "test" {
	 depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	 network_id = resource.meraki_network.test.network_id
	 bandwidth_limit_wan2_limit_up = 100
	 bandwidth_limit_wan2_limit_down = 600 
	 bandwidth_limit_cellular_limit_up = 101200
	 bandwidth_limit_cellular_limit_down = 101200	
	 bandwidth_limit_wan1_limit_up = 100
	 bandwidth_limit_wan1_limit_down = 600
 }
 `, orgId, serial)
	return result
}
