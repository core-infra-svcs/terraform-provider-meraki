package appliance

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceVpnSiteToSiteVpnDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			test_acc.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and Read Network.
				Config: testAccNetworksApplianceVpnSiteToSiteVpnDatasourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_appliance_vpn_site_to_site_vpn"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},
			{
				// Claim Network Device
				Config: testAccApplianceVpnSiteToSiteVpnDatasourceConfigClaimNetworksDevice(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				// Test case for reading Networks Appliance Vpn Site To Site Vpn.
				Config: testAccApplianceVpnSiteToSiteVpnDatasourceConfigReadNetworkApplianceVpnSiteToSiteVpn,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vpn_site_to_site_vpn.test", "mode", "hub"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vpn_site_to_site_vpn.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vpn_site_to_site_vpn.test", "subnets.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_networks_appliance_vpn_site_to_site_vpn.test", "subnets.0.local_subnet", "192.168.128.0/24"),
				),
			},
		},
	})
}

func testAccNetworksApplianceVpnSiteToSiteVpnDatasourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    name = "test_acc_network_appliance_vpn_site_to_site_vpn"
	product_types = ["appliance"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
    
}
`, orgId)
}

func testAccApplianceVpnSiteToSiteVpnDatasourceConfigClaimNetworksDevice(serial string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
   name = "test_acc_network_appliance_vpn_site_to_site_vpn"
	product_types = ["appliance"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"

}
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
	network_id = meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
`, serial)
}

const testAccApplianceVpnSiteToSiteVpnDatasourceConfigReadNetworkApplianceVpnSiteToSiteVpn = `
resource "meraki_network" "test" {
	name = "test_acc_network_appliance_vpn_site_to_site_vpn"
	product_types = ["appliance"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"

}

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
	network_id = meraki_network.test.network_id
}

resource "meraki_networks_appliance_vpn_site_to_site_vpn" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	network_id = resource.meraki_network.test.network_id
    mode = "hub"
    subnets = [
		{
			local_subnet = "192.168.128.0/24"
		}
	]	
}

data "meraki_networks_appliance_vpn_site_to_site_vpn" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test, meraki_networks_appliance_vpn_site_to_site_vpn.test]
	network_id = resource.meraki_network.test.network_id
	mode = resource.meraki_networks_appliance_vpn_site_to_site_vpn.test.mode
	id = resource.meraki_networks_appliance_vpn_site_to_site_vpn.test.id
	subnets = resource.meraki_networks_appliance_vpn_site_to_site_vpn.test.subnets
}
`
