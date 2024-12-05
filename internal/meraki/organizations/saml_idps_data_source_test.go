package organizations_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSamlIdpsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_saml_idps"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_saml_idps"),
			},

			// Enable SAML on test Organization
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigUpdateOrganizationSaml(),
				Check:  OrganizationsSamlTestChecks(true),
			},

			// Create and Read OrganizationsSamlIdp
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp(),
				Check:  OrganizationsSamlIdpTestChecks("https://sbuxforyou.com", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"),
			},

			// Read OrganizationsSamlIdps
			{
				Config: testAccOrganizationsSamlIdpsDataSourceConfigRead(),
				Check:  OrganizationsSamlIdpsDataSourceTestChecks(),
			},
		},
	})
}

// testAccOrganizationsSamlIdpsDataSourceConfigUpdateOrganizationSaml returns the configuration for enabling SAML on the organization
func testAccOrganizationsSamlIdpsDataSourceConfigUpdateOrganizationSaml() string {
	return `
	resource "meraki_organization" "test" {
	}

	resource "meraki_organization_saml" "test" {
		depends_on = [resource.meraki_organization.test]
		id = resource.meraki_organization.test.organization_id
		enabled = true
	}
	`
}

// testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp returns the configuration for creating the SAML IdP resource
func testAccOrganizationsSamlIdpsDataSourceConfigCreateSamlIdp() string {
	return `
	resource "meraki_organization" "test" {
	}
	resource "meraki_organizations_saml_idp" "test" {
		depends_on = [resource.meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		slo_logout_url = "https://sbuxforyou.com"
		x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"
	}
	`
}

// testAccOrganizationsSamlIdpsDataSourceConfigRead returns the configuration for reading the SAML IdP data source
func testAccOrganizationsSamlIdpsDataSourceConfigRead() string {
	return `
	resource "meraki_organization" "test" {}

	data "meraki_organizations_saml_idps" "test" {
		organization_id = resource.meraki_organization.test.organization_id
	}
	`
}

// OrganizationsSamlIdpsDataSourceTestChecks returns the test check functions for verifying the SAML IdPs data source
func OrganizationsSamlIdpsDataSourceTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.#":                             "1",
		"list.0.x_509_cert_sha1_fingerprint": "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24",
		"list.0.slo_logout_url":              "https://sbuxforyou.com",
	}

	return utils.ResourceTestCheck("data.meraki_organizations_saml_idps.test", expectedAttrs)
}
