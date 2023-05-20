package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceStaticRoutesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_appliance_static_routes"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateNetwork,
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

			// Create and Read Networks Appliance Static Routes.
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateNetworksApplianceStaticRoutes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "name", "My route"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "subnet", "192.168.o.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_ip", "192.168.0.0"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_static_routes.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Update testing
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigUpdateNetworksApplianceStaticRoutes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "name", "My route"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_ip", "192.168.1.0"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_vlan_id", "1.2.3.5"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "enable", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_address", "22:33:44:55:66:77"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_name", "Some client name"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.comment", "A reserved IP range"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.end", "192.168.1.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.start", "192.168.1.0"),
				),
			},
		},
	})
}

const testAccNetworksApplianceStaticRoutesResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
    name = "test_acc_meraki_organizations_networks_appliance_static_routes"
    api_enabled = true
}
`
const testAccNetworksApplianceStaticRoutesResourceConfigCreateNetwork = `
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

const testAccNetworksApplianceStaticRoutesResourceConfigCreateNetworksApplianceStaticRoutes = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
    network_id = "N_784752235069339433"  
	name = "My route"
    subnet = "192.168.0.0/24"
    gateway_ip = "192.168.0.0"
	reserved_ip_ranges = []
	
}
`

const testAccNetworksApplianceStaticRoutesResourceConfigUpdateNetworksApplianceStaticRoutes = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id    
	name = "My route"
    subnet = "192.168.1.0/24"
    gateway_ip = "192.168.1.0"
	gateway_vlan_id = "1.2.3.5"
	enabled = true
	fixed_ip_assignments_mac_address = "22:33:44:55:66:77"
	fixed_ip_assignments_mac_ip_address = "1.2.3.4"
	fixed_ip_assignments_mac_name = "Some client name"
	reserved_ip_ranges = [
        {
            start = "192.168.1.0"
            end = "192.168.1.1"
            comment = "A reserved IP range"
        }
    ]
}
`
