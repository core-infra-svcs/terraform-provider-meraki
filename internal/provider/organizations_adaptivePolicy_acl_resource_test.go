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
				Check:  resource.ComposeAggregateTestCheckFunc(

				// TODO - Check return data matches expected result
				// TODO - Example: resource.TestCheckResourceAttr("meraki_organizations_adaptive_policy_acl.test", "name", "testOrg1"),
				),
			},

			{
				Config: `
				resource "meraki_organizations_adaptivePolicy_acl" "testAcl" {
					id          = "784752235069308981"
					name        = "Block12345467868900123 sensitive web traffic"
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
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("meraki_organizations_adaptivePolicy_acl.testAcl", "description", "Blocks sensitive web traffic sets"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationsAdaptivepolicyAclResourceConfig = `
resource "meraki_organizations_adaptivePolicy_acl" "testAcl" {
	id          = "784752235069308981"
	name        = "Block12345467868900123 sensitive web traffic"
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
