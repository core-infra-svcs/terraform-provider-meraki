package vpn_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceVpnSiteToSiteVpnDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

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

			{
				// Test case for reading Networks Appliance Vpn Site To Site Vpn.
				Config: ApplianceVpnSiteToSiteVpnDatasourceConfigRead(),
				Check:  ApplianceVpnSiteToSiteVpnDatasourceConfigReadChecks(),
			},
		},
	})
}

func ApplianceVpnSiteToSiteVpnDatasourceConfigRead() string {
	return fmt.Sprintf(`
	%s

resource "meraki_networks_appliance_vpn_site_to_site_vpn" "test" {
    depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
    mode = "hub"
    subnets = [
		{
			local_subnet = "192.168.128.0/24"
		}
	]	
}

data "meraki_networks_appliance_vpn_site_to_site_vpn" "test" {
	depends_on = [resource.meraki_network.test, meraki_networks_appliance_vpn_site_to_site_vpn.test]
	network_id = resource.meraki_network.test.network_id
	mode = resource.meraki_networks_appliance_vpn_site_to_site_vpn.test.mode
	id = resource.meraki_networks_appliance_vpn_site_to_site_vpn.test.id
	subnets = resource.meraki_networks_appliance_vpn_site_to_site_vpn.test.subnets
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_network_appliance_vpn_site_to_site_vpn"),
	)
}

// ApplianceVpnSiteToSiteVpnDatasourceConfigReadChecks returns the test check functions for ApplianceVpnSiteToSiteVpnDatasourceConfigRead
func ApplianceVpnSiteToSiteVpnDatasourceConfigReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"mode":                   "hub",
		"subnets.0.local_subnet": "192.168.128.0/24",
	}
	return utils.ResourceTestCheck("data.meraki_networks_appliance_vpn_site_to_site_vpn.test", expectedAttrs)
}
