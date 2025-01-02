package idps_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestAccOrganizationsSamlIdpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_saml_idp"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_saml_idp"),
			},

			// Enable SAML on organization test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigSaml(),
				Check:  OrganizationsSamlTestChecks(true),
			},

			// Create and Read IdP test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigCreateIdp(),
				Check:  OrganizationsSamlIdpTestChecks("https://sbuxforyou.com", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"),
			},

			// Update IdP test
			{
				Config: testAccOrganizationsSamlIdpResourceConfigUpdate(),
				Check:  OrganizationsSamlIdpTestChecks("https://sbuxforyouandme.com", "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_saml_idp.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// testAccOrganizationsSamlIdpResourceConfigSaml returns the configuration for enabling SAML on the organization
func testAccOrganizationsSamlIdpResourceConfigSaml() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organization_saml" "test" {
		depends_on = [resource.meraki_organization.test]
		id = resource.meraki_organization.test.organization_id
		enabled = true
	}
	`
}

// testAccOrganizationsSamlIdpResourceConfigCreateIdp returns the configuration for creating the SAML IdP resource
func testAccOrganizationsSamlIdpResourceConfigCreateIdp() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_saml_idp" "test" {
		depends_on = [resource.meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		slo_logout_url = "https://sbuxforyou.com"
		x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:24"
	}
	`
}

// testAccOrganizationsSamlIdpResourceConfigUpdate returns the configuration for updating the SAML IdP resource
func testAccOrganizationsSamlIdpResourceConfigUpdate() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_saml_idp" "test" {
		depends_on = [resource.meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		slo_logout_url = "https://sbuxforyouandme.com"
		x_509_cert_sha1_fingerprint = "00:11:22:33:44:55:66:77:88:99:00:11:22:33:44:55:66:77:88:55"
	}
	`
}

// OrganizationsSamlTestChecks returns the test check functions for verifying the organization SAML state
func OrganizationsSamlTestChecks(enabled bool) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"enabled": strconv.FormatBool(enabled),
	}

	return utils.ResourceTestCheck("meraki_organization_saml.test", expectedAttrs)
}

// OrganizationsSamlIdpTestChecks returns the test check functions for verifying the SAML IdP resource
func OrganizationsSamlIdpTestChecks(sloLogoutUrl, fingerprint string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"slo_logout_url":              sloLogoutUrl,
		"x_509_cert_sha1_fingerprint": fingerprint,
	}

	return utils.ResourceTestCheck("meraki_organizations_saml_idp.test", expectedAttrs)
}
