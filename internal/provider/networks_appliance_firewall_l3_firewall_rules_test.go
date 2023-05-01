package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceFirewallL3FirewallRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_appliance_firewall_l3_firewall_rules"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_firewall_l3_firewall_rules.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigCreateNetwork,
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

			// Update and Read Network Settings.
			{
				Config: testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigUpdateNetworksApplianceFirewallL3FirewallRules,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.src_cidr", "Any"),
				),
			},
		},
	})
}

const testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
    name = "test_acc_meraki_organizations_networks_appliance_firewall_l3_firewall_rules"
    api_enabled = true
}
`
const testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigCreateNetwork = `
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

const testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigUpdateNetworksApplianceFirewallL3FirewallRules = `
        resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_firewall_l3_firewall_rules" "test" {
    depends_on = [resource.meraki_organization.test,
    resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    syslog_default_rule = false
    rules = [
    {
        comment =  "Allow TCP traffic to subnet with HTTP servers."
        policy = "allow"
        protocol = "tcp"
        dest_port = "443"
        dest_cidr = "192.168.1.0/24"
        src_port = "Any"
        src_cidr = "Any"
        syslog_enabled = false
    },
    {
        comment =  "Default rule"
        policy = "allow"
        protocol = "Any"
        dest_port = "Any"
        dest_cidr = "Any"
        src_port = "Any"
        src_cidr = "Any"
        syslog_enabled = false
    }

    ]

}
`
