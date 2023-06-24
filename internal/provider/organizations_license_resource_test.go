package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccOrganizationsLicenseResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsLicenseResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccOrganizationsLicenseResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_license"),
				),
			},

			// Update and Read OrganizationsLicense
			{
				Config: testAccOrganizationsLicenseResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_MX_LICENCE_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_license.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_license.test", "license_id", os.Getenv("TF_ACC_MERAKI_MX_LICENCE_ID")),
					resource.TestCheckResourceAttr("meraki_organizations_license.test", "device_serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				),
			},
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

// testAccOrganizationsLicenseResourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccOrganizationsLicenseResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_organizations_license"
 	api_enabled = true
 }
 `

// testAccOrganizationsLicenseResourceConfigUpdate is a constant string that defines the configuration for updating a organizations_license resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsLicenseResourceConfigUpdate(licenceId, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_organization" "test" {}

resource "meraki_organizations_license" "test" {
	depends_on = [resource.meraki_organization.test]
    organization_id = resource.meraki_organization.test.organization_id
    license_id = "%s"
    device_serial = "%s"
}
`, licenceId, serial)
	return result
}
