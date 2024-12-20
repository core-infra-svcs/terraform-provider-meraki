package device_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/device"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"testing"
)

func TestAccResource(t *testing.T) {
	t.Run("Validate Schema-Model Consistency", func(t *testing.T) {
		validateResourceSchemaModelConsistency(t)
	})

	t.Run("Test Resource Lifecycle", func(t *testing.T) {
		testResourceLifecycle(t)
	})
}

// Validate schema-model consistency for the resource
func validateResourceSchemaModelConsistency(t *testing.T) {
	testutils.ValidateResourceSchemaModelConsistency(
		t, device.GetResourceSchema.Attributes, device.ResourceModel{},
	)
}

// Test the full resource lifecycle
func testResourceLifecycle(t *testing.T) {
	mxSerial, msSerial, mrSerial, orgId := getResourceTestEnvironmentVariables(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create network
			{
				Config: createNetworkConfigResource(orgId),
				Check:  networkCheckResource("test_acc_device"),
			},
			// Step 2: Claim devices
			{
				Config: claimDevicesConfigResource(mxSerial, msSerial, mrSerial),
				Check:  claimDevicesCheckResource(),
			},
			// Step 3: Update devices
			{
				Config: updateDeviceConfigResource(mxSerial, msSerial, mrSerial),
				Check:  updateDeviceCheckResource(),
			},
			// Step 4: Import and validate state
			{
				ResourceName:      "meraki_devices.test_mx",
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_devices.test_mx", "serial", mxSerial),
				),
			},
		},
	})
}

// Retrieve required environment variables for tests
func getResourceTestEnvironmentVariables(t *testing.T) (string, string, string, string) {
	mxSerial := os.Getenv("TF_ACC_MERAKI_MX_SERIAL")
	msSerial := os.Getenv("TF_ACC_MERAKI_MS_SERIAL")
	mrSerial := os.Getenv("TF_ACC_MERAKI_MR_SERIAL")
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")

	if mxSerial == "" || msSerial == "" || mrSerial == "" || orgId == "" {
		t.Fatal("Environment variables TF_ACC_MERAKI_MX_SERIAL, TF_ACC_MERAKI_MS_SERIAL, TF_ACC_MERAKI_MR_SERIAL, and TF_ACC_MERAKI_ORGANIZATION_ID must be set for acceptance tests")
	}

	return mxSerial, msSerial, mrSerial, orgId
}

// Create network configuration
func createNetworkConfigResource(orgId string) string {
	return utils.CreateNetworkOrgIdConfig(orgId, "test_acc_device")
}

// Validate network creation
func networkCheckResource(networkName string) resource.TestCheckFunc {
	return utils.NetworkOrgIdTestChecks(networkName)
}

// Claim devices configuration
func claimDevicesConfigResource(mxSerial, msSerial, mrSerial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = ["%s", "%s", "%s"]
}
`, createNetworkConfigResource(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")), mxSerial, msSerial, mrSerial)
}

// Validate device claiming
func claimDevicesCheckResource() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources["meraki_networks_devices_claim.test"]
		if !ok {
			return fmt.Errorf("resource not found in state")
		}

		serialCountStr := rs.Primary.Attributes["serials.#"]
		if serialCountStr == "" {
			return fmt.Errorf("serial count not found in state")
		}

		serialCount, err := strconv.Atoi(serialCountStr)
		if err != nil {
			return fmt.Errorf("invalid serial count: %s", serialCountStr)
		}

		actualSerials := make(map[string]bool)
		for i := 0; i < serialCount; i++ {
			serialKey := fmt.Sprintf("serials.%d", i)
			actualSerial := rs.Primary.Attributes[serialKey]
			if actualSerial == "" {
				return fmt.Errorf("serial %d not found in state", i)
			}
			actualSerials[actualSerial] = true
		}

		expectedSerials := []string{
			os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
			os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
			os.Getenv("TF_ACC_MERAKI_MR_SERIAL"),
		}

		for _, expected := range expectedSerials {
			if !actualSerials[expected] {
				return fmt.Errorf("expected serial %s not found in actual serials: %+v", expected, actualSerials)
			}
		}

		return nil
	}
}

// Update devices configuration
func updateDeviceConfigResource(mxSerial, msSerial, mrSerial string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
	product_types   = ["wireless", "switch", "appliance"]
}

resource "meraki_networks_devices_claim" "test" {
	depends_on = [meraki_network.test]
	network_id = meraki_network.test.network_id
	serials    = ["%s", "%s", "%s"]
}

resource "meraki_devices" "test_mx" {
	depends_on = [meraki_networks_devices_claim.test]
	network_id = meraki_network.test.network_id
	serial     = "%s"
	name       = "Updated MX Device"
	tags       = ["updated"]
	address    = "123 MX Street"
}

resource "meraki_devices" "test_ms" {
	depends_on = [meraki_networks_devices_claim.test]
	network_id = meraki_network.test.network_id
	serial     = "%s"
	name       = "Updated MS Device"
	tags       = ["updated"]
	address    = "123 MS Street"
}

resource "meraki_devices" "test_mr" {
	depends_on = [meraki_networks_devices_claim.test]
	network_id = meraki_network.test.network_id
	serial     = "%s"
	name       = "Updated MR Device"
	tags       = ["updated"]
	address    = "123 MR Street"
}
`, os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), mxSerial, msSerial, mrSerial, mxSerial, msSerial, mrSerial)
}

// Validate device updates
func updateDeviceCheckResource() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_devices.test_mx", "name", "Updated MX Device"),
		resource.TestCheckResourceAttr("meraki_devices.test_mx", "tags.0", "updated"),
		resource.TestCheckResourceAttr("meraki_devices.test_mx", "address", "123 MX Street"),

		resource.TestCheckResourceAttr("meraki_devices.test_ms", "name", "Updated MS Device"),
		resource.TestCheckResourceAttr("meraki_devices.test_ms", "tags.0", "updated"),
		resource.TestCheckResourceAttr("meraki_devices.test_ms", "address", "123 MS Street"),

		resource.TestCheckResourceAttr("meraki_devices.test_mr", "name", "Updated MR Device"),
		resource.TestCheckResourceAttr("meraki_devices.test_mr", "tags.0", "updated"),
		resource.TestCheckResourceAttr("meraki_devices.test_mr", "address", "123 MR Street"),
	)
}
