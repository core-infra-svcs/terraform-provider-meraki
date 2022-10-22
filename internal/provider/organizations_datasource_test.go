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
			// Read testing
			{
				Config: testAccOrganizationsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations.list", "id", "example-id"),

					// this checks the number of organizations available to the API user.
					resource.TestCheckResourceAttr("data.meraki_organizations.list", "list.#", "1"),

					// Verify data inside a returned meraki organization by attribute element value
					//resource.TestCheckResourceAttr("data.meraki_organizations.list", "list.0.name", "DextersLab"),
					//resource.TestCheckResourceAttr("data.meraki_organizations.list", "list.0.api_enabled", "true"),
				),
			},
		},
	})
}

const testAccOrganizationsDataSourceConfig = `
data "meraki_organizations" "list" {
}
`
