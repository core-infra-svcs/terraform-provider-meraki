package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAdministeredIdentitiesMeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAdministeredIdentitiesMeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					// TODO - Check return data matches expected result
					resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "id", "example-id"),
					//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "email", "kiran.surapathi@gmail.com"),

					//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "name", "kiran kumar surapathi"),
					//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "authentication.saml.enabled", "false"),
					//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "list.0.api_enabled", "true"),
				),
			},
		},
	})
}

const testAccAdministeredIdentitiesMeDataSourceConfig = `
data "meraki_administered_identities_me" "test" {
    
}
`
