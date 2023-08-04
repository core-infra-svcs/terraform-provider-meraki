package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksNetFlowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksNetFlowResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_netflow"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksNetFlowResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Update and Read Networks NetFlow.
			{
				Config: testAccNetFlowResourceConfigUpdateNetworkNetFlowSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_netflow.test", "reporting_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_netflow.test", "eta_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_networks_netflow.test", "collector_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("meraki_networks_netflow.test", "collector_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_netflow.test", "eta_dst_port", "443"),
				),
			},
		},
	})
}

const testAccNetworksNetFlowResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_networks_netflow"
 	api_enabled = true
 } 
 `
const testAccNetworksNetFlowResourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}
 resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_network"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
 `

const testAccNetFlowResourceConfigUpdateNetworkNetFlowSettings = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
	
}
resource "meraki_networks_netflow" "test" {
	  depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
      network_id = resource.meraki_network.test.network_id
	  reporting_enabled = false     
      eta_enabled = false   
	  collector_ip = "1.2.3.4"
      collector_port = 443 
	  eta_dst_port = 443	  
}
`
