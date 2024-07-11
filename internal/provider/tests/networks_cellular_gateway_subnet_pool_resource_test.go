package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksCellularGatewaySubnetPoolResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksCellularGatewaySubnetPoolResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksCellularGatewaySubnetPoolResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_cellular_gateway_subnet_pool"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "4"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "cellularGateway"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.3", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// TODO: Update and Read NetworksCellularGatewaySubnetPool
			{
				Config: testAccNetworksCellularGatewaySubnetPoolResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_networks_cellular_gateway_subnet_pool.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_cellular_gateway_subnet_pool.test", "deployment_mode", "routed"),
					resource.TestCheckResourceAttr("meraki_networks_cellular_gateway_subnet_pool.test", "cidr", "192.168.0.0/22"),
					resource.TestCheckResourceAttr("meraki_networks_cellular_gateway_subnet_pool.test", "mask", "24"),
				),
			},
		},
		// ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_networks_cellular_gateway_subnet_pool.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890",
		   },
		*/

	})
}

// testAccNetworksCellularGatewaySubnetPoolResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksCellularGatewaySubnetPoolResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless", "cellularGateway"]
	tags = ["tag1"]
	name = "test_acc_networks_cellular_gateway_subnet_pool"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccNetworksCellularGatewaySubnetPoolResourceConfigUpdate is a constant string that defines the configuration for updating a networks_cellularGateway_subnetPool resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksCellularGatewaySubnetPoolResourceConfigUpdate = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless", "cellularGateway"]
}

resource "meraki_networks_cellular_gateway_subnet_pool" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
    cidr = "192.168.0.0/22"
    mask = 24    
}
`
