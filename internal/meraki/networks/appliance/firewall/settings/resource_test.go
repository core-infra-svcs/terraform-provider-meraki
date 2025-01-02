package settings_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceFirewallSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_settings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_firewall_settings"),
			},

			// Update and Read Networks Appliance Firewall Settings.
			{
				Config: NetworkApplianceFirewallSettingsResourceConfigCreate(),
				Check:  NetworkApplianceFirewallSettingsResourceConfigCreateChecks(),
			},
		},
	})
}

func NetworkApplianceFirewallSettingsResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_firewall_settings" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    spoofing_protection = {
        ip_source_guard =  {
            mode = "block"
        }
    }
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_firewall_settings"),
	)
}

// NetworkApplianceFirewallSettingsResourceConfigCreateChecks returns the test check functions for NetworkApplianceFirewallSettingsResourceConfigCreate
func NetworkApplianceFirewallSettingsResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"spoofing_protection.ip_source_guard.mode": "block",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_firewall_settings.test", expectedAttrs)
}
