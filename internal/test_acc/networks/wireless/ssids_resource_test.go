package wireless

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccNetworksWirelessSsidsResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	ssids := 10 // Number of SSIDs to create, Meraki max is 15
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
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

			// Create and Read testing without encryption
			{
				Config: testAccNetworksWirelessSsidsResourceConfigBasic(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "number", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "name", "My SSID"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "auth_mode", "psk"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "psk", "deadbeef"),
				),
			},

			// Create and Read testing with encryption
			{
				Config: testAccNetworksWirelessSsidsResourceConfigBasic(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "number", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "name", "My SSID"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "auth_mode", "psk"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "psk", "deadbeef"),
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
func testAccNetworksWirelessSsidsResourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	name = "test_acc_networks_wireless_ssids_resource"
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
  

`, orgId)
}

func testAccNetworksWirelessSsidsResourceConfigBasic(encryption bool) string {
	if encryption {
		return `
provider "meraki" {
  encryption_key = "my_secret_encryption_key"
}

resource "meraki_network" "test" {
  organization_id = "%s"
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
	} else {
		return `
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
	}
}

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

/*

	// Active Directory Authentication Test
		{
			Config: testAccNetworksWirelessSsidsResourceConfigAD,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.active_directory", "number", "1"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.active_directory", "name", "AD_SSID"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.active_directory", "enabled", "true"),
			),
		},

		// VLAN and Bandwidth Limits Test
		{
			Config: testAccNetworksWirelessSsidsResourceConfigVLANBandwidth,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.vlan_bandwidth", "number", "2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.vlan_bandwidth", "name", "VLANBandwidth"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.vlan_bandwidth", "enabled", "true"),
			),
		},

		// Full Configuration Test
		{
			Config: testAccNetworksWirelessSsidsResourceConfigFullConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.full_config", "number", "3"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.full_config", "name", "FullConfigSSID"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.full_config", "enabled", "true"),
			),
		},

		// Guest Access and Walled Garden Test
		{
			Config: testAccNetworksWirelessSsidsResourceConfigGuestAccess,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.guest_access", "number", "4"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.guest_access", "name", "GuestAccess"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.guest_access", "enabled", "true"),
			),
		},



*/

/*
const testAccNetworksWirelessSsidsResourceConfigAD = `
resource "meraki_network" "test" {
  product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "active_directory" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 1
	name = "AD_SSID"
	auth_mode = "psk"
	encryption_mode = "wpa"
	psk = "deadbeef"
	wpa_encryption_mode = "WPA2 only"
	splash_page = "Password-protected with Active Directory"
	active_directory = {
		credentials = {
			login_name = "user@example.com"
			password = "password"
		}
		servers = [
			{
				host = "192.168.1.1"
				port = 389
			}
		]
	}
	enabled = true
}
`

const testAccNetworksWirelessSsidsResourceConfigVLANBandwidth = `
resource "meraki_network" "test" {
  product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "vlan_bandwidth" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 2
	name = "VLANBandwidth"
	auth_mode = "psk"
	use_vlan_tagging = true
	vlan_id = 100
	enabled = true
	per_client_bandwidth_limit_down = 5000  // 5 Mbps
	per_client_bandwidth_limit_up   = 1000  // 1 Mbps
}
`

const testAccNetworksWirelessSsidsResourceConfigFullConfig = `
resource "meraki_network" "test" {
  product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "full_config" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 3
	name = "FullConfigSSID"
	auth_mode = "8021x-radius"
	radius_servers = {
		host = "radius.example.com"
		secret = "radiussecret"
	}
	ldap = {
		base_distinguished_name = "dc=example,dc=com"
		credentials = {
			distinguished_name = "cn=admin,dc=example,dc=com"
			password = "ldappassword"
		}
	}
	dns_rewrite = {
		enabled = true
		dns_custom_name_servers = ["8.8.8.8", "8.8.4.4"]
	}
	enabled = true
}
`

const testAccNetworksWirelessSsidsResourceConfigGuestAccess = `
resource "meraki_network" "test" {
  product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "guest_access" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 4
	name = "GuestAccess"
	splash_page = "Click-through splash page"
	walled_garden_enabled = true
	walled_garden_ranges = ["192.168.100.0/24", "www.example.com"]
	enabled = true
}
`
*/
