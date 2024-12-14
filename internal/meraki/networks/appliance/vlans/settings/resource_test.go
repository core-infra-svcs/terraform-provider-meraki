package settings_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceVlansSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_vlans_settings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_vlans_settings"),
			},

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: NetworkApplianceVlansSettingsResourceConfigUpdate(),
				Check:  NetworkApplianceVlansSettingsResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworkApplianceVlansSettingsResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_vlans_settings" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  vlans_enabled = true
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_vlans_settings"),
	)
}

// NetworkApplianceVlansSettingsResourceConfigUpdateChecks returns the test check functions for NetworkApplianceVlansSettingsResourceConfigUpdate
func NetworkApplianceVlansSettingsResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlans_enabled": "true",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_vlans_settings.test", expectedAttrs)
}
