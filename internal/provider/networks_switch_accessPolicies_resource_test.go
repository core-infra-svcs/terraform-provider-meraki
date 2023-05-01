package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksSwitchAccesspoliciesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: testAccNetworksSwitchAccesspoliciesResourceConfigCreateOrganization,
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
		},
	})
}

const testAccNetworksSwitchAccesspoliciesResourceConfigCreateOrganization = `
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

resource "meraki_networks_switch_access_policies" "test" {
  depends_on                = ["meraki_organization.test", "meraki_network.test"]
  network_id                = resource.meraki_network.test.network_id
  access_policy_type        = "Hybrid authentication"
  dot_1x_control_direction  = "inbound"
  host_mode                 = "Single-Host"
  increase_access_speed     = false
  name                      = "Access policy #1"
  radius_accounting_enabled = false
  radius_accounting_servers = []
  radius_coa_support_enabled = false
  radius_group_attribute     = "11"
  radius_servers             = [
    {
      host   = "1.2.3.4"
      port   = 22
      secret = "secret"
    }
  ]
  radius_testing_enabled = false
  url_redirect_walled_garden_enabled = false
  voice_vlan_clients                 = false
}
 `
