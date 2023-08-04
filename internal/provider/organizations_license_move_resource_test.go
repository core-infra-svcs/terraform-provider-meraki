package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccOrganizationsLicenseMoveResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsLicenseMoveResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Destination Organization.
			{
				Config: testAccOrganizationsLicenseMoveResourceConfigCreateOrganizationDestination,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test_destination", "name", "test_acc_meraki_organizations_move_license_destination"),
				),
			},

			/* TODO - Get a Valid Licence to Move

			// Claim Licence into Destination Organization from Source Organization
			{
				Config: testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToDestinationOrganization(os.Getenv("TF_ACC_MERAKI_ORGANIZATION"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_destination_move", "organization_id", os.Getenv("TF_ACC_MERAKI_ORGANIZATION")),
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_destination_move", "license_ids.#", "1"),
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_destination_move", "license_ids.0", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				),
			},

			// Claim Licence into Source Organization from Destination Organization
			{
				Config: testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToSourceOrganization(os.Getenv("TF_ACC_MERAKI_ORGANIZATION"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_source_move", "dest_organization_id", os.Getenv("TF_ACC_MERAKI_ORGANIZATION")),
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_source_move", "license_ids.#", "1"),
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_source_move", "license_ids.0", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				),
			},


			*/
		},

		// ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_organizations_license.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890, 0987654321",
		   },
		*/

	})
}

// testAccOrganizationsLicenseMoveResourceConfigCreateOrganizationDestination is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccOrganizationsLicenseMoveResourceConfigCreateOrganizationDestination = `
 resource "meraki_organization" "test_destination" {
 	name = "test_acc_meraki_organizations_move_license_destination"
 	api_enabled = true
 }
 `

// testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToDestinationOrganization is a constant string that defines the configuration for updating a organizations_license resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToDestinationOrganization(organizationId, licenceId string) string {
	result := fmt.Sprintf(`
resource "meraki_organization" "test_destination" {}

resource "meraki_organizations_license" "test_destination_move" {
	depends_on = [resource.meraki_organization.test_destination]
	organization_id = "%s"
    dest_organization_id = resource.meraki_organization.test_destination.organization_id
    license_ids = ["%s"]
}
`, organizationId, licenceId)
	return result
}

// testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToSourceOrganization is a constant string that defines the configuration for updating a organizations_license resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToSourceOrganization(organizationId, licenceId string) string {
	result := fmt.Sprintf(`
resource "meraki_organization" "test_destination" {}

resource "meraki_organizations_license" "test_source_move" {
	depends_on = [resource.meraki_organization.test_destination]
	organization_id = resource.meraki_organization.test_destination.organization_id
    dest_organization_id = "%s"
    license_ids = ["%s"]
}
`, organizationId, licenceId)
	return result
}
