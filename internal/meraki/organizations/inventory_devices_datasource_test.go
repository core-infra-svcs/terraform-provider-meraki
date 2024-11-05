package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccOrganizationsInventoryDevicesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
func TestAccOrganizationsInventoryDevicesDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Read OrganizationsInventoryDevices
			{
				Config: testAccOrganizationsInventoryDevicesDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check:  OrganizationsInventoryDevicesDataSourceTestChecks(),
			},
		},
	})
}

// testAccOrganizationsInventoryDevicesDataSourceConfigRead returns the configuration for reading organizations inventory devices data source
func testAccOrganizationsInventoryDevicesDataSourceConfigRead(orgID string) string {
	return fmt.Sprintf(`
	data "meraki_organizations_inventory_devices" "test" {
  		organization_id = "%s"
	}
	`, orgID)
}

// OrganizationsInventoryDevicesDataSourceTestChecks returns the test check functions for verifying inventory devices data source
func OrganizationsInventoryDevicesDataSourceTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		// Add any expected attributes to check, for example:
		// "list.#": "10",
	}

	return utils.ResourceTestCheck("data.meraki_organizations_inventory_devices.test", expectedAttrs)
}
