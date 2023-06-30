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
					resource.TestCheckResourceAttr("devices_management_interface.test", "serial", "serial"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan1_wan_enabled", "not configured"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan1_using_static_ip", "true"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan1_static_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan1_static_subnet_mask", "255.255.255.0"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan1_static_gateway_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan1_vlan", "7"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan2_wan_enabled", "enabled"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan2_using_static_ip", "false"),
					resource.TestCheckResourceAttr("devices_management_interface.test", "wan2_vlan", "2"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccDevicesManagementinterfaceResourceConfigCreate = `
resource "meraki_organization" "test" {
	name = "meraki_devices_management_interface"
	api_enabled = true
} 

resource "meraki_devices_management_interface" "test" {
	serial = "Q2HY-A497-YBMG"
    wan1_wan_enabled = "not configured"
	wan1_using_static_ip = false
    wan2_wan_enabled = "not configured"
	wan2_using_static_ip = false
}
`
