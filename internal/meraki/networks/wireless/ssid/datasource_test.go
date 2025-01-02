package ssid_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// TestAccNetworksWirelessSsidsDataSource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksWirelessSsidsDataSource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			//// Create and Read a Network.
			//{
			//	Config: testAccNetworksWirelessSsidsDataSourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_wireless_ssids"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
			//		resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
			//	),
			//},

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_wireless_ssids_resource"),
			},

			// Create and Read SSID without encryption
			{
				Config: NetworksWirelessSsidsResourceConfigBasic(false),
				Check:  NetworksWirelessSsidsResourceConfigBasicChecks(),
			},

			{
				Config: NetworksWirelessSsidsDataSourceConfigRead(),
				Check:  NetworksWirelessSsidsDataSourceConfigReadChecks(),
			},
		},
	})
}

// testAccNetworksWirelessSsidsDataSourceConfigCreateNetwork is a constant string that defines the configuration for creating a network resource in your tests.
// It depends on the organization resource.
func testAccNetworksWirelessSsidsDataSourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = %s
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_wireless_ssids"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
	return result
}

// testAccNetworksWirelessSsidsDataSourceConfigCreate is a constant string that defines the configuration for creating and updating a networks__test_acc_networks_wireless_ssids resource in your tests.
// It depends on both the organization and network resources.
const testAccNetworksWirelessSsidsDataSourceConfigCreate = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

data "meraki_networks_wireless_ssids" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
}
`

func NetworksWirelessSsidsDataSourceConfigRead() string {
	return fmt.Sprintf(`
	%s
data "meraki_networks_wireless_ssids" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_wireless_ssids_resource"),
	)
}

// NetworksWirelessSsidsDataSourceConfigReadChecks returns the test check functions for NetworksWirelessSsidsDataSourceConfigRead
func NetworksWirelessSsidsDataSourceConfigReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.0.number":    "0",
		"list.0.name":      "My SSID",
		"list.0.enabled":   "true",
		"list.0.auth_mode": "psk",
	}
	return utils.ResourceTestCheck("data.meraki_networks_wireless_ssids.test", expectedAttrs)
}
