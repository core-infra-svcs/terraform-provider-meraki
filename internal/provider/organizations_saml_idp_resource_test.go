package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlIdpResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccOrganizationsSamlIdpResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "organization_id", "1239794"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "slo_logout_url", "https://sbuxforyou.com"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "x_509cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationsSamlIdpResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "organization_id", "1239794"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "slo_logout_url", "https://sbuxforyouandme.com"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "x_509cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsSamlIdpResourceConfig = `
 resource "meraki_organizations_saml_idp" "test" {
	organization_id = "1239794"
	slo_logout_url = "https://sbuxforyou.com"
	x_509cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"
}
`

const testAccOrganizationsSamlIdpResourceConfigUpdate = `
resource "meraki_organizations_saml_idp" "test" {
	organization_id = "1239794"
	slo_logout_url = "https://sbuxforyouandme.com"
	x_509cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"
}
`
