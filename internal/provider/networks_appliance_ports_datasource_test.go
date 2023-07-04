package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkAppliancePortsDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_ports.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network .
			{
				Config: testAccNetworkAppliancePortsDatasourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			//TODO To List Network Appliance Ports VLANs enabled Network Needed for Testing.
			//  List Network Appliance Ports.
			{
				Config: testAccNetworkAppliancePortsDatasourceConfigListNetworkAppliancePorts(os.Getenv("TF_ACC_MERAKI_MX_SERIAL"), os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.#", "3"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.0.allowed_vlans", "all"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.0.drop_untagged_traffic", "true"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.0.number", "3"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.0.type", "trunk"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_ports.test", "list.0.vlan", "0"),
				),
			},
		},
	})
}

func testAccNetworkAppliancePortsDatasourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
 resource "meraki_network" "test" {	
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
 `, orgId)
	return result
}

// testAccNetworkAppliancePortsDatasourceConfigListNetworkAppliancePorts is a constant string that defines the configuration for reading a networks_appliance_ports datasource in your tests.
// It depends on both the organization and network resources.
func testAccNetworkAppliancePortsDatasourceConfigListNetworkAppliancePorts(serial string, orgId string) string {
	result := fmt.Sprintf(`
	resource "meraki_network" "test" {
		organization_id = "%s"
		product_types = ["appliance"]			
	}
	resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
    ]
    }
	resource "meraki_networks_appliance_vlans_settings" "test" {
		depends_on = [resource.meraki_network.test, meraki_networks_devices_claim.test]
		network_id = resource.meraki_network.test.network_id
		vlans_enabled = true
	}  
    data "meraki_networks_appliance_ports" "test" {
	depends_on = [resource.meraki_network.test, meraki_networks_devices_claim.test, meraki_networks_appliance_vlans_settings.test]	
	network_id = resource.meraki_network.test.network_id	
    }	
`, orgId, serial)
	return result
}
