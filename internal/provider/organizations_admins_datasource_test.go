package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsAdminsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationsAdminsDataSourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_admins"),
				),
			},

			/*
				// TODO - Create test admin
					{
						Config: testAccOrganizationsAdminsDataSourceConfigCreateAdmin,
						Check: resource.ComposeAggregateTestCheckFunc(
							//resource.TestCheckResourceAttr("meraki_organization_admin", "name", "test_admin"),
							//resource.TestCheckResourceAttr("meraki_organization_admin", "email", "test_admin@example.com"),
						),
					},
			*/

			// Read test admin
			{
				Config: testAccOrganizationsAdminsDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "id", "example-id"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "organization_id", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.#", "1"),

					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.id", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.name", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.email", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.org_access", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.account_status", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.two_factor_auth_enabled", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.has_api_key", "true"),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.last_active", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.tags.0.tag", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.tags.0.access", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.networks.0.id", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.networks.0.access", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.authentication_method", ""),

				),
			},
		},
	})
}

const testAccOrganizationsAdminsDataSourceConfigCreateOrg = `
resource "meraki_organization" "test" {
	name = "test_meraki_organizations_admins"
	api_enabled = true
}
`

/* TODO - Create admin to complete test coverage
const testAccOrganizationsAdminsDataSourceConfigCreateAdmin = `
resource "meraki_organization_admin" "test" {
	name = "test_admin"
	email = "test_admin@example.com"
}
`
*/

const testAccOrganizationsAdminsDataSourceConfigRead = `
resource "meraki_organization" "test" {
}

data "meraki_organizations_admins" "test" {
   organization_id = resource.meraki_organization.test.organization_id
}
`
