package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: debug device unclaiming issue.
/*
func TestAccOrganizationsCellularGatewayUplinkStatusesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read a Network
			{
				Config: testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check:  utils.NetworkTestChecks("test_acc_organizations_cellular_gateway_uplink_statuses"),
			},

			// Claim and Read NetworksDevicesClaim
			{
				Config: testAccOrganizationsCellularGatewayUplinkStatusesNetworksDevicesClaimResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},

			// Read OrganizationsCellularGatewayUplinkStatuses
			{
				Config: testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check:  OrganizationsCellularGatewayUplinkStatusesDataSourceTestChecks(os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
			},

			// Unclaim Device to allow Organization deletion
			{
				Config: testAccOrganizationsUnclaimDevice(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
*/

// Unclaim device from organization
func testAccOrganizationsUnclaimDevice(orgId string, serial string) string {
	return fmt.Sprintf(`
	resource "meraki_organizations_devices_unclaim" "test" {
		organization_id = "%s"
		serials = ["%s"]
	}
	`, orgId, serial)
}

// testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigCreateNetwork returns the configuration for creating a network resource
func testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
	%s
	`, utils.CreateNetworkConfig(
		"test_acc_meraki_organizations_cellular_gateway_uplink_statuses",
		"test_acc_organizations_cellular_gateway_uplink_statuses",
	))
}

// testAccOrganizationsCellularGatewayUplinkStatusesNetworksDevicesClaimResourceConfigCreate returns the configuration for creating and reading a networks_devices_claim resource
func testAccOrganizationsCellularGatewayUplinkStatusesNetworksDevicesClaimResourceConfigCreate(orgId string, serial string) string {
	return fmt.Sprintf(`
	%s
	resource "meraki_networks_devices_claim" "test" {
		depends_on = [resource.meraki_network.test]
		network_id = resource.meraki_network.test.network_id
		serials = ["%s"]
	}
	`,
		utils.CreateNetworkConfig("test_acc_meraki_organizations_cellular_gateway_uplink_statuses", "test_acc_organizations_cellular_gateway_uplink_statuses"),
		serial,
	)
}

// testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigRead returns the configuration for reading cellular gateway uplink statuses
func testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigRead(orgId string, serial string) string {
	return fmt.Sprintf(`
	%s
	resource "meraki_networks_devices_claim" "test" {
		depends_on = [resource.meraki_network.test]
		network_id = resource.meraki_network.test.network_id
		serials = ["%s"]
	}

	data "meraki_organizations_cellular_gateway_uplink_statuses" "test" {
		organization_id = "%s"
		serials = ["%s"]
	}
	`,
		utils.CreateNetworkConfig("test_acc_meraki_organizations_cellular_gateway_uplink_statuses", "test_acc_organizations_cellular_gateway_uplink_statuses"),
		serial,
		orgId,
		serial,
	)
}

// OrganizationsCellularGatewayUplinkStatusesDataSourceTestChecks returns the test check functions for verifying cellular gateway uplink statuses data source
func OrganizationsCellularGatewayUplinkStatusesDataSourceTestChecks(serial string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.#":        "1",
		"list.0.serial": serial,
		"list.0.model":  "MG21-NA",
	}

	return utils.ResourceTestCheck("data.meraki_organizations_cellular_gateway_uplink_statuses.test", expectedAttrs)
}
