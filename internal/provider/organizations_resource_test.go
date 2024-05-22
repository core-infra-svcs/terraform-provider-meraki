package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func init() {
	resource.AddTestSweepers("meraki_organization", &resource.Sweeper{
		Name: "meraki_organization",
		F: func(organization string) error {
			return sweepMerakiOrganization(organization)
		},
	})
}

func TestAccOrganizationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organization"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organization_update"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organization.test", "management_details_name", "MSP ID"),
					resource.TestCheckResourceAttr("meraki_organization.test", "management_details_value", "123456"),
				),
			},

			// Clone Organization testing
			{
				Config: testAccOrganizationResourceConfigClone,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.clone", "name", "test_acc_meraki_organization_clone"),
					resource.TestCheckResourceAttr("meraki_organization.clone", "api_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organization.clone", "management_details_name", "MSP ID"),
					resource.TestCheckResourceAttr("meraki_organization.clone", "management_details_value", "123456"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationResourceConfig = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organization"
	api_enabled = true
}
`

const testAccOrganizationResourceConfigUpdate = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organization_update"
	api_enabled = true
	management_details_name = "MSP ID"
	management_details_value = "123456"
}
`

const testAccOrganizationResourceConfigClone = `
resource "meraki_organization" "test" {}

resource "meraki_organization" "clone" {
	depends_on = [meraki_organization.test]
	clone_organization_id = resource.meraki_organization.test.organization_id
	name = "test_acc_meraki_organization_clone"
	api_enabled = true
	
}
`
