package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkAppliancePortsDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworkAppliancePortsDatasourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_networks_appliance_ports"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_ports.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network .
			{
				Config: testAccNetworkAppliancePortsDatasourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
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

			/*
				//TODO To List Network Appliance Ports VLANs enabled Network Needed for Testing.
				//  List Network Appliance Ports.
				{
					Config: testAccNetworkAppliancePortsDatasourceConfigListNetworkAppliancePorts,
					Check:  resource.ComposeAggregateTestCheckFunc(

						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "access_policy", "access_policy"),
						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "allowed_vlans", "allowed_vlans"),
						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "drop_untagged_traffic", "true"),
						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "number", "4"),
						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "type", "access"),
						resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "vlan", "12345"),

					),
				},
			*/

		},
	})
}

const testAccNetworkAppliancePortsDatasourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_networks_appliance_ports"
 	api_enabled = true
 } 
 `

const testAccNetworkAppliancePortsDatasourceConfigCreateNetwork = `
 resource "meraki_organization" "test" {}
 resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
 `

const testAccNetworkAppliancePortsDatasourceConfigListNetworkAppliancePorts = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]	
	product_types = ["appliance", "switch", "wireless"]	
}
data "meraki_networks_appliance_ports" "test" {
	  depends_on = [resource.meraki_organization.test,
	  resource.meraki_network.test]	
      network_id = resource.meraki_network.test.network_id
	  
}
`
