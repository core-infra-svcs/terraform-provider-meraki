package rules_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceFirewallL7FirewallRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_firewall_l7_firewall_rules.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_l7_firewall_rules"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_firewall_l7_firewall_rules"),
			},

			// Create and Read Networks Appliance Firewall L7 Firewall Rules.
			{
				Config: NetworksApplianceL7FirewallRulesResourceConfigCreate(),
				Check:  NetworksApplianceL7FirewallRulesResourceConfigCreateChecks(),
			},

			// Update and Read Networks Appliance Firewall L7 Firewall Rules.
			{
				Config: NetworksApplianceL7FirewallRulesResourceConfigUpdate(),
				Check:  NetworksApplianceL7FirewallRulesResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworksApplianceL7FirewallRulesResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_firewall_l7_firewall_rules" "test" {
    depends_on = [resource.meraki_network.test]
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
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_l7_firewall_rules"),
	)
}

// NetworksApplianceL7FirewallRulesResourceConfigCreateChecks returns the test check functions for NetworksApplianceL7FirewallRulesResourceConfigCreate
func NetworksApplianceL7FirewallRulesResourceConfigCreateChecks() resource.TestCheckFunc {
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
	return utils.ResourceTestCheck("meraki_networks_appliance_firewall_l7_firewall_rules.test", expectedAttrs)
}

func NetworksApplianceL7FirewallRulesResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_firewall_l7_firewall_rules" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    rules = [
    {
        policy = "deny",
        type = "host"
        value = "amazon.com"
    },
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
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_l7_firewall_rules"),
	)
}

// NetworksApplianceL7FirewallRulesResourceConfigUpdateChecks returns the test check functions for NetworksApplianceL7FirewallRulesResourceConfigUpdate
func NetworksApplianceL7FirewallRulesResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.0.policy": "deny",
		"rules.0.type":   "host",
		"rules.0.value":  "amazon.com",

		"rules.1.policy": "deny",
		"rules.1.type":   "host",
		"rules.1.value":  "google.com",

		"rules.2.policy": "deny",
		"rules.2.type":   "ipRange",
		"rules.2.value":  "10.11.12.00/24",

		"rules.3.policy": "deny",
		"rules.3.type":   "ipRange",
		"rules.3.value":  "10.11.12.00/24:5555",

		"rules.4.policy": "deny",
		"rules.4.type":   "port",
		"rules.4.value":  "23",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_firewall_l7_firewall_rules.test", expectedAttrs)
}
