package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksCellularGatewayDhcpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksCellularGatewayDhcpResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_cellular_gateway_dhcp"),
				),
			},

			// Create and Read testing
			{
				Config: testAccNetworksCellularGatewayDhcpResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "4"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "cellularGateway"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.3", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksCellularGatewayDhcpResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_cellular_gateway_dhcp.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_cellular_gateway_dhcp.test", "dns_custom_name_servers.0", "1.1.1.1"),
				),
			},
		},
	})
}

const testAccNetworksCellularGatewayDhcpResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_cellular_gateway_dhcp"
 	api_enabled = true
 }
 `

const testAccNetworksCellularGatewayDhcpResourceConfigCreate = `
resource "meraki_organization" "test" {}
 resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless", "cellularGateway"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccNetworksCellularGatewayDhcpResourceConfigUpdate = `
resource "meraki_organization" "test" {}
 resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless", "cellularGateway"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_cellular_gateway_dhcp" "test" {
    depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
	dhcp_lease_time =  "1 hour"
	dns_name_servers = "custom"
	dns_custom_name_servers = ["1.1.1.1", "8.8.8.8"]
}
`
