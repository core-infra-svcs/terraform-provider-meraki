package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdaptivepolicyAclsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrganizationsAdaptivepolicyAclsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.meraki_organizations_adaptivePolicy_acls.list", "id", "1232821"),
				),
			},
		},
	})
}

const testAccOrganizationsAdaptivepolicyAclsDataSourceConfig = `
data "meraki_organizations_adaptivePolicy_acls" "list" {
    id = "1232821"
}
`
