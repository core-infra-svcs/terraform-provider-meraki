package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceFirewallSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworkApplianceFirewallSettingsResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_network_appliance_firewall_settings"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworkApplianceFirewallSettingsResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
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

const testAccNetworkApplianceFirewallSettingsResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
    name = "test_acc_meraki_organizations_network_appliance_firewall_settings"
    api_enabled = true
}
`
const testAccNetworkApplianceFirewallSettingsResourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    organization_id = resource.meraki_organization.test.organization_id
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_network"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`

const testAccNetworkApplianceFirewallSettingsResourceConfigUpdateNetworkApplianceFirewallSettings = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_firewall_settings" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
    network_id = resource.meraki_network.test.network_id
    spoofing_protection = {
        ip_source_guard =  {
            mode = "block"
        }
    }
}
`
