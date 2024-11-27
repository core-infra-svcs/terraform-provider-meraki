package vpn_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceVpnSiteToSiteVpnResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			//// Create and Read Network.
			//{
			//	Config: testAccNetworksApplianceVpnSiteToSiteVpnResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_appliance_vpn_site_to_site_vpn"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
			//	),
			//},

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_appliance_vpn_site_to_site_vpn"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_network_appliance_vpn_site_to_site_vpn"),
			},

			// Claim Network Device
			{
				Config: ApplianceVpnSiteToSiteVpnResourceConfigClaimNetworksDevice(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},

			// Update and Read Networks Appliance Vpn Site To Site Vpn.
			{
				Config: ApplianceVpnSiteToSiteVpnResourceConfigUpdate(),
				Check:  ApplianceVpnSiteToSiteVpnResourceConfigUpdateChecks(),
			},
		},
	})
}

func ApplianceVpnSiteToSiteVpnResourceConfigClaimNetworksDevice(serial string) string {
	return fmt.Sprintf(`
	%s

resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_appliance_vpn_site_to_site_vpn"),
		serial,
	)
}

func ApplianceVpnSiteToSiteVpnResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_vpn_site_to_site_vpn" "test" {
    depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
    mode = "hub"
    subnets = [{
		local_subnet = "192.168.128.0/24"
	}]
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_appliance_vpn_site_to_site_vpn"),
	)
}

// ApplianceVpnSiteToSiteVpnResourceConfigUpdateChecks returns the test check functions for ApplianceVpnSiteToSiteVpnResourceConfigUpdate
func ApplianceVpnSiteToSiteVpnResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"mode":                   "hub",
		"subnets.0.local_subnet": "192.168.128.0/24",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_vpn_site_to_site_vpn.test", expectedAttrs)
}
