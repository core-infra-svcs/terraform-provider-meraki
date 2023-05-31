package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccOrganizationsClaimResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsClaimResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccOrganizationsClaimResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_claim"),
				),
			},

			// TODO: Create and Read OrganizationsClaim
			{
				Config: testAccOrganizationsClaimResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("OrganizationsClaim.test", "id", "example-id"),

					resource.TestCheckResourceAttr("organizations_claim.test", "orders.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "orders.0", os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),

					resource.TestCheckResourceAttr("organizations_claim.test", "serials.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),

					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.key", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.mode", "addDevices"),
				),
			},

			// TODO: Update and Read OrganizationsClaim
			{
				Config: testAccOrganizationsClaimResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("OrganizationsClaim.test", "id", "example-id"),

					resource.TestCheckResourceAttr("organizations_claim.test", "orders.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "orders.0", os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),

					resource.TestCheckResourceAttr("organizations_claim.test", "serials.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),

					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.key", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.mode", "addDevices"),
				),
			},
		},

		// The resource.Test function automatically tests the Delete operation.
	})
}

// testAccOrganizationsClaimResourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccOrganizationsClaimResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_organizations_claim"
 	api_enabled = true
 }
 `

// testAccOrganizationsClaimResourceConfigCreate is a constant string that defines the configuration for creating and reading a organizations_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsClaimResourceConfigCreate(order, serial, licence string) string {
	result := fmt.Sprintf(`
	resource "meraki_organization" "test" {}
	
	resource "meraki_organizations_claim" "test" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = [%s]
		serials = [%s]
		licences = [
			{
				"key": %s,
				"mode": "addDevices"
			}
		]
	
	}
`, order, serial, licence)
	return string(result)
}

// TODO: Make a change to the configuration to test
// testAccOrganizationsClaimResourceConfigUpdate is a constant string that defines the configuration for updating a organizations_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsClaimResourceConfigUpdate(order, serial, licence string) string {
	result := fmt.Sprintf(`
	resource "meraki_organization" "test" {}
	
	resource "meraki_organizations_claim" "test" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = [%s]
		serials = [%s]
		licences = [
			{
				"key": %s,
				"mode": "addDevices"
			}
		]
	
	}
`, order, serial, licence)
	return result
}
