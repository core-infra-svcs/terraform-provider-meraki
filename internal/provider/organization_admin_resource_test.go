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
					resource.TestCheckResourceAttr("meraki_organizations_admin.testAdmin", "name", "testAdmin12345"),
				),
			},

			{
				Config: `
				resource "meraki_organizations_admin" "testAdmin" {
					id          = "784752235069308981"
					name        = "testAdmin123456"
					email       = "kirankumar6002700924560166612891@gmail.com"
					orgaccess   = "none"
					tags        = [
						{
							tag = "west"
							access = "read-only"
						}
					]
					  
					
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("meraki_organizations_admin.testAdmin", "name", "testAdmin123456"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsAdminResourceConfig = `
resource "meraki_organizations_admin" "testAdmin" {
	id          = "784752235069308981"
	name        = "testAdmin12345"
	email       = "kirankumar6002700924560166612891@gmail.com"
	orgaccess   = "none"
	tags        = [
        {
            tag = "west"
            access = "read-only"
        }
    ]
     
	
  }
`
