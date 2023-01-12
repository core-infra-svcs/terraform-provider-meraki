package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TODO - Setup Testing Pipeline with dedicated email address that can be discovered by sweeper after tests run.
func TestAccOrganizationsAdminsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsAdminsDataSourceConfigCreateOrganizations,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_admins"),
				),
			},

			// TODO uncomment these tests once dedicated test pipeline enabled
			/*
				// Create test admin
				{
						Config: testAccOrganizationsAdminsDataSourceConfigCreateAdmin,
						Check:  resource.ComposeAggregateTestCheckFunc(
						//resource.TestCheckResourceAttr("meraki_organization_admin", "name", "test_admin"),
						//resource.TestCheckResourceAttr("meraki_organization_admin", "email", "test_admin@example.com"),
						),
					},
			*/

			// Read test admin
			{
				Config: testAccOrganizationsAdminsDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(

					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "organization_id", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.#", "1"),

					// TODO - uncomment these tests once dedicated test pipeline enabled
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.id", ""),
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.name", "test_admin"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.0.email", "meraki_organizations_admin_datasource_test@example.com"),
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

const testAccOrganizationsAdminsDataSourceConfigCreateOrganizations = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organizations_admins"
	api_enabled = true
}
`

const testAccOrganizationsAdminsDataSourceConfigCreateAdmin = `
resource "meraki_organization" "test" {
}

resource "meraki_organizations_admin" "test" {
depends_on = [
			resource.meraki_organization.test
		]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "testAdmin"
	email       = "meraki_organizations_admin_datasource_test@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			  {
			   tag = "west"
			   access = "read-only"
			  }]
    networks    = [{
                  id = "N_784752235069332413"
                  access = "read-only"
                }]
}
`

const testAccOrganizationsAdminsDataSourceConfigRead = `
resource "meraki_organization" "test" {
}

data "meraki_organizations_admins" "test" {
	depends_on = [
			resource.meraki_organization.test
		]
   organization_id = resource.meraki_organization.test.organization_id
}
`
