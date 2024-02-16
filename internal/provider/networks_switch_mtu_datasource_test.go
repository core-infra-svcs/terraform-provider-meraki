package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSwitchMtuDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkSwitchMtuDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_switch_mtu"),
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

			// Create and Read Networks Switch MTU Rules.
			{
				Config: testAccNetworkSwitchMtuDataSourceConfigCreateNetworkSwitchMtu,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "default_mtu_size", "9578"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switches.#", "3"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switches.0", "Q234-ABCD-0001"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switches.1", "Q234-ABCD-0002"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switches.2", "Q234-ABCD-0003"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switchProfiles.#", "2"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switchProfiles.0", "1284392014819"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.switchProfiles.1", "2983092129865"),
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu", "overrides.0.mtu_size", "1500"),
				),
			},

			// Read testing
			{
				Config: testAccNetworkSwitchMtuDataSourceRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "default_mtu_size", "9578"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switches.#", "3"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switches.0", "Q234-ABCD-0001"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switches.1", "Q234-ABCD-0002"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switches.2", "Q234-ABCD-0003"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switchProfiles.#", "2"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switchProfiles.0", "1284392014819"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switchProfiles.1", "2983092129865"),
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.mtu_size", "1500"),
				),
			},
		},
	})
}

func testAccNetworkSwitchMtuDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_network_switch_mtu"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworkSwitchMtuDataSourceConfigCreateNetworkSwitchMtu = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_switch_mtu" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    overrides {
        switches = ["Q234-ABCD-0001", "Q234-ABCD-0002", "Q234-ABCD-0003"]
        switchProfiles = ["1284392014819", "2983092129865"]
        mtu_size = 1500
    }
}
`

const testAccNetworkSwitchMtuDataSourceRead = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_switch_mtu" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    overrides {
        switches = ["Q234-ABCD-0001", "Q234-ABCD-0002", "Q234-ABCD-0003"]
        switchProfiles = ["1284392014819", "2983092129865"]
        mtu_size = 1500
    }
}

data "meraki_networks_switch_mtu" "test" {
    depends_on = [meraki_network.test, meraki_networks_switch_mtu.test]
    network_id = meraki_network.test.network_id
}
`
