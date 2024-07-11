package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDevicesSwitchPortsCycleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccDevicesSwitchPortsCycleResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_meraki_devices_switch_ports_cycle"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			/*
				// Claim A Device To A Network
				{
					Config: testAccDevicesSwitchPortsCycleResourceConfigClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "id", "example-id"),
					),
				},

				// Create and Read testing
				{
					Config: testAccDevicesSwitchPortsCycleResourceConfigCycle(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_devices_switch_ports_cycle.test", "id", "example-id"),
					),
				},
			*/

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDevicesSwitchPortsCycleResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["switch"]
	tags = ["tag1"]
	name = "test_meraki_devices_switch_ports_cycle"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// TODO: Response Body: {"errors":["Temporary overload, please try again"]} (Need a switch that is online for this to work
/*
func testAccDevicesSwitchPortsCycleResourceConfigClaimNetworkDevice(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["switch"]
	tags = ["tag1"]
	name = "test_meraki_devices_switch_ports_cycle"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
`, orgId, serial)
	return result
}


// It seems like this is a bug or a response given when the switch is not online.

func testAccDevicesSwitchPortsCycleResourceConfigCycle(orgId string, serial string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["switch"]
		tags = ["tag1"]
		name = "test_meraki_devices_switch_ports_cycle"
		timezone = "America/Los_Angeles"
		notes = "Additional description of the network"
}

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

resource "meraki_devices_switch_ports_cycle" "test" {
	depends_on = [resource.meraki_network.test, meraki_networks_devices_claim.test]
	serial = "%s"
	ports = ["3"]
}
`, orgId, serial, serial)
	return result
}

*/
