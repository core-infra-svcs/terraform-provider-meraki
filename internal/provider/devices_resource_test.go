package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesResourceSecurityAppliance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDevicesResourceConfig(os.Getenv("TF_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(

					// TODO - Check return data matches expected result
					resource.TestCheckResourceAttr("meraki_devices.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_devices.test", "serial", os.Getenv("TF_MERAKI_MX_SERIAL")),
				),
			},
			// Update testing
			{
				Config: testAccDevicesResourceConfigUpdate(os.Getenv("TF_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices.test", "name", "My AP"),
					resource.TestCheckResourceAttr("meraki_devices.test", "serial", os.Getenv("TF_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices.test", "model", "MX67C-NA"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDevicesResourceConfig(serialNumber string) string {
	return fmt.Sprintf(`
resource "meraki_devices" "test" {
 	serial = "%s"
 }
 `, serialNumber)
}

func testAccDevicesResourceConfigUpdate(serialNumber string) string {
	return fmt.Sprintf(`
resource "meraki_devices" "test" {
 	serial = "%s"
	name = "My AP"
	tags = ["sfo", "ca"]
 }
 `, serialNumber)
}
