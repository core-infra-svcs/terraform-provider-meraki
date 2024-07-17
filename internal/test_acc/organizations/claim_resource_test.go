package organizations

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccOrganizationsClaimResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsClaimResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccOrganizationsClaimResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_claim"),
				),
			},

			// Claim a Device by Serial into the Organization
			{
				Config: testAccOrganizationsClaimResourceConfigClaimSerial(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_serial", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_serial", "serials.#", "1"),
					resource.TestCheckResourceAttr("meraki_organizations_claim.test_serial", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				),
			},

			// Claim an Order into the Organization
			/* TODO - Need a valid Order SsidNumber
			{
					Config: testAccOrganizationsClaimResourceConfigClaimOrder(os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_organizations_claim", "id", "example-id"),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_order", "orders.#", "1"),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_order", "orders.0", os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),
					),
				},
			*/

			// Claim a Licence into the Organization
			/* TODO - Need a valid Licence
			{
					Config: testAccOrganizationsClaimResourceConfigClaimLicence(os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "id", "example-id"),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "licenses.#", "1"),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "licenses.0.key", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
						resource.TestCheckResourceAttr("meraki_organizations_claim.test_licence", "licenses.0.mode", "addDevices"),
					),
				},
			*/
		},
	})
}

// testAccOrganizationsClaimResourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccOrganizationsClaimResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_claim"
 	api_enabled = true
 }
 `

// testAccOrganizationsClaimResourceConfigClaimSerial is a constant string that defines the configuration for creating and reading a organizations_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsClaimResourceConfigClaimSerial(serial string) string {
	result := fmt.Sprintf(`
	resource "meraki_organization" "test" {}
	
	resource "meraki_organizations_claim" "test_serial" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = []
		licences = []
		serials = ["%s"]
	}
`, serial)
	return result
}

/*
// testAccOrganizationsClaimResourceConfigClaimOrder is a constant string that defines the configuration for creating and reading a organizations_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsClaimResourceConfigClaimOrder(order string) string {
	result := fmt.Sprintf(`
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_claim" "test_order" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = ["%s"]
		serials = []
		licences = []

	}
`, order)
	return result
}

// testAccOrganizationsClaimResourceConfigClaimLicence is a constant string that defines the configuration for creating and reading a organizations_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsClaimResourceConfigClaimLicence(licence string) string {
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
*/
