package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsApplianceVpnVpnFirewallRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_appliance_vpn_vpn_firewall_rules"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_appliance_vpn_vpn_firewall_rules"),
			},

			// Create and Read Network
			{
				Config: utils.CreateNetworkConfig("test_acc_meraki_organizations_appliance_vpn_vpn_firewall_rules", "test_acc_organizations_appliance_vpn_vpn_firewall_rules"),
				Check:  utils.NetworkTestChecks("test_acc_organizations_appliance_vpn_vpn_firewall_rules"),
			},

			// Update and Read Organizations Appliance Vpn Vpn Firewall Rules
			{
				Config: OrganizationsApplianceVpnVpnFirewallRulesResourceConfigUpdate(),
				Check:  OrganizationsApplianceVpnVpnFirewallRulesTestChecks(),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_appliance_vpn_vpn_firewall_rules.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// OrganizationsApplianceVpnVpnFirewallRulesResourceConfigUpdate returns the configuration string for updating appliance VPN firewall rules
func OrganizationsApplianceVpnVpnFirewallRulesResourceConfigUpdate() string {
	rules := `[
		{
			"comment": "Allow TCP traffic to subnet with HTTP servers.",
			"policy": "allow",
			"protocol": "tcp",
			"dest_port": "443",
			"dest_cidr": "192.168.1.0/24",
			"src_port": "Any",
			"src_cidr": "Any",
			"syslog_enabled": false
		},
		{
			"comment": "Default rule",
			"policy": "allow",
			"protocol": "Any",
			"dest_port": "Any",
			"dest_cidr": "Any",
			"src_port": "Any",
			"src_cidr": "Any",
			"syslog_enabled": true
		}
	]`

	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_appliance_vpn_vpn_firewall_rules" "test" {
		depends_on = [
			resource.meraki_organization.test,
			resource.meraki_network.test
		]
		organization_id = resource.meraki_organization.test.organization_id
		rules           = %s
	}
	`,
		utils.CreateNetworkConfig("test_acc_meraki_organizations_appliance_vpn_vpn_firewall_rules", "test_acc_organizations_appliance_vpn_vpn_firewall_rules"),
		rules,
	)
}

// OrganizationsApplianceVpnVpnFirewallRulesTestChecks returns the test check functions for verifying appliance VPN firewall rules
func OrganizationsApplianceVpnVpnFirewallRulesTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.comment":        "Allow TCP traffic to subnet with HTTP servers.",
		"rules.0.policy":         "allow",
		"rules.0.protocol":       "tcp",
		"rules.0.dest_port":      "443",
		"rules.0.dest_cidr":      "192.168.1.0/24",
		"rules.0.src_port":       "Any",
		"rules.0.src_cidr":       "Any",
		"rules.0.syslog_enabled": "false",

		"rules.1.comment":        "Default rule",
		"rules.1.policy":         "allow",
		"rules.1.protocol":       "Any",
		"rules.1.dest_port":      "Any",
		"rules.1.dest_cidr":      "Any",
		"rules.1.src_port":       "Any",
		"rules.1.src_cidr":       "Any",
		"rules.1.syslog_enabled": "true",
	}

	return utils.ResourceTestCheck("meraki_organizations_appliance_vpn_vpn_firewall_rules.test", expectedAttrs)
}
