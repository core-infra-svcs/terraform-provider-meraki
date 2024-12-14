package rules_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
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
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_firewall_l3_firewall_rules"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_wireless_ssids_firewall_l3_firewall_rules"),
			},

			// Create and Read NetworksWirelessSsidsFirewallL3FirewallRules
			{
				Config: NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreate(),
				Check:  NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreateChecks(),
			},

			// Update and Read NetworksWirelessSsidsFirewallL3FirewallRules
			{
				Config: NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdate(),
				Check:  NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdateChecks(),
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

func NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
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
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_firewall_l3_firewall_rules"),
	)
}

// NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreateChecks returns the test check functions for NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreate
func NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.comment":   "Allow TCP traffic to subnet with HTTP servers.",
		"rules.0.policy":    "allow",
		"rules.0.protocol":  "tcp",
		"rules.0.dest_port": "443",
		"rules.0.dest_cidr": "Any",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", expectedAttrs)
}

func NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
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
            comment = "Allow TCP traffic to subnet with HTTP servers duplicate.",
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
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_firewall_l3_firewall_rules"),
	)
}

// NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdateChecks returns the test check functions for NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdate
func NetworksWirelessSsidsFirewallL3FirewallRulesResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.comment":   "Allow TCP traffic to subnet with HTTP servers duplicate.",
		"rules.0.policy":    "allow",
		"rules.0.protocol":  "tcp",
		"rules.0.dest_port": "443",
		"rules.0.dest_cidr": "Any",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids_firewall_l3_firewall_rules.test", expectedAttrs)
}
