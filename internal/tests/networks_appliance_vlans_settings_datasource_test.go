package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceVlanSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkApplianceVlanSettingsDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_vlan_settings"),
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

			// Create and Read Networks Appliance Vlans Settings.
			{
				Config: testAccNetworkApplianceVlanSettingsDataSourceConfigUpdateNetworkApplianceVlanSettingsFalse,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans_settings.test", "vlans_enabled", "false"),
				),
			},

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: testAccNetworkApplianceVlanSettingsDataSourceConfigUpdateNetworkApplianceVlanSettingsTrue,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans_settings.test", "vlans_enabled", "true"),
				),
			},

			// Read Networks Appliance Vlans Settings.
			{
				Config: testAccNetworkApplianceVlanSettingsDataSourceReadNetworkApplianceVlanSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans_settings_enabled.test", "vlans_enabled", "true"),
				),
			},
		},
	})
}

func testAccNetworkApplianceVlanSettingsDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
 resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_appliance_vlan_settings"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
 `, orgId)
	return result
}

const testAccNetworkApplianceVlanSettingsDataSourceConfigUpdateNetworkApplianceVlanSettingsFalse = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_vlans_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  vlans_enabled = false
}
`
const testAccNetworkApplianceVlanSettingsDataSourceConfigUpdateNetworkApplianceVlanSettingsTrue = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_vlans_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  vlans_enabled = true
}
`

const testAccNetworkApplianceVlanSettingsDataSourceReadNetworkApplianceVlanSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

data "meraki_networks_appliance_vlans_settings_enabled" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
}
`
