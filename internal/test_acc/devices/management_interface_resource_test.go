package devices

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesManagementInterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesManagementInterfaceResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_device_management_interface"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim device to Network
			{
				Config: testAccDevicesManagementInterfaceResourceConfigCreate(
					os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
					os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
					os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "wan1.wan_enabled", "disabled"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "wan1.vlan", "2"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "wan1.using_static_ip", "false"),
				),
			},

			{
				Config: testAccDevicesManagementInterfaceResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "wan1.wan_enabled", "enabled"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "wan1.vlan", "2"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mx", "wan1.using_static_ip", "false"),
				),
			},

			{
				Config: testAccDevicesManagementInterfaceResourceConfigCreateMS(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.ms", "serial", os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.ms", "wan1.vlan", "2"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.ms", "wan1.using_static_ip", "false"),
				),
			},

			{
				Config: testAccDevicesManagementInterfaceResourceConfigCreateMR(os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mr", "serial", os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mr", "wan1.wan_enabled", "not configured"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mr", "wan1.vlan", "2"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.mr", "wan1.using_static_ip", "false"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDevicesManagementInterfaceResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
   organization_id = "%s"
   product_types = ["appliance", "switch", "wireless"]
   tags = ["tag1"]
   name = "test_acc_device_management_interface"
   timezone = "America/Los_Angeles"
   notes = "Additional description of the network"
}
`, orgId)
	return result
}

func testAccDevicesManagementInterfaceResourceConfigCreate(mxSerial, msSerial, mrSerial string) string {
	result := fmt.Sprintf(`resource "meraki_network" "test" {
        product_types = ["appliance", "switch", "wireless"]
}    

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
`, mxSerial, msSerial, mrSerial, mxSerial)
	return result
}

func testAccDevicesManagementInterfaceResourceConfigUpdate(serial string) string {
	result := fmt.Sprintf(`resource "meraki_network" "test" {
        product_types = ["appliance", "switch", "wireless"]
}    

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
`, serial, serial)
	return result
}

func testAccDevicesManagementInterfaceResourceConfigCreateMS(serial string) string {
	result := fmt.Sprintf(`resource "meraki_network" "test" {
        product_types = ["appliance", "switch", "wireless"]
}    

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
`, serial, serial)
	return result
}

func testAccDevicesManagementInterfaceResourceConfigCreateMR(serial string) string {
	result := fmt.Sprintf(`resource "meraki_network" "test" {
        product_types = ["appliance", "switch", "wireless"]
}    

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
`, serial, serial)
	return result
}
