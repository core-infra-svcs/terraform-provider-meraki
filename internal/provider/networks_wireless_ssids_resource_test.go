package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccNetworksWirelessSsidsResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	ssids := 14 // Number of SSIDs to create, meraki max is 15

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_wireless_ssids_resource"),
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

			// Create and Read testing
			{
				Config: testAccNetworksWirelessSsidsResourceConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "number", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "name", "My SSID"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "auth_mode", "open"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "adult_content_filtering_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "availability_tags.#", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "available_on_all_aps", "true"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "band_selection", "Dual band operation"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dns_rewrite.dns_custom_nameservers.#", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dns_rewrite.enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ip_assignment_mode", "NAT mode"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "mandatory_dhcp_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "min_bit_rate", "11"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_client_bandwidth_limit_down", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_client_bandwidth_limit_up", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_ssid_bandwidth_limit_down", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_ssid_bandwidth_limit_up", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "speed_burst.enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "splash_page", "None"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ssid_admin_accessible", "false"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "visible", "true"),
				),
			},

			// Test the creation of multiple SSIDs.
			{
				Config: testAccNetworksWirelessSsidsResourceConfigMultiplePolicies(orgId, ssids),
				Check: func(s *terraform.State) error {
					var checks []resource.TestCheckFunc
					// Dynamically generate checks for each SSID
					for i := 1; i <= ssids; i++ {
						resourceName := fmt.Sprintf("meraki_networks_wireless_ssids.test%d", i)
						expectedNumber := fmt.Sprintf("%d", i-1) // Assuming numbering starts from 0
						checks = append(checks,
							resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("SSID %d", i)),
							resource.TestCheckResourceAttr(resourceName, "number", expectedNumber),
							resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
							resource.TestCheckResourceAttr(resourceName, "auth_mode", "open"),
							resource.TestCheckResourceAttr(resourceName, "adult_content_filtering_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "availability_tags.#", "0"),
							resource.TestCheckResourceAttr(resourceName, "available_on_all_aps", "true"),
							resource.TestCheckResourceAttr(resourceName, "band_selection", "Dual band operation"),
							resource.TestCheckResourceAttr(resourceName, "dns_rewrite.dns_custom_nameservers.#", "0"),
							resource.TestCheckResourceAttr(resourceName, "dns_rewrite.enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "ip_assignment_mode", "NAT mode"),
							resource.TestCheckResourceAttr(resourceName, "mandatory_dhcp_enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "min_bit_rate", "11"),
							resource.TestCheckResourceAttr(resourceName, "per_client_bandwidth_limit_down", "0"),
							resource.TestCheckResourceAttr(resourceName, "per_client_bandwidth_limit_up", "0"),
							resource.TestCheckResourceAttr(resourceName, "per_ssid_bandwidth_limit_down", "0"),
							resource.TestCheckResourceAttr(resourceName, "per_ssid_bandwidth_limit_up", "0"),
							resource.TestCheckResourceAttr(resourceName, "speed_burst.enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "splash_page", "None"),
							resource.TestCheckResourceAttr(resourceName, "ssid_admin_accessible", "false"),
							resource.TestCheckResourceAttr(resourceName, "visible", "true"),
						)
					}
					return resource.ComposeAggregateTestCheckFunc(checks...)(s)
				},
			},

			//// Import test
			//{
			//	ResourceName:      "meraki_networks_wireless_ssids.test",
			//	ImportState:       true,
			//	ImportStateVerify: false,
			//	ImportStateId:     "1234567890, 0987654321",
			//},
		},
	})
}

// testAccNetworksWirelessSsidsResourceConfigCreateNetwork is a function which returns a string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_resource"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksWirelessSsidsResourceConfigBasic = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
	number = 0
    name = "My SSID"
    enabled = true
    auth_mode = "open"
	gre = {
		concentrator = {
			host = "Test Host"
		}
		key = 0
	}
    adult_content_filtering_enabled = false
    availability_tags = []
    available_on_all_aps = true
    band_selection = "Dual band operation"
    dns_rewrite = {
      dns_custom_nameservers = []
      enabled = false
    }
    ip_assignment_mode = "NAT mode"
    mandatory_dhcp_enabled = false
    min_bit_rate = 11
    per_client_bandwidth_limit_down = 0
    per_client_bandwidth_limit_up = 0
    per_ssid_bandwidth_limit_down = 0
    per_ssid_bandwidth_limit_up = 0
    speed_burst = {
      enabled = false
    }
    splash_page = "None"
    ssid_admin_accessible = false
    visible = true
}
`

func testAccNetworksWirelessSsidsResourceConfigMultiplePolicies(orgId string, ssids int) string {
	config := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_resource"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)

	// Append each ssid configuration
	for i := 1; i <= ssids; i++ {
		config += fmt.Sprintf(`
resource "meraki_networks_wireless_ssids" "test%d" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.id
	number = %d
    name = "SSID %d"
	enabled = true
    auth_mode = "open"
    adult_content_filtering_enabled = false
    availability_tags = []
    available_on_all_aps = true
    band_selection = "Dual band operation"
	gre = {
		concentrator = {
			host = "Test Host"
		}
		key = 0
	}
    dns_rewrite = {
      dns_custom_nameservers = []
      enabled = false
    }
    ip_assignment_mode = "NAT mode"
    mandatory_dhcp_enabled = false
    min_bit_rate = 11
    per_client_bandwidth_limit_down = 0
    per_client_bandwidth_limit_up = 0
    per_ssid_bandwidth_limit_down = 0
    per_ssid_bandwidth_limit_up = 0
    speed_burst = {
      enabled = false
    }
    splash_page = "None"
    ssid_admin_accessible = false
    visible = true
}
`, i, i-1, i)
	}

	return config
}
