package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdminsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrganizationsAdminsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					// TODO - Check return data matches expected result
					resource.TestCheckResourceAttr("data.meraki_organizations_admins.list", "id", "784752235069308980"),

					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.list", "list.0.email", "kiran.surapathi@gmail.com"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.list", "list.0.name", "kiran kumar surapathi"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_admins.list", "list.0.id", "784752235069314949"),
				),
			},
		},
	})
}

const testAccOrganizationsAdminsDataSourceConfig = `
data "meraki_organizations_admins" "list" {
    id = "784752235069308980"
}
`
