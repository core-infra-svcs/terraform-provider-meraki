package provider

// TODO: This test is only valid for devices that support this feature. MS120's do not.
/*
func TestAccNetworkStormControlResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworkStormControlResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"),
					os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_switch_storm_control"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read Networks Switch Qos Rules.
			{
				Config: testAccNetworkStormControlResourceConfigCreateNetworkStormControl(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"),
					os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "broadcast_threshold", "90"),
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "multicast_threshold", "90"),
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "unknown_unicast_threshold", "90"),
				),
			},

			//Update Networks Switch Qos Rules.
			{
				Config: testAccNetworkSwiStormControlResourceConfigUpdateNetworkStormControl(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"),
					os.Getenv("TF_ACC_MERAKI_MS_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "broadcast_threshold", "40"),
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "multicast_threshold", "40"),
					resource.TestCheckResourceAttr("meraki_networks_storm_control.test", "unknown_unicast_threshold", "40"),
				),
			},

			// Import testing
			{
				ResourceName:      "meraki_networks_storm_control.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkStormControlResourceConfigCreateNetwork(orgId, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["switch"]
    tags = ["tag1"]
    name = "test_acc_network_switch_storm_control"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
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

func testAccNetworkStormControlResourceConfigCreateNetworkStormControl(orgId, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["switch"]
}

resource "meraki_networks_storm_control" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 90
	multicast_threshold = 90
	unknown_unicast_threshold = 90
}

resource "meraki_devices_switch_port" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_storm_control.test]
	serial = "%s"
	storm_control_enabled = true
	port_id = 1
}

`, orgId, serial)
	return result
}

func testAccNetworkSwiStormControlResourceConfigUpdateNetworkStormControl(orgId, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["switch"]
}


resource "meraki_networks_storm_control" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
	broadcast_threshold = 40
	multicast_threshold = 40
	unknown_unicast_threshold = 40
}

resource "meraki_devices_switch_port" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_storm_control.test]
	serial = "%s"
	storm_control_enabled = true
	port_id = 1
}

`, orgId, serial)
	return result
}
*/
