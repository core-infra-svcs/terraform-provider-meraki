package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsApplianceVpnVpnFirewallRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsApplianceVpnVpnFirewallRulesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_appliance_vpn_vpn_firewall_rules"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccOrganizationsApplianceVpnVpnFirewallRulesResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_organizations_appliance_vpn_vpn_firewall_rules"),
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

			// Update and Read Organizations Appliance Vpn Vpn FirewallRules.
			{
				Config: testAccOrganizationsApplianceVpnVpnFirewallRulesResourceConfigUpdateOrganizationsApplianceVpnVpnFirewallRules,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.src_port", "Any"),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.src_cidr", "Any"),
					resource.TestCheckResourceAttr("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", "rules.0.syslog_enabled", "false"),
				),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_appliance_vpn_vpn_firewall_rules.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{
					// Add any attributes you want to ignore during import verification
				},
			},
		},
	})
}

const testAccOrganizationsApplianceVpnVpnFirewallRulesResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
    name = "test_acc_meraki_organizations_appliance_vpn_vpn_firewall_rules"
    api_enabled = true
}
`
const testAccOrganizationsApplianceVpnVpnFirewallRulesResourceConfigCreateNetwork = `
        resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    organization_id = resource.meraki_organization.test.organization_id
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_organizations_appliance_vpn_vpn_firewall_rules"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`

const testAccOrganizationsApplianceVpnVpnFirewallRulesResourceConfigUpdateOrganizationsApplianceVpnVpnFirewallRules = `
        resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}



resource "meraki_organizations_appliance_vpn_vpn_firewall_rules" "test" {
    depends_on = [resource.meraki_organization.test,
    resource.meraki_network.test]
    organization_id = resource.meraki_organization.test.organization_id
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
        syslog_enabled = true
    }
    ]
}
`
