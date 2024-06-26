package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksDevicesClaimResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksDevicesClaimResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Claim a Device by Serial into the Organization
			/*
				{
						Config: testAccNetworksDevicesClaimResourceConfigClaimOrgSerial(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"),
							os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_organizations_claim.test_serial", "id", "example-id"),
							resource.TestCheckResourceAttr("meraki_organizations_claim.test_serial", "serials.#", "1"),
							resource.TestCheckResourceAttr("meraki_organizations_claim.test_serial", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
						),
					},
			*/

			// Create and Read a Network. If a network with the same name already exists this will not create.
			{
				Config: testAccNetworksDevicesClaimResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_devices_claim"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim and Read NetworksDevicesClaim
			{
				Config: testAccNetworksDevicesClaimResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				),
			},

			// Import Test
			{
				ResourceName:      "meraki_networks_devices_claim.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

/*
func testAccNetworksDevicesClaimResourceConfigClaimOrgSerial(orgId, serial string) string {
	result := fmt.Sprintf(`

	resource "meraki_organizations_claim" "test_serial" {
		organization_id = %s
		orders = []
		licences = []
		serials = ["%s"]
	}
`, orgId, serial)
	return result
}
*/

// testAccNetworksDevicesClaimResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksDevicesClaimResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_network_devices_claim"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccNetworksDevicesClaimResourceConfigCreate is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccNetworksDevicesClaimResourceConfigCreate(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        product_types = ["appliance"]
}    

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}	
`, serial)
	return result
}
