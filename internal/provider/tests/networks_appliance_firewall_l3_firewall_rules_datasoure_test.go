package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksApplianceFirewallL3FirewallRulesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksApplianceFirewallL3FirewallRulesDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksApplianceFirewallL3FirewallRulesDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_firewall_l3_firewall_rules_datasources"),
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
				Config: testAccNetworksApplianceFirewallL3FirewallRulesDataResourceConfigUpdateNetworksApplianceFirewallL3FirewallRules,
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

			// Read L3 Firewall Rules
			{
				Config: testAccNetworksApplianceFirewallL3FirewallRulesDataResourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.src_port", "Any"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.src_cidr", "Any"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", "rules.0.syslog_enabled", "false"),
				),
			},
		},
	})
}

// testAccNetworksApplianceFirewallL3FirewallRulesDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksApplianceFirewallL3FirewallRulesDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_appliance_firewall_l3_firewall_rules_datasources"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksApplianceFirewallL3FirewallRulesDataResourceConfigUpdateNetworksApplianceFirewallL3FirewallRules = `

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

// testAccNetworksApplianceFirewallL3FirewallRulesDataResourceConfigRead is a constant string that defines the configuration for creating and updating a devices_switch_ports_dataSource resource in your tests.
// It depends on both the organization and network resources.
func testAccNetworksApplianceFirewallL3FirewallRulesDataResourceConfigRead(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	    organization_id = "%s"
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
data "meraki_networks_appliance_firewall_l3_firewall_rules" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_appliance_firewall_l3_firewall_rules.test]
	network_id = resource.meraki_network.test.network_id
}
`, orgId)
	return result
}
