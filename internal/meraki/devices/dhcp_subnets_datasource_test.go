package devices_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccDevicesApplianceDhcpSubnetsDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccDevicesApplianceDhcpSubnetsDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
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

			// Claim and Read NetworksDevicesClaim
			{
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigClaim(os.Getenv("TF_ACC_MERAKI_MX_SERIAL"), os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},

			// Update Network VLAN Settings
			{
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigVlanSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_devices_appliance_dhcp_subnets"),
				),
			},

			// Update and Read DevicesApplianceDhcpSubnets
			{
				Config: testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_devices_appliance_dhcp_subnets.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				),
			},
		},

		// The resource.Test function automatically tests the Delete operation.
	})
}

// testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead is a constant string that defines the configuration for creating and updating a devices__appliance_dhcp_subnets resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesApplianceDhcpSubnetsDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types = ["appliance"]
	tags = ["tag1"]
	name = "test_acc_devices_appliance_dhcp_subnets"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccDevicesApplianceDhcpSubnetsDataSourceConfigClaim is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesApplianceDhcpSubnetsDataSourceConfigClaim(serial string, orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
        organization_id = "%s"
        product_types = ["appliance"]
}    
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}	
`, orgId, serial)
	return result
}

// testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead is a constant string that defines the configuration for creating and updating a devices__appliance_dhcp_subnets resource in your tests.
// It depends on both the organization and network resources.
const testAccDevicesApplianceDhcpSubnetsDataSourceConfigVlanSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	vlans_enabled = true
}
`

// testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead is a constant string that defines the configuration for creating and updating a devices__appliance_dhcp_subnets resource in your tests.
// It depends on both the organization and network resources.
func testAccDevicesApplianceDhcpSubnetsDataSourceConfigRead(serialID string) string {
	return fmt.Sprintf(
		`
resource "meraki_network" "test" {
	product_types = ["appliance"]
}

resource "meraki_networks_appliance_vlans_settings" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	vlans_enabled = true
}

data "meraki_devices_appliance_dhcp_subnets" "test" {
	depends_on = [meraki_networks_appliance_vlans_settings.test]
	serial = "%s"
}

`, serialID)
}
