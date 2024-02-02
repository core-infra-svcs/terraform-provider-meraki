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

			// Create and Read a Network.
			{
				Config: testAccDevicesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_device"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim Device to Network
			{
				Config: testAccDevicesResourceConfigDeviceClaim(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(

					// Claim A Device To A Network
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				),
			},

			// Update and Read Device Attributes
			{
				Config: testAccDevicesResourceConfigUpdateDevice(os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices.test", "id", os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices.test", "serial", os.Getenv("TF_ACC_MERAKI_MR_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices.test", "lat", "37.418095"),
					resource.TestCheckResourceAttr("meraki_devices.test", "lng", "-122.09853"),
					resource.TestCheckResourceAttr("meraki_devices.test", "address", "new address"),
					resource.TestCheckResourceAttr("meraki_devices.test", "name", "test_acc_test_device"),
					resource.TestCheckResourceAttr("meraki_devices.test", "notes", "test notes"),
					resource.TestCheckResourceAttr("meraki_devices.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_devices.test", "tags.0", "recently-added"),
					//resource.TestCheckResourceAttr("meraki_devices.test", "beacon_id_params.%", "3"),
					//resource.TestCheckResourceAttr("meraki_devices.test", "beacon_id_params.uuid", "00000000-0000-0000-0000-000000000000"),
					//resource.TestCheckResourceAttr("meraki_devices.test", "beacon_id_params.major", "5"),
					//resource.TestCheckResourceAttr("meraki_devices.test", "beacon_id_params.minor", "3"),
				),
			},

			// ImportState test case.
			{
				ResourceName:      "meraki_devices.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccNetworksDevicesClaimResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource. This will not work if the network already exists
func testAccDevicesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	name = "test_acc_network_device"
	product_types = ["wireless"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

func testAccDevicesResourceConfigDeviceClaim(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["wireless"]
}

resource "meraki_networks_devices_claim" "test" {
	depends_on = ["resource.meraki_network.test"]
	network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}


`, orgId, serial)
	return result
}

// testAccDevicesResourceConfigUpdateDevice is a constant string that defines the configuration for updating a devices' resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesResourceConfigUpdateDevice(serial string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	product_types = ["wireless"]
}

resource "meraki_networks_devices_claim" "test" {
	depends_on = ["resource.meraki_network.test"]
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_devices" "test" {
	depends_on = ["resource.meraki_network.test", "resource.meraki_networks_devices_claim.test"]
	network_id = resource.meraki_network.test.network_id
  	serial = "%s"
    lat = 37.418095
    lng = -122.09853
    address = "new address"
    name = "test_acc_test_device"
    notes = "test notes"
    tags = ["recently-added"]
}	
`, serial)
	return result
}
