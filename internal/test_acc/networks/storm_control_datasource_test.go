package networks

/*
import (
	"fmt"
)


func TestAccNetworkStormControlDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkStormControlDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_switch_storm_control_data"),
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

			// Claim Device
			{
				Config: testAccNetworkStormControlDataSourceClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_switch_storm_control_data"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
				),
			},

			// Create and Read Networks Switch Qos Rules.
			{
				Config: testAccNetworkStormControlDataSourceConfigCreateNetworkStormControl(os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "broadcast_threshold", "90"),
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "multicast_threshold", "90"),
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "unknown_unicast_threshold", "90"),
				),
			},

			// Read Datasource
			{
				Config: testAccNetworkStormControlDataSourceConfigReadNetworkStormControl,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_storm_control.test", "broadcast_threshold", "90"),
					resource.TestCheckResourceAttr("data.meraki_networks_storm_control.test", "multicast_threshold", "90"),
					resource.TestCheckResourceAttr("data.meraki_networks_storm_control.test", "unknown_unicast_threshold", "90"),
				),
			},
		},
	})
}


func testAccNetworkStormControlDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_network_switch_storm_control_data"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}

`, orgId)
	return result
}

// testAccNetworkStormControlDataSourceClaimNetworkDevice is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccNetworkStormControlDataSourceClaimNetworkDevice(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_devices_claim" "test" {
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
`, serial)
	return result
}

func testAccNetworkStormControlDataSourceConfigCreateNetworkStormControl(serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_networks_storm_control" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 90
	multicast_threshold = 90
	unknown_unicast_threshold = 90
}

resource "meraki_devices_switch_port" "test" {
	depends_on = [resource.meraki_networks_storm_control.test]
	serial = "%s"
	storm_control_enabled = true
	port_id = 1
}

`, serial)
	return result
}

const testAccNetworkStormControlDataSourceConfigReadNetworkStormControl = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}

resource "meraki_networks_storm_control" "test" {
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 90
	multicast_threshold = 90
	unknown_unicast_threshold = 90
}

resource "meraki_devices_switch_port" "test" {
	storm_control_enabled = true
	port_id = 1
}

data "meraki_networks_storm_control" "test" {
	network_id = resource.meraki_network.test.network_id
}

`
*/
