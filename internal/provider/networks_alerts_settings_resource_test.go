package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNetworksAlertsSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksAlertsSettingsResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_meraki_networks_alerts_settings"),
				),
			},

			//Create test Network
			{
				Config: testAccNetworksAlertsSettingsResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "id", "example-id"),
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

			// Create and Read Testing
			{
				Config: testAccNetworksAlertsSettingsResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "default_destinations.emails.0", "meraki_organizations_admin_test2@example.com"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "default_destinations.snmp", "true"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "default_destinations.all_admins", "true"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "default_destinations.http_server_ids.0", "aHR0cHM6Ly93d3cuZXhhbXBsZS5jb20vd2ViaG9va3M="),

					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.type", "gatewayDown"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.alert_destinations.emails", "meraki_organizations_admin_test2@example.com"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.alert_destinations.snmp", "true"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.alert_destinations.all_admins", "true"),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.alert_destinations.all_admins", "aHR0cHM6Ly93d3cuZXhhbXBsZS5jb20vd2ViaG9va3M="),
					resource.TestCheckResourceAttr("meraki_networks_alerts_settings.test", "alerts.0.filters.timeout", "60"),
				),
			},
		},
	})
}

const testAccNetworksAlertsSettingsResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_meraki_networks_alerts_settings"
 	api_enabled = true
 }
`

const testAccNetworksAlertsSettingsResourceConfigCreateNetwork = `
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

const testAccNetworksAlertsSettingsResourceConfigCreate = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_alerts_settings" "test" {
	depends_on = [resource.meraki_organization.test, resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	default_destinations = {
		emails = ["miles@meraki.com"]
		snmp = true
		all_admins = true
		http_server_ids = []
	}
	alerts = [
		{
			type = "gatewayDown"
			enabled = true
			alert_destinations = {
				emails = ["miles@meraki.com"]
				snmp = true
				all_admins = true
				http_server_ids = []
			}
			filters = {
				timeout = 60
			}
		}
	]
}
`
