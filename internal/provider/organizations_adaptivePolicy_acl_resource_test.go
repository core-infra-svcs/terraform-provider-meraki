package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdaptivepolicyAclResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationsAdaptivepolicyAclResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("meraki_organizations_adaptivePolicy_acl.testAcl", "name", "testacl"),
				),
			},

			// Update testing
			{
				Config: testUpdatedAccOrganizationsAdaptivepolicyAclResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("meraki_organizations_adaptivePolicy_acl.testAcl", "description", "Blocks sensitive web traffic sets"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "meraki_organizations_adaptivePolicy_acl.testAcl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "784752235069308981,testacl",
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

// Config for create and read
const testAccOrganizationsAdaptivepolicyAclResourceConfig = `
resource "meraki_organizations_adaptivePolicy_acl" "testAcl" {
	id          = "784752235069308981"
	name        = "testacl"
	description = "Blocks sensitive web traffic"
	ipversion   = "ipv6"
	rules       = [      
                  {
                   policy = "allow"
                   protocol = "any"
                   srcport = "any"
                   dstport = "any"
                  }
                  ]  
  }
`

// Config for update
const testUpdatedAccOrganizationsAdaptivepolicyAclResourceConfig = `
resource "meraki_organizations_adaptivePolicy_acl" "testAcl" {
	id          = "784752235069308981"
	name        = "testacl"
	description = "Blocks sensitive web traffic sets"
	ipversion   = "ipv6"
	rules       = [
				  {
					policy = "allow"
					protocol = "any"
					srcport = "any"
					dstport = "any"
				  }
				  ] 
  
  }
`
