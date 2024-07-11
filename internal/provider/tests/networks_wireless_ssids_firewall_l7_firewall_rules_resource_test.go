package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksWirelessSsidsFirewallL7FirewallRulesResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksWirelessSsidsFirewallL7FirewallRulesResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_wireless_ssids_firewall_l7_firewall_rules"),
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

			// Create and Read NetworksWirelessSsidsFirewallL7FirewallRules
			{
				Config: testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.1.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.1.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.1.value", "10.11.12.00/24"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.2.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.2.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.2.value", "10.11.12.00/24:5555"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.3.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.3.type", "port"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.3.value", "23"),
				),
			},

			// Update and Read NetworksWirelessSsidsFirewallL7FirewallRules
			{
				Config: testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.value", "yahoo.com"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.1.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.1.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.1.value", "10.11.13.00/24"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.2.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.2.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.2.value", "10.13.12.00/24:5555"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.3.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.3.type", "port"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.3.value", "43"),
				),
			},
		},
		// ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890, 0987654321",
		   },
		*/

	})
}

// testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_firewall_l7_firewall_rules"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreate is a constant string that defines the configuration for creating and reading a networks_wireless_ssids_firewall_l7FirewallRules resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreate = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids_firewall_l7_firewall_rules" "test" {
    depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
    number = "0"
    rules = [
        {
            policy =  "deny"
            type = "host"
            value = "google.com"
        },
        {
            policy = "deny"
            type = "port"
            value = "23"
        },
        {
            policy = "deny"
            type = "ipRange"
            value = "10.11.12.00/24"
        },
        {
            policy = "deny",
            type = "ipRange"
            value = "10.11.12.00/24:5555"
        }
        ]
}
`

// testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdate is a constant string that defines the configuration for updating a networks_wireless_ssids_firewall_l7FirewallRules resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdate = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids_firewall_l7_firewall_rules" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
    number = "0"
    rules = [
        {
            policy =  "deny"
            type = "host"
            value = "yahoo.com"
        },
        {
            policy = "deny"
            type = "port"
            value = "43"
        },
        {
            policy = "deny"
            type = "ipRange"
            value = "10.11.13.00/24"
        },
        {
            policy = "deny",
            type = "ipRange"
            value = "10.13.12.00/24:5555"
        }
        ]
}
`
