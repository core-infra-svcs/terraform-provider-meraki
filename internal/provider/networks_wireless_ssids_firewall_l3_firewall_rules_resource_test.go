package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksWirelessSsidsFirewallL3FirewallRulesResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksWirelessSsidsFirewallL3FirewallRulesResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_wireless_ssids_firewall_l3_firewall_rules"),
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

			// Create and Read NetworksWirelessSsidsFirewallL3FirewallRules
			{
				Config: testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.dest_cidr", "Any"),
				),
			},

			// Update and Read NetworksWirelessSsidsFirewallL3FirewallRules
			{
				Config: testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", "rules.0.dest_cidr", "Any"),
				),
			},
		},

		// ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890, 0987654321",
		   },
		*/

	})
}

// testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_firewall_l3_firewall_rules"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreate is a constant string that defines the configuration for creating and reading a networks_wireless_ssids_firewall_l3FirewallRules resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreate = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_wireless_ssids_firewall_l3_firewall_rules" "test" {
    depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
    number = "0"
    rules = [
        {
            comment = "Allow TCP traffic to subnet with HTTP servers.",
            policy = "allow",
			ip_ver = "ipv4",
            protocol = "tcp",
            dest_port = "443",
            dest_cidr = "Any"
        },
		{
            comment = "Wireless clients accessing LAN",
            policy = "deny",
			ip_ver = "ipv4",
            protocol = "Any",
            dest_port = "Any",
            dest_cidr = "Local LAN"
        },
		{
            comment = "Default rule",
            policy = "allow",
			ip_ver = "ipv4",
            protocol = "Any",
            dest_port = "Any",
            dest_cidr = "Any"
        }
    ]
    }

`

// testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdate is a constant string that defines the configuration for updating a networks_wireless_ssids_firewall_l3FirewallRules resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdate = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids_firewall_l3_firewall_rules" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
    number = "0"
      rules = [
        {
            comment = "Allow TCP traffic to subnet with HTTP servers.",
            policy = "allow",
			ip_ver = "ipv4",
            protocol = "tcp",
            dest_port = "443",
            dest_cidr = "Any"
        },
		{
            comment = "Wireless clients accessing LAN",
            policy = "deny",
			ip_ver = "ipv4",
            protocol = "Any",
            dest_port = "Any",
            dest_cidr = "Local LAN"
        },
		{
            comment = "Default rule",
            policy = "allow",
			ip_ver = "ipv4",
            protocol = "Any",
            dest_port = "Any",
            dest_cidr = "Any"
        }
    ]
    }
`
