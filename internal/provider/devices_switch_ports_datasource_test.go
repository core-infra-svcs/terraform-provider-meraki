package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesSwitchPortsDataSource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	serial := os.Getenv("TF_ACC_MERAKI_MS_SERIAL")
	ports := 24 // The number of ports to include in the test, make sure to check the max of your switch

	// Configuration for claiming a device
	claimConfig := testAccDevicesSwitchPortsDataSourceConfigClaimDevice(orgId, serial)
	networkConfig := testAccDevicesSwitchPortsDataSourceConfigCreateNetwork(orgId)
	read := testAccDevicesSwitchPortsDataSourceRead(serial)

	// Generate switch port configurations for each port
	portConfigs := ""
	for i := 1; i <= ports; i++ {
		portConfigs += fmt.Sprintf(`
resource "meraki_devices_switch_port" "test_%[1]d" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    serial = "%[2]s"
    port_id = %[1]d
    enabled = true
    type = "access"
    poe_enabled = true
    isolation_enabled = false
    rstp_enabled = true
    stp_guard = "disabled"
    link_negotiation = "Auto negotiate"
    udld = "Alert only"
    dai_trusted = false
    vlan = 10
    voice_vlan = 20
    allowed_vlans = "all"
    profile = { 
        enabled = false
        iname = ""
        id="0"
    }
}
`, i, serial)
	}

	fullConfig := networkConfig + claimConfig + portConfigs + read

	// Prepare a slice to hold all check functions
	var checks []resource.TestCheckFunc
	for i := 1; i <= ports; i++ {
		prefix := fmt.Sprintf("meraki_devices_switch_port.test_%d", i)
		checks = append(checks,
			resource.TestCheckResourceAttr(prefix, "port_id", fmt.Sprintf("%d", i)),
			resource.TestCheckResourceAttr(prefix, "enabled", "true"),
			resource.TestCheckResourceAttr(prefix, "type", "access"),
			resource.TestCheckResourceAttr(prefix, "poe_enabled", "true"),
			resource.TestCheckResourceAttr(prefix, "isolation_enabled", "false"),
			resource.TestCheckResourceAttr(prefix, "rstp_enabled", "true"),
			resource.TestCheckResourceAttr(prefix, "stp_guard", "disabled"),
			resource.TestCheckResourceAttr(prefix, "link_negotiation", "Auto negotiate"),
			resource.TestCheckResourceAttr(prefix, "udld", "Alert only"),
			resource.TestCheckResourceAttr(prefix, "dai_trusted", "false"),
			resource.TestCheckResourceAttr(prefix, "vlan", "10"),
			resource.TestCheckResourceAttr(prefix, "voice_vlan", "20"),
			resource.TestCheckResourceAttr(prefix, "allowed_vlans", "all"),
			resource.TestCheckResourceAttr(prefix, "profile.enabled", "false"),
			resource.TestCheckResourceAttr(prefix, "profile.iname", ""),
			resource.TestCheckResourceAttr(prefix, "profile.id", "0"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

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

			// Tests the datasource by reading a blank switch
			{
				Config: testAccDevicesSwitchPortsDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.poe_enabled", "false"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.type", "trunk"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.vlan", "1"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.isolation_enabled", "false"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.rstp_enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.stp_guard", "disabled"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.link_negotiation", "Auto negotiate"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.udld", "Alert only"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.access_policy_type", "Open"),
					resource.TestCheckResourceAttr("data.meraki_devices_switch_ports.test", "list.0.dai_trusted", "false"),
				),
			},

			// Tests the datasource by first adding configuration to n number of ports and then reading the configured switch
			{
				Config: fullConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

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

// Tests the datasource by reading a blank/unconfigured switch
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

func testAccDevicesSwitchPortsDataSourceConfigClaimDevice(orgId string, serial string) string {
	return fmt.Sprintf(`
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = ["%s"]
}
`, serial)
}

func testAccDevicesSwitchPortsDataSourceRead(serial string) string {
	result := fmt.Sprintf(`
data "meraki_devices_switch_ports" "test" {
	serial = "%s"
}
`, serial)
	return result
}
