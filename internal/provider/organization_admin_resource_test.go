package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdminResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationsAdminResourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organization_admin"),
				),
			},

			// Create and Read testing
			{
				Config: testAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "name", "testAdmin"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "email", "meraki_organizations_admin_test@example.com"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "org_access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "authentication_method", "Email"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "west"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.id", "N_784752235069332413"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
				),
			},

			// Update testing
			{
				Config: testUpdatedAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.1.tag", "east"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.1.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.1.id", "N_784752235069332414"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.1.access", "read-only"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded organizationId + adminId.
			/*
				{
						ResourceName:      "meraki_organizations_admin.test",
						ImportState:       true,
						ImportStateVerify: false,
						ImportStateId:     "657525545596096508, 657525545596237587",
					},
			*/

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsAdminResourceConfigCreateOrg = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_organization_admin"
 	api_enabled = true
 }
 `

const testAccOrganizationsAdminResourceConfig = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_admin" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	name        = "testAdmin"
	email       = "meraki_organizations_admin_test@example.com"
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
const testUpdatedAccOrganizationsAdminResourceConfig = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_admin" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	name        = "testAdmin"
	email       = "meraki_organizations_admin_test@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			{
				tag = "west"
				access = "read-only"
			},
			{
				tag = "east"
				access = "read-only"
			  }]
    networks    = [{
                  id = "N_784752235069332413"
                  access = "read-only"
                },
				{
                  id = "N_784752235069332414"
                  access = "read-only"
                }]
}
`
