package vlan_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceVlansDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_meraki_networks_appliance_vlans"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_meraki_networks_appliance_vlans"),
			},

			// Create and Read a VLAN
			{
				Config: NetworksApplianceVlansDataSourceConfigCreate(),
				Check:  NetworksApplianceVlansDataSourceConfigCreateChecks(),
			},

			// Read Ports
			{
				Config: NetworksApplianceVlansDataSourceConfigRead(),
				Check:  NetworksApplianceVlansDataSourceConfigReadChecks(),
			},
		},
	})
}

func NetworksApplianceVlansDataSourceConfigCreate() string {
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
    subnet = "192.168.2.0/24"
    appliance_ip = "192.168.2.2"
    dhcp_handling = "Run a DHCP server"
    dhcp_lease_time = "12 hours"
    dhcp_boot_options_enabled = true
    dhcp_boot_next_server = "2.3.4.5"
    dhcp_boot_filename = "updated.file"
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
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_meraki_networks_appliance_vlans"),
	)
}

// NetworksApplianceVlansDataSourceConfigCreateChecks returns the test check functions for NetworksApplianceVlansDataSourceConfigCreate
func NetworksApplianceVlansDataSourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlan_id":                      "10",
		"name":                         "My VLAN",
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

func NetworksApplianceVlansDataSourceConfigRead() string {
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
    subnet = "192.168.2.0/24"
    appliance_ip = "192.168.2.2"
    dhcp_handling = "Run a DHCP server"
    dhcp_lease_time = "12 hours"
    dhcp_boot_options_enabled = true
    dhcp_boot_next_server = "2.3.4.5"
    dhcp_boot_filename = "updated.file"
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

data "meraki_networks_appliance_vlans" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_meraki_networks_appliance_vlans"),
	)
}

// NetworksApplianceVlansDataSourceConfigReadChecks returns the test check functions for NetworksApplianceVlansDataSourceConfigRead
func NetworksApplianceVlansDataSourceConfigReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.0.vlan_id":                   "1",
		"list.0.name":                      "Default",
		"list.0.appliance_ip":              "192.168.128.1",
		"list.0.subnet":                    "192.168.128.0/24",
		"list.0.dns_nameservers":           "upstream_dns",
		"list.0.dhcp_handling":             "Run a DHCP server",
		"list.0.dhcp_lease_time":           "1 day",
		"list.0.dhcp_boot_options_enabled": "false",

		"list.1.vlan_id":                      "10",
		"list.1.name":                         "My VLAN",
		"list.1.subnet":                       "192.168.2.0/24",
		"list.1.appliance_ip":                 "192.168.2.2",
		"list.1.dhcp_handling":                "Run a DHCP server",
		"list.1.dhcp_lease_time":              "12 hours",
		"list.1.dhcp_boot_options_enabled":    "true",
		"list.1.dhcp_boot_next_server":        "2.3.4.5",
		"list.1.dhcp_boot_filename":           "updated.file",
		"list.1.reserved_ip_ranges.0.start":   "192.168.2.0",
		"list.1.reserved_ip_ranges.0.end":     "192.168.2.1",
		"list.1.reserved_ip_ranges.0.comment": "A newly reserved IP range",
		"list.1.dns_nameservers":              "upstream_dns",
		"list.1.dhcp_options.0.code":          "6",
		"list.1.dhcp_options.0.type":          "text",
		"list.1.dhcp_options.0.value":         "six",
	}
	return utils.ResourceTestCheck("data.meraki_networks_appliance_vlans.test", expectedAttrs)
}
