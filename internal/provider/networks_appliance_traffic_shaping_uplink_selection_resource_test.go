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
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_appliance_traffic_shaping_uplink_selection"),
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

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworkApplianceVlansSettings(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans_settings.test", "vlans_enabled", "true"),
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
	name = "test_acc_network_appliance_traffic_shaping_uplink_selection"
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
		name = "test_acc_network_appliance_traffic_shaping_uplink_selection"
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

func testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworkApplianceVlansSettings(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
		tags = ["tag1"]
		name = "test_acc_network_appliance_traffic_shaping_uplink_selection"
		timezone = "America/Los_Angeles"
		notes = "Additional description of the network"
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  vlans_enabled = true
}
`, orgId)
	return result
}

func testAccNetworksApplianceTrafficShapingUplinkSelectionResourceConfigUpdateNetworksApplianceTrafficShapingUplinkSelection(orgId string, serial string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
		tags = ["tag1"]
		name = "test_acc_network_appliance_traffic_shaping_uplink_selection"
		timezone = "America/Los_Angeles"
		notes = "Additional description of the network"
}

resource "meraki_networks_appliance_traffic_shaping_uplink_selection" "test" {
	depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
	active_active_auto_vpn_enabled = false
    default_uplink = "wan1"
    load_balancing_enabled = true
    failover_and_failback = {
        immediate = {
            enabled = true
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
                            vlan = 10
                            host = 254
                        }
                        destination = {
                            port = "8080"
                            cidr = "192.168.10.0/24"
                        }
                    }
                }
            ]
            preferred_uplink = "wan1"
        }
    ]
    vpn_traffic_uplink_preferences = [
        {
            traffic_filters = [
                {
                    type = "applicationCategory"
                    value = {
						id = "meraki:layer7/category/1"
                        protocol = "tcp"
                        source = {
                            port = "any"
                            cidr = "192.168.1.0/24"
                            network = "L_23456789"
                            host = 200
                        }
                        destination = {
                            port = "1-1024"
                            cidr = "0.0.0.0/0"
                            network = "L_12345678"
							vlan = 1
							fqdn = "www.google.com"
                        }
                    }
                }
            ]
            preferred_uplink = "bestForVoIP"
            fail_over_criterion = "poorPerformance"
            performance_class = {
                type = "custom"
                builtin_performance_class_name = "VoIP"
                custom_performance_class_id = "123456"
            }
        }
    ]
}	
`, orgId)
	return result
}
