package subnets_test

import (
	"fmt"
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

	t.Run("Read DHCP Subnets Data Source", func(t *testing.T) {
		testReadApplianceDhcpSubnetsDataSource(t)
	})
}

func testCreateAndReadNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccDevicesApplianceDhcpSubnetsPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_devices_appliance_dhcp_subnets"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
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
					os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"),
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
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigVlanSettings,
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
					resource.TestCheckResourceAttr(
						"data.meraki_devices_appliance_dhcp_subnets.test",
						"serial",
						os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
					),
				),
			},
		},
	})
}

func testReadApplianceDhcpSubnetsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccDevicesApplianceDhcpSubnetsPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplianceDhcpSubnetsDataSourceConfig(os.Getenv("TF_ACC_MERAKI_DEVICE_SERIAL")),
				Check:  testAccDevicesApplianceDhcpSubnetsDataSourceCheck(),
			},
			{
				ResourceName:      "data.meraki_devices_appliance_dhcp_subnets.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDevicesApplianceDhcpSubnetsPreCheck(t *testing.T) {
	if v := os.Getenv("TF_ACC_MERAKI_DEVICE_SERIAL"); v == "" {
		t.Fatal("TF_ACC_MERAKI_DEVICE_SERIAL must be set for acceptance tests")
	}
	if v := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"); v == "" {
		t.Fatal("TF_ACC_MERAKI_ORGANIZATION_ID must be set for acceptance tests")
	}
}

func testAccDevicesApplianceDhcpSubnetsDataSourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_devices_appliance_dhcp_subnets"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
}

func testAccDevicesApplianceDhcpSubnetsDataSourceConfigClaim(serial string, orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
}    
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
    ]
}	
`, orgId, serial)
}

const testAccDevicesApplianceDhcpSubnetsDataSourceConfigVlanSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [meraki_network.test]
	network_id = meraki_network.test.network_id
	vlans_enabled = true
}
`

func testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead(serialID string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	product_types = ["appliance"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [meraki_network.test]
	network_id = meraki_network.test.network_id
	vlans_enabled = true
}

data "meraki_devices_appliance_dhcp_subnets" "test" {
	depends_on = [meraki_networks_appliance_vlans_settings.test]
	serial = "%s"
}
`, serialID)
}

func testAccApplianceDhcpSubnetsDataSourceConfig(serial string) string {
	return fmt.Sprintf(`
data "meraki_devices_appliance_dhcp_subnets" "test" {
  serial = "%s"
}
`, serial)
}

func testAccDevicesApplianceDhcpSubnetsDataSourceCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.meraki_devices_appliance_dhcp_subnets.test",
			"serial",
			os.Getenv("TF_ACC_MERAKI_DEVICE_SERIAL"),
		),
		resource.TestCheckResourceAttrSet("data.meraki_devices_appliance_dhcp_subnets.test", "id"),
		resource.TestCheckResourceAttrSet("data.meraki_devices_appliance_dhcp_subnets.test", "list.#"),
	)
}
