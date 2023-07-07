package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDevicesManagementinterfaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDevicesManagementinterfaceResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", "Q2HY-BHEX-TLTC"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "1"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccDevicesManagementinterfaceResourceConfigCreate = `
resource "meraki_devices_management_interface" "test" {
	serial = "Q2HY-BHEX-TLTC"
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
`
