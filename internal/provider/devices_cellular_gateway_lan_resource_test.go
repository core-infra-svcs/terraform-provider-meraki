package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDevicesCellulargatewayLanResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccDevicesCellulargatewayLanResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_devices_cellular_gateway_lan"),
				),
			},

			// Create and Read testing
			{
				Config: testAccDevicesCellulargatewayLanResourceConfigCreate,
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
				Config: testAccDevicesCellulargatewayLanResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("merak_devices_cellular_gateway_lan.test", "id", "example-id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccDevicesCellulargatewayLanResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_devices_cellular_gateway_lan"
 	api_enabled = true
 }
 `

const testAccDevicesCellulargatewayLanResourceConfigCreate = `
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

const testAccDevicesCellulargatewayLanResourceConfigUpdate = `
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

resource "meraki_devices_cellular_gateway_lan" "test" {
	serial = "Q2GX-4F25-JFM4"
	device_name = "name of the MG"
	fixed_ip_assignments = [
		{
			mac = "0b:00:00:00:00:ac"
			name = "server 1"
			ip = "192.168.0.10"
		},
		{
			mac = "0b:00:00:00:00:ab"
			name = "server 2"
			ip = "192.168.0.20"
		}
	]
	reserved_ip_ranges = [
		{
			start = "192.168.1.0"
			end = "192.168.1.1"
			comment = "A reserved IP range"
		}
	]
}
`
