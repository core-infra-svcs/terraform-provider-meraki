package move_test

/* TODO - Get a Valid License to Move
import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccOrganizationsLicenseMoveResource function is used to test the CRUD operations of the Terraform resource you are developing.
func TestAccOrganizationsLicenseMoveResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{


			// Create and Read a Destination Organization.
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_move_license_destination"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_move_license_destination"),
			},


			// Claim License into Destination Organization from Source Organization
			{
				Config: testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToDestinationOrganization(os.Getenv("TF_ACC_MERAKI_ORGANIZATION"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check:  OrganizationsLicenseMoveResourceTestChecks("destination", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
			},

			// Claim License into Source Organization from Destination Organization
			{
				Config: testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToSourceOrganization(os.Getenv("TF_ACC_MERAKI_ORGANIZATION"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check:  OrganizationsLicenseMoveResourceTestChecks("source", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
			},


		},
	})
}


// testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToDestinationOrganization returns the configuration for moving a license to the destination organization
func testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToDestinationOrganization(organizationId, licenceId string) string {
	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_license" "test_destination_move" {
		depends_on = [resource.meraki_organization.test_destination]
		organization_id = "%s"
		dest_organization_id = resource.meraki_organization.test_destination.organization_id
		license_ids = ["%s"]
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_move_license_destination"),
		organizationId, licenceId,
	)
}

// testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToSourceOrganization returns the configuration for moving a license back to the source organization
func testAccOrganizationsLicenseMoveResourceConfigMoveLicenceToSourceOrganization(organizationId, licenceId string) string {
	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_license" "test_source_move" {
		depends_on = [resource.meraki_organization.test_destination]
		organization_id = resource.meraki_organization.test_destination.organization_id
		dest_organization_id = "%s"
		license_ids = ["%s"]
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_move_license_destination"),
		organizationId, licenceId,
	)
}

// OrganizationsLicenseMoveResourceTestChecks returns the test check functions for verifying license move
func OrganizationsLicenseMoveResourceTestChecks(moveType, licenceId string) resource.TestCheckFunc {
	return utils.ResourceTestCheck(fmt.Sprintf("meraki_organizations_claim.test_%s_move", moveType), map[string]string{
		"license_ids.#": "1",
		"license_ids.0": licenceId,
	})
}
*/
