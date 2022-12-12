package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations"),
				),
			},

			// Read OrganizationsDataSource
			{
				Config: testAccOrganizationsDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_organizations.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations.test", "list.1.name", "test_meraki_organizations"),
					resource.TestCheckResourceAttr("data.meraki_organizations.test", "list.1.api_enabled", "true"),
				),
			},
		},
	})
}

const testAccOrganizationsDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_organizations"
 	api_enabled = true
 }
 `

const testAccOrganizationsDataSourceConfigRead = `
data "meraki_organizations" "test" {
}
`
