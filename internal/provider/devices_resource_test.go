package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDevicesResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					// TODO - Check return data matches expected result
					resource.TestCheckResourceAttr("meraki_devices.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_devices.test", "serial", "Q2HY-6Y6T-X3HX"),
					resource.TestCheckResourceAttr("meraki_devices.test", "name", "testDevice1"),
				),
			},

			// TODO - Update+Read Test

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccDevicesResourceConfig = `
resource "meraki_devices" "test" {
    serial = "Q2HY-6Y6T-X3HX"
	name = "testDevice1"
}
`
