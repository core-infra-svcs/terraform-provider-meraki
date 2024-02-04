package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksAppliancePortResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworksAppliancePortResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_ports"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim Appliance To Network
			{
				Config: testAccNetworksAppliancePortResourceConfigClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: testAccNetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans_settings.test", "vlans_enabled", "true"),
				),
			},

			//  Configure Appliance Port
			{
				Config: testAccNetworksAppliancePortResourceConfigListNetworkAppliancePorts,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_ports.test", "port_id", "4"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_ports.test", "allowed_vlans", "all"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_ports.test", "drop_untagged_traffic", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_ports.test", "enabled", "true"),
					//resource.TestCheckResourceAttr("meraki_networks_appliance_ports.test", "number", "3"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_ports.test", "type", "trunk"),
				),
			},
		},
	})
}

// testAccNetworksAppliancePortResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource. This will not work if the network already exists
func testAccNetworksAppliancePortResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_networks_appliance_ports"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesResourceConfigClaimNetworkDevice is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccNetworksAppliancePortResourceConfigClaimNetworkDevice(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	product_types = ["appliance"]
}

resource "meraki_networks_devices_claim" "test" {
    depends_on = ["resource.meraki_network.test"]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
`, serial)
	return result
}

// testAccNetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance"]
}
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_networks_appliance_vlans_settings" "test" {
		depends_on = ["resource.meraki_network.test", "resource.meraki_networks_devices_claim.test"]
		network_id = resource.meraki_network.test.network_id
		vlans_enabled = true
}
`

// testAccNetworksAppliancePortResourceConfigListNetworkAppliancePort is a constant string that defines the configuration for reading a networks_appliance_ports datasource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksAppliancePortResourceConfigListNetworkAppliancePorts = `
resource "meraki_network" "test" {
	product_types = ["appliance"]
}
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}
resource "meraki_networks_appliance_vlans_settings" "test" {
	network_id = resource.meraki_network.test.network_id
	vlans_enabled = true
}
	
resource "meraki_networks_appliance_ports" "test" {
	depends_on = ["resource.meraki_network.test", "resource.meraki_networks_devices_claim.test", "resource.meraki_networks_appliance_vlans_settings.test"]
	network_id = resource.meraki_network.test.network_id
	port_id = 4
	allowed_vlans = "all"
	drop_untagged_traffic = true
	enabled = true	
	type = "trunk"	
    }
`
