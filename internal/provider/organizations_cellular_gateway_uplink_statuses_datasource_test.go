package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccOrganizationsCellularGatewayUplinkStatusesDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccOrganizationsCellularGatewayUplinkStatusesDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read a Network.
			{
				Config: testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "cellularGateway"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim and Read NetworksDevicesClaim
			{
				Config: testAccOrganizationsCellularGatewayUplinkStatusesNetworksDevicesClaimResourceConfigCreate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "id", "example-id"),
				),
			},

			// Read OrganizationsCellularGatewayUplinkStatuses
			{
				Config: testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigRead(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_organizations_cellular_gateway_uplink_statuses.test", "id", "example-id"),
					resource.TestCheckResourceAttr("data.meraki_organizations_cellular_gateway_uplink_statuses.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_organizations_cellular_gateway_uplink_statuses.test", "list.0.serial", os.Getenv("TF_ACC_MERAKI_MG_SERIAL")),
					resource.TestCheckResourceAttr("data.meraki_organizations_cellular_gateway_uplink_statuses.test", "list.0.model", "MG21-NA"),
				),
			},
		},
	})
}

// testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {	
	organization_id = "%s"
	product_types = ["cellularGateway"]
	tags = ["tag1"]
	name = "test_acc_network"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccNetworksDevicesClaimResourceConfigCreate is a constant string that defines the configuration for creating and reading a networks_devices_claim resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsCellularGatewayUplinkStatusesNetworksDevicesClaimResourceConfigCreate(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	    organization_id = "%s"
        product_types = ["cellularGateway"]
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

// testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigRead is a constant string that defines the configuration for creating and updating a organizations_cellularGateway_uplink_statuses resource in your tests.
// It depends on both the organization and network resources.
func testAccOrganizationsCellularGatewayUplinkStatusesDataSourceConfigRead(orgId string, serial string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	    organization_id = "%s"
        product_types = ["cellularGateway"]
}    
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

data "meraki_organizations_cellular_gateway_uplink_statuses" "test" {
	organization_id = "%s"
}	
`, orgId, serial, orgId)
	return result
}
