package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchQosRulesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkSwitchQosRulesResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_switch_qos_rules"),
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
				Config: testAccNetworkSwitchQosRulesResourceConfigCreateNetworkSwitchQosRules,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "vlan", "100"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "dst_port", "3000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "src_port", "2000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "dscp", "0"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "protocol", "TCP"),
				),
			},

			// Update Networks Switch Qos Rules.
			{
				Config: testAccNetworkSwitchQosRulesResourceConfigUpdateNetworkSwitchQosRules,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "vlan", "101"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "dst_port", "4000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "src_port", "3000"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "dscp", "0"),
					resource.TestCheckResourceAttr("meraki_networks_switch_qos_rules.test", "protocol", "UDP"),
				),
			},
		},
	})
}

func testAccNetworkSwitchQosRulesResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_network_switch_qos_rules"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworkSwitchQosRulesResourceConfigCreateNetworkSwitchQosRules = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_switch_qos_rules" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 100
    protocol = "TCP"
    src_port = 2000
    dst_port = 3000
    dscp = 0
}
`

const testAccNetworkSwitchQosRulesResourceConfigUpdateNetworkSwitchQosRules = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_switch_qos_rules" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan = 101
    protocol = "UDP"
    src_port = 3000
    dst_port = 4000
    dscp = 0
}
`
