package mtu_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksSwitchMtuResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_switch_mtu.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_mtu"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_switch_mtu"),
			},

			// Update and Read Networks Switch Mtu.
			{
				Config: SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettings(),
				Check:  SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettingsChecks(),
			},
		},
	})
}

func SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettings() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_mtu" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    default_mtu_size = 9578
    overrides = []
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_mtu"),
	)
}

// SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettingsChecks returns the test check functions for SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettings
func SwitchMtuResourceConfigUpdateNetworkSwitchMtuSettingsChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"default_mtu_size": "9578",
	}
	return utils.ResourceTestCheck("meraki_networks_switch_mtu.test", expectedAttrs)
}
