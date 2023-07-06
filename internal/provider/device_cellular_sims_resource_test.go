package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesCellularSimsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccDevicesCellularSimsResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID")),
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

			// Update testing
			{
				Config: testAccDevicesCellularSimsResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.#", "1"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.slot", "sim1"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.is_primary", "true"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.apns.#", "0"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.enabled", "false"),
				),
			},

			/*
				{
						ResourceName:      "meraki_devices_cellular_sims.test",
						ImportState:       true,
						ImportStateVerify: false,
						ImportStateId:     "1234567890",
					},
			*/

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDevicesCellularSimsResourceConfigCreate(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance", "switch", "wireless", "cellularGateway"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesCellularSimsResourceConfigUpdate is a constant string that defines the configuration for creating and updating a devices cellular sims resource config update resource in your tests.
// It depends on network resources.
func testAccDevicesCellularSimsResourceConfigUpdate(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance", "switch", "wireless", "cellularGateway"]
}
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
resource "meraki_devices_cellular_sims" "test" {
	serial = "%s"
	sims = [{
		slot = "sim1"
		apns = []
		is_primary = true
	}]
	sim_failover = {
		enabled = false
	}
	
}
`, orgId, serial, serial)
	return result
}
