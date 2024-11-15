package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNetworksDevicesClaimResource function is used to test the CRUD operations of the Terraform resource you are developing.
// It runs the test cases in order to create, read, update, and delete the resource and checks the state at each step.
func TestAccNetworksDevicesClaimResource(t *testing.T) {

	claimDevices := []string{
		os.Getenv("TF_ACC_MERAKI_MR_SERIAL"),
		os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
		os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
	}

	updateDevices := []string{
		os.Getenv("TF_ACC_MERAKI_MR_SERIAL"),
		os.Getenv("TF_ACC_MERAKI_MG_SERIAL"),
		os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
		os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
	}

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

			// Claim Device serials
			{
				Config: DevicesClaimResourceConfigDeviceClaimWithSerials(claimDevices),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.#", "3"),
					testCheckSerialsUnordered("meraki_networks_devices_claim.test", claimDevices),
				),
			},

			// Update and Claim additional Device serial
			{
				Config: DevicesClaimResourceConfigDeviceClaimWithSerials(updateDevices),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.#", "4"),
					testCheckSerialsUnordered("meraki_networks_devices_claim.test", updateDevices),
				),
			},

			// Import Test
			{
				ResourceName:      "meraki_networks_devices_claim.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Custom test check function to verify serials presence regardless of order
func testCheckSerialsUnordered(resourceName string, expectedSerials []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		serialCountAttr, ok := rs.Primary.Attributes["serials.#"]
		if !ok {
			return fmt.Errorf("serials attribute not found")
		}

		serialCount, err := strconv.Atoi(serialCountAttr)
		if err != nil {
			return fmt.Errorf("failed to parse serial count: %s", err)
		}

		if serialCount != len(expectedSerials) {
			return fmt.Errorf("expected %d serials but found %d", len(expectedSerials), serialCount)
		}

		actualSerials := make([]string, serialCount)
		for i := 0; i < serialCount; i++ {
			actualSerials[i] = rs.Primary.Attributes[fmt.Sprintf("serials.%d", i)]
		}

		sort.Strings(actualSerials)
		sort.Strings(expectedSerials)

		if !reflect.DeepEqual(expectedSerials, actualSerials) {
			return fmt.Errorf("expected serials %v but found %v", expectedSerials, actualSerials)
		}

		return nil
	}
}

func DevicesClaimResourceConfigDeviceClaimWithSerials(serials []string) string {
	serialsFormatted := ""
	for _, serial := range serials {
		serialsFormatted += fmt.Sprintf("\"%s\", ", serial)
	}

	// Remove the trailing comma and space
	serialsFormatted = serialsFormatted[:len(serialsFormatted)-2]

	return fmt.Sprintf(`
	%s

resource "meraki_networks_devices_claim" "test" {
    depends_on = ["resource.meraki_network.test"]
    network_id = resource.meraki_network.test.network_id
    serials = [%s]
}

`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "nw name"), serialsFormatted)
}
