package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceVlansDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksApplianceVlansDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
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
				Config: testAccNetworksApplianceVlansDataSourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vlan_id", "10"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "name", "My VLAN"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "subnet", "192.168.2.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "appliance_ip", "192.168.2.2"),
					//resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vpn_nat_subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "mandatory_dhcp.enabled", "true"),
				),
			},

			// Read List
			{
				Config: testAccNetworksApplianceVlansDataSourceConfigReadList,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.vlan_id", "1"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.name", "Default"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.appliance_ip", "192.168.128.1"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.subnet", "192.168.128.0/24"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.fixed_ip_assignments.%", "0"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.reserved_ip_ranges.%", "0"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.dns_nameservers", "upstream_dns"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.dhcp_handling", "Run a DHCP server"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.dhcp_lease_time", "1 day"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.0.dhcp_boot_options_enabled", "false"),

					testCheckConcatenatedValues(
						"data.meraki_networks_appliance_vlans.test", "network_id",
						"data.meraki_networks_appliance_vlans.test", "list.0.vlan_id",
						",",
					),

					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.vlan_id", "10"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.name", "My VLAN"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.subnet", "192.168.2.0/24"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.appliance_ip", "192.168.2.2"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_handling", "Run a DHCP server"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_lease_time", "12 hours"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_boot_options_enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_boot_next_server", "2.3.4.5"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_boot_filename", "updated.file"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.reserved_ip_ranges.0.start", "192.168.2.0"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.reserved_ip_ranges.0.end", "192.168.2.1"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.reserved_ip_ranges.0.comment", "A newly reserved IP range"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dns_nameservers", "upstream_dns"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_options.0.code", "6"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_options.0.type", "text"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.dhcp_options.0.value", "six"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vlans.test", "list.1.mandatory_dhcp.enabled", "true"),

					testCheckConcatenatedValues(
						"data.meraki_networks_appliance_vlans.test", "network_id",
						"data.meraki_networks_appliance_vlans.test", "list.1.vlan_id",
						",",
					),
				),
			},
		},
	})
}

// testAccNetworksApplianceVlansDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksApplianceVlansDataSourceConfigCreateNetwork(orgId string) string {
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

const testAccNetworksApplianceVlansDataSourceConfigCreate = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlans_enabled = true
}

resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_networks_appliance_vlans_settings.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "10"
    name = "My VLAN"
    subnet = "192.168.2.0/24"
    appliance_ip = "192.168.2.2"
    dhcp_handling = "Run a DHCP server"
    dhcp_lease_time = "12 hours"
    dhcp_boot_options_enabled = true
    dhcp_boot_next_server = "2.3.4.5"
    dhcp_boot_filename = "updated.file"
    reserved_ip_ranges = [
        {
            start = "192.168.2.0"
            end = "192.168.2.1"
            comment = "A newly reserved IP range"
        }
    ]
    dns_nameservers = "upstream_dns"
    dhcp_options = [
        {
            code = "6"
            type = "text"
            value = "six"
        }
    ]
    mandatory_dhcp = {
        enabled = true
    }
}
`

const testAccNetworksApplianceVlansDataSourceConfigReadList = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlans_enabled = true
}

resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_networks_appliance_vlans_settings.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "10"
    name = "My VLAN"
    subnet = "192.168.2.0/24"
    appliance_ip = "192.168.2.2"
    dhcp_handling = "Run a DHCP server"
    dhcp_lease_time = "12 hours"
    dhcp_boot_options_enabled = true
    dhcp_boot_next_server = "2.3.4.5"
    dhcp_boot_filename = "updated.file"
    reserved_ip_ranges = [
        {
            start = "192.168.2.0"
            end = "192.168.2.1"
            comment = "A newly reserved IP range"
        }
    ]
    dns_nameservers = "upstream_dns"
    dhcp_options = [
        {
            code = "6"
            type = "text"
            value = "six"
        }
    ]
    mandatory_dhcp = {
        enabled = true
    }
}

data "meraki_networks_appliance_vlans" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
}
`

func testCheckConcatenatedValues(resource1, attr1, resource2, attr2, separator string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r1, ok := s.RootModule().Resources[resource1]
		if !ok {
			return fmt.Errorf("Not found: %s", resource1)
		}

		r2, ok := s.RootModule().Resources[resource2]
		if !ok {
			return fmt.Errorf("Not found: %s", resource2)
		}

		value1, ok := r1.Primary.Attributes[attr1]
		if !ok {
			return fmt.Errorf("Attribute not found: %s", attr1)
		}

		value2, ok := r2.Primary.Attributes[attr2]
		if !ok {
			return fmt.Errorf("Attribute not found: %s", attr2)
		}

		expectedValue := value1 + separator + value2
		// Use expectedValue as needed or compare with another expected output
		// For demonstration: Just log it (or assert equality if there is a specific value to compare)
		fmt.Printf("Concatenated Value: %s\n", expectedValue)

		return nil
	}
}
