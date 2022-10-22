package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "testOrg1"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),

					// save postgres Id to variable
					resource.TestCheckResourceAttrWith("meraki_organization.test", "id",
						func(value string) error {
							id = value
							if len(id) < 1 {
								return fmt.Errorf("failed to save postgresId from state: %s", id)
							}
							return nil
						}),
				),
				//

			},

			/*
				// TODO - Figure out why update testing produces a 404 from the Meraki API. May require a custom test check.
					// Update testing
					{
						Config: testAccOrganizationResourceUpdate(id),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_organization.test", "name", "testOrg1"),
							resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
							resource.TestCheckResourceAttr("meraki_organization.test", "id", id),
						),
					},
			*/

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationResourceConfig = `
resource "meraki_organization" "test" {
	name = "testOrg1"
	api_enabled = true
}
`

/*
func testAccOrganizationResourceUpdate(postgresId string) string {
	return fmt.Sprintf(`
resource "meraki_organization" "test" {
	id = "%s"
	name = "testOrg1"
	api_enabled = true

}
`, postgresId)
}
*/
