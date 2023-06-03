package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDevicesResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccDevicesResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccDevicesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_devices"),
				),
			},

			// Update and Read Devices
			{
				Config: testAccDevicesResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_devices.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices.test", "lat", "37.418095"),
					resource.TestCheckResourceAttr("meraki_devices.test", "lng", "-122.09853"),
					resource.TestCheckResourceAttr("meraki_devices.test", "address", "new address"),
					resource.TestCheckResourceAttr("meraki_devices.test", "name", "test device"),
					resource.TestCheckResourceAttr("meraki_devices.test", "notes", "test notes"),
					resource.TestCheckResourceAttr("meraki_devices.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_devices.test", "tags.0", "recently-added"),
				),
			},
		},
		// ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_devices.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890, 0987654321",
		   },
		*/

	})
}

// testAccDevicesResourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccDevicesResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_devices"
 	api_enabled = true
 }
 `

// testAccDevicesResourceConfigUpdate is a constant string that defines the configuration for updating a devices resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesResourceConfigUpdate(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_organization" "test" {}
resource "meraki_devices" "test" {
	depends_on = [resource.meraki_organization.test]
  	serial = "%s"
    lat = 37.418095
    lng = -122.09853
    address = "new address"
    name = "test device"
    notes = "test notes"
    beacon_id_params = {}
    tags = ["recently-added"]
}	
`, serial)
	return result
}
