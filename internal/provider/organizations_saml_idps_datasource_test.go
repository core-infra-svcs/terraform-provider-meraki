package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlIdpsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfig("1239794"),
				Check: resource.ComposeAggregateTestCheckFunc(

					// TODO - Check return data matches expected result
					// TODO - Example:
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "id", "example-id"),
					// TODO - Example:
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.#", "1"),
					// TODO - Example:
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.x_509cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:99"),
					// TODO - Example:
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.slo_logout_url", "https://somewhereelseforyou.com"),
				),
			},
		},
	})
}

//const testAccOrganizationsSamlIdpsDataSourceConfig = `
//data "meraki_organizations_saml_idps" "test" {
//
//}
//`

func testAccOrganizationsSamlIdpsDataSourceConfig(organizationId string) string {
	return fmt.Sprintf(`
data "meraki_organizations_saml_idps" "test" {
      organization_id = "%s"
}
`, organizationId)

}
