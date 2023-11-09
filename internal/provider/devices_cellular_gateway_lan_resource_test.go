package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDevicesCellularGatewayLanResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccDevicesCellularGatewayLanResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_meraki_devices_cellular_gateway_lan"),
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

			// Update testing
			{
				Config: testAccDevicesCellularGatewayLanResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("merak_devices_cellular_gateway_lan.test", "id", "example-id"),
					resource.TestCheckResourceAttr("merak_devices_cellular_gateway_lan.test", "device_name", "name of the MG"),
					resource.TestCheckResourceAttr("merak_devices_cellular_gateway_lan.test", "fixed_ip_assignments.0.name", "server 1"),
					resource.TestCheckResourceAttr("merak_devices_cellular_gateway_lan.test", "reserved_ip_range.0.start", "192.168.1.0"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDevicesCellularGatewayLanResourceConfigCreate(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_meraki_devices_cellular_gateway_lan"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

func testAccDevicesCellularGatewayLanResourceConfigUpdate(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_devices_cellular_gateway_lan" "test" {
	depends_on = ["meraki_network.test"]
	serial = "%s"
	device_name = "name of the MG"
	fixed_ip_assignments = [
		{
			mac = "0b:00:00:00:00:ac"
			name = "server 1"
			ip = "192.168.0.10"
		},
		{
			mac = "0b:00:00:00:00:ab"
			name = "server 2"
			ip = "192.168.0.20"
		}
	]
	reserved_ip_ranges = [
		{
			start = "192.168.1.0"
			end = "192.168.1.1"
			comment = "A reserved IP range"
		}
	]
}
`, serial)
	return result
}
