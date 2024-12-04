package subnets_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesApplianceDhcpSubnetsDataSource(t *testing.T) {
	t.Run("Create and Read Network", func(t *testing.T) {
		testCreateAndReadNetwork(t)
	})

	t.Run("Claim and Read Network Devices", func(t *testing.T) {
		testClaimAndReadNetworkDevices(t)
	})

	t.Run("Update Network VLAN Settings", func(t *testing.T) {
		testUpdateNetworkVLANSettings(t)
	})

	t.Run("Update and Read DevicesApplianceDhcpSubnets", func(t *testing.T) {
		testUpdateAndReadDevicesApplianceDhcpSubnets(t)
	})
}

func testCreateAndReadNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccDevicesApplianceDhcpSubnetsPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_devices_appliance_dhcp_subnets"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_devices_appliance_dhcp_subnets"),
			},
		},
	})
}

func testClaimAndReadNetworkDevices(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccDevicesApplianceDhcpSubnetsPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigClaim(
					os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func testUpdateNetworkVLANSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccDevicesApplianceDhcpSubnetsPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: DevicesApplianceDhcpSubnetsDataSourceConfigVlanSettings(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_devices_appliance_dhcp_subnets"),
				),
			},
		},
	})
}

func testUpdateAndReadDevicesApplianceDhcpSubnets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccDevicesApplianceDhcpSubnetsPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "resources.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "resources.0.subnet", "192.168.128.0/24"),
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "resources.0.vlan_id", "1"),
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "resources.0.free_count", "253"),
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "resources.0.used_count", "0"),
				),
			},
		},
	})
}

func testAccDevicesApplianceDhcpSubnetsPreCheck(t *testing.T) {
	if v := os.Getenv("TF_ACC_MERAKI_MX_SERIAL"); v == "" {
		t.Fatal("TF_ACC_MERAKI_DEVICE_SERIAL must be set for acceptance tests")
	}
	if v := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"); v == "" {
		t.Fatal("TF_ACC_MERAKI_ORGANIZATION_ID must be set for acceptance tests")
	}
}

func testAccDevicesApplianceDhcpSubnetsDataSourceConfigClaim(serial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
    ]
}	
`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_devices_appliance_dhcp_subnets"),
		serial)
}

func DevicesApplianceDhcpSubnetsDataSourceConfigVlanSettings(serial string) string {
	return fmt.Sprintf(`
	%s

resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
    ]
}	

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	network_id = meraki_network.test.network_id
	vlans_enabled = true
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_devices_appliance_dhcp_subnets"),
		serial,
	)
}

func testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead(serial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
    ]
}	

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	network_id = meraki_network.test.network_id
	vlans_enabled = true
}

data "meraki_devices_appliance_dhcp_subnets" "test" {
	depends_on = [meraki_networks_appliance_vlans_settings.test]
	serial = "%s"
}

output "dhcp_subnets" {
    value = data.meraki_devices_appliance_dhcp_subnets.test
}
`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_devices_appliance_dhcp_subnets"),
		serial, serial)
}
