package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!
// TODO - Testing is meant to be atomic in that we give very specific instructions for how to create, read, update, and delete infrastructure across test steps.
// TODO - This is really useful for troubleshooting resources/data sources during development and provides a high level of confidence that our provider works as intended.
func TestAccOrganizationsClaimResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsClaimResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_claim"),
				),
			},

			// Create and Read testing
			{
				Config: testAccOrganizationsClaimResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER"),
					os.Getenv("TF_ACC_MERAKI_MX_SERIAL"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("organizations_claim.test", "id", "example-id"),

					resource.TestCheckResourceAttr("organizations_claim.test", "orders.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "orders.0", os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),

					resource.TestCheckResourceAttr("organizations_claim.test", "serials.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),

					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.key", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.mode", "addDevices"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationsClaimResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER"),
					os.Getenv("TF_ACC_MERAKI_MX_SERIAL"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("organizations_claim.test", "id", "example-id"),

					resource.TestCheckResourceAttr("organizations_claim.test", "orders.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "orders.0", os.Getenv("TF_ACC_MERAKI_ORDER_NUMBER")),

					resource.TestCheckResourceAttr("organizations_claim.test", "serials.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),

					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.#", "1"),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.key", os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
					resource.TestCheckResourceAttr("organizations_claim.test", "licenses.0.mode", "addDevices"),
				),
			},

			//	TODO - organizations_licences_seats (Remove Licenses)

			// TODO - organizations_inventory_claim (Remove Serials)

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsClaimResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_organizations_claim"
 	api_enabled = true
 }
 `

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
	return result
}

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

/*
const testAccOrganizationsClaimResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_claim" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	orders = ["4CXXXXXXX"]
	serials = ["Q234-ABCD-5678"]
	licences = [
		{
			"key": "Z2XXXXXXXXXX",
			"mode": "addDevices"
		}
	]

}
`

const testAccOrganizationsClaimResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_claim" "test" {
	organization_id  = resource.meraki_organization.test.organization_id
    orders = ["4CXXXXXXX"]
	serials = ["Q234-ABCD-5678"]
	licences = [
		{
			"key": "Z2XXXXXXXXXX",
			"mode": "addDevices"
		}
	]

}
`
*/
