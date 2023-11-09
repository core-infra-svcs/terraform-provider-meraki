package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDevicesSwitchPortResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccDevicesSwitchPortResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccDevicesSwitchPortResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_devices_switch_port_resource"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Read Devices Switch Port
			{
				Config: testAccDevicesSwitchPortResourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "port_id", "1"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "enabled", "true"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "poe_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "type", "access"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "isolation_enabled", "false"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "link_negotiation", "Auto negotiate"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "udld", "Alert only"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "vlan", "10"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "voice_vlan", "20"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "allowed_vlans", "all"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "sticky_mac_allow_list_limit", "5"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "dai_trusted", "false"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "access_policy_type", "Sticky MAC allow list"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "profile.enabled", "false"),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "profile.iname", ""),
					resource.TestCheckResourceAttr("meraki_devices_switch_port.test", "profile.id", "0"),
				),
			},
		},
	})
}

// testAccDevicesSwitchPortResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccDevicesSwitchPortResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["switch"]
	tags = ["tag1"]
	name = "test_acc_devices_switch_port_resource"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesSwitchPortResourceConfigRead is a constant string that defines the configuration for creating and updating a devices_switch_ports_dataSource resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesSwitchPortResourceConfigRead(orgId string, serial string) string {
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

resource "meraki_devices_switch_port" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	serial = "%s"
	name = "My switch port"	
	tags = ["tag1", "tag2"]
	port_id = 1
	enabled = true
	type = "access"
	poe_enabled = true
	isolation_enabled = false
	rstp_enabled = true
	stp_guard = "disabled"
	link_negotiation = "Auto negotiate"	
	udld = "Alert only"	
	dai_trusted = false	
	sticky_mac_allow_list_limit = 5
	access_policy_type = "Sticky MAC allow list"
	//access_policy_number = 2
	//mac_allow_list = ["34:56:fe:ce:8e:b0", "34:56:fe:ce:8e:b1"]	
	//sticky_mac_allow_list = ["34:56:fe:ce:8e:b0", "34:56:fe:ce:8e:b1"]
	//adaptive_policy_group_id = "123"
	//port_schedule_id = "1234"
	//peer_sgt_capable = false
	//flexible_stacking_enabled = false
	//storm_control_enabled = false
	vlan = 10
	voice_vlan = 20
	allowed_vlans = "all"
	profile = { 
		enabled = false
		iname = ""
		id="0"
	}
	
}
`, orgId, serial, serial)
	return result
}
