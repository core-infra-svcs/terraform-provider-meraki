package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!
// TODO - Testing is meant to be atomic in that we give very specific instructions for how to create, read, update, and delete infrastructure across test steps.
// TODO - This is really useful for troubleshooting resources/data sources during development and provides a high level of confidence that our provider works as intended.
func TestAccOrganizationsClaimResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - Usually the first step in building a resource is to create an organization or network to configure.
			/*
			   // Create test Organization
			               {
			                   Config: testAccOrganizationsClaimResourceConfigCreateOrganization,
			                   Check: resource.ComposeAggregateTestCheckFunc(
			                       resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
			                       resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_Organizations_claim"),
			                   ),
			               },
			*/

			// TODO - Next, run the create test step for the resource you are developing. It is important to validate every field returned by read.
			// Create and Read testing
			{
				Config: testAccOrganizationsClaimResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("organizations_claim.test", "id", "example-id"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "name", "Block sensitive web traffic"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "description", "Blocks sensitive web traffic"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "ip_version", "ipv6"),

					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.#", "1"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.policy", "deny"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.protocol", "tcp"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.src_port", "1,33"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.dst_port", "22-30"),
				),
			},

			// TODO - Once a resource has been created, we will test the ability to modify it. Make sure to test all values that are modifiable by the API call.
			// Update testing
			{
				Config: testAccOrganizationsClaimResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("organizations_claim.test", "id", "example-id"),

					// resource.TestCheckResourceAttr("data.OrganizationsClaims.test", "list.#", "2"),

					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.policy", "deny"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.protocol", "tcp"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.src_port", "1,33"),
					// resource.TestCheckResourceAttr("organizations_claim.test", "rules.0.dst_port", "22-30"),
					// resource.TestCheckResourceAttr("OrganizationsClaims.test", "list.1.rules.0.policy", "allow"),
					// resource.TestCheckResourceAttr("OrganizationsClaims.test", "list.1.rules.0.protocol", "any"),
					// resource.TestCheckResourceAttr("OrganizationsClaims.test", "list.1.rules.0.src_port", "any"),
					// resource.TestCheckResourceAttr("OrganizationsClaims.test", "list.1.rules.0.dst_port", "any"),
				),
			},

			// TODO - ImportState testing - An import statement should ONLY include the required attributes to make a Read func call (example: organizationId + networkId).
			// TODO - Currently This only works with hard-coded values so if you find a dynamic way to test please update these template.
			/*
				{
						ResourceName:      "meraki_Organizations_claim.test",
						ImportState:       true,
						ImportStateVerify: false,
						ImportStateId:     "1234567890, 0987654321",
					},
			*/

			// TODO - Check your test environment for dangling resources. During the early stages of development it is not uncommon to find organizations,
			// TODO - networks or admins which did not get picked up because the resource errored out before the delete stage.
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TODO - Usually we need to create an organization. Determine if this makes sense for your workflow.
/*
const testAccAccOrganizationsClaimResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_Organizations_claim"
 	api_enabled = true
 }
 `
*/

// TODO - Create your resource, make sure to include only the applicable attributes modifiable for CREATE.
const testAccAccOrganizationsClaimResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_Organizations_claim" "test" {
	organization_id = resource.meraki_organization.test.organization_id
       {
			"orders": [
				"4CXXXXXXX"
			],
			"serials": [
				"Q234-ABCD-5678"
			],
			"licenses": [
				{
					"key": "Z2XXXXXXXXXX",
					"mode": "addDevices"
				}
			]
		}
`

// TODO - Update the resource ensuring that all modifiable attributes are tested
/*
const testAccAccOrganizationsClaimResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_Organizations_claim" "test" {
	organization_id  = resource.meraki_organization.test.organization_id
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
*/
