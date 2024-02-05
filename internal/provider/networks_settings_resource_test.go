package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkSettingsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_settings"),
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

			// Update and Read Network Settings.
			{
				Config: testAccNetworkSettingsResourceConfigUpdateNetworkSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_settings.test", "local_status_page_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_settings.test", "remote_status_page_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_settings.test", "local_status_page.authentication.enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_settings.test", "local_status_page.authentication.username", "admin"),
					resource.TestCheckResourceAttr("meraki_networks_settings.test", "local_status_page.authentication.password", "testpassword"),
					resource.TestCheckResourceAttr("meraki_networks_settings.test", "secure_port_enabled", "false"),
				),
			},

			/*
				{
						ResourceName:      "meraki_networks_settings.test",
						ImportState:       true,
						ImportStateVerify: true,
					},
			*/
		},
	})
}

func testAccNetworkSettingsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_settings"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworkSettingsResourceConfigUpdateNetworkSettings = `

resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]	
}

resource "meraki_networks_settings" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
	local_status_page_enabled = false
	remote_status_page_enabled = false
	local_status_page = { 
		authentication = { 
			enabled = true
			username = "admin"
			password = "testpassword"
	}
	}
  	secure_port_enabled = false
	named_vlans_enabled = true
}
`
