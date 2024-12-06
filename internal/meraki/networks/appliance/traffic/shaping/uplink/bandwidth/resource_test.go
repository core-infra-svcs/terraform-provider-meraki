package bandwidth_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceTrafficShapingUplinkBandWidthResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_traffic_shaping_uplink_bandwidth"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_traffic_shaping_uplink_bandwidth"),
			},

			// Update Network Appliance Traffic Shaping UplinkBandWidth.
			{
				Config: NetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  NetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdate(serial string) string {
	return fmt.Sprintf(`
	%s
 resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

 resource "meraki_networks_appliance_traffic_shaping_uplink_bandwidth" "test" {
	 depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	 network_id = resource.meraki_network.test.network_id
	 bandwidth_limit_wan2_limit_up = 100
	 bandwidth_limit_wan2_limit_down = 600 
	 bandwidth_limit_cellular_limit_up = 101200
	 bandwidth_limit_cellular_limit_down = 101200	
	 bandwidth_limit_wan1_limit_up = 100
	 bandwidth_limit_wan1_limit_down = 600
 }
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_traffic_shaping_uplink_bandwidth"),
		serial,
	)
}

// NetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateChecks returns the test check functions for NetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdate
func NetworksApplianceTrafficShapingUplinkBandWidthResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"bandwidth_limit_wan1_limit_up":       "100",
		"bandwidth_limit_wan1_limit_down":     "600",
		"bandwidth_limit_wan2_limit_up":       "100",
		"bandwidth_limit_wan2_limit_down":     "600",
		"bandwidth_limit_cellular_limit_up":   "101200",
		"bandwidth_limit_cellular_limit_down": "101200",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_traffic_shaping_uplink_bandwidth.test", expectedAttrs)
}
