package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchAccesspoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccOrganizationsNetworksDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_organizations_networks"),
				),
			},

			// Create and Read Network
			{
				Config: testAccOrganizationsNetworksCreateAccessPolicyDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Delete testing automatically occurs in TestCase

			// Create access policy for switch network
			//{
			//	Config: testAccOrganizationsNetworksCreateAccessPolicyDataSourceConfig,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
			//	),
			//},
		},
	})
}

const testAccOrganizationsNetworksCreateAccessPolicyDataSourceConfig = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "Main Office"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}

data "meraki_networks_switch_access_policies" "test" {
    network_id = resource.meraki_network.test.network_id
}
`
