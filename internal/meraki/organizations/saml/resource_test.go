package saml_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestAccOrganizationSamlResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organization_saml"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organization_saml"),
			},

			// Enable SAML on organization test
			{
				Config: testAccOrganizationSamlResourceConfigSaml(),
				Check:  OrganizationsSamlTestChecks(true),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organization_saml.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// testAccOrganizationSamlResourceConfigSaml returns the configuration for enabling SAML on the organization
func testAccOrganizationSamlResourceConfigSaml() string {
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

// OrganizationsSamlTestChecks returns the test check functions for verifying the organization SAML state
func OrganizationsSamlTestChecks(enabled bool) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"enabled": strconv.FormatBool(enabled),
	}

	return utils.ResourceTestCheck("meraki_organization_saml.test", expectedAttrs)
}
