package rules_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchQosRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_qos_rules"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_network_switch_qos_rules"),
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

			// Read testing
			{
				Config: NetworkSwitchQosRulesDataSourceRead(),
				Check:  NetworkSwitchQosRulesDataSourceReadChecks(),
			},
		},
	})
}

func NetworkSwitchQosRulesDataSourceRead() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_qos_rule" "tcp" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 100
    protocol = "TCP"
    src_port = 2000
    dst_port = 3000
    dscp = 0
}

resource "meraki_networks_switch_qos_rule" "udp" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 101
    protocol = "UDP"
    src_port = 3000
    dst_port = 4000
    dscp = 0
}

data "meraki_networks_switch_qos_rules" "test" {
    depends_on = [
		resource.meraki_network.test, meraki_networks_switch_qos_rule.udp, meraki_networks_switch_qos_rule.tcp
	]
    network_id = resource.meraki_network.test.network_id
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_qos_rule"),
	)
}

// NetworkSwitchQosRulesDataSourceReadChecks returns the test check functions for NetworkSwitchQosRulesDataSourceRead
func NetworkSwitchQosRulesDataSourceReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.0.vlan":     "100",
		"list.0.dst_port": "3000",
		"list.0.src_port": "2000",
		"list.0.dscp":     "0",
		"list.0.protocol": "TCP",

		"list.1.vlan":     "101",
		"list.1.dst_port": "4000",
		"list.1.src_port": "3000",
		"list.1.dscp":     "0",
		"list.1.protocol": "UDP",
	}
	return utils.ResourceTestCheck("data.meraki_networks_switch_qos_rules.test", expectedAttrs)
}
