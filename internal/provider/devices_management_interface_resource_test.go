package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesManagementInterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesManagementInterfaceResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_device_management_interface"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim device to Network
			{
				Config: testAccDevicesManagementInterfaceResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "2"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
func testAccDevicesManagementInterfaceResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
   organization_id = %s
   product_types = ["appliance"]
   tags = ["tag1"]
   name = "test_acc_device_management_interface"
   timezone = "America/Los_Angeles"
   notes = "Additional description of the network"
}

`, orgId)
	return result
}

func testAccDevicesManagementInterfaceResourceConfigCreate(serial string) string {
	result := fmt.Sprintf(`resource "meraki_network" "test" {
        product_types = ["appliance"]
}    

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

resource "meraki_devices_management_interface" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
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
