package settings_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_switch_settings.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_settings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_switch_settings"),
			},

			// Update and Read Network Switch Settings.
			{
				Config: NetworkSwitchSettingsResourceConfigUpdate(),
				Check:  NetworkSwitchSettingsResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworkSwitchSettingsResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  vlan = 100
	  use_combined_power = true
	  power_exceptions = []
	 
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_settings"),
	)
}

// NetworkSwitchSettingsResourceConfigUpdateChecks returns the test check functions for NetworkSwitchSettingsResourceConfigUpdate
func NetworkSwitchSettingsResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlan":               "100",
		"use_combined_power": "true",
	}
	return utils.ResourceTestCheck("meraki_networks_switch_settings.test", expectedAttrs)
}
