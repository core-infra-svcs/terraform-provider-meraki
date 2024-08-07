package organizations

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
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

			// Import State testing
			{
				ResourceName:            "meraki_organization.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
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
