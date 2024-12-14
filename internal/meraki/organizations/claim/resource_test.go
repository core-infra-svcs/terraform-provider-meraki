package claim_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccOrganizationsClaimResource function is used to test the CRUD operations of the Terraform resource you are developing.
func TestAccOrganizationsClaimResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_claim"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_claim"),
			},

			// Claim a Device by Serial into the Organization
			{
				Config: utils.ClaimDeviceConfig(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  utils.ClaimDeviceTestChecks(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			},

			// Claim an Order into the Organization (commented out since it needs valid Order SsidNumber)
			/*
				{
					Config: utils.ClaimOrderConfig(os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),
					Check: utils.ClaimOrderTestChecks(os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),
				},
			*/

			// Claim a License into the Organization (commented out since it needs a valid License)
			/*
				{
					Config: utils.ClaimLicenseConfig(os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					Check: utils.ClaimLicenseTestChecks(os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				},
			*/
		},
	})
}
