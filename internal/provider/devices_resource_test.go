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
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim A Device To A Network
			{
				Config: testAccDevicesResourceConfigClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "id", "example-id"),
				),
			},

			// Update and Read Device Attributes
			{
				Config: testAccDevicesResourceConfigUpdateDevice(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_devices.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("meraki_devices.test", "lat", "37.418095"),
					resource.TestCheckResourceAttr("meraki_devices.test", "lng", "-122.09853"),
					resource.TestCheckResourceAttr("meraki_devices.test", "address", "new address"),
					resource.TestCheckResourceAttr("meraki_devices.test", "name", "test_acc_mx_device"),
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

// testAccNetworksDevicesClaimResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource. Thisd will not work if the network already exists
func testAccDevicesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_network_device"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesResourceConfigClaimNetworkDevice is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesResourceConfigClaimNetworkDevice(orgId string, serial string) string {
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
`, orgId, serial)
	return result
}

// testAccDevicesResourceConfigUpdateDevice is a constant string that defines the configuration for updating a devices' resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesResourceConfigUpdateDevice(orgId string, serial string) string {
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

resource "meraki_devices" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	network_id = resource.meraki_network.test.network_id
  	serial = "%s"
    lat = 37.418095
    lng = -122.09853
    address = "new address"
    name = "test_acc_mx_device"
    notes = "test notes"
    beacon_id_params = {}
    tags = ["recently-added"]
}	
`, orgId, serial, serial)
	return result
}
