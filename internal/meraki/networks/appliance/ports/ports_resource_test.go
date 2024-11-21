package ports_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksAppliancePortResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_ports"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_ports"),
			},

			// Claim Appliance To Network
			{
				Config: NetworksAppliancePortResourceClaimNetworkDevice(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
			},

			// Update and Read Networks Appliance Vlans Settings.
			{
				Config: NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings(),
				Check:  NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettingsChecks(),
			},

			//  Configure Appliance Port
			{
				Config: NetworksAppliancePortResourceConfigListNetworkAppliancePorts(),
				Check:  NetworksAppliancePortResourceConfigListNetworkAppliancePortsChecks(),
			},
		},
	})
}

func NetworksAppliancePortResourceClaimNetworkDevice(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = ["resource.meraki_network.test"]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_ports"),
		serial,
	)
}

func NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}
resource "meraki_networks_appliance_vlans_settings" "test" {
		depends_on = ["resource.meraki_network.test", "resource.meraki_networks_devices_claim.test"]
		network_id = resource.meraki_network.test.network_id
		vlans_enabled = true
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_ports"),
	)
}

// NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettingsChecks returns the test check functions for NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettings
func NetworksAppliancePortResourceConfigUpdateNetworkApplianceVlansSettingsChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"vlans_enabled": "true",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_vlans_settings.test", expectedAttrs)
}

func NetworksAppliancePortResourceConfigListNetworkAppliancePorts() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
	network_id = resource.meraki_network.test.network_id
}
resource "meraki_networks_appliance_vlans_settings" "test" {
	network_id = resource.meraki_network.test.network_id
	vlans_enabled = true
}
resource "meraki_networks_appliance_ports" "test" {
	depends_on = ["resource.meraki_network.test", "resource.meraki_networks_devices_claim.test", "resource.meraki_networks_appliance_vlans_settings.test"]
	network_id = resource.meraki_network.test.network_id
	port_id = 4
	allowed_vlans = "all"
	drop_untagged_traffic = true
	enabled = true	
	type = "trunk"	
    }
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "nw_name"),
	)
}

// NetworksAppliancePortResourceConfigListNetworkAppliancePortsChecks returns the test check functions for NetworksAppliancePortResourceConfigListNetworkAppliancePorts
func NetworksAppliancePortResourceConfigListNetworkAppliancePortsChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"port_id":               "4",
		"allowed_vlans":         "all",
		"drop_untagged_traffic": "true",
		"enabled":               "true",
		"type":                  "trunk",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_ports.test", expectedAttrs)
}
