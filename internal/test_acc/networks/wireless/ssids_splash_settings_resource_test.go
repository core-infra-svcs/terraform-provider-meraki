package wireless

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksWirelessSsidsSplashSettingsResource function is used to test the CRUD operations of the Terraform resource developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksWirelessSsidsSplashSettingsResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { test_acc.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: test_acc.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_wireless_ssids_splash_settings"),
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
				Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetworkSystemsManager(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.testhub", "name", "test_acc_networks_wireless_ssids_splash_settings_hub"),
					resource.TestCheckResourceAttr("meraki_network.testhub", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.testhub", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.testhub", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.testhub", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.testhub", "product_types.0", "systemsManager"),
					resource.TestCheckResourceAttr("meraki_network.testhub", "notes", "Additional description of the network"),
				),
			},

			// Create and Read NetworksWirelessSsidsSplashSettings
			{
				Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
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
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_image.extension", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_image.md5", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_logo.extension", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_logo.md5", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_prepaid_front.extension", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_prepaid_front.md5", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.systems_manager_network.id", os.Getenv("TF_ACC_MERAKI_NETWORK_ID")),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.enforced_systems.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.enforced_systems.0", "iOS"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.strength", "focused"),
				),
			},

			// Update and Read NetworksWirelessSsidsSplashSettings
			{
				Config: testAccNetworksWirelessSsidsSplashSettingsResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_url", "https://www.updatedcustom_splash_url.com"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "use_splash_url", "false"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_timeout", "1440"),
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
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_image.extension", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_image.md5", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_logo.extension", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_logo.md5", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_prepaid_front.extension", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "splash_prepaid_front.md5", ""),
					//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.systems_manager_network.id", os.Getenv("TF_ACC_MERAKI_NETWORK_ID")),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.enforced_systems.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.enforced_systems.0", "iOS"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids_splash_settings.test", "sentry_enrollment.strength", "focused"),
				),
			},
			// Import State testing
			{
				ResourceName:      "meraki_networks_wireless_ssids_splash_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"splash_image",
					"splash_logo",
					"splash_prepaid_front",
				},
			},
		},

		// The resource.Test function automatically tests the Delete operation.
	})
}

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" { 
	organization_id = "%s" 
	product_types = ["appliance", "switch", "wireless"] 
	tags = ["tag1"] 
	name = "test_acc_networks_wireless_ssids_splash_settings" 
	timezone = "America/Los_Angeles" 
	notes = "Additional description of the network" 
}
`, orgId)
	return result
}

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetworkSystemsManager is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreateNetworkSystemsManager(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "testhub" {
	organization_id = "%s"
	product_types = ["systemsManager"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_splash_settings_hub"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

`, orgId)
	return result
}

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreate is a constant string that defines the configuration for creating and reading a networks_wireless_ssids_splash_settings resource in your tests.
// It depends on both the organization and network resources.
func testAccNetworksWirelessSsidsSplashSettingsResourceConfigCreate(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_splash_settings"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_network" "testhub" {
	organization_id = "%s"
	product_types = ["systemsManager"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_splash_settings_hub"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_wireless_ssids_splash_settings" "test" {
	depends_on = [resource.meraki_network.test, resource.meraki_network.testhub]
	network_id = resource.meraki_network.test.network_id
	number = "0"
	splash_url = "https://www.custom_splash_url.com"
	use_splash_url = false
	splash_timeout = 1440
	welcome_message = "Welcome!"
	redirect_url = "https://example.com"
	use_redirect_url = false
	block_all_traffic_before_sign_on = false
	controller_disconnection_behavior = "default"
	allow_simultaneous_logins = false
	billing = {
		prepaid_access_fast_login_enabled = false
		reply_to_email_address = "user@email.com"
		free_access = {
			enabled = true
			duration_in_minutes = 120
		}
	}

	guest_sponsorship = {
		duration_in_minutes = 120
		guest_can_request_time_frame = false
	}

	splash_image = {
		image = {
			contents = "Q2lzY28gTWVyYWtp"
			format = "jpg"
		}
	}

	splash_logo = {
		image = {
			contents = "Q2lzY28gTWVyYWtp"
			format = "jpg"
		}
	}

	splash_prepaid_front = {
		image = {
			contents = "Q2lzY28gTWVyYWtp"
			format = "jpg"
		}
	}

	sentry_enrollment =  {
		strength = "click-through"
		systems_manager_network = {
			id = resource.meraki_network.testhub.network_id
		}
	
		enforced_systems = ["iOS"]
		strength = "focused"
	}
}
`, orgId, orgId)
	return result
}

// testAccNetworksWirelessSsidsSplashSettingsResourceConfigUpdate is a constant string that defines the configuration for updating a networks_wireless_ssids_splash_settings resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsSplashSettingsResourceConfigUpdate = `
resource "meraki_network" "test" {
product_types = ["appliance", "switch", "wireless"]
tags = ["tag1"]
}

resource "meraki_network" "testhub" {
product_types = ["systemsManager"]
tags = ["tag1"]
}

resource "meraki_networks_wireless_ssids_splash_settings" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	number = "0"
	splash_url = "https://www.updatedcustom_splash_url.com"
	use_splash_url = false
	splash_timeout = 1440
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
	splash_image = {
		image = {
			contents = "Q2lzY28gTWVyYWtp"
			format = "jpg"
		}
	}
	splash_logo = {
		image = {
			contents = "Q2lzY28gTWVyYWtp"
			format = "jpg"
		}
	}
	splash_prepaid_front = {
		image = {
			contents = "Q2lzY28gTWVyYWtp"
			format = "jpg"
		}
	}
	sentry_enrollment =  {
		strength = "click-through"
		systems_manager_network = {
			id = resource.meraki_network.testhub.network_id
		}
		enforced_systems = ["iOS"]
		strength = "focused"
	}
}
`
