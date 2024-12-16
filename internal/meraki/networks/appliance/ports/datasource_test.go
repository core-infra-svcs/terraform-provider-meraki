package ports_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkAppliancePortsDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_ports"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_ports"),
			},

			// Claim Appliance To Network
			{
				Config: NetworksAppliancePortResourceClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			},

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings(),
				Check:  NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettingsChecks(),
			},

			//  Ports Network Appliance Ports.
			{
				Config: NetworkAppliancePortsDatasourceReadAppliancePorts(),
				Check:  NetworkAppliancePortsDatasourceReadAppliancePortsChecks(),
			},
		},
	})
}

func NetworkAppliancePortsDatasourceReadAppliancePorts() string {
	return fmt.Sprintf(`
	%s

	data "meraki_networks_appliance_ports" "test" {
	depends_on = ["resource.meraki_network.test", "resource.meraki_networks_devices_claim.test", "resource.meraki_networks_appliance_vlans_settings.test"]
	network_id = resource.meraki_network.test.network_id
    }
	`,
		NetworksAppliancePortResourceConfigListNetworkAppliancePorts(),
	)
}

// NetworkAppliancePortsDatasourceReadAppliancePortsChecks returns the test check functions for NetworkAppliancePortsDatasourceReadAppliancePorts
func NetworkAppliancePortsDatasourceReadAppliancePortsChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.0.allowed_vlans":         "all",
		"list.0.drop_untagged_traffic": "true",
		"list.0.enabled":               "true",
		"list.0.type":                  "trunk",
		"list.0.vlan":                  "0",
	}
	return utils.ResourceTestCheck("data.meraki_networks_appliance_ports.test", expectedAttrs)
}
