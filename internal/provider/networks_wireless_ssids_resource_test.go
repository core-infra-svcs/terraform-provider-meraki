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
	ssids := 10 // Number of SSIDs to create, Meraki max is 15
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsResourceConfigCreateNetwork(orgId),
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
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "auth_mode", "psk"),
				),
			},

			// Import test
			{
				ResourceName:            "meraki_networks_wireless_ssids.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
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
							resource.TestCheckResourceAttr(resourceName, "auth_mode", "psk"),
						)
					}
					return resource.ComposeAggregateTestCheckFunc(checks...)(s)
				},
			},
		},
	})
}

// testAccNetworksWirelessSsidsResourceConfigCreateNetwork is a function which returns a string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsResourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
  name          = "test_acc_networks_wireless_ssids_resource"
  organization_id = "%s"
  timezone      = "America/Los_Angeles"
  tags          = ["tag1"]
  product_types = ["appliance", "switch", "wireless"]
  notes         = "Additional description of the network"
}
`, orgId)
}

const testAccNetworksWirelessSsidsResourceConfigBasic = `
resource "meraki_network" "test" {
  product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 0
	auth_mode = "psk"
	enabled = true
	encryption_mode = "wpa"
	name = "My SSID"
	psk = "deadbeef"
	wpa_encryption_mode = "WPA2 only"	
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
	network_id = resource.meraki_network.test.network_id
	number = %d
	name = "SSID %d"
	auth_mode = "psk"
	enabled = true
	encryption_mode = "wpa"
	psk = "deadbeef"
	wpa_encryption_mode = "WPA2 only"
}
`, i, i-1, i)
	}
	return config
}
