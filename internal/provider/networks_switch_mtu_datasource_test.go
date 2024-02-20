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

			// TODO: Create Switch Profile

			// TODO: Claim Network Device
			/*
				{
						Config: testAccNetworkSwitchMtuDataSourceConfigClaimDevice(os.Getenv("TF_ACC_MERAKI_MS_SERIAL"), os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.0", os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
						),
					},
			*/

			// Create and Read Networks Switch MTU.
			{
				Config: testAccNetworkSwitchMtuDataSourceConfigCreateNetworkSwitchMtu(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_switch_mtu.test", "default_mtu_size", "1500"),

					// TODO: Create Override Settings
					// resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switch_profiles.0", "meraki_switch_profile.test.id"),
					// resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.switches.0", os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
					// resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "overrides.0.mtu_size", "9578"),
				),
			},

			// Update Switch MTU
			{
				Config: testAccNetworkSwitchMtuDataSourceRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_switch_mtu.test", "default_mtu_size", "9578"),
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

// TODO: func/const for testAccNetworkSwitchMtuDataSourceConfigSwitchProfile

func testAccNetworkSwitchMtuDataSourceConfigClaimDevice(serial string, orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
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

// TODO: Uncomment overrides
func testAccNetworkSwitchMtuDataSourceConfigCreateNetworkSwitchMtu(serial string) string {
	result := fmt.Sprintf(`

	resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
	}

	resource "meraki_networks_switch_mtu" "test" {
    	depends_on     = [meraki_network.test]
    	network_id     = meraki_network.test.network_id
    	default_mtu_size = 1500
		overrides = [

		/*
		{
				switches = [
					"%s"
				],
				switch_profiles = [
					meraki_switch_profile.test.id
				],
				mtu_size = 1540
    		}
		*/

		]
		
}
`, serial)
	return result
}

const testAccNetworkSwitchMtuDataSourceRead = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_switch_mtu" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    default_mtu_size = 9578
    overrides = []
}

data "meraki_networks_switch_mtu" "test" {
    depends_on = [meraki_network.test, meraki_networks_switch_mtu.test]
    network_id = meraki_network.test.network_id
}
`
