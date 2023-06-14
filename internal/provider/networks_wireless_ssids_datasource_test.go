package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TestAccNetworksWirelessSsidsDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksWirelessSsidsDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccNetworksWirelessSsidsDataSourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_wireless_ssids"),
				),
			},

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsDataSourceConfigCreateNetwork,
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

			{
				Config: testAccNetworksWirelessSsidsDataSourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("networks_wireless_ssids.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_networks_wireless_ssids.test", "list.#", "0"),
				),
			},
		},
	})
}

// testAccNetworksWirelessSsidsDataSourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccNetworksWirelessSsidsDataSourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_wireless_ssids"
 	api_enabled = true
 }
 `

// testAccNetworksWirelessSsidsDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
const testAccNetworksWirelessSsidsDataSourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
    depends_on = [resource.meraki_organization.test]
    organization_id = resource.meraki_organization.test.organization_id
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "Main Office"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`

// testAccNetworksWirelessSsidsDataSourceConfigCreate is a constant string that defines the configuration for creating and updating a networks__wireless_ssids resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsDataSourceConfigCreate = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
}

data "meraki_networks_wireless_ssids" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
  	network_id = resource.meraki_network.test.network_id
}
`
