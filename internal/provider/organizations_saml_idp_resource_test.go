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

			// Create test Organization
			{
				Config: testAccOrganizationsSamlIdpResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test-acc-meraki-organizations-saml-idp"),
				),
			},

			// Enable SAML on organization test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigSaml,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "enabled", "true"),
				),
			},

			// Create and Read Idp test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigCreateIdp,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "slo_logout_url", "https://sbuxforyou.com"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "x_509_cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"),
				),
			},

			// Update Idp test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "slo_logout_url", "https://sbuxforyouandme.com"),
					resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "x_509_cert_sha1_fingerprint", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"),
					//resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "organization_id", ""),
					//resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "consumer_url", ""),
					//resource.TestCheckResourceAttr("meraki_organizations_saml_idp.test", "idp_id", ""),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsSamlIdpResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test-acc-meraki-organizations-saml-idp"
 	api_enabled = true
 }
 `

const testAccOrganizationsSamlIdpResourceConfigSaml = `
resource "meraki_organization" "test" {
}
resource "meraki_organization_saml" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	enabled = true
}
`

const testAccOrganizationsSamlIdpResourceConfigCreateIdp = `
resource "meraki_organization" "test" {
}
 resource "meraki_organizations_saml_idp" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	slo_logout_url = "https://sbuxforyou.com"
	x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"
}
`

const testAccOrganizationsSamlIdpResourceConfigUpdate = `
resource "meraki_organization" "test" {
}

resource "meraki_organizations_saml_idp" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	slo_logout_url = "https://sbuxforyouandme.com"
	x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"
}
`
