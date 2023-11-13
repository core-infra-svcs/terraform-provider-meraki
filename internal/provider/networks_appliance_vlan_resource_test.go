package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksApplianceVlanResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksApplianceVlanResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
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

			// Update testing
			{
				Config: testAccNetworksApplianceVlanResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "name", "My VLAN"),
				),
			},
		},
	})
}

// testAccNetworksApplianceVlanResourceConfigCreate is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksApplianceVlanResourceConfigCreate(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
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

const testAccNetworksApplianceVlanResourceConfigUpdate = `

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
	vlan_id = "123"
    name = "My VLAN"
    subnet = "192.168.1.0/24"
    appliance_ip = "192.168.1.2"
    template_vlan_type = "same"
    cidr = "192.168.1.0/24"
    mask = 24

	reserved_ip_ranges = []

	dhcp_options = []

	fixed_ip_assignments = {}

    ipv6 = {
        enabled = false,
        prefix_assignments = []
    }

	mandatory_dhcp = {
		enabled = false
	}

}
`
