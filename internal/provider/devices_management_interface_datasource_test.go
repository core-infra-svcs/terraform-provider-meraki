package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesManagementInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Network.
			{
				Config: testAccDevicesManagementInterfaceDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
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

			// Claim Appliance To Network
			{
				Config: testAccDevicesManagementInterfaceDatasourceConfigClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MX_SERIAL"), os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					//resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					//resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
					//resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
					//resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "1023"),
					//
					resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					//resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
					//resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
					//resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.vlan", "1023"),
				),
			},
			//
			//// Create and Read Networks Switch Qos Rules.
			//{
			//	Config: testAccDevicesManagementInterfaceDataSourceConfigCreate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "1023"),
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.static_dns.#", "2"),
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.static_dns.0", "1.2.3.2"),
			//		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.static_dns.1", "1.2.3.3"),
			//	),
			//},
			//
			//// Read testing
			//{
			//	Config: testAccDevicesManagementInterfaceDataSourceRead(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.vlan", "1023"),
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.static_dns.#", "2"),
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.static_dns.0", "1.2.3.2"),
			//		resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "wan1.static_dns.1", "1.2.3.3"),
			//	),
			//},
		},
	})
}

func testAccDevicesManagementInterfaceDataSourceConfigCreateNetwork(orgId string) string {
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

func testAccDevicesManagementInterfaceDatasourceConfigClaimNetworkDevice(serial string, orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
}    
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

data "meraki_devices_management_interface" "test" {
	serial = "%s"
}
`, orgId, serial, serial)
	return result
}

//func testAccDevicesManagementInterfaceDataSourceConfigCreate(serial string) string {
//	result := fmt.Sprintf(`
//resource "meraki_devices_management_interface" "test" {
// serial = "%s"
//	wan1 = {
//		wan_enabled = "enabled"
//		using_static_ip = false
//		vlan 		  	= 1023
//		static_dns = [
//			"1.2.3.2",
//			"1.2.3.3"
// 	]
//	}
//}
//`, serial)
//	return result
//}
//
//func testAccDevicesManagementInterfaceDataSourceRead(serial string) string {
//	result := fmt.Sprintf(`
//
//data "meraki_devices_management_interface" "test" {
//	serial = "%s"
//}
//`, serial)
//	return result
//}
