package organizations

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationPolicyObjectResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create Organization
			{
				Config: testAccOrganizationPolicyObjectResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_policy_object"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
				),
			},

			// Create and Read testing
			{
				Config: testAccOrganizationPolicyObjectResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_policy_object.test", "name", "test_acc_meraki_organizations_policy_object"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationPolicyObjectResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_policy_object.test", "name", "test_acc_meraki_organizations_policy_object"),
				),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_policy_object.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationPolicyObjectResourceConfig = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organizations_policy_object"
	api_enabled = true
}
`

const testAccOrganizationPolicyObjectResourceConfigCreate = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organization_policy_object"
}

resource "meraki_organizations_policy_object" "test" {
	depends_on = [meraki_organization.test]
	organization_id = meraki_organization.test.organization_id
	name = "test_acc_meraki_organizations_policy_object"
	category = "network"
	type = "cidr"
	cidr = "10.0.0.0/24"
    ip = "1.2.3.4"
    group_ids = []
}
`

const testAccOrganizationPolicyObjectResourceConfigUpdate = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organization_policy_object"
}

resource "meraki_organizations_policy_object" "test" {
	depends_on = [meraki_organization.test]
	organization_id = meraki_organization.test.organization_id
	name = "Web Servers - Datacenter 10"
	category = "network"
	type = "cidr"
	cidr = "10.0.0.0/24"
	fqdn = "example.com"
	mask = "255.255.255.0"
    ip = "1.2.3.4"
    group_ids = [8]
}
`
