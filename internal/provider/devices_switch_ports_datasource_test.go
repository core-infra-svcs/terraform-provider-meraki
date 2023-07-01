package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDevicesSwitchPortsDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccDevicesSwitchPortsDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccDevicesSwitchPortsDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_devices_switch_ports"),
				),
			},

			// Create and Read a Network.
			{
				Config: testAccDevicesSwitchPortsDataSourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
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

			// Read Devices Switch Ports
			{
				Config: testAccDevicesSwitchPortsDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_MS_SERIAL"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.#", "52"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.port_id", "49"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.poe_enabled", "false"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.type", "trunk"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.vlan", "1"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.allowed_vlans", "all"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.isolation_enabled", "false"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.rstp_enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.stp_guard", "disabled"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.link_negotiation", "Auto negotiate"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.udld", "Alert only"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.access_policy_type", "Open"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.dai_trusted", "false"),
				),
			},
		},
	})
}

// testAccDevicesSwitchPortsDataSourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccDevicesSwitchPortsDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_devices_switch_ports"
 	api_enabled = true
 }
 `

// testAccDevicesSwitchPortsDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
const testAccDevicesSwitchPortsDataSourceConfigCreateNetwork = `
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

// testAccDevicesSwitchPortsDataSourceConfigRead is a constant string that defines the configuration for creating and updating a devices_switch_ports_dataSource resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesSwitchPortsDataSourceConfigRead(serial1 string, serial2 string) string {
	result := fmt.Sprintf(`
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
        depends_on = [resource.meraki_organization.test]
        product_types = ["appliance", "switch", "wireless", "cellularGateway"]
}
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
data "meraki_devices_switch_ports" "test" {
	serial = "%s"
}
`, serial1, serial2)
	return result
}
