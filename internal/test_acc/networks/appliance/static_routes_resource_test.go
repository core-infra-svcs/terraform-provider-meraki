package appliance

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceStaticRoutesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			test_acc.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_appliance_static_routes"),
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
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "subnet", "192.168.129.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_ip", "192.168.128.1"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksApplianceStaticRoutesResourceConfigUpdateNetworksApplianceStaticRoutes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "name", "My route"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "subnet", "192.168.129.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "gateway_ip", "192.168.128.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "enable", "true"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_address", "22:33:44:55:66:77"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_ip_address", "192.168.128.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "fixed_ip_assignments_mac_name", "Some client name"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.comment", "A reserved IP range"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.start", "192.168.128.1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_static_routes.test", "reserved_ip_ranges.0.end", "192.168.128.2"),
				),
			},
		},
	})
}

func testAccNetworksApplianceStaticRoutesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_appliance_static_routes"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksApplianceStaticRoutesResourceConfigCreateNetworksApplianceStaticRoutes = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id  
    name = "My route"
    subnet = "192.168.129.0/24"
    gateway_ip = "192.168.128.1"
	reserved_ip_ranges = []
	
}
`

const testAccNetworksApplianceStaticRoutesResourceConfigUpdateNetworksApplianceStaticRoutes = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id    
	name = "My route"
    subnet = "192.168.129.0/24"
	fixed_ip_assignments_mac_address = "22:33:44:55:66:77"
	fixed_ip_assignments_mac_ip_address = "192.168.128.1"
	fixed_ip_assignments_mac_name = "Some client name"   
	reserved_ip_ranges = [
        {
            start = "192.168.128.1"
            end = "192.168.128.2"
            comment = "A reserved IP range"
        }
    ]
}
`
