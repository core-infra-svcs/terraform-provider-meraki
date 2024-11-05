package _switch_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksSwitchDscpToCosMappingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworksSwitchDscpToCosMappingsResourceNetworkCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_switch_dscp_to_cos_mappings"),
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

			// Create and Read Test
			{
				Config: testAccNetworksSwitchDscpToCosMappingsResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.0.dscp", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.0.cos", "1"),
				),
			},

			// Update and Read Test
			{
				Config: testAccNetworksSwitchDscpToCosMappingsResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.0.dscp", "63"),
					resource.TestCheckResourceAttr("meraki_networks_switch_dscp_to_cos_mappings.test", "mappings.0.cos", "5"),
				),
			},
		},
	})
}

func testAccNetworksSwitchDscpToCosMappingsResourceNetworkCreate(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
  organization_id = %s
  product_types   = ["appliance", "switch", "wireless"]
  tags            = ["tag1", "tag2"]
  name            = "test_acc_networks_switch_dscp_to_cos_mappings"
  timezone        = "America/Los_Angeles"
  notes           = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksSwitchDscpToCosMappingsResourceConfigCreate = `
resource "meraki_network" "test" {
  product_types   = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_switch_dscp_to_cos_mappings" "test" {
  depends_on                = [meraki_network.test]
  network_id                = resource.meraki_network.test.network_id
  mappings = [
	{
		dscp = 1
		cos = 1
	}
  ]
}
`

const testAccNetworksSwitchDscpToCosMappingsResourceConfigUpdate = `
resource "meraki_network" "test" {
  product_types   = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_switch_dscp_to_cos_mappings" "test" {
  depends_on                = [meraki_network.test]
  network_id                = resource.meraki_network.test.network_id
  mappings = [
	{
		dscp = 63
		cos = 5
	}
  ]
}
`
