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
			// Create and Read testing
			{
				Config: testAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					// TODO - Check return data matches expected result
					resource.TestCheckResourceAttr("meraki_organizations_admin.testAdmin", "name", "testAdmin1234666889"),
				),
			},

			// Update testing
			{
				Config: testUpdatedAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("meraki_organizations_admin.testAdmin", "name", "testAdmin1234666889"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "meraki_organizations_admin.testAdmin",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "1232821,test20226668890101@gmail.com",
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsAdminResourceConfig = `
resource "meraki_organizations_admin" "testAdmin" {
	id          = "1232821"
	name        = "testAdmin1234666889"
	email       = "test20226668890101@gmail.com"
	orgaccess   = "read-only"
    tags        = [
		          {
			       tag = "west"
			       access = "read-only"
		          },
                  {
			        tag = "east"
			        access = "read-only"
		          }
	              ]
    networks    = [{
                  id = "N_784752235069332413"
                  access = "read-only"
                }]
}
`
const testUpdatedAccOrganizationsAdminResourceConfig = `
resource "meraki_organizations_admin" "testAdmin" {
	id          = "1232821"
	name        = "testAdmin1234666889"
	email       = "test20226668890101@gmail.com"
	orgaccess   = "read-only"

	networks = [{
	  id = "N_784752235069332413"
	  access = "read-only"
	}]
	
	tags        = [
		{
			tag = "west"
			access = "read-only"
		},
		{
			 tag = "east"
			 access = "read-only"
		}
	     ]
				
	}
`
