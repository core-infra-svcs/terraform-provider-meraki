package vlans_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksApplianceVlansResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_meraki_networks_appliance_vlan"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_meraki_networks_appliance_vlan"),
			},

			// Create and Read a VLAN
			{
				Config: NetworksApplianceVlansResourceConfigCreate(),
				Check:  NetworksApplianceVlansResourceConfigCreateChecks(),
			},

			// Update testing
			{
				Config: NetworksApplianceVlansResourceConfigUpdate(),
				Check:  NetworksApplianceVlansResourceConfigUpdateChecks(),
			},

			// Import testing
			{
				ResourceName:      "meraki_networks_appliance_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func NetworksApplianceVlansResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_vlans_settings" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlans_enabled = true
}

resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_networks_appliance_vlans_settings.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "10"
    name = "My VLAN"
    subnet = "192.168.1.0/24"
	appliance_ip = "192.168.1.2"
	cidr = "192.168.1.0/24"
	mask = 28
    //vpn_nat_subnet = "192.168.1.0/24"
	dhcp_relay_server_ips = ["192.168.1.0/24", "192.168.128.0/24"]
    mandatory_dhcp = {
        enabled = true
    }

}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_meraki_networks_appliance_vlan"),
	)
}

// NetworksApplianceVlansResourceConfigCreateChecks returns the test check functions for NetworksApplianceVlansResourceConfigCreate
func NetworksApplianceVlansResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlan_id":                 "10",
		"name":                    "My VLAN",
		"subnet":                  "192.168.1.0/24",
		"appliance_ip":            "192.168.1.2",
		"dhcp_relay_server_ips.0": "192.168.1.0/24",
		"dhcp_relay_server_ips.1": "192.168.128.0/24",
		"mandatory_dhcp.enabled":  "true",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_vlan.test", expectedAttrs)
}

func NetworksApplianceVlansResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_vlan" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    vlan_id = "10"
    name = "My Updated VLAN"
    subnet = "192.168.2.0/24"
    appliance_ip = "192.168.2.2"
    dhcp_handling = "Run a DHCP server"
    dhcp_lease_time = "12 hours"
    dhcp_boot_options_enabled = true
    dhcp_boot_next_server = "2.3.4.5"
    dhcp_boot_filename = "updated.file"
	fixed_ip_assignments = {
		"22:33:44:55:66:77": {
			"ip": "192.168.2.10",
			"name": "Some client name"
	  	}
	}
    reserved_ip_ranges = [
        {
            start = "192.168.2.0"
            end = "192.168.2.1"
            comment = "A newly reserved IP range"
        }
    ]
    dns_nameservers = "upstream_dns"
    dhcp_options = [
        {
            code = "6"
            type = "text"
            value = "six"
        }
    ]
    mandatory_dhcp = {
        enabled = true
    }
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_meraki_networks_appliance_vlan"),
	)
}

// NetworksApplianceVlansResourceConfigUpdateChecks returns the test check functions for NetworksApplianceVlansResourceConfigUpdate
func NetworksApplianceVlansResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlan_id":                      "10",
		"name":                         "My Updated VLAN",
		"subnet":                       "192.168.2.0/24",
		"appliance_ip":                 "192.168.2.2",
		"dhcp_handling":                "Run a DHCP server",
		"dhcp_lease_time":              "12 hours",
		"dhcp_boot_options_enabled":    "true",
		"dhcp_boot_next_server":        "2.3.4.5",
		"dhcp_boot_filename":           "updated.file",
		"reserved_ip_ranges.0.start":   "192.168.2.0",
		"reserved_ip_ranges.0.end":     "192.168.2.1",
		"reserved_ip_ranges.0.comment": "A newly reserved IP range",
		"dns_nameservers":              "upstream_dns",
		"dhcp_options.0.code":          "6",
		"dhcp_options.0.type":          "text",
		"dhcp_options.0.value":         "six",
		"mandatory_dhcp.enabled":       "true",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_vlan.test", expectedAttrs)
}
