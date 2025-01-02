package organization_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organization"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organization"),
			},

			// Update testing
			{
				Config: testAccOrganizationResourceConfigUpdate(),
				Check:  testAccOrganizationUpdateTestChecks(),
			},

			// Clone Organization testing
			{
				Config: testAccOrganizationResourceConfigClone(),
				Check:  testAccOrganizationCloneTestChecks(),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organization.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// testAccOrganizationResourceConfigUpdate returns the configuration string for updating the organization
func testAccOrganizationResourceConfigUpdate() string {
	return `
	resource "meraki_organization" "test" {
		name = "test_acc_meraki_organization_update"
		api_enabled = true
		management_details_name = "MSP ID"
		management_details_value = "123456"
	}
	`
}

// testAccOrganizationUpdateTestChecks returns the test check functions for the organization update
func testAccOrganizationUpdateTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":                     "test_acc_meraki_organization_update",
		"api_enabled":              "true",
		"management_details_name":  "MSP ID",
		"management_details_value": "123456",
	}

	return utils.ResourceTestCheck("meraki_organization.test", expectedAttrs)
}

// testAccOrganizationResourceConfigClone returns the configuration for cloning the organization
func testAccOrganizationResourceConfigClone() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organization" "clone" {
		depends_on = [meraki_organization.test]
		clone_organization_id = resource.meraki_organization.test.organization_id
		name = "test_acc_meraki_organization_clone"
		api_enabled = true
	}
	`
}

// testAccOrganizationCloneTestChecks returns the test check functions for the cloned organization
func testAccOrganizationCloneTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":                     "test_acc_meraki_organization_clone",
		"api_enabled":              "true",
		"management_details_name":  "MSP ID",
		"management_details_value": "123456",
	}

	return utils.ResourceTestCheck("meraki_organization.clone", expectedAttrs)
}
