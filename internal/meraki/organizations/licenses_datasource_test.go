package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TestAccOrganizationsLicensesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
func TestAccOrganizationsLicensesDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_move_license_source"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_move_license_source"),
			},

			// Claim a License into the Organization (commented out since it needs a valid License)
			/*
				{
					Config: utils.ClaimLicenseConfig(os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					Check: utils.ClaimLicenseTestChecks(os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				},
			*/

			// Read Organizations Licenses
			/*
				{
					Config: testAccOrganizationsLicensesDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_MX_LICENCE_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					Check: OrganizationsLicensesDataSourceTestChecks(os.Getenv("TF_ACC_MERAKI_MX_LICENCE_ID")),
				},
			*/
		},
	})
}

// testAccOrganizationsLicensesDataSourceConfigRead returns the configuration for reading licenses data source
func testAccOrganizationsLicensesDataSourceConfigRead(licenceID, serial string) string {
	return fmt.Sprintf(`
	%s

	data "meraki_organizations_licenses" "test" {
		organization_id = resource.meraki_organization.test.organization_id
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_move_license_source"),
	)
}

// OrganizationsLicensesDataSourceTestChecks returns the test check functions for verifying licenses data source
func OrganizationsLicensesDataSourceTestChecks(licenseID string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.#":            "1",
		"list.0.license_id": licenseID,
	}

	return utils.ResourceTestCheck("data.meraki_organizations_licenses.test", expectedAttrs)
}
