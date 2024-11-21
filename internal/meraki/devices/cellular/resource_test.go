package cellular_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDevicesCellularSimsResource tests the full lifecycle of the devices cellular sims resource.
func TestAccDevicesCellularSimsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesCellularSimsResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check:  testAccDevicesCellularSimsCheckCreate(),
			},
			{
				Config: testAccDevicesCellularSimsResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check:  testAccDevicesCellularSimsCheckUpdate(),
			},
			{
				ResourceName:      "meraki_devices_cellular_sims.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     os.Getenv("TF_ACC_MERAKI_MG_SERIAL"),
				Check:             testAccDevicesCellularSimsCheckImport(),
			},
		},
	})
}

// testAccDevicesCellularSimsCheckCreate validates the attributes after creation.
func testAccDevicesCellularSimsCheckCreate() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "serial", os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.#", "1"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.slot", "sim1"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.is_primary", "true"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.enabled", "false"),
		resource.TestCheckNoResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.timeout"),
	)
}

// testAccDevicesCellularSimsCheckUpdate validates the attributes after update.
func testAccDevicesCellularSimsCheckUpdate() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.#", "1"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.slot", "sim1"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.is_primary", "true"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.enabled", "true"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.timeout", "300"),
	)
}

// testAccDevicesCellularSimsCheckImport validates the attributes after import.
func testAccDevicesCellularSimsCheckImport() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "serial", os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.#", "1"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sims.0.slot", "sim1"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.enabled", "true"),
		resource.TestCheckResourceAttr("meraki_devices_cellular_sims.test", "sim_failover.timeout", "300"),
	)
}

// testAccDevicesCellularSimsResourceConfigCreate generates the test configuration for creating a network with a device cellular SIMs resource.
func testAccDevicesCellularSimsResourceConfigCreate(orgID string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types   = ["cellularGateway"]
	tags            = ["tag1"]
	name            = "test_acc_device_cellular_sims"
	timezone        = "America/Los_Angeles"
	notes           = "Additional description of the network"
}

resource "meraki_devices_cellular_sims" "test" {
	serial = "%s"
	sims = [{
		slot       = "sim1"
		apns       = []
		is_primary = true
	}]
	sim_failover = {
		enabled = false
	}
}
`, orgID, os.Getenv("TF_ACC_MERAKI_MG_SERIAL"))
}

// testAccDevicesCellularSimsResourceConfigUpdate generates the test configuration for updating the device cellular SIMs resource.
func testAccDevicesCellularSimsResourceConfigUpdate(orgID, serial string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types   = ["cellularGateway"]
}

resource "meraki_networks_devices_claim" "test" {
	depends_on = [meraki_network.test]
	network_id = meraki_network.test.network_id
	serials    = ["%s"]
}

resource "meraki_devices_cellular_sims" "test" {
	depends_on  = [meraki_network.test, meraki_networks_devices_claim.test]
	serial      = "%s"
	sims = [{
		slot       = "sim1"
		apns       = []
		is_primary = true
	}]
	sim_failover = {
		enabled = true
		timeout = 300
	}
}
`, orgID, serial, serial)
}
