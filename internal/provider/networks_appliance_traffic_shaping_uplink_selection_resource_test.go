package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceTrafficShapingUplinkSelectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_device"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim A Device To A Network
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "id", "example-id"),
				),
			},

			// Update and Read Network Settings.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworksApplianceTrafficShapingUplinkSelection(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_selection.test", "active_active_auto_vpn_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_traffic_shaping_uplink_selection.test", "default_uplink", "wan1"),
				),
			},
		},
	})
}

func testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_network_device"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesResourceConfigClaimNetworkDevice is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigClaimNetworkDevice(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
		tags = ["tag1"]
		name = "test_acc_network_device"
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
`, orgId, serial)
	return result
}

func testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworksApplianceTrafficShapingUplinkSelection(orgId string, serial string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
		tags = ["tag1"]
		name = "test_acc_network_device"
		timezone = "America/Los_Angeles"
		notes = "Additional description of the network"
}

resource "meraki_networks_appliance_traffic_shaping_uplink_selection" "test" {
	depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
	active_active_auto_vpn_enabled = true
	default_uplink = "wan1"
	load_balancing_enabled = true
	failover_and_failback = {
		immediate =  {
		   enabled = false
		}
    }
	wan_traffic_uplink_preferences = [
		{
			traffic_filters = [
				{
					type = "custom"
					value = {
						protocol = "tcp"
						source = {
							port = "1-1024"					
							cidr = "192.168.1.0/24"
							vlan = 10
							host = 254
						}
						destination = {
							cidr = "any"
							port = "any"
						}
					}
				}
			]
			preferred_uplink = "wan1"
		}
	]
	
	vpn_traffic_uplink_preferences = [
		traffic_filters = [
				{
					type = "custom"
					value = {
						protocol = "tcp"
						source = {
							port = "1-1024"					
							cidr = "192.168.1.0/24"
							vlan = 10
							host = 254
						}
						destination = {
							port = "1-1024"					
							cidr = "192.168.1.0/24"
							vlan = 10
							host = 254
							fqdn = "example.com"
						}
					}
				}
		]
	]
}	
`, orgId)
	return result
}

/*
{
		  trafficFilters = [
			{
			  type = "applicationCategory"
			  value = {
				protocol = "tcp"
				source = {
				  port = "any"
				  cidr = "192.168.1.0/24"
				  network = "L_23456789"
				  vlan = 20
				  host = 200
				}
				destination = {
				  port = "1-1024"
				  cidr = "any"
				  network = "L_12345678"
				  vlan = 10
				  host = 254
				  fqdn = "www.google.com"
				}
			  }
			}
		  ]
		  preferredUplink = "bestForVoIP"
		  failOverCriterion = "poorPerformance"
		  performanceClass = {
			type = "custom"
			builtinPerformanceClassName = "VoIP"
			customPerformanceClassId = "123456"
		  }
		}
*/
