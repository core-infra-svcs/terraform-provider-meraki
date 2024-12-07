package control_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"

	"testing"
)

// TODO: This test is only valid for devices that support this feature. MS120's do not.

func TestAccNetworkStormControlResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_network_switch_storm_control"),
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

			// Update and Read Storm control
			{
				Config: NetworkStormControlResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check:  NetworkStormControlResourceConfigUpdateChecks(),
			},

			// Import testing
			{
				ResourceName:      "meraki_networks_storm_control.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func NetworkStormControlResourceClaimNetworkDevice(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control"),
		serial,
	)
}

// NetworkStormControlResourceClaimNetworkDeviceCheck returns the test check functions for NetworkStormControlResourceClaimNetworkDevice
func NetworkStormControlResourceClaimNetworkDeviceCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":            "test_acc_network_switch_storm_control",
		"product_types.0": "appliance",
		"product_types.1": "cellularGateway",
		"product_types.2": "switch",
		"product_types.3": "wireless",
		"tags.0":          "tag1",
	}
	return utils.ResourceTestCheck("meraki_network.test", expectedAttrs)
}

func NetworkStormControlResourceConfigCreate(serial string) string {
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
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control"),
		serial,
	)
}

// NetworkStormControlResourceConfigCreateChecks returns the test check functions for NetworkStormControlResourceConfigCreate
func NetworkStormControlResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"broadcast_threshold":       "90",
		"multicast_threshold":       "90",
		"unknown_unicast_threshold": "90",
	}
	return utils.ResourceTestCheck("meraki_networks_storm_control.test", expectedAttrs)
}

func NetworkStormControlResourceConfigUpdate(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_networks_storm_control" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 40
	multicast_threshold = 40
	unknown_unicast_threshold = 40
}

resource "meraki_devices_switch_port" "test" {
	depends_on = [resource.meraki_networks_storm_control.test]
	serial = "%s"
	storm_control_enabled = true
	port_id = 1
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_switch_storm_control"),
		serial,
	)
}

// NetworkStormControlResourceConfigUpdateChecks returns the test check functions for NetworkStormControlResourceConfigUpdate
func NetworkStormControlResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"broadcast_threshold":       "40",
		"multicast_threshold":       "40",
		"unknown_unicast_threshold": "40",
	}
	return utils.ResourceTestCheck("meraki_networks_storm_control.test", expectedAttrs)
}
