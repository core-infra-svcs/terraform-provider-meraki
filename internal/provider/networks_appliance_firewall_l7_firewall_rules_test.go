package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceFirewallL7FirewallRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksApplianceFirewallL7FirewallRulesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_appliance_firewall_l7_firewall_rules"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_firewall_l7_firewall_rules.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceFirewallL7FirewallRulesResourceConfigCreateNetwork,
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

			// Create and Read Networks Appliance Firewall L7 Firewall Rules.
			{
				Config: testAccNetworksApplianceFirewallL7FirewallRulesResourceConfigCreateNetworksApplianceFirewallL7FirewallRules,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.1.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.1.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.1.value", "10.11.12.00/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.2.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.2.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.2.value", "10.11.12.00/24:5555"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.3.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.3.type", "port"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l7_firewall_rules.test", "rules.3.value", "23"),
				),
			},
		},
	})
}

const testAccNetworksApplianceFirewallL7FirewallRulesResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_networks_appliance_firewall_l7_firewall_rules"
 	api_enabled = true
 }
 `
const testAccNetworksApplianceFirewallL7FirewallRulesResourceConfigCreateNetwork = `
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

const testAccNetworksApplianceFirewallL7FirewallRulesResourceConfigCreateNetworksApplianceFirewallL7FirewallRules = `
 resource "meraki_organization" "test" {}
 
 resource "meraki_network" "test" {
	 depends_on = [resource.meraki_organization.test]	
	 product_types = ["appliance", "switch", "wireless"]	
 }
 
 resource "meraki_networks_appliance_firewall_l7_firewall_rules" "test" {
	 depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
	 network_id = resource.meraki_network.test.network_id
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
