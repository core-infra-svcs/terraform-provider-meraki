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

			// Create test Organization
			{
				Config: testAccOrganizationsAdminsDataSourceConfigCreateOrganizations,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_admins"),
				),
			},

			// Create test admin
			{
				Config: testAccOrganizationsAdminsDataSourceConfigCreateAdmin,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "name", "test_acc_admin"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "email", "test_acc_meraki_organizations_admin_datasource_test1@example.com"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "org_access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "authentication_method", "Email"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "has_api_key", "false"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "west"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
				),
			},

			// Read test admin
			{
				Config: testAccOrganizationsAdminsDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(

					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "organization_id", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.#", "2"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.id", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.name", "test_acc_admin"),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.email", "test_acc_meraki_organizations_admin_datasource_test1@example.com"),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.org_access", "read-only"),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.account_status", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.two_factor_auth_enabled", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.has_api_key", "false"),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.last_active", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.tags.0.tag", "west"),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.tags.0.access", "read-only"),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.networks.0.id", ""),
					// resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.networks.0.access", ""),
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.authentication_method", "Email"),
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
resource "meraki_organization" "test" {}

resource "meraki_organizations_admin" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "test_acc_admin"
	email       = "test_acc_meraki_organizations_admin_datasource_test1@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
	tags = [{
				   tag = "west"
				   access = "read-only"
				  }]
	networks    = []
}
`

const testAccOrganizationsAdminsDataSourceConfigRead = `
resource "meraki_organization" "test" {}

data "meraki_organizations_admins" "test" {
	depends_on = [resource.meraki_organization.test]
   organization_id = resource.meraki_organization.test.organization_id
}
`
