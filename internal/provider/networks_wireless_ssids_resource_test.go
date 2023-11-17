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

			// Create and Read testing
			{
				Config: testNetworkWirelessSSID(os.Getenv("TF_ACC_MERAKI_NETWORK_ID"), os.Getenv("TF_ACC_MERAKI_MX_LICENCE")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_wireless_ssids.test", "id", "example-id"),
				),
			},

			{
				ResourceName:      "meraki_networks_wireless_ssids.test",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "1234567890, 0987654321",
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testNetworkWirelessSSID(networkID string, number string) string {
	return fmt.Sprintf(testAccNetworksWirelessSsidsResourceConfigCreate, networkID, number)
}

const testAccNetworksWirelessSsidsResourceConfigCreate = `
resource "meraki_networks_wireless_ssids" "test" {
  network_id = "%s"
  number = "%s"
  ssid = {
    name = "My SSID"
    auth_mode = "8021x-radius"
    enterprise_admin_access = "access enabled"
    encryption_mode = "wpa"
    psk = "deadbeef"
    wpa_encryption_mode = "WPA2 only"
    splash_page = "Click-through splash page"
    radius_called_station_id = "00-11-22-33-44-55:AP1"
    radius_authentication_nas_id = "00-11-22-33-44-55:AP1"
    radius_failover_policy = "Deny access"
    radius_load_balancing_policy = "Round robin"
    radius_attribute_for_group_policies = "Filter-Id"
    ip_assignment_mode = "NAT mode"
    concentrator_network_id = "N_24329156"
    secondary_concentrator_network_id = "disabled"
    band_selection = "5 GHz band only"
    radius_server_timeout = 5
    radius_server_attempts_limit = 5
    radius_accounting_interim_interval = 5
    vlan_id = 10
    default_vlan_id = 1
    per_client_bandwidth_limit_up = 0
    per_client_bandwidth_limit_down = 0
    per_ssid_bandwidth_limit_up = 0
    per_ssid_bandwidth_limit_down = 0
    radius_guest_vlan_id = 1
    min_bitrate = 5.5
    use_vlan_tagging = false
    disassociate_clients_on_vpn_failover = false
    radius_override = false
    radius_guest_vlan_enabled = true
	enabled = true
	radius_proxy_enabled = false
    radius_testing_enabled = true
    radius_fallback_enabled = true
    radius_coa_enabled = true
    radius_accounting_enabled = true
    lan_isolation_enabled = true
    visible = true
    available_on_all_aps = false
    mandatory_dhcp_enabled = false
    adult_content_filtering_enabled = false
    walled_garden_enabled = true
    dot11w = {
      enabled = true
      required = false
    }
    dot11r = {
      enabled = true
      adaptive = true
    }
	local_radius = {
      cache_timeout = 60
      password_authentication = {
        enabled = false
      }
      certificate_authentication = {
        enabled = true
        use_ldap = false
        use_ocsp = true
        ocsp_responder_url = "http://ocsp-server.example.com"
        client_root_ca_certificate = {
          contents = "test"
        }
      }
    }
 	ldap = {
      servers = [
        {
          host = "127.0.0.1"
          port = 389
        }
      ]
      credentials = {
        distinguished_name = "cn=user,dc=example,dc=com"
        password = "password"
      }
      base_distinguished_name = "dc=example,dc=com"
      server_ca_certificate = {
          contents = "test"
      }
    }
    active_directory = {
      servers = [
        {
          host = "127.0.0.1"
          port = 3268
        }
      ]
      credentials = {
        logon_name = "user"
        password = "password"
      }
    }
	dns_rewrite = {
      enabled = true
      dns_custom_nameservers = ["8.8.8.8", "8.8.4.4"]
    }
    speed_burst = {
      enabled = true
    }
	gre = {
      concentrator = {
        host = "192.168.1.1"
      }
      key = 5
    }
	oauth = {
      allowed_domains = ["example.com"]
    }
	splash_guest_sponsor_domains = ["example.com"]
	walled_garden_ranges = ["example.com", "1.1.1.1/32"]
    availability_tags = ["tag1", "tag2"]
    radius_servers = [
      {
        host = "0.0.0.0"
        secret = "secret-string"
		ca_certificate = "test"
        port = 3000
        open_roaming_certificate_id = 2
        radsec_enabled = true
      }
    ]
    radius_accounting_servers = [
      {
        host = "0.0.0.0"
        secret = "secret-string"
		ca_certificate = "test"
        port = 3000
        radsec_enabled = true
      }
    ]
    ap_tags_and_vlan_ids = [
      {
        tags = ["tag1", "tag2"]
        vlan_id = 100
      }
    ]
  }
}`
