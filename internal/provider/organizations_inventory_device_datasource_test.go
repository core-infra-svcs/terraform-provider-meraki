package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccOrganizationsInventoryDevicesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsInventoryDeviceDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccOrganizationsInventoryDeviceDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_inventory_device"),
				),
			},

			// Create and Read a Network.
			{
				Config: testAccOrganizationsInventoryDeviceDataSourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
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

			// TODO: Create and Read OrganizationsInventoryDevices
			//{
			//	Config: testAccOrganizationsInventoryDevicesDataSourceConfigCreate,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("organizations_inventory_devices.test", "id", "example-id"),
			//	),
			//},

			// Read OrganizationsInventoryDevices
			{
				Config: testAccOrganizationsInventoryDeviceDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations_inventory_device.test", "id", "example-id"),
				),
			},
		},
	})
}

// testAccOrganizationsInventoryDeviceDataSourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccOrganizationsInventoryDeviceDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_organizations_inventory_device"
 	api_enabled = true
 }
 `

// testAccOrganizationsInventoryDeviceDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
const testAccOrganizationsInventoryDeviceDataSourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

// TODO: Make a change to the configuration to test
// testAccOrganizationsInventoryDevicesDataSourceConfigCreate is a constant string that defines the configuration for creating and updating a organizations_{organizationId}_inventory_devices resource in your tests.
// It depends on both the organization and network resources.
//const testAccOrganizationsInventoryDevicesDataSourceConfigCreate = `
//resource "meraki_organization" "test" {}
//resource "meraki_network" "test" {
//	depends_on = [resource.meraki_organization.test]
//	product_types = ["appliance", "switch", "wireless"]
//}
//
//resource "meraki_organizations_inventory_device" "test" {
//	depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
//  	network_id = resource.meraki_network.test.network_id
//}
//`

// testAccOrganizationsInventoryDeviceDataSourceConfigRead is a constant string that defines the configuration for creating and updating a organizations_{organizationId}_inventory_devices resource in your tests.
// It depends on both the organization and network resources.
var testAccOrganizationsInventoryDeviceDataSourceConfigRead = func(orgID string) string {
	return fmt.Sprintf(`
data "meraki_organizations_inventory_device" "test" {
  	organization_id = "%s"
    serial = "Q2KD-QF5S-L52G"
}
`, orgID)
}
