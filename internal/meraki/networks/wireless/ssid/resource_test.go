package ssid_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccNetworksWirelessSsidsResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	ssids := 10 // Number of SSIDs to create, Meraki max is 15
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_wireless_ssids_resource"),
			},

			// Create and Read SSID without encryption
			{
				Config: NetworksWirelessSsidsResourceConfigBasic(false),
				Check:  NetworksWirelessSsidsResourceConfigBasicChecks(),
			},

			// Create and Read SSID with encryption
			{
				Config: NetworksWirelessSsidsResourceConfigBasic(true),
				Check:  NetworksWirelessSsidsResourceConfigBasicChecks(),
			},

			//TODO: ImportState test case.
			{
				ResourceName:      "meraki_networks_wireless_ssids.test",
				ImportState:       true,
				ImportStateVerify: true,
			},

			//Test RADIUS servers creation
			{
				Config: NetworksWirelessSsidsResourceConfigRadiusServers(),
				Check:  NetworksWirelessSsidsResourceConfigRadiusServersChecks(),
			},

			// Test RADIUS updating
			{
				Config: NetworksWirelessSsidsResourceConfigRadiusServersUpdate(),
				Check:  NetworksWirelessSsidsResourceConfigRadiusServersUpdateChecks(),
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
							resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
							resource.TestCheckResourceAttr(resourceName, "auth_mode", "psk"),
						)
					}
					return resource.ComposeAggregateTestCheckFunc(checks...)(s)
				},
			},
		},
	})
}

func NetworksWirelessSsidsResourceConfigBasic(encryption bool) string {
	if encryption {
		return fmt.Sprintf(`
provider "meraki" {
  encryption_key = "my_secret_encryption_key"
}

	%s

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
	radius_servers = []
}
	
	`,
			utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
		)
	} else {
		return fmt.Sprintf(`
	%s

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
	radius_servers = []
}
	
	`,
			utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
		)
	}
}

// NetworksWirelessSsidsResourceConfigBasicChecks returns the test check functions for NetworksWirelessSsidsResourceConfigBasic
func NetworksWirelessSsidsResourceConfigBasicChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"number":    "0",
		"name":      "My SSID",
		"enabled":   "true",
		"auth_mode": "psk",
		"psk":       "deadbeef",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids.test", expectedAttrs)
}

func NetworksWirelessSsidsResourceConfigRadiusServers() string {
	return fmt.Sprintf(`
provider "meraki" {
  encryption_key = "my_secret_encryption_key"
}

	%s

resource "meraki_networks_wireless_ssids" "test_radius" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 1
	auth_mode = "8021x-radius"
	enabled = true
	encryption_mode = "wpa-eap"
	name = "My Radius SSID TEST"
	wpa_encryption_mode = "WPA2 only"
	radius_servers = [{
		host = "radius.example.com"
		port = 1812
		secret = "radius_secret"
		rad_sec_enabled = true
		ca_certificate = "ca_cert_value"
	}]
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
	)
}

// NetworksWirelessSsidsResourceConfigRadiusServersChecks returns the test check functions for NetworksWirelessSsidsResourceConfigRadiusServers
func NetworksWirelessSsidsResourceConfigRadiusServersChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"radius_servers.0.host":            "radius.example.com",
		"radius_servers.0.port":            "1812",
		"radius_servers.0.secret":          "radius_secret",
		"radius_servers.0.rad_sec_enabled": "true",
		"radius_servers.0.ca_certificate":  "ca_cert_value",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids.test_radius", expectedAttrs)
}

func NetworksWirelessSsidsResourceConfigRadiusServersUpdate() string {
	return fmt.Sprintf(`
provider "meraki" {
  encryption_key = "my_secret_encryption_key"
}

	%s

resource "meraki_networks_wireless_ssids" "test_radius" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = 1
	auth_mode = "8021x-radius"
	enabled = true
	encryption_mode = "wpa-eap"
	name = "My Radius SSID"
	wpa_encryption_mode = "WPA2 only"
	radius_servers = [{
		host = "radius.example.com"
		port = 1812
		secret = "new_radius_secret"
		rad_sec_enabled = true
		ca_certificate = "new_ca_cert_value"
	}]
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
	)
}

// NetworksWirelessSsidsResourceConfigRadiusServersUpdateChecks returns the test check functions for NetworksWirelessSsidsResourceConfigRadiusServersUpdate
func NetworksWirelessSsidsResourceConfigRadiusServersUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"radius_servers.0.host":            "radius.example.com",
		"radius_servers.0.port":            "1812",
		"radius_servers.0.secret":          "new_radius_secret",
		"radius_servers.0.rad_sec_enabled": "true",
		"radius_servers.0.ca_certificate":  "new_ca_cert_value",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids.test_radius", expectedAttrs)
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
	enabled = false
	encryption_mode = "wpa"
	psk = "deadbeef"
	wpa_encryption_mode = "WPA2 only"
}
`, i, i-1, i)
	}
	return config
}
