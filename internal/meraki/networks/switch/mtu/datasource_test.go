package mtu_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchMtuDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_mtu"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_switch_mtu"),
			},

			// Update Networks Switch Mtu.
			{
				Config: SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettings(),
				Check:  SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettingsChecks(),
			},

			// Update Switch MTU
			{
				Config: NetworkSwitchMtuDataSourceRead(),
				Check:  NetworkSwitchMtuDataSourceReadChecks(),
			},
		},
	})
}

func NetworkSwitchMtuDataSourceRead() string {
	return fmt.Sprintf(`
	%s

data "meraki_networks_switch_mtu" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_mtu"),
	)
}

// NetworkSwitchMtuDataSourceReadChecks returns the test check functions for NetworkSwitchMtuDataSourceRead
func NetworkSwitchMtuDataSourceReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"default_mtu_size": "9578",
	}
	return utils.ResourceTestCheck("data.meraki_networks_switch_mtu.test", expectedAttrs)
}
