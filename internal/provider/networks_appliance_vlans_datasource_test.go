package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TestAccNetworksApplianceVlansDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksApplianceVlansDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccNetworksApplianceVlansDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_appliance_vlans"),
				),
			},

			// Create and Read a Network.
			{
				Config: testAccNetworksApplianceVlansDataSourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
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

			// TODO: Create and Read NetworksApplianceVlans
			{
				Config: testAccNetworksApplianceVlansDataSourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "id", "example-id"),
					// TODO: Check the type and naming of the attribute "ApplianceIp".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "appliance_ip", "example-string"),
					// TODO: Check the type and naming of the attribute "Cidr".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "cidr", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpBootFilename".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_boot_filename", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpBootNextServer".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_boot_next_server", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpBootOptionsEnabled".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_boot_options_enabled", "true"),
					// TODO: Check the type and naming of the attribute "DhcpHandling".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_handling", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpLeaseTime".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_lease_time", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpOptions".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_options", "example-array"),
					// TODO: Check the type and naming of the attribute "DhcpRelayServerIps".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_relay_server_ips", "example-array"),
					// TODO: Check the type and naming of the attribute "DnsNameservers".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dns_nameservers", "example-string"),
					// TODO: Check the type and naming of the attribute "FixedIpAssignments".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "fixed_ip_assignments", "example-object"),
					// TODO: Check the type and naming of the attribute "GroupPolicyId".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "group_policy_id", "example-string"),
					// TODO: Check the type and naming of the attribute "Id".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "id", "example-string"),
					// TODO: Check the type and naming of the attribute "Ipv6Enabled".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "ipv6_enabled", "true"),
					// TODO: Check the type and naming of the attribute "Ipv6PrefixAssignments".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "ipv6_prefix_assignments", "example-array"),
					// TODO: Check the type and naming of the attribute "Mask".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "mask", "123"),
					// TODO: Check the type and naming of the attribute "Name".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "name", "example-string"),
					// TODO: Check the type and naming of the attribute "ReservedIpRanges".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "reserved_ip_ranges", "example-array"),
					// TODO: Check the type and naming of the attribute "Subnet".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "subnet", "example-string"),
					// TODO: Check the type and naming of the attribute "TemplateVlanType".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "template_vlan_type", "example-string"),
					// TODO: Check the type and naming of the attribute "VpnNatSubnet".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "vpn_nat_subnet", "example-string"),
				),
			},

			// TODO: Update and Read NetworksApplianceVlans
			{
				Config: testAccNetworksApplianceVlansDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "id", "example-id"),
					// TODO: Check the type and naming of the attribute "ApplianceIp".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "appliance_ip", "example-string"),
					// TODO: Check the type and naming of the attribute "Cidr".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "cidr", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpBootFilename".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_boot_filename", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpBootNextServer".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_boot_next_server", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpBootOptionsEnabled".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_boot_options_enabled", "true"),
					// TODO: Check the type and naming of the attribute "DhcpHandling".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_handling", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpLeaseTime".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_lease_time", "example-string"),
					// TODO: Check the type and naming of the attribute "DhcpOptions".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_options", "example-array"),
					// TODO: Check the type and naming of the attribute "DhcpRelayServerIps".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dhcp_relay_server_ips", "example-array"),
					// TODO: Check the type and naming of the attribute "DnsNameservers".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "dns_nameservers", "example-string"),
					// TODO: Check the type and naming of the attribute "FixedIpAssignments".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "fixed_ip_assignments", "example-object"),
					// TODO: Check the type and naming of the attribute "GroupPolicyId".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "group_policy_id", "example-string"),
					// TODO: Check the type and naming of the attribute "Id".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "id", "example-string"),
					// TODO: Check the type and naming of the attribute "Ipv6Enabled".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "ipv6_enabled", "true"),
					// TODO: Check the type and naming of the attribute "Ipv6PrefixAssignments".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "ipv6_prefix_assignments", "example-array"),
					// TODO: Check the type and naming of the attribute "Mask".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "mask", "123"),
					// TODO: Check the type and naming of the attribute "Name".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "name", "example-string"),
					// TODO: Check the type and naming of the attribute "ReservedIpRanges".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "reserved_ip_ranges", "example-array"),
					// TODO: Check the type and naming of the attribute "Subnet".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "subnet", "example-string"),
					// TODO: Check the type and naming of the attribute "TemplateVlanType".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "template_vlan_type", "example-string"),
					// TODO: Check the type and naming of the attribute "VpnNatSubnet".
					resource.TestCheckResourceAttr("networks_appliance_vlans.test", "vpn_nat_subnet", "example-string"),
				),
			},
		},

		// TODO: Finally, make sure there are no dangling resources in your test environment.
		// The resource.Test function automatically tests the Delete operation.
	})
}

// testAccNetworksApplianceVlansDataSourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccNetworksApplianceVlansDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_appliance_vlans"
 	api_enabled = true
 }
 `

// testAccNetworksApplianceVlansDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
const testAccNetworksApplianceVlansDataSourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

// TODO: Make a change to the configuration to test
// testAccNetworksApplianceVlansDataSourceConfigCreate is a constant string that defines the configuration for creating and updating a networks__appliance_vlans resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksApplianceVlansDataSourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
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
    mask = 28
	reserved_ip_ranges = [
		{
			start = "192.168.1.0"
			end = "192.168.1.1"
			comment = "A reserved IP range"
      	}
	]
	dhcp_options = [
		{
			code = "5"
			type = "text"
			value = "five"
      	}
    ]
	fixed_ip_assignments = {
	}
    ipv6 = {
        enabled = true
        prefix_assignments = []
    }
	mandatory_dhcp = {
		enabled = true
	}
}
`

// TODO: Make a change to the configuration to test
// testAccNetworksApplianceVlansDataSourceConfigRead is a constant string that defines the configuration for creating and updating a networks__appliance_vlans resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksApplianceVlansDataSourceConfigRead = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
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
    mask = 28
	reserved_ip_ranges = [
		{
			start = "192.168.1.0"
			end = "192.168.1.1"
			comment = "A reserved IP range"
      	}
	]
	dhcp_options = [
		{
			code = "5"
			type = "text"
			value = "five"
      	}
    ]
	fixed_ip_assignments = {
	}
    ipv6 = {
        enabled = true
        prefix_assignments = []
    }
	mandatory_dhcp = {
		enabled = true
	}
}

data "meraki_networks_appliance_vlans" "test" {
    depends_on = [resource.meraki_networks_appliance_vlans.test]
    network_id = resource.meraki_network.test.network_id
}
`
