package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
