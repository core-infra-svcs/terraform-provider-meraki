package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"

	"testing"
)

func TestAccNetworkStormControlDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control_data"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_network_switch_storm_control_data"),
			},

			// Claim Device
			{
				Config: NetworkStormControlDataSourceClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check:  NetworkStormControlDataSourceClaimNetworkDeviceChecks(),
			},

			// Create and Read Networks Switch Qos Rules.
			{
				Config: NetworkStormControlDataSourceConfigCreate(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check:  NetworkStormControlDataSourceConfigCreateChecks(),
			},

			// Read Datasource
			{
				Config: NetworkStormControlDataSourceRead(),
				Check:  NetworkStormControlDataSourceReadChecks(),
			},
		},
	})
}

func NetworkStormControlDataSourceClaimNetworkDevice(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control_data"),
		serial,
	)
}

// NetworkStormControlDataSourceClaimNetworkDeviceChecks returns the test check functions for NetworkStormControlDataSourceClaimNetworkDevice
func NetworkStormControlDataSourceClaimNetworkDeviceChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":            "test_acc_network_switch_storm_control_data",
		"product_types.0": "appliance",
		"product_types.1": "cellularGateway",
		"product_types.2": "switch",
		"product_types.3": "wireless",
		"tags.0":          "tag1",
	}
	return utils.ResourceTestCheck("meraki_network.test", expectedAttrs)
}

func NetworkStormControlDataSourceConfigCreate(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_networks_storm_control" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 90
	multicast_threshold = 90
	unknown_unicast_threshold = 90
}

resource "meraki_devices_switch_port" "test" {
	depends_on = [resource.meraki_networks_storm_control.test]
	serial = "%s"
	storm_control_enabled = true
	port_id = 1
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control_data"),
		serial,
	)
}

// NetworkStormControlDataSourceConfigCreateChecks returns the test check functions for NetworkStormControlDataSourceConfigCreate
func NetworkStormControlDataSourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"broadcast_threshold":       "90",
		"multicast_threshold":       "90",
		"unknown_unicast_threshold": "90",
	}
	return utils.ResourceTestCheck("meraki_networks_storm_control.test", expectedAttrs)
}

//const testAccNetworkStormControlDataSourceConfigReadNetworkStormControl = `
//resource "meraki_network" "test" {
//    product_types = ["appliance", "switch", "wireless"]
//}
//
//resource "meraki_networks_devices_claim" "test" {
//	network_id = resource.meraki_network.test.network_id
//}
//
//resource "meraki_networks_storm_control" "test" {
//    network_id = resource.meraki_network.test.network_id
//	broadcast_threshold = 90
//	multicast_threshold = 90
//	unknown_unicast_threshold = 90
//}
//
//
//data "meraki_networks_storm_control" "test" {
//	network_id = resource.meraki_network.test.network_id
//}
//
//`

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
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "nw name"),
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
