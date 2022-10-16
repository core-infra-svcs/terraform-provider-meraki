package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_organization.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("scaffolding_organization.test", "id", "Organization-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "scaffolding_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// Organization code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute"},
			},
			// Update and Read testing
			{
				Config: testAccOrganizationResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_organization.test", "configurable_attribute", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "scaffolding_organization" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}
