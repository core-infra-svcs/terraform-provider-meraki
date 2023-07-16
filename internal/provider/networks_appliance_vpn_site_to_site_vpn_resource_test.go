package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceVpnSiteToSiteVpnResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network.
			{
				Config: testAccNetworksApplianceVpnSiteToSiteVpnResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Site1"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},
			// Update and Read Networks Appliance Vpn Site To Site Vpn.
			{
				Config: testAccApplianceVpnSiteToSiteVpnResourceConfigUpdateNetworkApplianceVpnSiteToSiteVpnSettings(os.Getenv("TF_ACC_MERAKI_ORGANZIATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL"), os.Getenv("TF_ACC_MAIN_OFFICE_SUB_TEST_NETWORK_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vpn_site_to_site_vpn.test", "mode", "hub"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vpn_site_to_site_vpn.test", "hubs.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vpn_site_to_site_vpn.test", "hubs.0.hub_id", os.Getenv("TF_ACC_MAIN_OFFICE_SUB_TEST_NETWORK_ID")),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vpn_site_to_site_vpn.test", "subnets.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_appliance_vpn_site_to_site_vpn.test", "subnets.0.local_subnet", "192.168.128.0/24"),
				),
			},
		},
	})
}

func testAccNetworksApplianceVpnSiteToSiteVpnResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance"]
    tags = ["tag1"]
    name = "Site1"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

func testAccApplianceVpnSiteToSiteVpnResourceConfigUpdateNetworkApplianceVpnSiteToSiteVpnSettings(orgId string, serial string, hub_id string) string {
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
resource "meraki_networks_appliance_vpn_site_to_site_vpn" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	network_id = resource.meraki_network.test.network_id
    mode = "hub"
    hubs = [{
		hub_id = "%s"
		
		}]
    subnets = [{
		local_subnet = "192.168.128.0/24"
	}]
}
`, orgId, serial, hub_id)
	return result
}
