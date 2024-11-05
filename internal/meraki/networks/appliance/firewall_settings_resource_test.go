package appliance_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceFirewallSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkApplianceFirewallSettingsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_firewall_settings"),
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

			// Update and Read Networks Appliance Firewall Settings.
			{
				Config: testAccNetworkApplianceFirewallSettingsResourceConfigUpdateNetworkApplianceFirewallSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_settings.test", "spoofing_protection.ip_source_guard.mode", "block"),
				),
			},
		},
	})
}

func testAccNetworkApplianceFirewallSettingsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_appliance_firewall_settings"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworkApplianceFirewallSettingsResourceConfigUpdateNetworkApplianceFirewallSettings = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_firewall_settings" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    spoofing_protection = {
        ip_source_guard =  {
            mode = "block"
        }
    }
}
`
