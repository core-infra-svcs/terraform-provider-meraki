package provider

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!
// TODO - Testing is meant to be atomic in that we give very specific instructions for how to create, read, update, and delete infrastructure across test steps.
// TODO - This is really useful for troubleshooting resources/data sources during development and provides a high level of confidence that our provider works as intended.
func TestAccNetworks{Networkid}WirelessSsids{Number}VpnResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

          // TODO - Usually the first step in building a resource is to create an organization or network to configure.
        /*
        // Create test Organization
                    {
                        Config: testAccNetworks{Networkid}WirelessSsids{Number}VpnResourceConfigCreateOrganization,
                        Check: resource.ComposeAggregateTestCheckFunc(
                            resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
                            resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_{network_id}_wireless_ssids_{number}_vpn"),
                        ),
                    },
        */

            // TODO - Next, run the create test step for the resource you are developing. It is important to validate every field returned by read.
			// Create and Read testing
            			{
                            Config: testAccNetworks{Networkid}WirelessSsids{Number}VpnResourceConfigCreate,
                            Check: resource.ComposeAggregateTestCheckFunc(
                                resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "id", "example-id"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "name", "Block sensitive web traffic"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "description", "Blocks sensitive web traffic"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "ip_version", "ipv6"),

                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.#", "1"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.policy", "deny"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.protocol", "tcp"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.src_port", "1,33"),
                                // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.dst_port", "22-30"),
                            ),
                        },

            // TODO - Once a resource has been created, we will test the ability to modify it. Make sure to test all values that are modifiable by the API call.
			// Update testing
            			{
            				Config: testAccNetworks{Networkid}WirelessSsids{Number}VpnResourceConfigUpdate,
            				Check: resource.ComposeAggregateTestCheckFunc(
            					resource.TestCheckResourceAttr("Networks{Networkid}WirelessSsids{Number}Vpn.test", "id", "example-id"),

                               // resource.TestCheckResourceAttr("data.Networks{Networkid}WirelessSsids{Number}Vpns.test", "list.#", "2"),

                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.policy", "deny"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.protocol", "tcp"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.src_port", "1,33"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpn.test", "rules.0.dst_port", "22-30"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpns.test", "list.1.rules.0.policy", "allow"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpns.test", "list.1.rules.0.protocol", "any"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpns.test", "list.1.rules.0.src_port", "any"),
                               // resource.TestCheckResourceAttr("networks_{network_id}_wireless_ssids_{number}_vpns.test", "list.1.rules.0.dst_port", "any"),
            				),
            			},

			// TODO - ImportState testing - An import statement should ONLY include the required attributes to make a Read func call (example: organizationId + networkId).
			// TODO - Currently This only works with hard-coded values so if you find a dynamic way to test please update these template.
			/*
				{
						ResourceName:      "meraki_networks_{network_id}_wireless_ssids_{number}_vpn.test",
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
const testAccNetworks{Networkid}WirelessSsids{Number}VpnResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_{network_id}_wireless_ssids_{number}_vpn"
 	api_enabled = true
 }
 `
*/

// TODO - Create your resource, make sure to include only the applicable attributes modifiable for CREATE.
const testAccNetworks{Networkid}WirelessSsids{Number}VpnResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_networks_{network_id}_wireless_ssids_{number}_vpn" "test" {
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
            }
        ]
    }
`

// TODO - Update the resource ensuring that all modifiable attributes are tested
/*
const testAccNetworks{Networkid}WirelessSsids{Number}VpnResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_networks_{network_id}_wireless_ssids_{number}_vpn" "test" {
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


