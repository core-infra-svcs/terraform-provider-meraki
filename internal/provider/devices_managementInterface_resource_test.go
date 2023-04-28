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
					resource.TestCheckResourceAttr("devices_management_interface.test", "id", "example-id"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "name", "Block sensitive web traffic"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "description", "Blocks sensitive web traffic"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "ip_version", "ipv6"),

					// resource.TestCheckResourceAttr("devices_management_interface.test", "rules.#", "1"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "rules.0.policy", "deny"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "rules.0.protocol", "tcp"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "rules.0.src_port", "1,33"),
					// resource.TestCheckResourceAttr("devices_management_interface.test", "rules.0.dst_port", "22-30"),
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
	serial = "serial"
    wan1_wan_enabled = "not configured"
	wan1_using_static_ip = true
	wan1_static_ip = "1.2.3.4"
	wan1_static_subnet_mask = "255.255.255.0"
	wan1_static_gateway_ip = "1.2.3.1"
	wan1_vlan = 7
	
	wan2_wan_enabled = "enabled"
	wan2_using_static_ip = false
	wan2_vlan = 2
}
`
