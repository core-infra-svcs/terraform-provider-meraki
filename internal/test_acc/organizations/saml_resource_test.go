package organizations

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationSamlResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationSamlResourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organization_saml"),
				),
			},

			// Enable SAML on organization test
			{
				Config: testAccOrganizationSamlResourceConfigSaml,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization_saml.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization_saml.test", "enabled", "true"),
				),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organization_saml.test",
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

const testAccOrganizationSamlResourceConfigCreateOrg = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organization_saml"
 	api_enabled = true
 }
 `

const testAccOrganizationSamlResourceConfigSaml = `
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
