package firewall_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksApplianceFirewallL3FirewallRulesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksApplianceFirewallL3FirewallRulesDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_l3_firewall_rules_datasources"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_firewall_l3_firewall_rules_datasources"),
			},

			// Create and Read Network Settings.
			{
				Config: NetworksApplianceL3FirewallRulesResourceConfigCreate(),
				Check:  NetworksApplianceL3FirewallRulesResourceConfigCreateChecks(),
			},

			// Read L3 Firewall Rules
			{
				Config: NetworksApplianceFirewallL3FirewallRulesDataSourceConfigRead(),
				Check:  NetworksApplianceFirewallL3FirewallRulesDataSourceConfigReadChecks(),
			},
		},
	})
}

func NetworksApplianceFirewallL3FirewallRulesDataSourceConfigRead() string {
	return fmt.Sprintf(`
	%s
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
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_l3_firewall_rules_datasources"),
	)
}

// NetworksApplianceFirewallL3FirewallRulesDataSourceConfigReadChecks returns the test check functions for NetworksApplianceFirewallL3FirewallRulesDataSourceConfigRead
func NetworksApplianceFirewallL3FirewallRulesDataSourceConfigReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.comment":        "Allow TCP traffic to subnet with HTTP servers.",
		"rules.0.policy":         "allow",
		"rules.0.protocol":       "tcp",
		"rules.0.dest_port":      "443",
		"rules.0.dest_cidr":      "192.168.1.0/24",
		"rules.0.src_port":       "Any",
		"rules.0.src_cidr":       "Any",
		"rules.0.syslog_enabled": "false",
	}
	return utils.ResourceTestCheck("data.meraki_networks_appliance_firewall_l3_firewall_rules.test", expectedAttrs)
}
