package device_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDevicesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			/*
				{
						Config: testAccDevicesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_device"),
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
					{
						Config: testAccDevicesResourceConfigDeviceClaim(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MR_SERIAL"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.#", "3"),
						),
					},
					{
						Config: testAccDevicesResourceConfigUpdateDevice(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MR_SERIAL"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_devices.test_mr", "tags.0", "recently-added"),
							resource.TestCheckResourceAttr("meraki_devices.test_mr", "name", "test_acc_mr_device"),
							resource.TestCheckResourceAttr("meraki_devices.test_mr", "address", "new address"),

							resource.TestCheckResourceAttr("meraki_devices.test_ms", "tags.0", "recently-added"),
							resource.TestCheckResourceAttr("meraki_devices.test_ms", "name", "test_acc_ms_device"),
							resource.TestCheckResourceAttr("meraki_devices.test_ms", "address", "new address"),

							resource.TestCheckResourceAttr("meraki_devices.test_mx", "tags.0", "recently-added"),
							resource.TestCheckResourceAttr("meraki_devices.test_mx", "name", "test_acc_mx_device"),
							resource.TestCheckResourceAttr("meraki_devices.test_mx", "address", "new address"),
						),
					},
			*/

			{
				Config: testAccDevicesResourceConfigUpdateDevice(os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices.test_mr", "tags.0", "recently-added"),
					resource.TestCheckResourceAttr("meraki_devices.test_mr", "name", "test_acc_mr_device"),
					resource.TestCheckResourceAttr("meraki_devices.test_mr", "address", "new address"),
				),
			},

			// import test
			{
				ResourceName:      "meraki_devices.test_mr",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDevicesResourceConfigUpdateDevice(serialMR string) string {
	return fmt.Sprintf(`

resource "meraki_devices" "test_mr" {
	serial = "%s"
	name = "test_acc_mr_device"
	tags = ["recently-added"]
	address = "new address"
}

`, serialMR)
}

/*
func testAccDevicesResourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	name = "test_acc_network_device"
	product_types = ["wireless", "switch", "appliance"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
}

func testAccDevicesResourceConfigDeviceClaim(orgId, serialMR, serialMS, serialMX string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["wireless", "switch", "appliance"]
}

resource "meraki_networks_devices_claim" "test" {
	depends_on = ["resource.meraki_network.test"]
	network_id = resource.meraki_network.test.network_id
	serials = ["%s", "%s", "%s"]
}

`, orgId, serialMR, serialMS, serialMX)
}

func testAccDevicesResourceConfigUpdateDevice(orgId, serialMR, serialMS, serialMX string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["wireless", "switch", "appliance"]
}

resource "meraki_devices" "test_mr" {
	depends_on = ["resource.meraki_network.test"]
	network_id = resource.meraki_network.test.network_id
	serial = "%s"
	name = "test_acc_mr_device"
	tags = ["recently-added"]
	address = "new address"
}

resource "meraki_devices" "test_ms" {
	depends_on = ["resource.meraki_network.test"]
	network_id = resource.meraki_network.test.network_id
	serial = "%s"
	name = "test_acc_ms_device"
	address = "new address"
	tags = ["recently-added"]
}

resource "meraki_devices" "test_mx" {
	depends_on = ["resource.meraki_network.test"]
	network_id = resource.meraki_network.test.network_id
	serial = "%s"
	name = "test_acc_mx_device"
	address = "new address"
	tags = ["recently-added"]
}
`, orgId, serialMR, serialMS, serialMX)
}

*/
