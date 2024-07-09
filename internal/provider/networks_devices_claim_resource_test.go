package provider

import (
	"fmt"
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		// Steps is a slice of TestStep where each TestStep represents a test case.
		Steps: []resource.TestStep{

			// create network
			{
				Config: testAccDevicesClaimResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_claim_device"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Claim Device serials
			{
				Config: testAccDevicesClaimResourceConfigDeviceClaimWithSerials(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), claimDevices),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_devices_claim.test", "serials.#", "3"),
					testCheckSerialsUnordered("meraki_networks_devices_claim.test", claimDevices),
				),
			},

			// Update and Claim additional Device serial
			{
				Config: testAccDevicesClaimResourceConfigDeviceClaimWithSerials(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), updateDevices),
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

func testAccDevicesClaimResourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	name = "test_acc_network_claim_device"
	product_types = ["wireless", "switch", "appliance"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
}

func testAccDevicesClaimResourceConfigDeviceClaimWithSerials(orgId string, serials []string) string {
	serialsFormatted := ""
	for _, serial := range serials {
		serialsFormatted += fmt.Sprintf("\"%s\", ", serial)
	}

	// Remove the trailing comma and space
	serialsFormatted = serialsFormatted[:len(serialsFormatted)-2]

	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
 	name = "test_acc_network_claim_device"
	product_types = ["wireless", "switch", "appliance", "cellularGateway"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}

resource "meraki_networks_devices_claim" "test" {
    depends_on = ["resource.meraki_network.test"]
    network_id = resource.meraki_network.test.network_id
    serials = [%s]
}

`, orgId, serialsFormatted)
}
