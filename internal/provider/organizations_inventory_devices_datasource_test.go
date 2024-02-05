package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccOrganizationsInventoryDevicesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsInventoryDevicesDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Read OrganizationsInventoryDevices
			{
				Config: testAccOrganizationsInventoryDevicesDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check:  resource.ComposeAggregateTestCheckFunc(
				//resource.TestCheckResourceAttr("data.meraki_organizations_inventory_devices.test", "list.#", "10"),
				),
			},
		},
	})
}

// testAccOrganizationsInventoryDevicesDataSourceConfigRead is a constant string that defines the configuration for creating and updating a organizations_{organizationId}_inventory_devices resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsInventoryDevicesDataSourceConfigRead(orgID string) string {
	return fmt.Sprintf(`
data "meraki_organizations_inventory_devices" "test" {
  	organization_id = "%s"
}
`, orgID)
}
