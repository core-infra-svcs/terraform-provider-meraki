package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdaptivePolicyAclResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationsAdaptivePolicyAclResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_adaptive_policy_acl"),
				),
			},

			// Create and Read testing
			{
				Config: testAccOrganizationsAdaptivePolicyAclResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "name", "Block sensitive web traffic"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "description", "Blocks sensitive web traffic"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "ip_version", "ipv6"),

					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.src_port", "1,33"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.0.dst_port", "22-30"),
				),
			},

			// TODO - ImportState testing
			/*
				{
						ResourceName:      "meraki_organizations_adaptive_policy_acl.test",
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "1234567890, 00000111111",
					},
			*/

			// Update testing
			{
				Config: testAccOrganizationsAdaptivePolicyAclResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.#", "2"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.1.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.1.protocol", "any"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.1.src_port", "any"),
					resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "rules.1.dst_port", "any"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsAdaptivePolicyAclResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_adaptive_policy_acl"
 	api_enabled = true
 }
 `

// Config for create and read
const testAccOrganizationsAdaptivePolicyAclResourceConfig = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_adaptive_policy_acl" "test" {
	depends_on = [
    	resource.meraki_organization.test
  	]
	organization_id          = resource.meraki_organization.test.organization_id
	name = "Block sensitive web traffic"
	description = "Blocks sensitive web traffic"
	ip_version   = "ipv6"
	rules = [      
		{
			"policy": "deny",
			"protocol": "tcp",
			"src_port": "1,33",
			"dst_port": "22-30"
		}
	]  
  }
`

// Config for update
const testAccOrganizationsAdaptivePolicyAclResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_organizations_adaptive_policy_acl" "test" {
	depends_on = [
    	resource.meraki_organization.test
  	]
	organization_id = resource.meraki_organization.test.organization_id
	name = "Block sensitive web traffic"
	description = "Blocks sensitive web traffic"
	ip_version   = "ipv6"
	rules = [      
		{
			"policy": "deny",
			"protocol": "tcp",
			"src_port": "1,33",
			"dst_port": "22-30"
		},
		{
			"policy": "allow",
			"protocol": "any",
			"src_port": "any",
			"dst_port": "any"
		}
	]  
  }
`
