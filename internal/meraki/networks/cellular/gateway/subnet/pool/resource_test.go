package pool_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksCellularGatewaySubnetPoolResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksCellularGatewaySubnetPoolResource(t *testing.T) {

	// The resource.Test function is used to run the test cases.
	resource.Test(t, resource.TestCase{
		// PreCheck function is used to do the necessary setup before running the test cases.
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_cellular_gateway_subnet_pool"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_cellular_gateway_subnet_pool"),
			},

			{
				Config: NetworksCellularGatewaySubnetPoolResourceConfigCreate(),
				Check:  NetworksCellularGatewaySubnetPoolResourceConfigCreateChecks(),
			},

			{
				Config: NetworksCellularGatewaySubnetPoolResourceConfigUpdate(),
				Check:  NetworksCellularGatewaySubnetPoolResourceConfigUpdateChecks(),
			},

			//// ImportState test case.
			//{
			//	ResourceName:      "meraki_networks_cellular_gateway_subnet_pool.test",
			//	ImportState:       true,
			//	ImportStateVerify: false,
			//	ImportStateId:     "1234567890",
			//},
		},
	})
}

// NetworksCellularGatewaySubnetPoolResourceConfigCreate is a constant string that defines the configuration for updating a networks_cellularGateway_subnetPool resource in your tests.
// It depends on both the organization and network resources.
func NetworksCellularGatewaySubnetPoolResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s

	resource "meraki_networks_cellular_gateway_subnet_pool" "test" {
		depends_on = [resource.meraki_network.test]
		id = resource.meraki_network.test.network_id
		cidr = "192.168.0.0/22"
		mask = 24    
	}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_cellular_gateway_subnet_pool"),
	)
}

// NetworksCellularGatewaySubnetPoolResourceConfigCreateChecks returns the aggregated test check functions for the cellular gateway subnet pool resource
func NetworksCellularGatewaySubnetPoolResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"deployment_mode": "routed",
		"cidr":            "192.168.0.0/22",
		"mask":            "24",
	}
	return utils.ResourceTestCheck("meraki_networks_cellular_gateway_subnet_pool.test", expectedAttrs)
}

// NetworksCellularGatewaySubnetPoolResourceConfigUpdate is a constant string that defines the configuration for updating a networks_cellularGateway_subnetPool resource in your tests.
// It depends on both the organization and network resources.
func NetworksCellularGatewaySubnetPoolResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s

	resource "meraki_networks_cellular_gateway_subnet_pool" "test" {
		depends_on = [resource.meraki_network.test]
		id = resource.meraki_network.test.network_id
		cidr = "10.0.0.0/22"
		mask = 24    
	}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_cellular_gateway_subnet_pool"),
	)
}

// NetworksCellularGatewaySubnetPoolResourceConfigUpdateChecks returns the aggregated test check functions for the cellular gateway subnet pool resource
func NetworksCellularGatewaySubnetPoolResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"deployment_mode": "routed",
		"cidr":            "10.0.0.0/22",
		"mask":            "24",
	}
	return utils.ResourceTestCheck("meraki_networks_cellular_gateway_subnet_pool.test", expectedAttrs)
}
