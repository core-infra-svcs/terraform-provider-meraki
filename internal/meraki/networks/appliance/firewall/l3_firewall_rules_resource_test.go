package firewall_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceFirewallL3FirewallRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_firewall_l3_firewall_rules"),
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
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.src_port", "Any"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.src_cidr", "Any"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.syslog_enabled", "false"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.

			{
				ResourceName:      "meraki_networks_appliance_firewall_l3_firewall_rules.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_appliance_firewall_l3_firewall_rules"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksApplianceFirewallL3FirewallRulesResourceConfigUpdateNetworksApplianceFirewallL3FirewallRules = `

resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_appliance_firewall_l3_firewall_rules" "test" {
    depends_on = [resource.meraki_network.test]
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
