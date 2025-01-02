package port_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDevicesSwitchPortResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	serial := os.Getenv("TF_ACC_MERAKI_MS_SERIAL")
	ports := 2

	// Configuration for claiming a device and creating a network
	claimConfig := testAccDevicesSwitchPortResourceConfigClaimDevice(orgId, serial)
	networkConfig := testAccDevicesSwitchPortResourceConfigCreateNetwork(orgId)

	// Generate access port test configuration
	accessConfig, accessChecks := generateAccessPortConfig(serial, ports)
	fullAccessConfig := networkConfig + claimConfig + accessConfig

	// Generate access port test configuration update
	accessConfigUpdate, accessChecksUpdate := generateAccessPortConfigUpdate(serial, ports)
	fullAccessConfigUpdate := networkConfig + claimConfig + accessConfigUpdate

	// Generate trunk port test configuration
	trunkConfig, trunkChecks := generateTrunkPortConfig(serial, ports)
	fullTrunkConfig := networkConfig + claimConfig + trunkConfig

	// Generate trunk port test configuration update
	trunkConfigUpdate, trunkChecksUpdate := generateTrunkPortConfigUpdate(serial, ports)
	fullTrunkConfigUpdate := networkConfig + claimConfig + trunkConfigUpdate

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// bulk access port test
			{
				Config: fullAccessConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(accessChecks...),
			},

			// bulk access port update test
			{
				Config: fullAccessConfigUpdate,
				Check:  resource.ComposeAggregateTestCheckFunc(accessChecksUpdate...),
			},

			// bulk trunk port test
			{
				Config: fullTrunkConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(trunkChecks...),
			},

			// bulk trunk port update test
			{
				Config: fullTrunkConfigUpdate,
				Check:  resource.ComposeAggregateTestCheckFunc(trunkChecksUpdate...),
			},
		},
	})
}

func testAccDevicesSwitchPortResourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = "%s"
    product_types = ["switch"]
    tags = ["tag1"]
    name = "test_acc_devices_switch_port_resource"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
}

func testAccDevicesSwitchPortResourceConfigClaimDevice(orgId string, serial string) string {
	return fmt.Sprintf(`
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = ["%s"]
}
`, serial)
}

func generateAccessPortConfig(serial string, ports int) (string, []resource.TestCheckFunc) {
	portConfig := ""
	var checks []resource.TestCheckFunc

	for i := 1; i <= ports; i++ {
		portConfig += fmt.Sprintf(`
resource "meraki_devices_switch_port" "test_%[1]d" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    serial = "%[2]s"
    port_id = %[1]d
    enabled = "true"
    type = "access"
}
`, i, serial)

		prefix := fmt.Sprintf("meraki_devices_switch_port.test_%d", i)
		checks = append(checks,
			// attribute checks
			resource.TestCheckResourceAttr(prefix, "enabled", "true"),
			resource.TestCheckResourceAttr(prefix, "type", "access"),
		)
	}
	return portConfig, checks
}

func generateAccessPortConfigUpdate(serial string, ports int) (string, []resource.TestCheckFunc) {
	portConfig := ""
	var checks []resource.TestCheckFunc

	for i := 1; i <= ports; i++ {
		portConfig += fmt.Sprintf(`
resource "meraki_devices_switch_port" "test_%[1]d" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    serial = "%[2]s"
    port_id = %[1]d
    enabled = "false"
    type = "access"
  
}
`, i, serial)

		prefix := fmt.Sprintf("meraki_devices_switch_port.test_%d", i)
		checks = append(checks,
			// attribute checks
			resource.TestCheckResourceAttr(prefix, "enabled", "false"),
			resource.TestCheckResourceAttr(prefix, "type", "access"),
		)
	}
	return portConfig, checks
}

func generateTrunkPortConfig(serial string, ports int) (string, []resource.TestCheckFunc) {
	portConfig := ""
	var checks []resource.TestCheckFunc

	for i := 1; i <= ports; i++ {
		portConfig += fmt.Sprintf(`
resource "meraki_devices_switch_port" "test_%[1]d" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    serial = "%[2]s"
    port_id = %[1]d
    enabled = "true"
    type = "trunk"

}
`, i, serial)

		prefix := fmt.Sprintf("meraki_devices_switch_port.test_%d", i)
		checks = append(checks,
			// attribute checks
			resource.TestCheckResourceAttr(prefix, "enabled", "true"),
			resource.TestCheckResourceAttr(prefix, "type", "trunk"),
		)
	}
	return portConfig, checks
}

func generateTrunkPortConfigUpdate(serial string, ports int) (string, []resource.TestCheckFunc) {
	portConfig := ""
	var checks []resource.TestCheckFunc

	for i := 1; i <= ports; i++ {
		portConfig += fmt.Sprintf(`
resource "meraki_devices_switch_port" "test_%[1]d" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    serial = "%[2]s"
    port_id = %[1]d
    enabled = "false"
    type = "trunk"
 
}
`, i, serial)

		prefix := fmt.Sprintf("meraki_devices_switch_port.test_%d", i)
		checks = append(checks,
			// attribute checks
			resource.TestCheckResourceAttr(prefix, "enabled", "false"),
			resource.TestCheckResourceAttr(prefix, "type", "trunk"),
		)
	}
	return portConfig, checks
}
