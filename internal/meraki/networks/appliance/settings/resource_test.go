package settings_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_settings.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_settings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_settings"),
			},

			// Update and Read Network Appliance Settings.
			{
				Config: NetworkApplianceSettingsResourceConfigUpdate(),
				Check:  NetworkApplianceSettingsResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworkApplianceSettingsResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  client_tracking_method = "MAC address"
	  deployment_mode = "routed"
	  dynamic_dns_prefix = "test"
	  dynamic_dns_enabled = true
	  
	 
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_settings"),
	)
}

// NetworkApplianceSettingsResourceConfigUpdateChecks returns the test check functions for NetworkApplianceSettingsResourceConfigUpdate
func NetworkApplianceSettingsResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"client_tracking_method": "MAC address",
		"deployment_mode":        "routed",
		"dynamic_dns_prefix":     "test",
		"dynamic_dns_enabled":    "true",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_settings.test", expectedAttrs)
}
