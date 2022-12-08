package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdaptivePolicyAclsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create test organization
			{
				Config: testAccOrganizationsAdaptivePolicyAclsDataSourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_admin_adaptive_policy_acls"),
				),
			},

			// Read testing
			{
				Config: testAccOrganizationsAdaptivePolicyAclsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "id", "example-id"),
					//resource.TestCheckResourceAttr("data.meraki_organizations_adaptive_policy_acls.test", "list.#", "1"),
				),
			},
		},
	})
}

const testAccOrganizationsAdaptivePolicyAclsDataSourceConfigCreateOrg = `
resource "meraki_organization" "test" {
	name = "test_meraki_organizations_admin_adaptive_policy_acls"
	api_enabled = true
}
`

const testAccOrganizationsAdaptivePolicyAclsDataSourceConfig = `
resource "meraki_organization" "test" {}

data "meraki_organizations_adaptive_policy_acls" "test" {
    organization_id = resource.meraki_organization.test.organization_id
}
`
