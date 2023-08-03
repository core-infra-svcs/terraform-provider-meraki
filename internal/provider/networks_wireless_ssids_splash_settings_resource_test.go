package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO - Requires additional licences to complete integration testing

// TestAccNetworksWirelessSsidsSplashSettingsResource function is used to test the CRUD operations of the Terraform resource developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksWirelessSsidsSplashSettingsResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read an Organization.
			{
				Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_networks_wireless_ssids_splash_settings"),
				),
			},
			/*
				// Create and Read a Network.
				{
					Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetwork,
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

				// Create and Read a SystemsManager Network.
				{
					Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetworkSystemsManager,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_network.test_systems_manager", "name", "SM"),
						resource.TestCheckResourceAttr("meraki_network.test_systems_manager", "timezone", "America/Los_Angeles"),
						resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
						resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
						resource.TestCheckResourceAttr("meraki_network.test_systems_manager", "product_types.#", "1"),
						resource.TestCheckResourceAttr("meraki_network.test_systems_manager", "product_types.0", "systemsManager"),
						resource.TestCheckResourceAttr("meraki_network.test_systems_manager", "notes", "Additional description of the network"),
					),
				},

				// TODO: Create and Read NetworksWirelessSsidsSplashSettings
				{
					Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreate,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "id", "example-id"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_url", "https://www.custom_splash_url.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "use_splash_url", "false"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_timeout", "1440"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "welcome_message", "Welcome!"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "redirect_url", "https://example.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "use_redirect_url", "false"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "block_all_traffic_before_sign_on", "false"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "controller_disconnection_behavior", "default"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "allow_simultaneous_logins", "false"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.reply_to_email_address", "user@email.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.prepaid_access_fast_login_enabled", "false"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.free_access.enabled", "true"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.free_access.duration_in_minutes", "120"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "guest_sponsorship.guest_can_request_time_frame", "false"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "guest_sponsorship.duration_in_minutes", "120"),
					),
				},
			*/

			/*
				// TODO: Update and Read NetworksWirelessSsidsSplashSettings
					{
						Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigUpdate,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "id", "example-id"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_url", "https://www.updatedcustom_splash_url.com"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "use_splash_url", "false"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_timeout", "1450"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "welcome_message", "Welcome hii!"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "redirect_url", "https://updatedexample.com"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "use_redirect_url", "false"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "block_all_traffic_before_sign_on", "true"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "controller_disconnection_behavior", "open"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "allow_simultaneous_logins", "true"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.reply_to_email_address", "updateduser@email.com"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.prepaid_access_fast_login_enabled", "false"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.free_access.enabled", "false"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "billing.free_access.duration_in_minutes", "60"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "guest_sponsorship.guest_can_request_time_frame", "true"),
							resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "guest_sponsorship.duration_in_minutes", "60"),
						),
					},
			*/
		},
		// TODO: Optionally, you can add an ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_networks_wireless_ssids_splash_settings.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890, 0987654321",
		   },
		*/

		// The resource.Test function automatically tests the Delete operation.
	})
}

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateOrganization is a constant string that defines the configuration for creating an organization resource in your tests.
const testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_networks_wireless_ssids_splash_settings"
 	api_enabled = true
 }
 `

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetworkSystemsManager is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
const testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetworkSystemsManager = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test_systems_manager" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["systemsManager"]
	name = "SM"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
const testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetwork = `
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

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreate is a constant string that defines the configuration for creating and reading a networks_wireless_ssids_splash_settings resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
}


resource "meraki_networks_wireless_ssids_splash_settings" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_network.test_systems_manager,resource.meraki_organization.test]
  	network_id = resource.meraki_network.test.network_id
    number = "0"
	splash_url = "https://example.com"
    use_splash_url = false
    splash_timeout = 1440
    redirect_url = "https://example.com"
    use_redirect_url = false
    welcome_message = "Welcome!"
    splash_logo = {
		image = {}
	}
    splash_image = {
		image = {}
	}
    splash_prepaid_front = {
		image = {}
	}

    block_all_traffic_before_sign_on = false
    controller_disconnection_behavior = "default"
    allow_simultaneous_logins = false
    guest_sponsorship = {
        duration_in_minutes = 30
        guest_can_request_time_frame = false
    }
    billing = {
        free_access = {
            enabled = true
            duration_in_minutes = 120
        }
        prepaid_access_fast_login_enabled = true
        reply_to_email_address = "user@email.com"
    }
    sentry_enrollment = {
        systems_manager_network = {
			id = resource.meraki_organization.test_systems_manager.organization_id
		}
        strength = "focused"
        enforced_systems = [
            "iOS"
        ]
    }
}
`

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigUpdate is a constant string that defines the configuration for updating a networks_wireless_ssids_splash_settings resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsSplashSettingsResourceConfigUpdate = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless", "systemsManager"]
}

resource "meraki_network" "test_systems_manager" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["systemsManager"]
}

resource "meraki_networks_wireless_ssids_splash_settings" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_network.test_systems_manager, resource.meraki_organization.test]
  	network_id = resource.meraki_network.test.network_id
    number = "0"
	use_splash_url = false
	splash_timeout = 1450
	welcome_message = "Welcome hii!"
	redirect_url = "https://updatedexample.com"
	use_redirect_url = false
	block_all_traffic_before_sign_on = true
	controller_disconnection_behavior = "open"
	allow_simultaneous_logins = true
	billing = {
		prepaid_access_fast_login_enabled = false
        reply_to_email_address = "updateduser@email.com"
		free_access = {
			enabled = false
            duration_in_minutes = 60
		}
	}
	guest_sponsorship = {
        duration_in_minutes = 60
        guest_can_request_time_frame = true
    }

	splash_image = {}
	splash_logo = {}
	splash_prepaid_front = {} 
	sentry_enrollment =  {
        systems_manager_network = {
		} 
		enforced_systems = ["iOS"]
		strength = "focused"
    }
}
`
