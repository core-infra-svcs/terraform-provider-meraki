package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!
// TODO - Testing is meant to be atomic in that we give very specific instructions for how to create, read, update, and delete infrastructure across test steps.
// TODO - This is really useful for troubleshooting resources/data sources during development and provides a high level of confidence that our provider works as intended.
func TestAccDevicesSwitchPortsCycleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccDevicesSwitchPortsCycleResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_devices_switch_ports_cycle"),
				),
			},

			//Create and Read testing
			{
				Config: testAccDevicesSwitchPortsCycleResourceConfigCreate,
				Check:  resource.ComposeAggregateTestCheckFunc(
				//resource.TestCheckResourceAttr("meraki_devices_switch_ports_cycle.test", "id", "example-id"),
				),
			},

			//
			//{
			//	ResourceName:      "meraki_devices_switch_ports_cycle.test",
			//	ImportState:       true,
			//	ImportStateVerify: false,
			//	ImportStateId:     "1234567890, 0987654321",
			//},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccDevicesSwitchPortsCycleResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_devices_switch_ports_cycle"
 	api_enabled = true
 }
 `

const testAccDevicesSwitchPortsCycleResourceConfigCreate = `
resource "meraki_devices_switch_ports_cycle" "test" {
	serial = "Q2GX-AQWX-FZW9"
	ports = ["1"]
}
`
