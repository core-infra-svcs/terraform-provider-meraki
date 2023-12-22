package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksApplianceVlansResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksApplianceVlansResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_meraki_networks_appliance_vlans"),
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

			// Create and Read a VLAN
			{
				Config: testAccNetworksApplianceVlansResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans.test", "name", "My VLAN"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksApplianceVlansResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans.test", "name", "My VLAN 2"),
				),
			},
		},
	})
}

// testAccNetworksApplianceVlansResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksApplianceVlansResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
organization_id = %s
product_types = ["appliance", "switch", "wireless"]
tags = ["tag1"]
name = "test_acc_meraki_networks_appliance_vlans"
timezone = "America/Los_Angeles"
notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksApplianceVlansResourceConfigCreate = `

resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	vlans_enabled = true
}

resource "meraki_networks_appliance_vlans" "test" {
	depends_on = [resource.meraki_networks_appliance_vlans_settings.test]
	network_id = resource.meraki_network.test.network_id
	vlan_id = "10"
    name = "My VLAN"
    subnet = "192.168.1.0/24"
    appliance_ip = "192.168.1.2"
}
`

const testAccNetworksApplianceVlansResourceConfigUpdate = `

resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_vlans" "test" {
	network_id = resource.meraki_network.test.network_id
	vlan_id = "10"
    name = "My VLAN 2"
    subnet = "192.168.1.0/24"
    appliance_ip = "192.168.1.2"
}
`
