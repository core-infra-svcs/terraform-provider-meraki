package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDevicesSwitchPortResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	serial := os.Getenv("TF_ACC_MERAKI_MS_SERIAL")
	ports := 24 // The number of ports to include in the test

	// Configuration for claiming a device
	claimConfig := testAccDevicesSwitchPortResourceConfigClaimDevice(orgId, serial)
	networkConfig := testAccDevicesSwitchPortResourceConfigCreateNetwork(orgId)

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

	fullConfig := networkConfig + claimConfig + portConfigs

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
				Config: testAccDevicesSwitchPortResourceConfigCreateNetwork(orgId),
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

			{
				Config: fullConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
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
