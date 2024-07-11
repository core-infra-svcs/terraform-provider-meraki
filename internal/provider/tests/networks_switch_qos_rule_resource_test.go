package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchQosRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkSwitchQosRuleResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_switch_qos_rule"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read Networks Switch Qos Rules.
			{
				Config: testAccNetworkSwitchQosRuleResourceConfigCreateNetworkSwitchQosRule,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "vlan", "100"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "dst_port", "3000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "src_port", "2000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "dscp", "0"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "protocol", "TCP"),
				),
			},

			// Update Networks Switch Qos Rules.
			{
				Config: testAccNetworkSwitchQosRuleResourceConfigUpdateNetworkSwitchQosRule,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "vlan", "101"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "dst_port", "4000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "src_port", "3000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "dscp", "0"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rule.test", "protocol", "UDP"),
				),
			},

			// Import testing
			{
				ResourceName:      "meraki_networks_switch_qos_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkSwitchQosRuleResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_network_switch_qos_rule"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworkSwitchQosRuleResourceConfigCreateNetworkSwitchQosRule = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_switch_qos_rule" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 100
    protocol = "TCP"
    src_port = 2000
    dst_port = 3000
    dscp = 0
}
`

const testAccNetworkSwitchQosRuleResourceConfigUpdateNetworkSwitchQosRule = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_switch_qos_rule" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 101
    protocol = "UDP"
    src_port = 3000
    dst_port = 4000
    dscp = 0
}
`
