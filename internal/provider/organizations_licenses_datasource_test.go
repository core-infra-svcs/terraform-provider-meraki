package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccOrganizationsLicensesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsLicensesDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccOrganizationsLicensesDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_move_license_source"),
				),
			},

			// Claim a Licence into the Organization
			/* TODO - Need a valid Licence
			{
					Config: testAccOrganizationsLicensesDataSourceConfigClaimLicence(os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "id", "example-id"),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "licenses.#", "1"),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "licenses.0.key", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "licenses.0.mode", "addDevices"),
					),
				},


			// Read Organizations Licenses
			{
				Config: testAccOrganizationsLicensesDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_MX_LICENCE_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations_licenses.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_organizations_licenses.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_licenses.test", "list.0.license_id", os.Getenv("TF_ACC_MERAKI_MX_LICENCE_ID")),
				),
			},

			*/
		},
	})
}

// testAccOrganizationsLicensesDataSourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccOrganizationsLicensesDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_move_license_source"
 	api_enabled = true
 }
 `

/*
// testAccOrganizationsLicensesDataSourceConfiggClaimLicence is a constant string that defines the configuration for creating and reading a organizations_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsLicensesDataSourceConfigClaimLicence(licence string) string {
	result := fmt.Sprintf(`
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_claim" "test_licence" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = []
		serials = []
		licences = [
			{
				key = "%s"
				mode = "addDevices"
			}
		]

	}
`, licence)
	return result
}

// testAccOrganizationsLicensesDataSourceConfigRead is a constant string that defines the configuration for creating and updating a organizations_licenses resource in your tests.
// It depends on both the organization and network resources.
const testAccOrganizationsLicensesDataSourceConfigRead = `
resource "meraki_organization" "test" {}
resource "meraki_organizations_license" "test_licence" {}

data "meraki_organizations_licenses" "test" {
    depends_on = ["resource.meraki_organization.test", "resource.meraki_organizations_claim.test_licence"]
	organization_id = resource.meraki_organization.test.organization_id
}
`
*/
