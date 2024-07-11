package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_settings.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworkApplianceSettingsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_settings"),
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

			// Update and Read Network Appliance Settings.
			{
				Config: testAccNetworkApplianceSettingsResourceConfigUpdateNetworkApplianceSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_settings.test", "client_tracking_method", "MAC address"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_settings.test", "deployment_mode", "routed"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_settings.test", "dynamic_dns_prefix", "test"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_settings.test", "dynamic_dns_enabled", "true"),
				),
			},
		},
	})
}

func testAccNetworkApplianceSettingsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
 resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_appliance_settings"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
 `, orgId)
	return result
}

const testAccNetworkApplianceSettingsResourceConfigUpdateNetworkApplianceSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  client_tracking_method = "MAC address"
	  deployment_mode = "routed"
	  dynamic_dns_prefix = "test"
	  dynamic_dns_enabled = true
	  
	 
}
`
