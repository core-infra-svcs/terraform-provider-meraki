package organizations

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlIdpsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_saml_idps"),
				),
			},

			// Enable SAML on test Organization
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigUpdateOrganizationSaml,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization_saml.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "enabled", "true"),
				),
			},

			// Create and Read OrganizationsSamlIdp
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "slo_logout_url", "https://sbuxforyou.com"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "x_509_cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"),
				),
			},

			// Read OrganizationsSamlIdps
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(

					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "id", "example-id"),

					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.x_509_cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"),
					resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.slo_logout_url", "https://sbuxforyou.com"),

					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.idp_id", ""),
					//resource.TestCheckResourceAttr("data.meraki_organizations_saml_idps.test", "list.0.consumer_url", ""),
				),
			},
		},
	})
}

const testAccOrganizationsSamlIdpsDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_saml_idps"
 	api_enabled = true
 }
 `

const testAccOrganizationsSamlIdpsDataSourceConfigUpdateOrganizationSaml = `
resource "meraki_organization" "test" {
}
resource "meraki_organization_saml" "test" {
	depends_on = [
    	resource.meraki_organization.test
  	]
	organization_id = resource.meraki_organization.test.organization_id
	enabled = true
}
`

const testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp = `
resource "meraki_organization" "test" {
}
 resource "meraki_organizations_saml_idp" "test" {
	depends_on = [
    	resource.meraki_organization.test
  	]
	organization_id = resource.meraki_organization.test.organization_id
	slo_logout_url = "https://sbuxforyou.com"
	x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"
}
`
const testAccOrganizationsSamlIdpsDataSourceConfigRead = `
resource "meraki_organization" "test" {}

data "meraki_organizations_saml_idps" "test" {
      organization_id = resource.meraki_organization.test.organization_id
}
`
