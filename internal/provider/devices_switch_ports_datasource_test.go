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

			// Create and Read a Network.
			{
				Config: testAccDevicesSwitchPortsDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_device_switch_ports"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Read Devices Switch Ports
			{
				Config: testAccDevicesSwitchPortsDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "id", "example-id"),
					//resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.#", "52"),
					//resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.port_id", "49"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.poe_enabled", "false"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.type", "trunk"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.vlan", "1"),
					//resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.allowed_vlans", "all"),
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

// testAccDevicesSwitchPortsDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccDevicesSwitchPortsDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["switch"]
	tags = ["tag1"]
	name = "test_acc_device_switch_ports"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesSwitchPortsDataSourceConfigRead is a constant string that defines the configuration for creating and updating a devices_switch_ports_dataSource resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesSwitchPortsDataSourceConfigRead(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	    organization_id = "%s"
        product_types = ["switch"]
}
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
data "meraki_devices_switch_ports" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	serial = "%s"
}
`, orgId, serial, serial)
	return result
}
