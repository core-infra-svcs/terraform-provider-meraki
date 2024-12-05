package _interface_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesManagementInterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_device_management_interface"),
			},

			// Claim device to Network
			{
				Config: DevicesManagementInterfaceResourceConfigCreate(
					os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
					os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
					os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check: DevicesManagementInterfaceResourceConfigCreateCheck(),
			},

			{
				Config: DevicesManagementInterfaceResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  DevicesManagementInterfaceResourceConfigUpdateCheck(),
			},

			{
				Config: DevicesManagementInterfaceResourceConfigCreateMS(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check:  DevicesManagementInterfaceResourceConfigCreateMSCheck(),
			},

			{
				Config: DevicesManagementInterfaceResourceConfigCreateMR(os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check:  DevicesManagementInterfaceResourceConfigCreateMRCheck(),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func DevicesManagementInterfaceResourceConfigCreate(mxSerial, msSerial, mrSerial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s", "%s", "%s"
  ]
}

resource "meraki_devices_management_interface" "mx" {
    depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	serial = "%s"
	wan1 = {
		wan_enabled = "disabled"
		vlan = 2
		using_static_ip = false
	}
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
		mxSerial, msSerial, mrSerial, mxSerial,
	)
}

// DevicesManagementInterfaceResourceConfigCreateCheck returns the test check functions for DevicesManagementInterfaceResourceConfigCreate
func DevicesManagementInterfaceResourceConfigCreateCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
		"wan1.wan_enabled":     "disabled",
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	}
	return utils.ResourceTestCheck("meraki_devices_management_interface.mx", expectedAttrs)
}

func DevicesManagementInterfaceResourceConfigUpdate(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

resource "meraki_devices_management_interface" "mx" {
    depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	serial = "%s"
	wan1 = {
		wan_enabled = "enabled"
		vlan = 2
		using_static_ip = false
	}
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
		serial, serial,
	)
}

// DevicesManagementInterfaceResourceConfigUpdateCheck returns the test check functions for DevicesManagementInterfaceResourceConfigUpdate
func DevicesManagementInterfaceResourceConfigUpdateCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
		"wan1.wan_enabled":     "enabled",
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	}
	return utils.ResourceTestCheck("meraki_devices_management_interface.mx", expectedAttrs)
}

func DevicesManagementInterfaceResourceConfigCreateMS(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

resource "meraki_devices_management_interface" "ms" {
    depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	serial = "%s"
	wan1 = {
		wan_enabled = null
		vlan = 2
		using_static_ip = false
	}
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
		serial, serial,
	)
}

// DevicesManagementInterfaceResourceConfigCreateMSCheck returns the test check functions for DevicesManagementInterfaceResourceConfigCreateMS
func DevicesManagementInterfaceResourceConfigCreateMSCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	}
	return utils.ResourceTestCheck("meraki_devices_management_interface.ms", expectedAttrs)
}

func DevicesManagementInterfaceResourceConfigCreateMR(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

resource "meraki_devices_management_interface" "mr" {
    depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	serial = "%s"
	wan1 = {
		wan_enabled = "not configured"
		vlan = 2
		using_static_ip = false
	}
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
		serial, serial,
	)
}

// DevicesManagementInterfaceResourceConfigCreateMRCheck returns the test check functions for DevicesManagementInterfaceResourceConfigCreateMR
func DevicesManagementInterfaceResourceConfigCreateMRCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MR_SERIAL"),
		"wan1.wan_enabled":     "not configured",
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	}
	return utils.ResourceTestCheck("meraki_devices_management_interface.mr", expectedAttrs)
}
