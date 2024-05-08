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
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_meraki_networks_appliance_vlan"),
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
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vlan_id", "10"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "name", "My VLAN"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "appliance_ip", "192.168.1.2"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_relay_server_ips.#", "2"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_relay_server_ips.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_relay_server_ips.1", "192.168.128.0/24"),
					//resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vpn_nat_subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "mandatory_dhcp.enabled", "true"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksApplianceVlansResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vlan_id", "10"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "name", "My Updated VLAN"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "subnet", "192.168.2.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "appliance_ip", "192.168.2.2"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_handling", "Run a DHCP server"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_lease_time", "12 hours"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_boot_options_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_boot_next_server", "2.3.4.5"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_boot_filename", "updated.file"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "list.0.fixed_ip_assignments.%", "0"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "reserved_ip_ranges.0.start", "192.168.2.0"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "reserved_ip_ranges.0.end", "192.168.2.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "reserved_ip_ranges.0.comment", "A newly reserved IP range"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dns_nameservers", "upstream_dns"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_options.0.code", "6"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_options.0.type", "text"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "dhcp_options.0.value", "six"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "mandatory_dhcp.enabled", "true"),
				),
			},

			/*
				// Create and Read a VLAN IPv6
					{
						Config: testAccNetworksApplianceVlansResourceConfigCreateIPv6,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vlan_id", "20"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "name", "My IPv6 VLAN"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.enabled", "true"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.autonomous", "false"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.static_prefix", "2001:db8:3c4d:15::/64"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.static_appliance_ip6", "2001:db8:3c4d:15::1"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.origin.type", "internet"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.origin.interfaces.0", "wan1"),
						),
					},

					// Update testing IPv6
					{
						Config: testAccNetworksApplianceVlansResourceConfigUpdateIPv6,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "vlan_id", "20"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "name", "My Updated IPv6 VLAN"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.enabled", "true"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.autonomous", "true"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.static_prefix", "2001:db8:3c4d:16::/64"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.static_appliance_ip6", "2001:db8:3c4d:16::1"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.origin.type", "internet"),
							resource.TestCheckResourceAttr("meraki_networks_appliance_vlan.test", "ipv6.prefix_assignments.0.origin.interfaces.0", "wan1"),
						),
					},


			*/

			// Import testing
			{
				ResourceName:      "meraki_networks_appliance_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
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
name = "test_acc_meraki_networks_appliance_vlan"
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

resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_networks_appliance_vlans_settings.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "10"
    name = "My VLAN"
    subnet = "192.168.1.0/24"
	appliance_ip = "192.168.1.2"
	cidr = "192.168.1.0/24"
	mask = 28
    //vpn_nat_subnet = "192.168.1.0/24"
	dhcp_relay_server_ips = ["192.168.1.0/24", "192.168.128.0/24"]
    mandatory_dhcp = {
        enabled = true
    }

}
`

/*
// TODO: Figure out IPv6 dependencies
const testAccNetworksApplianceVlansResourceConfigCreateIPv6 = `
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
    vlan_id = "20"
    name = "My IPv6 VLAN"
    ipv6 = {
        enabled = true
        prefix_assignments = [
            {
                autonomous = false
                static_prefix = "2001:db8:3c4d:15::/64"
                static_appliance_ip6 = "2001:db8:3c4d:15::1"
                origin = {
                    type = "internet"
                    interfaces = ["wan1"]
                }
            }
        ]
    }

}
`
*/

const testAccNetworksApplianceVlansResourceConfigUpdate = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "10"
    name = "My Updated VLAN"
    subnet = "192.168.2.0/24"
    appliance_ip = "192.168.2.2"
    dhcp_handling = "Run a DHCP server"
    dhcp_lease_time = "12 hours"
    dhcp_boot_options_enabled = true
    dhcp_boot_next_server = "2.3.4.5"
    dhcp_boot_filename = "updated.file"
	fixed_ip_assignments = {
		"22:33:44:55:66:77": {
			"ip": "192.168.2.10",
			"name": "Some client name"
	  	}
	}
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

/*
const testAccNetworksApplianceVlansResourceConfigUpdateIPv6 = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "20"
    name = "My Updated IPv6 VLAN"
    ipv6 = {
        enabled = true
        prefix_assignments = [
            {
                autonomous = true
                static_prefix = "2001:db8:3c4d:16::/64"
                static_appliance_ip6 = "2001:db8:3c4d:16::1"
                origin = {
                    type = "internet"
                    interfaces = ["wan0"]
                }
            }
        ]
    }
}
`
*/
