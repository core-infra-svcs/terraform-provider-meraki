package rules_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchQosRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_qos_rule"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_network_switch_qos_rule"),
			},

			// Create and Read Networks Switch Qos Rules.
			{
				Config: NetworkSwitchQosRuleResourceConfigCreate(),
				Check:  NetworkSwitchQosRuleResourceConfigCreateChecks(),
			},

			// Update Networks Switch Qos Rules.
			{
				Config: NetworkSwitchQosRuleResourceConfigUpdate(),
				Check:  NetworkSwitchQosRuleResourceConfigUpdateChecks(),
			},

			// Import testing
			{
				ResourceName:      "meraki_networks_switch_qos_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func NetworkSwitchQosRuleResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_qos_rule" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 100
    protocol = "TCP"
    src_port = 2000
    dst_port = 3000
    dscp = 0
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_qos_rule"),
	)
}

// NetworkSwitchQosRuleResourceConfigCreateChecks returns the test check functions for NetworkSwitchQosRuleResourceConfigCreate
func NetworkSwitchQosRuleResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlan":     "100",
		"dst_port": "3000",
		"src_port": "2000",
		"dscp":     "0",
		"protocol": "TCP",
	}
	return utils.ResourceTestCheck("meraki_networks_switch_qos_rule.test", expectedAttrs)
}

func NetworkSwitchQosRuleResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_qos_rule" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 101
    protocol = "UDP"
    src_port = 3000
    dst_port = 4000
    dscp = 0
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_qos_rule"),
	)
}

// NetworkSwitchQosRuleResourceConfigUpdateChecks returns the test check functions for NetworkSwitchQosRuleResourceConfigUpdate
func NetworkSwitchQosRuleResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlan":     "101",
		"dst_port": "4000",
		"src_port": "3000",
		"dscp":     "0",
		"protocol": "UDP",
	}
	return utils.ResourceTestCheck("meraki_networks_switch_qos_rule.test", expectedAttrs)
}
