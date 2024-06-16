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
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_saml_idp"),
				),
			},

			// Enable SAML on organization test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigSaml,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "enabled", "true"),
				),
			},

			// Create and Read Idp test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigCreateIdp,
				Check: resource.ComposeAggregateTestCheckFunc(
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
				),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_saml_idp.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{
					// Add any attributes you want to ignore during import verification
				},
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsSamlIdpResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organizations_saml_idp"
	api_enabled = true
}
`

const testAccOrganizationsSamlIdpResourceConfigSaml = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organizations_saml_idp"
	api_enabled = true
}
resource "meraki_organization_saml" "test" {
	depends_on = [
		resource.meraki_organization.test
	]
	organization_id = resource.meraki_organization.test.organization_id
	enabled = true
}
`

const testAccOrganizationsSamlIdpResourceConfigCreateIdp = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organizations_saml_idp"
	api_enabled = true
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

const testAccOrganizationsSamlIdpResourceConfigUpdate = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organizations_saml_idp"
	api_enabled = true
}
resource "meraki_organizations_saml_idp" "test" {
	depends_on = [
		resource.meraki_organization.test
	]
	organization_id = resource.meraki_organization.test.organization_id
	slo_logout_url = "https://sbuxforyouandme.com"
	x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"
}
`
