package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchDscptocosmappingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworksSwitchDscptocosmappingsResourceOrganizationCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_networks_switch_dscp_to_cos_mappings"),
				),
			},
			// Create and Read Network.
			{
				Config: testAccNetworksSwitchDscptocosmappingsResourceNetworkCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},
			{
				Config: testAccNetworksSwitchDscptocosmappingsResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.0.dscp", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.0.cos", "1"),
				),
			},
		},
	})
}

const testAccNetworksSwitchDscptocosmappingsResourceOrganizationCreate = `
resource "meraki_organization" "test" {
  name = "test_acc_meraki_networks_switch_dscp_to_cos_mappings"
  api_enabled = true
}
`

const testAccNetworksSwitchDscptocosmappingsResourceNetworkCreate = `
resource "meraki_organization" "test" {
  name = "test_meraki_organizations"
  api_enabled = true
}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "Main Office"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}
`

const testAccNetworksSwitchDscptocosmappingsResourceConfigCreate = `
resource "meraki_organization" "test" {
  name = "test_meraki_organizations"
  api_enabled = true
}

resource "meraki_network" "test" {
  depends_on      = ["meraki_organization.test"]
  organization_id = resource.meraki_organization.test.organization_id
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "Main Office"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}

resource "meraki_networks_switch_dscp_to_cos_mappings" "test" {
  depends_on                = ["meraki_organization.test", "meraki_network.test"]
  network_id                = resource.meraki_network.test.network_id
  mappings = [
	{
		dscp = 1
		cos = 1
	}
  ]
}
`
