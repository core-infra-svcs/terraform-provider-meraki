package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlIdpsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_saml_idps"),
				),
			},

			/*
				// TODO - Create test admin
					{
						Config: testAccOrganizationsSamlIdpsDataSourceConfigCreateAdmin,
						Check: resource.ComposeAggregateTestCheckFunc(
							//resource.TestCheckResourceAttr("meraki_organization_admin", "name", "test_admin"),
							//resource.TestCheckResourceAttr("meraki_organization_admin", "email", "test_admin@example.com"),
						),
					},
			*/

			/*
					// TODO - Create test saml idp
				 			{
				 				Config: testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp,
				 				Check: resource.ComposeAggregateTestCheckFunc(
				 					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "id", "example-id"),
				 					// resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "organization_id", ""),
				 					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "slo_logout_url", "https://sbuxforyou.com"),
				 					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "x_509cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:99"),
				 				),
				 			},
			*/

			// TODO - Read testing
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "id", "example-id"),

					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.#", "1"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.x_509_cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:99"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.slo_logout_url", "https://somewhereelseforyou.com"),

					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.idp_id", ""),
					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.consumer_url", ""),
				),
			},
		},
	})
}

const testAccOrganizationsSamlIdpsDataSourceConfigCreateOrg = `
resource "meraki_organization" "test" {
	name = "test_meraki_organizations_saml_idps"
	api_enabled = true
}
`

/* TODO - Create admin to complete test coverage
const testAccOrganizationsSamlIdpsDataSourceConfigCreateAdmin = `
resource "meraki_organization_admin" "test" {
	name = "test_admin"
	email = "test_admin@example.com"
}
`
*/

/* TODO -Enable SAML to complete test coverage
const testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp = `
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_saml_idp" "test" {
 	organization_id = resource.meraki_organization.test.organization_id
 	slo_logout_url = "https://sbuxforyou.com"
 	x_509cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:99"
 }
 `
*/

const testAccOrganizationsSamlIdpsDataSourceConfigRead = `
resource "meraki_organization" "test" {}

data "meraki_organizations_saml_idps" "test" {
      organization_id = resource.meraki_organization.test.organization_id
}
`
