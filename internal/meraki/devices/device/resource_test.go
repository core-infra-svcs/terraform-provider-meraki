package device_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/device"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccResource runs schema validation and resource lifecycle tests.
func TestAccResource(t *testing.T) {
	t.Run("Validate Schema-Model Consistency", func(t *testing.T) {
		validateDeviceResourceSchemaModelConsistency(t)
	})

	t.Run("Test Resource Lifecycle", func(t *testing.T) {
		testDeviceResourceLifecycle(t)
	})
}

// Validate schema-model consistency for the device resource
func validateDeviceResourceSchemaModelConsistency(t *testing.T) {
	testutils.ValidateResourceSchemaModelConsistency(
		t, device.GetResourceSchema.Attributes, device.ResourceModel{},
	)
}

// Test the full lifecycle of the device resource
func testDeviceResourceLifecycle(t *testing.T) {
	mxSerial, orgId := getDeviceResourceTestEnvVars(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createNetworkConfigResource(orgId),
				Check:  networkCheckResource("test_acc_device"),
			},
			{
				Config: claimDeviceConfigResource(mxSerial),
				Check:  claimDeviceCheckResource(),
			},
			{
				Config: updateDeviceConfigResource(mxSerial),
				Check:  updateDeviceCheckResource(),
			},
		},
	})
}

// Retrieve required environment variables for tests
func getDeviceResourceTestEnvVars(t *testing.T) (string, string) {
	mxSerial := os.Getenv("TF_ACC_MERAKI_MX_SERIAL")
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")

	if mxSerial == "" || orgId == "" {
		t.Fatal("Environment variables TF_ACC_MERAKI_MX_SERIAL and TF_ACC_MERAKI_ORGANIZATION_ID must be set for acceptance tests")
	}

	return mxSerial, orgId
}

// Create network configuration
func createNetworkConfigResource(orgId string) string {
	return utils.CreateNetworkOrgIdConfig(orgId, "test_acc_device")
}

// Validate network creation
func networkCheckResource(networkName string) resource.TestCheckFunc {
	return utils.NetworkOrgIdTestChecks(networkName)
}

// Claim device configuration
func claimDeviceConfigResource(mxSerial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_networks_devices_claim" "test_claim" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = ["%s"]
}
`, createNetworkConfigResource(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")), mxSerial)
}

// Validate device claim
func claimDeviceCheckResource() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources["meraki_networks_devices_claim.test_claim"]
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

		if serialCount != 1 {
			return fmt.Errorf("expected 1 serial, got %d", serialCount)
		}

		expectedSerial := os.Getenv("TF_ACC_MERAKI_MX_SERIAL")
		actualSerial := rs.Primary.Attributes["serials.0"]

		if actualSerial != expectedSerial {
			return fmt.Errorf("expected serial %s, got %s", expectedSerial, actualSerial)
		}

		return nil
	}
}

// Update device configuration with device claim
func updateDeviceConfigResource(mxSerial string) string {
	return fmt.Sprintf(`
%s

resource "meraki_networks_devices_claim" "test_claim" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials    = ["%s"]
}

resource "meraki_devices" "test_device" {
    depends_on = [meraki_networks_devices_claim.test_claim]
    serial = "%s"
    name   = "test_acc_device"
    tags   = ["test"]
    address = "test"
    notes   = "test"
}
`, createNetworkConfigResource(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")), mxSerial, mxSerial)
}

// Validate device update
func updateDeviceCheckResource() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_devices.test_device", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
		resource.TestCheckResourceAttr("meraki_devices.test_device", "name", "test_acc_device"),
		resource.TestCheckResourceAttr("meraki_devices.test_device", "tags.#", "1"),
		resource.TestCheckResourceAttr("meraki_devices.test_device", "address", "test"),
		resource.TestCheckResourceAttr("meraki_devices.test_device", "notes", "test"),
	)
}
