package control_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"

	"testing"
)

func TestAccNetworkStormControlDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control_data"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_network_switch_storm_control_data"),
			},

			// Claim Device
			{
				Config: NetworkStormControlResourceClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check:  NetworkStormControlResourceClaimNetworkDeviceCheck(),
			},

			// Create and Read Storm control
			{
				Config: NetworkStormControlResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check:  NetworkStormControlResourceConfigCreateChecks(),
			},

			// Read Datasource
			{
				Config: NetworkStormControlDataSourceRead(),
				Check:  NetworkStormControlDataSourceReadChecks(),
			},
		},
	})
}

func NetworkStormControlDataSourceRead() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_networks_storm_control" "test" {
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 90
	multicast_threshold = 90
	unknown_unicast_threshold = 90
}


data "meraki_networks_storm_control" "test" {
	network_id = resource.meraki_network.test.network_id
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control_data"),
	)
}

// NetworkStormControlDataSourceReadChecks returns the test check functions for NetworkStormControlDataSourceRead
func NetworkStormControlDataSourceReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"broadcast_threshold":       "90",
		"multicast_threshold":       "90",
		"unknown_unicast_threshold": "90",
	}
	return utils.ResourceTestCheck("data.meraki_networks_storm_control.test", expectedAttrs)
}
