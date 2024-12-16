package uplink_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksCellularGatewayUplinkResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksCellularGatewayUplinkResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_cellular_gateway_uplink"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_cellular_gateway_uplink"),
			},

			// Update and Read NetworksCellularGatewayUplink
			{
				Config: NetworksCellularGatewayUplinkResourceConfigUpdate(),
				Check:  NetworksCellularGatewayUplinkResourceConfigUpdateChecks(),
			},
		},
		// ImportState test case.
		/*
		   {
		       ResourceName:      "meraki_networks_cellular_gateway_uplink.test",
		       ImportState:       true,
		       ImportStateVerify: false,
		       ImportStateId:     "1234567890",
		   },
		*/

	})
}

// testAccNetworksCellularGatewayUplinkResourceConfigUpdate is a constant string that defines the configuration for updating a networks_cellularGateway_uplink resource in your tests.
// It depends on both the organization and network resources.
func NetworksCellularGatewayUplinkResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_cellular_gateway_uplink" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
    bandwidth_limits = {
        limit_up = 51200
        limit_down = 51200
    }

}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_cellular_gateway_uplink"),
	)
}

// testAccNetworksCellularGatewayUplinkResourceConfigUpdateChecks returns the aggregated test check functions for the cellular gateway uplink redource
func NetworksCellularGatewayUplinkResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"bandwidth_limits.limit_up":   "51200",
		"bandwidth_limits.limit_down": "51200",
	}
	return utils.ResourceTestCheck("meraki_networks_cellular_gateway_uplink.test", expectedAttrs)
}
