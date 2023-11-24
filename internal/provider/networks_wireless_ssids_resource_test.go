package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksWirelessSsidsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccNetworksWirelessSsidsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_wireless_ssids_resource"),
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

			// Create and Read testing
			{
				Config: testAccNetworksWirelessSsidsResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "number", "0"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "name", "My SSID"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "auth_mode", "open"),
				),
			},

			// Available attributes:
			/*
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "enterprise_admin_access", "access enabled"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "encryption_mode", "wpa"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "psk", "YourSecureP@ssw0rd"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "wpa_encryption_mode", "WPA2 only"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dot11w.enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dot11w.required", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dot11r.enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dot11r.adaptive", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "splash_page", "Password-protected with custom RADIUS"),
				//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "splash_guest_sponsor_domains.0", "sponsor1.com"),
				//resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "splash_guest_sponsor_domains.1", "sponsor2.org"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "oauth.allowed_domains.0", "example.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "oauth.allowed_domains.1", "example.org"),

				// Local radius checks
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "local_radius.cache_timeout", "3600"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "local_radius.password_authentication.enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "local_radius.certificate_authentication.enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "local_radius.certificate_authentication.use_ldap", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "local_radius.certificate_authentication.use_ocsp", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "local_radius.certificate_authentication.ocsp_responder_url", "http://ocsp.example.com"),


						// LDAP checks
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.servers.#", "2"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.servers.0.host", "ldap.example.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.servers.0.port", "389"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.servers.1.host", "ldap2.example.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.servers.1.port", "389"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.credentials.distinguished_name", "cn=readonly,dc=example,dc=com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.credentials.password", "readonly_password"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.base_distinguished_name", "dc=example,dc=com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ldap.server_ca_certificate.contents", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"),

						// Active Directory checks
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.servers.#", "2"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.servers.0.host", "ad.example.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.servers.0.port", "389"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.servers.1.host", "backup-ad.example.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.servers.1.port", "389"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.credentials.logon_name", "ldap_user@example.com"),
						resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "active_directory.credentials.password", "ldap_user_password"),


					// Radius Servers checks
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_servers.#", "2"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_servers.0.host", "server1.example.com"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_servers.0.port", "1812"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_servers.1.host", "server2.example.com"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_servers.1.port", "1812"),


				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.#", "2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.host", "server1.example.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.port", "1813"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.host", "server2.example.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.port", "1813"),

				// Radius configuration checks
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_proxy_enabled", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_testing_enabled", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_called_station_id", "station-id"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_authentication_nas_id", "0"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_server_timeout", "5"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_server_attempts_limit", "3"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_fallback_enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_coa_enabled", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_fail_over_policy", "Deny access"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_load_balancing_policy", "Strict priority order"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_enabled", "true"),

				// Radius Accounting Servers configuration checks
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.#", "2"),

				// Check the first Radius Accounting Server
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.host", "server1.example.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.port", "1813"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.ca_certificate", ""),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.id", "1111111111111111111"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.0.open_roaming_certificate_id", ""),

				// Check the second Radius Accounting Server
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.host", "server2.example.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.port", "1813"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.ca_certificate", ""),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.id", "2222222222222222222"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_servers.1.open_roaming_certificate_id", ""),

				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_accounting_interim_interval", "600"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_attribute_for_group_policies", "Reply-Message"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ip_assignment_mode", "Bridge mode"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "use_vlan_tagging", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "concentrator_network_id", "N_1234567890abcdef"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "secondary_concentrator_network_id", "N_abcdef1234567890"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "disassociate_clients_on_vpn_fail_over", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "vlan_id", "3"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "default_vlan_id", "1"),

				// Check for every item in ap_tags_and_vlan_ids
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ap_tags_and_vlan_ids.#", "1"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ap_tags_and_vlan_ids.0.tags.#", "2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ap_tags_and_vlan_ids.0.tags.0", "tag1"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ap_tags_and_vlan_ids.0.tags.1", "tag2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "ap_tags_and_vlan_ids.0.vlan_id", "4"),

				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "walled_garden_enabled", "true"),

				// Check for every item in walled_garden_ranges
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "walled_garden_ranges.#", "2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "walled_garden_ranges.0", "walled-garden.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "walled_garden_ranges.1", "www.walled-garden.edu"),

				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "gre.concentrator.host", "server1.example.com"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "gre.concentrator.key", "5678"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_override", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_guest_vlan_enabled", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "radius_guest_vlan_id", "2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "min_bit_rate", "12"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "band_selection", "Dual band operation with Band Steering"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_client_bandwidth_limit_down", "5120"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_client_bandwidth_limit_up", "5120"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_ssid_bandwidth_limit_down", "56250"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "per_ssid_bandwidth_limit_up", "56250"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "lan_isolation_enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "visible", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "available_on_all_aps", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "availability_tags.#", "0"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "mandatory_dhcp_enabled", "false"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "adult_content_filtering_enabled", "false"),

					// DNSRewrite checks
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dns_rewrite.enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dns_rewrite.dns_custom_nameservers.#", "2"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dns_rewrite.dns_custom_nameservers.0", "server2.example.com"),
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "dns_rewrite.dns_custom_nameservers.1", "server1.example.com"),


				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "speed_burst.enabled", "true"),

				// Named Vlans checks
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.default_vlan_name", "vlan-name"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.by_ap_tags.#", "1"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.by_ap_tags.0.tags.#", "2"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.by_ap_tags.0.tags.0", "tag3"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.by_ap_tags.0.tags.1", "tag4"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.tagging.by_ap_tags.0.vlan_name", "vlan-name-1"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.radius.guest_vlan.enabled", "true"),
				resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "named_vlans.radius.guest_vlan.name", "guest-vlan"),

			*/

			// {
			//				ResourceName:      "meraki_networks_wireless_ssids.test",
			//				ImportState:       true,
			//				ImportStateVerify: false,
			//				ImportStateId:     "1234567890, 0987654321",
			//			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccNetworksWirelessSsidsResourceConfigCreateNetwork is a function which returns a string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_wireless_ssids_resource"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworksWirelessSsidsResourceConfigCreate = `

resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_networks_wireless_ssids" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    number = 0
	name = "My SSID"
	enabled = true
	auth_mode = "open"
}

`

/*

//enterprise_admin_access = "access enabled"  # Possible values could be 'access enabled' or 'access disabled'
	encryption_mode = "wpa"                    # Possible values: 'wep', 'wpa'
	//psk = "YourSecureP@ssw0rd"
	//wpa_encryption_mode = "WPA2 only"           # Possible values: 'WPA1 only', 'WPA1 and WPA2', 'WPA2 only', etc.

	ip_assignment_mode = "Bridge mode"
	use_vlan_tagging = true
	//concentrator_network_id = "N_1234567890abcdef"
	//secondary_concentrator_network_id = "N_abcdef1234567890"
	disassociate_clients_on_vpn_fail_over = false
	vlan_id = 3
	default_vlan_id = 1
	ap_tags_and_vlan_ids = [
        {
            tags = ["tag1", "tag2"]
            vlan_id = 4
        }
    ]
	walled_garden_enabled = true
    walled_garden_ranges = [
        "walled-garden.com",
        "www.walled-garden.edu"
    ]
	gre = {
        concentrator = {
            host = "server1.example.com"
        }
		key = 5678
    }

	min_bit_rate = 12
	band_selection = "Dual band operation with Band Steering"
	per_client_bandwidth_limit_down = 5120
    per_client_bandwidth_limit_up = 5120
    per_ssid_bandwidth_limit_down = 56250
    per_ssid_bandwidth_limit_up = 56250
	lan_isolation_enabled = true
	visible = true
	available_on_all_aps = true
	availability_tags = []
	mandatory_dhcp_enabled = false
	adult_content_filtering_enabled = false

	speed_burst = {
        enabled = true
    }
	named_vlans = {
        tagging = {
            enabled = true
            default_vlan_name = "vlan-name"
            by_ap_tags = [
                {
                    tags = ["tag3", "tag4"]
                    vlan_name = "vlan-name-1"
                }
            ]
        }

    }

dot11r = {
		enabled = true
		adaptive = true
	}

dot11w = {
		enabled = true
		required = false
	}

radius = {
            guest_vlan = {
                enabled = true
                name = "guest-vlan"
            }
        }

radius_override = true
	radius_guest_vlan_enabled = false
	radius_guest_vlan_id = 2

splash_page = "Password-protected with custom RADIUS"
		splash_guest_sponsor_domains = [
		"sponsor1.com",
		"sponsor2.org"
	]


	local_radius = {
		cache_timeout = 3600  # Timeout in seconds
		password_authentication = {
			enabled = true
		}
		certificate_authentication = {
			enabled = true
			use_ldap = false
			use_ocsp = false
			ocsp_responder_url = "http://ocsp.example.com"
			client_root_ca_certificate = {
				contents = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
			}
		}
	}



	radius_proxy_enabled = false
	radius_testing_enabled = false
	radius_called_station_id = "station-id"
    //radius_authentication_nas_id = "0"
	radius_server_timeout = 5
    radius_server_attempts_limit = 3
	radius_fallback_enabled = true
    radius_coa_enabled = false
	radius_fail_over_policy = "Deny access"
    radius_load_balancing_policy = "Strict priority order"
	radius_accounting_enabled = true

	radius_accounting_interim_interval = 600
	radius_attribute_for_group_policies = "Reply-Message"

radius_accounting_servers = [
        {
            ca_certificate = null
            host = "server1.example.com"
            id = "1111111111111111111"
            //open_roaming_certificate_id = null
            port = 1813
        },
        {
            ca_certificate = null
            host = "server2.example.com"
            id = "2222222222222222222"
            //open_roaming_certificate_id = null
            port = 1813
        }
    ]

	radius_servers = [
			{
				ca_certificate = null
				host = "server1.example.com"
				id = "1111111111111111111"
				//open_roaming_certificate_id = 2
				port = 1812
			},
			{
				ca_certificate = null
				host = "server2.example.com"
				id = "2222222222222222222"
				//open_roaming_certificate_id = 2
				port = 1812
			}
		]



	active_directory = {
		servers = [
			{
				host = "ad.example.com"
				port = 389
			},
			{
				host = "backup-ad.example.com"
				port = 389
			}
		]
		credentials = {
			logon_name = "ldap_user@example.com"
			password = "ldap_user_password"
		}
	}

	oauth = {
		allowed_domains = ["example.com", "example.org"]
	}


	ldap = {
		servers = [
			{
				host = "ldap.example.com"
				port = 389
			},
			{
				host = "ldap2.example.com"
				port = 389
			}
		]
		credentials = {
			distinguished_name = "cn=readonly,dc=example,dc=com"
			password = "readonly_password"
		}
		base_distinguished_name = "dc=example,dc=com"
		server_ca_certificate = {
			contents = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
		}
	}

	dns_rewrite = {
        enabled = true
        dns_custom_nameservers = ["server2.example.com", "server1.example.com"]
    }

*/
