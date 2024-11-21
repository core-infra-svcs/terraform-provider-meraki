package cellular_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDevicesCellularSimsResource tests the creation, update, and deletion of the devices cellular sims resource.
func TestAccDevicesCellularSimsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesCellularSimsResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_device_cellular_sims"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "cellularGateway"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},
			{
				Config: testAccDevicesCellularSimsResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.#", "1"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.slot", "sim1"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.is_primary", "true"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.apns.#", "0"),
					resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.enabled", "false"),
				),
			},
			{
				ResourceName:      "meraki_devices_cellular_sims.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     os.Getenv("TF_ACC_MERAKI_MG_SERIAL"),
			},
		},
	})
}

// testAccDevicesCellularSimsResourceConfigCreate generates the test configuration for creating a network.
func testAccDevicesCellularSimsResourceConfigCreate(orgID string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["cellularGateway"]
	tags = ["tag1"]
	name = "test_acc_device_cellular_sims"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgID)
}

// testAccDevicesCellularSimsResourceConfigUpdate generates the test configuration for updating a device cellular SIMs resource.
func testAccDevicesCellularSimsResourceConfigUpdate(orgID, serial string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["cellularGateway"]
}

resource "meraki_networks_devices_claim" "test" {
	depends_on = [meraki_network.test]
	network_id = meraki_network.test.network_id
	serials    = ["%s"]
}

resource "meraki_devices_cellular_sims" "test" {
	depends_on  = [meraki_network.test, meraki_networks_devices_claim.test]
	serial      = "%s"
	sims        = [{
		slot       = "sim1"
		apns       = []
		is_primary = true
	}]
	sim_failover = {
		enabled = false
	}
}
`, orgID, serial, serial)
}
