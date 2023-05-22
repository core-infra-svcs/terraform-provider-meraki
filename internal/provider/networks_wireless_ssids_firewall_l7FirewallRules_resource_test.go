package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksWirelessSsidsFirewallL7firewallrulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_wireless_ssids_firewall_l7_firewall_rules"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read testing
			//{
			//	Config: testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigCreate,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "id", "example-id"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "name", "Block sensitive web traffic"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "description", "Blocks sensitive web traffic"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "ip_version", "ipv6"),
			//
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.#", "1"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.policy", "deny"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.protocol", "tcp"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.src_port", "1,33"),
			//		// resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.dst_port", "22-30"),
			//	),
			//},

			// TODO - Once a resource has been created, we will test the ability to modify it. Make sure to test all values that are modifiable by the API call.
			// Update testing
			//			{
			//				Config: testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigUpdate,
			//				Check: resource.ComposeAggregateTestCheckFunc(
			//					resource.TestCheckResourceAttr("NetworksWirelessSsidsFirewallL7firewallrules.test", "id", "example-id"),
			//
			//                   // resource.TestCheckResourceAttr("data.NetworksWirelessSsidsFirewallL7firewallruless.test", "list.#", "2"),
			//
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.policy", "deny"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.protocol", "tcp"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.src_port", "1,33"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_rules.test", "rules.0.dst_port", "22-30"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_ruless.test", "list.1.rules.0.policy", "allow"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_ruless.test", "list.1.rules.0.protocol", "any"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_ruless.test", "list.1.rules.0.src_port", "any"),
			//                   // resource.TestCheckResourceAttr("networks_wireless_ssids_firewall_l7_firewall_ruless.test", "list.1.rules.0.dst_port", "any"),
			//				),
			//			},

			// TODO - ImportState testing - An import statement should ONLY include the required attributes to make a Read func call (example: organizationId + networkId).
			// TODO - Currently This only works with hard-coded values so if you find a dynamic way to test please update these template.
			/*
				{
						ResourceName:      "meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test",
						ImportState:       true,
						ImportStateVerify: false,
						ImportStateId:     "1234567890, 0987654321",
					},
			*/

			// TODO - Check your test environment for dangling resources. During the early stages of development it is not uncommon to find organizations,
			// TODO - networks or admins which did not get picked up because the resource errored out before the delete stage.
			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_wireless_ssids_firewall_l7_firewall_rules"
 	api_enabled = true
 }
 `

const testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "Main Office"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}
`

// TODO - Create your resource, make sure to include only the applicable attributes modifiable for CREATE.
const testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_networks_wireless_ssids_firewall_l7_firewall_rules" "test" {
	organization_id = resource.meraki_organization.test.organization_id
        name = "Block sensitive web traffic"
        description = "Blocks sensitive web traffic"
        ip_version   = "ipv6"
        rules = [
            {
                "policy": "deny",
                "protocol": "tcp",
                "src_port": "1,33",
                "dst_port": "22-30"
            }
        ]
    }
`

// TODO - Update the resource ensuring that all modifiable attributes are tested
/*
const testAccNetworksWirelessSsidsFirewallL7firewallrulesResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_networks_wireless_ssids_firewall_l7_firewall_rules" "test" {
	organization_id  = resource.meraki_organization.test.organization_id
    name = "Block sensitive web traffic"
    description = "Blocks sensitive web traffic"
    ip_version   = "ipv6"
    rules = [
        {
            "policy": "deny",
            "protocol": "tcp",
            "src_port": "1,33",
            "dst_port": "22-30"
        },
        {
            "policy": "allow",
            "protocol": "any",
            "src_port": "any",
            "dst_port": "any"
        }
    ]
  }
`
*/
