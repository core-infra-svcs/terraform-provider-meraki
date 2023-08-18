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
				Config: testAccDevicesManagementinterfaceResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "1"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDevicesManagementinterfaceResourceConfigCreate(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_devices_management_interface" "test" {
	serial = "%s"
    wan1 = {
		wan_enabled = "enabled"
		vlan = 1
		using_static_ip = false
	}
	wan2 = {
		wan_enabled= "enabled"
		vlan = 1
		using_static_ip = false
	}
}
`, serial)
	return result
}
