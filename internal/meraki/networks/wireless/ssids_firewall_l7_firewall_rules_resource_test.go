package wireless_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
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
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_firewall_l7_firewall_rules"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_wireless_ssids_firewall_l7_firewall_rules"),
			},

			// Create and Read NetworksWirelessSsidsFirewallL7FirewallRules
			{
				Config: NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreate(),
				Check:  NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreateChecks(),
			},

			// Update and Read NetworksWirelessSsidsFirewallL7FirewallRules
			{
				Config: NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdate(),
				Check:  NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdateChecks(),
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

func NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
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
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_firewall_l7_firewall_rules"),
	)
}

// NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreateChecks returns the test check functions for NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreate
func NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.policy": "deny",
		"rules.0.type":   "host",
		"rules.0.value":  "google.com",

		"rules.1.policy": "deny",
		"rules.1.type":   "ipRange",
		"rules.1.value":  "10.11.12.00/24",

		"rules.2.policy": "deny",
		"rules.2.type":   "ipRange",
		"rules.2.value":  "10.11.12.00/24:5555",

		"rules.3.policy": "deny",
		"rules.3.type":   "port",
		"rules.3.value":  "23",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", expectedAttrs)
}

func NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
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
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_firewall_l7_firewall_rules"),
	)
}

// NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdateChecks returns the test check functions for NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdate
func NetworksWirelessSsidsFirewallL7FirewallRulesResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.policy": "deny",
		"rules.0.type":   "host",
		"rules.0.value":  "yahoo.com",

		"rules.1.policy": "deny",
		"rules.1.type":   "ipRange",
		"rules.1.value":  "10.11.13.00/24",

		"rules.2.policy": "deny",
		"rules.2.type":   "ipRange",
		"rules.2.value":  "10.13.12.00/24:5555",

		"rules.3.policy": "deny",
		"rules.3.type":   "port",
		"rules.3.value":  "43",
	}
	return utils.ResourceTestCheck("meraki_networks_wireless_ssids_firewall_l7_firewall_rules.test", expectedAttrs)
}
