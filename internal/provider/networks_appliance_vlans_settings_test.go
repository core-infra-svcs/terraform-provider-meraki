package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkApplianceVlansSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworkApplianceVlansSettingsResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_network_appliance_vlans_settings"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_appliance_vlans_settings.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworkApplianceVlansSettingsResourceConfigCreateNetwork,
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

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: testAccNetworkApplianceVlansSettingsResourceConfigUpdateNetworkApplianceVlansSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_appliance_vlans_settings.test", "vlans_enabled", "true"),
				),
			},
		},
	})
}

const testAccNetworkApplianceVlansSettingsResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_network_appliance_vlans_settings"
 	api_enabled = true
 } 
 `
const testAccNetworkApplianceVlansSettingsResourceConfigCreateNetwork = `
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

const testAccNetworkApplianceVlansSettingsResourceConfigUpdateNetworkApplianceVlansSettings = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_appliance_vlans_settings" "test" {
	  depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
      network_id = resource.meraki_network.test.network_id
	  vlans_enabled = true
	  
	  
	 
}
`
