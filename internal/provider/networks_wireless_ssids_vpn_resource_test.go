package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!
// TODO - Testing is meant to be atomic in that we give very specific instructions for how to create, read, update, and delete infrastructure across test steps.
// TODO - This is really useful for troubleshooting resources/data sources during development and provides a high level of confidence that our provider works as intended.
func TestAccNetworksWirelessSsidsVpnResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksWirelessSsidsVpnResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_wireless_ssids_vpn"),
				),
			},

			//Create and Read testing
			{
				Config: testAccNetworksWirelessSsidsVpnResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
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

			// Update testing
			{
				Config: testAccNetworksWirelessSsidsVpnResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_vpn.test", "id", "example-id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccNetworksWirelessSsidsVpnResourceConfigCreateOrganization = `
resource "meraki_organization" "test" {
	name = "test_meraki_networks_wireless_ssids_vpn"
	api_enabled = true
}
`

const testAccNetworksWirelessSsidsVpnResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccNetworksWirelessSsidsVpnResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_wireless_ssids_vpn" "test" {
	network_id  = resource.meraki_network.test.network_id
	number = "123"
	concentrator = {
		network_id = resource.meraki_network.test.network_id
		vlan_id = 44
		name = "some concentrator name"
	}
	failover = {
		request_ip = "1.1.1.1"
		idle_timeout = 10
		heartbeat_interval = 30
	}
	split_tunnel = {
		enabled = true
		rules = [
			{
				protocol = "Any"
				dest_cidr = "1.1.1.1/32"
				dest_port = "any"
				policy = "allow"
				comment = "split tunnel rule 1"
			},
			{
				protocol = "Any"
				dest_cidr = "foo.com"
				dest_port = "any"
				policy = "deny"
				comment = "split tunnel rule 2"
			}
		]
	}
}
`
