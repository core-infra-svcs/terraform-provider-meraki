package _interface_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesManagementInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
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
				Config: testAccDevicesManagementInterfaceDatasourceConfigClaimNetworkDevice,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", "Q3FA-RGA5-FZJF"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "2"),
					resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),

					resource.TestCheckResourceAttr("data.meraki_devices_management_interface.test", "serial", "Q3FA-RGA5-FZJF"),
				),
			},
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

const testAccDevicesManagementInterfaceDatasourceConfigClaimNetworkDevice = `
resource "meraki_network" "test" {
        product_types = ["appliance"]
}    

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "Q3FA-RGA5-FZJF"
  ]
}

resource "meraki_devices_management_interface" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	serial = "Q3FA-RGA5-FZJF"
	wan1 = {
		wan_enabled = "enabled"
		vlan = 2
		using_static_ip = false
	}
}

data "meraki_devices_management_interface" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test, resource.meraki_devices_management_interface.test]
	serial = "Q3FA-RGA5-FZJF"
}
`
